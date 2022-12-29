package main

import (
	"github.com/bwmarrin/discordgo"
)

func commandSchedule(m *discordgo.MessageCreate, args []string) {
	// Create the user in the database if it does not already exist.
	var user *User
	if v, err := userGet(m.Author); err != nil {
		msg := "Failed to get the user from the database: " + err.Error()
		log.Error(msg)
		discordSend(m.ChannelID, msg)
		return
	} else {
		user = v
	}

	timezone := user.GetTimezone()

	var channelIDs []string
	if v, err := modals.Races.GetAllScheduled(); err != nil {
		log.Fatal("Failed to get the scheduled races: " + err.Error())
		return
	} else {
		channelIDs = v
	}

	if len(channelIDs) == 0 {
		msg := "There are no races currently scheduled for this week."
		discordSend(m.ChannelID, msg)
		return
	}

	msg := ""
	for _, channelID := range channelIDs {
		var race *Race
		if v, err := getRace(channelID); err != nil {
			msg := "Failed to get the race from the database: " + err.Error()
			log.Error(msg)
			discordSend(m.ChannelID, msg)
			return
		} else {
			race = v
		}

		msg += "**" + getDate(race.DatetimeScheduled.Time, timezone) + "**\n"
		msg += matchGetDescription(race) + "\n"
	}

	discordSend(m.ChannelID, msg)
}
