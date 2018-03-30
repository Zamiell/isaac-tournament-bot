package main

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// Valid rulesets are "seeded" and "unseeded"
type Tournament struct {
	Name              string
	ChallongeURL      string
	ChallongeID       float64
	Ruleset           string
	DiscordCategoryID string
}

var (
	challongeUsername string
	challongeAPIKey   string
	tournaments       = make(map[string]Tournament)

	// We don't want to use the default http.Client structure because it has no default timeout set
	myHTTPClient = &http.Client{
		Timeout: 10 * time.Second,
	}
)

func challongeInit() {
	// Read the Challonge configuration from the environment variables
	challongeUsername = os.Getenv("CHALLONGE_USERNAME")
	if len(challongeUsername) == 0 {
		log.Fatal("The \"CHALLONGE_USERNAME\" environment variable is blank. Set it in the \".env\" file.")
		return
	}

	challongeAPIKey = os.Getenv("CHALLONGE_API_KEY")
	if len(challongeAPIKey) == 0 {
		log.Fatal("The \"CHALLONGE_API_KEY\" environment variable is blank. Set it in the \".env\" file.")
		return
	}

	tournamentURLsString := os.Getenv("TOURNAMENT_CHALLONGE_URLS")
	if len(tournamentURLsString) == 0 {
		log.Fatal("The \"TOURNAMENT_CHALLONGE_URLS\" environment variable is blank. Set it in the \".env\" file.")
		return
	}
	tournamentURLs := strings.Split(tournamentURLsString, ",")

	tournamentRulesetsString := os.Getenv("TOURNAMENT_RULESETS")
	if len(tournamentRulesetsString) == 0 {
		log.Fatal("The \"TOURNAMENT_RULESETS\" environment variable is blank. Set it in the \".env\" file.")
		return
	}
	tournamentRulesets := strings.Split(tournamentRulesetsString, ",")

	tournamentDiscordCategoryIDsString := os.Getenv("TOURNAMENT_DISCORD_CATEGORY_IDS")
	if len(tournamentDiscordCategoryIDsString) == 0 {
		log.Fatal("The \"TOURNAMENT_DISCORD_CATEGORY_IDS\" environment variable is blank. Set it in the \".env\" file.")
		return
	}
	tournamentDiscordCategoryIDs := strings.Split(tournamentDiscordCategoryIDsString, ",")

	// Get all of the Challonge user's tournaments
	apiURL := "https://api.challonge.com/v1/tournaments.json?"
	apiURL += "api_key=" + challongeAPIKey
	var raw []byte
	if v, err := challongeGetJSON("GET", apiURL, nil); err != nil {
		log.Fatal("Failed to get the tournament from Challonge:", err)
		return
	} else {
		raw = v
	}

	jsonTournaments := make([]interface{}, 0)
	if err := json.Unmarshal(raw, &jsonTournaments); err != nil {
		log.Fatal("Failed to unmarshal the Challonge JSON:", err)
	}

	// Figure out the ID for all the tournaments listed in the environment variable
	for i, tournamentURL := range tournamentURLs {
		found := false
		for _, v := range jsonTournaments {
			vMap := v.(map[string]interface{})
			jsonTournament := vMap["tournament"].(map[string]interface{})
			// log.Info("We are an admin of Challonge tournament: " + jsonTournament["name"].(string))
			if jsonTournament["url"] == tournamentURL {
				found = true
				tournaments[tournamentURLs[i]] = Tournament{
					Name:              jsonTournament["name"].(string),
					ChallongeURL:      tournamentURL,
					ChallongeID:       jsonTournament["id"].(float64),
					Ruleset:           tournamentRulesets[i],
					DiscordCategoryID: tournamentDiscordCategoryIDs[i],
				}
				break
			}
		}
		if !found {
			log.Fatal("Failed to find the \"" + tournamentURL + "\" tournament in this Challonge user's tournament list.")
		}
	}
}

func challongeGetJSON(method string, apiURL string, data io.Reader) ([]byte, error) {
	//log.Info("Making a "+method+" request to Challonge:", apiURL) // Uncomment when debugging

	var req *http.Request
	if v, err := http.NewRequest(method, apiURL, data); err != nil {
		return nil, err
	} else {
		req = v
	}

	var resp *http.Response
	if v, err := myHTTPClient.Do(req); err != nil {
		return nil, err
	} else {
		resp = v
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("Bad return status: " + strconv.Itoa(resp.StatusCode))
	}

	var raw []byte
	if v, err := ioutil.ReadAll(resp.Body); err != nil {
		return nil, err
	} else {
		raw = v
	}

	return raw, nil
}

func challongeGetParticipantName(tournament map[string]interface{}, participantID float64) string {
	// Go through all of the participants in this tournament
	for _, v := range tournament["participants"].([]interface{}) {
		vMap := v.(map[string]interface{})
		participant := vMap["participant"].(map[string]interface{})
		if participant["id"] == participantID {
			return participant["name"].(string)
		}
	}

	return "Unknown"
}
