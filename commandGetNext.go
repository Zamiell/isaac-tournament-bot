package main

import (
	"database/sql"

	"github.com/bwmarrin/discordgo"
)

func commandGetNext(m *discordgo.MessageCreate, args []string) {
	// Create the user in the database if it does not already exist
	var user *User
	if v, err := userGet(m.Author); err != nil {
		msg := "Failed to get the user from the database: " + err.Error()
		log.Error(msg)
		discordSend(m.ChannelID, msg)
		return
	} else {
		user = v
	}

	var channelID string
	if v, err := modals.Races.GetNext(); err == sql.ErrNoRows {
		msg := "There are no races currently scheduled for this week."
		discordSend(m.ChannelID, msg)
		return
	} else if err != nil {
		msg := "Failed to get the next race from the database: " + err.Error()
		log.Error(msg)
		discordSend(m.ChannelID, msg)
		return
	} else {
		channelID = v
	}

	var race *Race
	if v, err := raceGet(channelID); err != nil {
		msg := "Failed to get the race from the database: " + err.Error()
		log.Error(msg)
		discordSend(m.ChannelID, msg)
		return
	} else {
		race = v
	}

	var timezone string
	if user.Timezone.Valid {
		timezone = user.Timezone.String
	} else {
		timezone = "UTC"
	}
	msg := "The next scheduled match is on:\n"
	msg += getDate(race.DatetimeScheduled.Time, timezone) + "\n"
	msg += matchGetDescription(race)
	discordSend(m.ChannelID, msg)
}
