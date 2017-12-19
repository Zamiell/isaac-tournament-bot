package main

import (
	"github.com/Zamiell/isaac-tournament-bot/src/models"
	"github.com/bwmarrin/discordgo"
)

func commandReschedule(m *discordgo.MessageCreate, args []string) {
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
		discordSend(m.ChannelID, "Only \""+race.Racer1.Username+"\" and \""+race.Racer2.Username+"\" can reschedule this match.")
		return
	}

	if !race.DatetimeScheduled.Valid {
		discordSend(m.ChannelID, "There is not a time scheduled for this match yet, so there is no need to reschedule.")
		return
	}

	// Set the scheduled time to null
	if err := db.Races.UnsetDatetimeScheduled(m.ChannelID); err != nil {
		msg := "Failed to unset the scheduled time: " + err.Error()
		log.Error(msg)
		discordSend(m.ChannelID, msg)
		return
	}

	discordSend(m.ChannelID, "The currently scheduled time has been deleted. Please suggest a new time with the `!schedule` command.")
}
