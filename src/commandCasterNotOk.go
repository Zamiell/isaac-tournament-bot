package main

import (
	"github.com/Zamiell/isaac-tournament-bot/src/models"
	"github.com/bwmarrin/discordgo"
)

func commandCasterNotOk(m *discordgo.MessageCreate, args []string) {
	var race models.Race
	if v, err := raceGet(m.ChannelID); err != nil {
		msg := "Failed to get the race from the database: " + err.Error()
		log.Error(msg)
		discordSend(m.ChannelID, msg)
		return
	} else {
		race = v
	}

	// Check to see if this person is one of the two racers
	if m.Author.ID != race.Racer1.DiscordID && m.Author.ID != race.Racer2.DiscordID {
		discordSend(m.ChannelID, "You cannot submit disapproval for a match that you are not participanting in.")
		return
	}

	// Check to see if someone is casting this match
	if !race.CasterID.Valid {
		discordSend(m.ChannelID, "No-one has volunteered to cast this match, so there is no need to submit disapproval.")
		return
	}

	// Unset the caster
	if err := db.Races.UnsetCaster(m.ChannelID); err != nil {
		msg := "Failed to unset the caster in the database: " + err.Error()
		log.Error(msg)
		discordSend(m.ChannelID, msg)
		return
	}

	msg := m.Author.Username + " has denied permission for " + race.Caster.Mention() + " to rebroadcast the race. They have been removed as the registered caster for this match."
	discordSend(m.ChannelID, msg)
}
