package main

import (
	"database/sql"

	"github.com/Zamiell/isaac-tournament-bot/src/models"
	"github.com/bwmarrin/discordgo"
)

func commandCastCancel(m *discordgo.MessageCreate, args []string) {
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

	// Check to see if someone is already casting this match
	if !race.CasterID.Valid {
		discordSend(m.ChannelID, "No-one is currently marked as casting this match, so there is no need to cancel anything.")
		return
	}

	// Check to see if this person is one of the two racers
	if m.Author.ID == race.Racer1.DiscordID || m.Author.ID == race.Racer2.DiscordID {
		discordSend(m.ChannelID, "If you don't want "+race.Caster.Username+" to cast your match, use the `!casternotok` command.")
		return
	}

	// Check to see if this person is the one who volunteered
	if m.Author.ID != race.Caster.DiscordID {
		discordSend(m.ChannelID, "Only "+race.Caster.Username+" can cancel their request to cast this match.")
		return
	}

	// Unset the caster
	if err := db.Races.UnsetCaster(m.ChannelID); err != nil {
		msg := "Failed to unset the caster in the database: " + err.Error()
		log.Error(msg)
		discordSend(m.ChannelID, msg)
		return
	}

	msg := "You have been removed as the designated caster for this match."
	discordSend(m.ChannelID, msg)

}
