package main

import (
	"strconv"
	"strings"

	"github.com/Zamiell/isaac-tournament-bot/src/models"
	"github.com/bwmarrin/discordgo"
	timezone "github.com/tkuchiki/go-timezone"
)

func commandTimezone(m *discordgo.MessageCreate, args []string) {
	if len(args) == 0 {
		commandTimezonePrint(m)
		return
	}

	// Create the user in the database if it does not already exist
	if _, err := racerGet(m.Author); err != nil {
		msg := "Failed to get the racer from the database: " + err.Error()
		log.Error(msg)
		discordSend(m.ChannelID, msg)
		return
	}

	// See if the submitted timezone is valid
	newTimezone := strings.ToUpper(args[0])
	var offset int
	if v, err := timezone.GetOffset(newTimezone); err != nil {
		discordSend(m.ChannelID, "That is not a valid timezone. For example, try using \"EST\", \"GMT-5\", or \"GMT+3\".")
		return
	} else {
		offset = v
	}
	hours := offset / -3600

	// Set the new timezone
	if err := db.Racers.SetTimeZone(m.Author.ID, hours); err != nil {
		msg := "Failed to update the timezone: " + err.Error()
		log.Error(msg)
		discordSend(m.ChannelID, msg)
		return
	}

	msg := m.Author.Mention() + ", your timezone has been set to: **GMT"
	if hours < 0 {
		msg += "-"
	} else {
		msg += "+"
	}
	msg += strconv.Itoa(hours) + "**"
	discordSend(m.ChannelID, msg)
}

func commandTimezonePrint(m *discordgo.MessageCreate) {
	var racer models.Racer
	if v, err := racerGet(m.Author); err != nil {
		msg := "Failed to get the racer from the database: " + err.Error()
		log.Error(msg)
		discordSend(m.ChannelID, msg)
		return
	} else {
		racer = v
	}

	msg := m.Author.Mention() + ", your timezone is "
	if racer.Timezone.Valid {
		msg += "currently set to: **" + getTimezone(racer.Timezone.Int64) + "**\n\n"
	} else {
		msg += "**not currently set**.\n\n"
	}
	msg += "Set your timezone with:\n"
	msg += "```\n"
	msg += "!timezone [timezone]\n"
	msg += "```\n"
	msg += "For example:\n"
	msg += "```\n"
	msg += "!timezone EST\n"
	msg += "```\n"
	msg += "Or:\n"
	msg += "```\n"
	msg += "!timezone GMT+3\n"
	msg += "```"
	discordSend(m.ChannelID, msg)
}

func getTimezone(timezone int64) string {
	timezoneString := "GMT"
	if timezone < 0 {
		timezoneString += "-"
	} else {
		timezoneString += "+"
	}
	timezoneString += strconv.Itoa(int(timezone))
	return timezoneString
}
