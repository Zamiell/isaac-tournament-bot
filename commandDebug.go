package main

import (
	"database/sql"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func commandDebug(m *discordgo.MessageCreate, args []string) {
	if !isAdmin(m) {
		return
	}

	/*
		log.Info("Tournaments:")
		for _, tournament := range tournaments { // This is a map indexed by "ChallongeURL"
			log.Info(tournament.Name + " - " + tournament.ChallongeURL + " - " + floatToString(tournament.ChallongeID) + " - " + tournament.Ruleset)
		}
	*/

	// Check to see if this is a race channel (and get the race from the database)
	var race *Race
	if v, err := getRace(m.ChannelID); err == sql.ErrNoRows {
		discordSend(m.ChannelID, "You can only use that command in a race channel.")
		return
	} else if err != nil {
		msg := "Failed to get the race from the database: " + err.Error()
		log.Error(msg)
		discordSend(m.ChannelID, msg)
		return
	} else {
		race = v
	}

	msg := "DEBUG:\n"
	msg += fmt.Sprintf("State: %v\n", race.State)
	msg += fmt.Sprintf("ActiveRacer: %v\n", race.ActiveRacer)
	msg += fmt.Sprintf("Racer1ID: %v\n", race.Racer1ID)
	msg += fmt.Sprintf("Racer2ID: %v\n", race.Racer2ID)
	msg += fmt.Sprintf("Characters: %v\n", race.Characters)
	msg += fmt.Sprintf("CharactersRemaining: %v\n", race.CharactersRemaining)
	msg += fmt.Sprintf("Builds: %v\n", race.Builds)
	msg += fmt.Sprintf("BuildsRemaining: %v\n", race.BuildsRemaining)
	discordSend(m.ChannelID, msg)
}
