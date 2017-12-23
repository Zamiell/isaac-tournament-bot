package main

import (
	"os"
	"strconv"
	"time"

	"github.com/Zamiell/isaac-tournament-bot/src/models"
)

var (
	bansNum int
)

func matchInit() {
	// Read the OAuth secret from the environment variable
	bansNumString := os.Getenv("BANS_NUM")
	if len(bansNumString) == 0 {
		log.Fatal("The \"BANS_NUM\" environment variable is blank. Set it in the \".env\" file.")
		return
	}

	if v, err := strconv.Atoi(bansNumString); err != nil {
		log.Fatal("The \"BANS_NUM\" environment variable is not a number.")
		return
	} else {
		bansNum = v
	}
}

func matchStart(race models.Race) {
	// Sleep until the match starts
	sleepDuration := race.DatetimeScheduled.Time.Sub(time.Now().UTC())
	if sleepDuration < 5*time.Minute {
		sleepDuration = 0
	} else {
		sleepDuration -= 5 * time.Minute
	}

	time.Sleep(sleepDuration)

	// Re-get the race from the database
	if v, err := raceGet(race.ChannelID); err != nil {
		msg := "Failed to re-get the race from the database: " + err.Error()
		log.Error(msg)
		discordSend(race.ChannelID, msg)
		return
	} else {
		race = v
	}

	// Check to see if this match has started already
	if race.State != 1 {
		log.Info("Reached the \"matchStart\" function when the state was " + strconv.Itoa(race.State) + ". Doing nothing.")
		return
	}

	// Set the state to 2
	if err := db.Races.SetState(race.ChannelID, 2); err != nil {
		msg := "Failed to set the state for race \"" + race.Name() + "\": " + err.Error()
		log.Error(msg)
		discordSend(race.ChannelID, msg)
		return
	}

	// Randomly decide who starts
	race.ActivePlayer = getRandom(1, 2)
	if err := db.Races.SetActivePlayer(race.ChannelID, race.ActivePlayer); err != nil {
		msg := "Failed to set the active player for race \"" + race.Name() + "\": " + err.Error()
		log.Error(msg)
		discordSend(race.ChannelID, msg)
		return
	}

	charactersStart(race)
}

func matchBansEnd(race models.Race) {

}
