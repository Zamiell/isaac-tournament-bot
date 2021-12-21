package main

import (
	"database/sql"
	"strconv"

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
	if v, err := raceGet(m.ChannelID); err == sql.ErrNoRows {
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
	msg += strconv.Itoa(len(race.Characters)) + "\n"
	msg += "A" + sliceToString(race.Characters) + "A\n"
	discordSend(m.ChannelID, msg)
}
