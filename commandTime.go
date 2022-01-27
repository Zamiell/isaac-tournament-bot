package main

import (
	"database/sql"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/kierdavis/dateparser"
)

func commandTime(m *discordgo.MessageCreate, args []string) {
	if len(args) == 0 {
		commandSchedulePrint(m)
		return
	}

	// Check to see if they might be mistakenly trying to do the "!timeok" command
	if len(args) == 1 && args[0] == "ok" {
		commandTimeOk(m, args)
		return
	}

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

	// Check to see if this is a race channel (and get the race from the database)
	var race *Race
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

	// Check to see if this person is one of the two racers
	var activeRacer int
	if m.Author.ID == race.Racer1.DiscordID {
		activeRacer = 1
	} else if m.Author.ID == race.Racer2.DiscordID {
		activeRacer = 2
	} else {
		discordSend(m.ChannelID, "Only \""+race.Racer1.Username+"\" and \""+race.Racer2.Username+"\" can schedule a time for this match.")
		return
	}

	// Check to see if this race has already been scheduled
	if race.State != RaceStateInitial {
		discordSend(m.ChannelID, "The race has already been scheduled. To delete this time and start over, use the `!timedelete` command.")
		return
	}

	// Check to see if this person has a timezone specified
	if !user.Timezone.Valid {
		discordSend(m.ChannelID, "You must specify a timezone with the `!timezone` command before you can suggest a time for the match.")
		return
	}

	// Check to see if this person has a stream specified
	if !user.StreamURL.Valid {
		discordSend(m.ChannelID, "You must specify a stream URL with the `!stream` command before you can suggest a time for the match.")
		return
	}

	// Check to see if this is a valid time
	input := strings.Join(args, " ")
	var datetime time.Time
	if v, err := dateparser.Parse(input); err != nil {
		msg := "Failed to parse the time: " + err.Error()
		discordSend(m.ChannelID, msg)
		return
	} else {
		datetime = v
	}

	// Get the timezone offset for this person
	// https://stackoverflow.com/questions/34975007/in-go-how-can-i-extract-the-value-of-my-current-local-time-offset
	loc, _ := time.LoadLocation(user.Timezone.String)
	t := time.Now().In(loc)
	_, offset := t.Zone()

	// Change the time to correspond to the local time zone
	datetime = datetime.Add(time.Second * time.Duration(offset) * -1)

	// Check to see if it is in the future
	difference := datetime.Sub(time.Now().UTC())
	if difference < 0 {
		discordSend(m.ChannelID, "You must schedule a date in the future.")
		return
	}

	// Set the new scheduled time
	if err := modals.Races.SetDatetimeScheduled(m.ChannelID, datetime, activeRacer); err != nil {
		msg := "Failed to update the scheduled time: " + err.Error()
		log.Error(msg)
		discordSend(m.ChannelID, msg)
		return
	}

	var racer1 *User
	var racer2 *User
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
		now := time.Now()
		_, offset1 := now.In(loc1).Zone()
		_, offset2 := now.In(loc2).Zone()
		if offset1 == offset2 {
			timezonesEqual = true
		}
	}

	msg := racer1.Mention() + " has suggested that the match be scheduled at: *"
	msg += getDate(datetime, racer1.Timezone.String) + "*\n"

	if !timezonesEqual && racer2.Timezone.Valid {
		msg += racer2.Mention() + ", this is equal to: *"
		msg += getDate(datetime, racer2.Timezone.String) + "*\n"
		msg += "If"
	} else {
		msg += racer2.Mention() + ", if"
	}
	msg += " this time is good for you, please use the `!timeok` command. Otherwise, suggest a new time with: `!time [date & time]`"
	discordSend(m.ChannelID, msg)
}

func commandSchedulePrint(m *discordgo.MessageCreate) {
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

	// Check to see if this is a race channel (and get the race from the database)
	var race *Race
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

	msg := ""
	if race.DatetimeScheduled.Valid {
		var timezone string
		if user.Timezone.Valid {
			timezone = user.Timezone.String
		} else {
			timezone = "UTC"
		}
		msg += "The currently scheduled time for the match is: *" + getDate(race.DatetimeScheduled.Time, timezone) + "*\n"
	} else {
		msg += "This match is not scheduled yet.\n"
	}

	if race.State != RaceStateInitial {
		msg += "Both racers have agreed to this time.\n"
		msg += "To delete this time and start over, use the `!timedelete` command."
	} else {
		msg += "You can suggest a new time with: `!time [date & time]`\n"
		msg += "e.g. `!time 6pm sat`"
	}
	discordSend(m.ChannelID, msg)
}
