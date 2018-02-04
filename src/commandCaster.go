package main

import (
	"database/sql"

	"github.com/Zamiell/isaac-tournament-bot/src/models"
	"github.com/bwmarrin/discordgo"
)

func commandCaster(m *discordgo.MessageCreate, args []string) {
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

	if !race.CasterID.Valid {
		discordSend(m.ChannelID, "No-one has volunteered to cast this match yet.")
		return
	}

	if race.CasterP1 && race.CasterP2 {
		msg := race.Caster.Username + " is registered to cast this match at: " + race.Caster.StreamURL.String + "\n"
		msg += "(Both racers have agreed already.)"
		discordSend(m.ChannelID, msg)
	} else {
		msg := race.Caster.Username + " has requested to cast this match at: " + race.Caster.StreamURL.String + "\n"
		if race.CasterP1 {
			msg += race.Racer2.Username + " still needs to okay this with the `!casterok` command."
		} else if race.CasterP2 {
			msg += race.Racer1.Username + " still needs to okay this with the `!casterok` command."
		}
		discordSend(m.ChannelID, msg)
	}
}
