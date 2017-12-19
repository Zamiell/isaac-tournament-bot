package main

import (
	"strings"
	"time"

	"github.com/Zamiell/isaac-tournament-bot/src/models"
	"github.com/araddon/dateparse"
	"github.com/bwmarrin/discordgo"
)

func commandSchedule(m *discordgo.MessageCreate, args []string) {
	if len(args) == 0 {
		commandSchedulePrint(m)
		return
	}

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
		discordSend(m.ChannelID, "Only \""+race.Racer1.Username+"\" and \""+race.Racer2.Username+"\" can schedule a time for this match.")
		return
	}

	// Check to see if this is a valid time
	datetime := strings.Join(args, " ")
	var datetimeScheduled time.Time
	if t, err := dateparse.ParseAny(datetime); err != nil {
		msg := "The datetime of \"" + datetime + "\" is not valid. Please see the following page for some working examples:\n"
		msg += "<https://github.com/araddon/dateparse/blob/master/example/main.go#L12>"
		discordSend(m.ChannelID, msg)
		return
	} else {
		datetimeScheduled = t
	}

	// Set the new scheduled time
	if err := db.Races.SetDatetimeScheduled(m.ChannelID, datetimeScheduled); err != nil {
		msg := "Failed to update the scheduled time: " + err.Error()
		log.Error(msg)
		discordSend(m.ChannelID, msg)
		return
	}

	var racer1 models.Racer
	var racer2 models.Racer
	if m.Author.ID == race.Racer1.DiscordID {
		racer1 = race.Racer1
		racer2 = race.Racer2
	} else if m.Author.ID == race.Racer2.DiscordID {
		racer1 = race.Racer2
		racer2 = race.Racer1
	}

	timezonesEqual := false
	if racer1.Timezone.Valid && racer2.Timezone.Valid && racer1.Timezone.Int64 == racer2.Timezone.Int64 {
		timezonesEqual = true
	}

	timeformatString := "Monday, January 2 @ 15:04"
	msg := racer1.Mention() + " has suggested that the match be scheduled at: **" + datetimeScheduled.Format(timeformatString)
	if racer1.Timezone.Valid && !timezonesEqual {
		msg += " (" + getTimezone(racer1.Timezone.Int64) + ")**\n"
	} else {
		msg += "**\n"
	}
	if racer2.Timezone.Valid && !timezonesEqual {
		msg += racer2.Mention() + ", this is equal to: **" + datetimeScheduled.Format(timeformatString)
		msg += " (" + getTimezone(racer2.Timezone.Int64) + ")**\nIf"
	} else {
		msg += racer2.Mention() + ", if"
	}
	msg += " this time is good for you, please use the `!confirm` command. Otherwise, suggest a new time with: `!schedule [date and time]`"
	discordSend(m.ChannelID, msg)
}

func commandSchedulePrint(m *discordgo.MessageCreate) {
	var race models.Race
	if v, err := raceGet(m.ChannelID); err != nil {
		msg := "Failed to get the race from the database: " + err.Error()
		log.Error(msg)
		discordSend(m.ChannelID, msg)
		return
	} else {
		race = v
	}

	msg := ""
	if race.DatetimeScheduled.Valid {
		msg += "The currently scheduled time is: " + race.DatetimeScheduled.Time.String() + "\n"
	} else {
		msg += "This match is not scheduled yet.\n"
	}

	if race.State > 0 {
		msg += "Both racers have agreed to this time.\n"
		msg += "To delete this time and start over, use the `!reschedule` command."
	} else {
		msg += "You can suggest a new time with:\n"
		msg += "```\n"
		msg += "!schedule [date and time]\n"
		msg += "```\n"
		msg += "For example:\n"
		msg += "```\n"
		msg += "!schedule 02/06/2018 22:00\n"
		msg += "```"
	}
	discordSend(m.ChannelID, msg)
}
