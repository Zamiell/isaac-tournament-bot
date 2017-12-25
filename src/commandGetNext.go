package main

import (
	"github.com/Zamiell/isaac-tournament-bot/src/models"
	"github.com/bwmarrin/discordgo"
)

func commandGetNext(m *discordgo.MessageCreate, args []string) {
	// Create the user in the database if it does not already exist
	var racer models.Racer
	if v, err := racerGet(m.Author); err != nil {
		msg := "Failed to get the racer from the database: " + err.Error()
		log.Error(msg)
		discordSend(m.ChannelID, msg)
		return
	} else {
		racer = v
	}

	var channelID string
	if v, err := db.Races.GetNext(); err != nil {
		msg := "Failed to get the next race from the database: " + err.Error()
		log.Error(msg)
		discordSend(m.ChannelID, msg)
		return
	} else {
		channelID = v
	}

	var race models.Race
	if v, err := raceGet(channelID); err != nil {
		msg := "Failed to get the race from the database: " + err.Error()
		log.Error(msg)
		discordSend(m.ChannelID, msg)
		return
	} else {
		race = v
	}

	msg := "The next scheduled match is: *" + race.Name() + "*\n"
	var timezone string
	if racer.Timezone.Valid {
		timezone = racer.Timezone.String
	} else {
		timezone = "UTC"
	}
	msg += "It is scheduled to begin at: *" + getDate(race.DatetimeScheduled.Time, timezone) + "*"
	discordSend(m.ChannelID, msg)
}
