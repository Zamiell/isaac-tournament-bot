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

	tournamentNamesString := os.Getenv("TOURNAMENT_NAMES")
	if len(tournamentNamesString) == 0 {
		log.Fatal("The \"TOURNAMENT_NAMES\" environment variable is blank. Set it in the \".env\" file.")
		return
	}
	tournamentNames := strings.Split(tournamentNamesString, ",")

	tournamentRulesetsString := os.Getenv("TOURNAMENT_RULESETS")
	if len(tournamentRulesetsString) == 0 {
		log.Fatal("The \"TOURNAMENT_RULESETS\" environment variable is blank. Set it in the \".env\" file.")
		return
	}
	tournamentRulesets := strings.Split(tournamentRulesetsString, ",")

	discordCategoryIDsString := os.Getenv("DISCORD_CATEGORY_IDS")
	if len(discordCategoryIDsString) == 0 {
		log.Fatal("The \"DISCORD_CATEGORY_IDS\" environment variable is blank. Set it in the \".env\" file.")
		return
	}
	discordCategoryIDs := strings.Split(discordCategoryIDsString, ",")

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
	for i, _ := range tournamentNames {
		found := false
		for _, v := range jsonTournaments {
			vMap := v.(map[string]interface{})
			jsonTournament := vMap["tournament"].(map[string]interface{})
			if jsonTournament["url"] == tournamentNames[i] {
				found = true
				tournaments[tournamentNames[i]] = Tournament{
					ChallongeID:       jsonTournament["id"].(float64),
					Ruleset:           tournamentRulesets[i],
					DiscordCategoryID: discordCategoryIDs[i],
				}
				break
			}
		}
		if !found {
			log.Fatal("Failed to find the \"" + tournamentNames[i] + "\" tournament in this Challonge user's tournament list.")
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
