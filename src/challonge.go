package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

var (
	challongeUsername       string
	challongeAPIKey         string
	challongeTournamentName string
	challongeTournamentID   float64

	// We don't want to use the default http.Client structure because it has no default timeout set
	myHTTPClient = &http.Client{
		Timeout: 10 * time.Second,
	}
)

func challongeInit() {
	// Read the Challonge API key from the environment variable
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

	challongeTournamentName = os.Getenv("CHALLONGE_TOURNAMENT_NAME")
	if len(challongeTournamentName) == 0 {
		log.Fatal("The \"CHALLONGE_TOURNAMENT_NAME\" environment variable is blank. Set it in the \".env\" file.")
		return
	}

	// Figure out the ID of this tournament by getting all of the Challonge user's tournaments
	apiURL := "https://api.challonge.com/v1/tournaments.json?"
	apiURL += "api_key=" + challongeAPIKey
	var raw []byte
	if v, err := challongeGetJSON(apiURL); err != nil {
		log.Fatal("Failed to get the tournament from Challonge:", err)
		return
	} else {
		raw = v
	}

	tournaments := make([]interface{}, 0)
	if err := json.Unmarshal(raw, &tournaments); err != nil {
		log.Fatal("Failed to unmarshal the Challonge JSON:", err)
	}

	found := false
	for _, v := range tournaments {
		vMap := v.(map[string]interface{})
		tournament := vMap["tournament"].(map[string]interface{})
		if tournament["url"] == challongeTournamentName {
			found = true
			challongeTournamentID = tournament["id"].(float64)
			break
		}
	}
	if !found {
		log.Fatal("Failed to find the \"" + challongeTournamentName + "\" tournament in this Challonge user's tournament list.")
	}
}

func challongeGetJSON(apiURL string) ([]byte, error) {
	// log.Info("Making a POST request to Challonge:", apiURL) // Uncomment when debugging

	resp, err := myHTTPClient.Get(apiURL)
	if err != nil {
		return nil, err
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
