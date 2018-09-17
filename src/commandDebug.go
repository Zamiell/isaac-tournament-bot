package main

import (
	"database/sql"
	"github.com/bwmarrin/discordgo"

	"github.com/Zamiell/isaac-tournament-bot/src/models"
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
	var race models.Race
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
	if race.Racer1.CasterAlwaysOk {
		msg += "Racer 1 \"" + race.Racer1.Username + "\" has enabled automatic caster approval.\n"
	}
	if race.Racer2.CasterAlwaysOk {
		msg += "Racer 2 \"" + race.Racer2.Username + "\" has enabled automatic caster approval.\n"
	}
	discordSend(m.ChannelID, msg)
}
