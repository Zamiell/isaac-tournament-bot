package main

import (
	"strings"
	"time"

	"github.com/Zamiell/isaac-tournament-bot/src/models"
	"github.com/bwmarrin/discordgo"
	"github.com/kierdavis/dateparser"
)

func commandTime(m *discordgo.MessageCreate, args []string) {
	if len(args) == 0 {
		commandSchedulePrint(m)
		return
	}

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

	// Check to see if this race has already been scheduled
	if race.State != 0 {
		discordSend(m.ChannelID, "The race has already been scheduled. To delete this time and start over, use the `!timedelete` command.")
		return
	}

	// Check to see if this person already has a timezone specified
	if !racer.Timezone.Valid {
		discordSend(m.ChannelID, "You must specify a timezone with the `!timezone` command before you can suggest a time for the match.")
		return
	}

	// Check to see if this is a valid time
	input := strings.Join(args, " ")
	//input += " " + getGMT(racer.Timezone.String)
	input += " PST"
	log.Info("input:", input)
	var datetime time.Time
	if t, err := dateparser.Parse(input); err != nil {
		msg := "Failed to parse the time: " + err.Error()
		discordSend(m.ChannelID, msg)
		return
	} else {
		datetime = t
	}

	// Change it to UTC
	datetimeUTC := datetime.UTC()

	// Set the new scheduled time
	if err := db.Races.SetDatetimeScheduled(m.ChannelID, datetimeUTC); err != nil {
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
	if racer1.Timezone.Valid && racer2.Timezone.Valid {
		// Get the short name of their time zone so that we can add it to their submitted time
		loc1, _ := time.LoadLocation(racer1.Timezone.String)
		loc2, _ := time.LoadLocation(racer2.Timezone.String)
		_, offset1 := time.Now().In(loc1).Zone()
		_, offset2 := time.Now().In(loc2).Zone()
		if offset1 == offset2 {
			timezonesEqual = true
		}
	}

	msg := racer1.Mention() + " has suggested that the match be scheduled at: "
	msg += getDate(datetimeUTC, racer1.Timezone.String) + "\n"

	if !timezonesEqual && racer2.Timezone.Valid {
		msg += racer2.Mention() + ", this is equal to: "
		msg += getDate(datetimeUTC, racer2.Timezone.String) + "\n"
		msg += "If"
	} else {
		msg += racer2.Mention() + ", if"
	}
	msg += " this time is good for you, please use the `!timeok` command. Otherwise, suggest a new time with: `!time [date & time]`"
	discordSend(m.ChannelID, msg)

	// Convert it to their timezone
	dateFormatString := "Monday, January 2 @ 15:04"

	log.Info("Racer \"" + m.Author.Username + "\" suggested the time of: " + datetimeUTC.Format(dateFormatString))
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
		msg += "To delete this time and start over, use the `!timedelete` command."
	} else {
		msg += "You can suggest a new time with: `!time [date & time]`\n"
		msg += "For example: `!time 6pm sat`"
	}
	discordSend(m.ChannelID, msg)
}
