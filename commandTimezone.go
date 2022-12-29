package main

import (
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	timezone "github.com/tkuchiki/go-timezone"
)

func commandTimezone(m *discordgo.MessageCreate, args []string) {
	if len(args) == 0 {
		commandTimezonePrint(m)
		return
	}

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

	// See if the submitted timezone is a short timezone.
	tz := timezone.New()
	newTimezone := args[0]
	if _, err := tz.GetTzAbbreviationInfo(strings.ToUpper(newTimezone)); err == nil {
		newTimezone = strings.ToUpper(newTimezone)
		msg := "That is not specific enough. Please use `!timezone [timezone]` and select from the following list of timezones inside " + newTimezone + ":\n"
		msg += "```\n"
		var timezones []string
		if v, err := tz.GetTimezones(newTimezone); err != nil {
			msg = "Failed to get the list of timezones for " + newTimezone + ": " + err.Error()
			discordSend(m.ChannelID, msg)
			return
		} else {
			timezones = v
		}
		for _, zone := range timezones {
			if strings.HasPrefix(zone, newTimezone) {
				continue
			}
			msg += zone + "\n"
		}
		msg += "```"
		discordSend(m.ChannelID, msg)
		return
	}

	// See if the submitted timezone is valid.
	if _, err := time.LoadLocation(newTimezone); err != nil {
		msg := "That is not a valid timezone. The submitted timezone has to exactly match the TZ column of the following page:\n"
		msg += "<https://en.wikipedia.org/wiki/List_of_tz_database_time_zones>\n"
		msg += "e.g. `!timezone America/New_York`"
		discordSend(m.ChannelID, msg)
		return
	}

	// Set the new timezone.
	if err := modals.Users.SetTimezone(m.Author.ID, newTimezone); err != nil {
		msg := "Failed to update the timezone: " + err.Error()
		log.Error(msg)
		discordSend(m.ChannelID, msg)
		return
	}

	msg := "The timezone for **" + user.Username + "** has been set to: **" + getTimezone(newTimezone) + "**"
	discordSend(m.ChannelID, msg)
}

func commandTimezonePrint(m *discordgo.MessageCreate) {
	var user *User
	if v, err := userGet(m.Author); err != nil {
		msg := "Failed to get the user from the database: " + err.Error()
		log.Error(msg)
		discordSend(m.ChannelID, msg)
		return
	} else {
		user = v
	}

	msg := m.Author.Mention() + ", your timezone is "
	if user.Timezone.Valid {
		msg += "currently set to: **" + getTimezone(user.Timezone.String) + "**\n\n"
	} else {
		msg += "**not currently set**.\n\n"
	}
	msg += "Set your timezone with: `!timezone [timezone]`\n"
	msg += "e.g. `!timezone America/New_York`\n"
	msg += "The submitted timezone has to exactly match the TZ column of the following page:\n"
	msg += "<https://en.wikipedia.org/wiki/List_of_tz_database_time_zones>"
	discordSend(m.ChannelID, msg)
}
