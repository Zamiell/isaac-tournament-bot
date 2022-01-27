package main

import (
	"database/sql"
	"math"
	"time"

	"github.com/bwmarrin/discordgo"
)

func commandStatus(m *discordgo.MessageCreate, args []string) {
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

	announceStatus(m, race, false)
}

func announceStatus(m *discordgo.MessageCreate, race *Race, shouldPing bool) {
	if race.State == RaceStateInitial {
		printStatusInitial(m, race, shouldPing)
	} else if race.State == RaceStateScheduled {
		printStatusScheduled(m, race)
	} else if race.State == RaceStateVetoCharacters {

	} else if race.State == RaceStateVetoBuilds {

	} else if race.State == RaceStateInProgress {

	} else if race.State == RaceStateCompleted {

	}
}

func printStatusInitial(m *discordgo.MessageCreate, race *Race, shouldPing bool) {
	racer1 := race.Racer1
	racer2 := race.Racer2

	var members []*discordgo.Member
	if v, err := getDiscordMembers(); err != nil {
		log.Error(err)
		discordSend(m.ChannelID, err.Error())
		return
	} else {
		members = v
	}

	discordUser1 := getDiscordUserByID(members, racer1.DiscordID)
	discordUser2 := getDiscordUserByID(members, racer2.DiscordID)

	var racer1NameToUse string
	var racer2NameToUse string
	if shouldPing {
		racer1NameToUse = discordUser1.Mention()
		racer2NameToUse = discordUser2.Mention()
	} else {
		racer1NameToUse = racer1.Username
		racer2NameToUse = racer2.Username
	}

	// Find out if the racers have set their timezone
	msg := ""
	if racer1.Timezone.Valid {
		msg += racer1NameToUse + " has a timezone of: " + getTimezone(racer1.Timezone.String) + "\n"
	} else {
		msg += racer1NameToUse + ", your timezone is **not currently set**. Please set one with: `!timezone [timezone]`\n"
	}
	if racer2.Timezone.Valid {
		msg += racer2NameToUse + " has a timezone of: " + getTimezone(racer2.Timezone.String) + "\n"
	} else {
		msg += racer2NameToUse + ", your timezone is **not currently set**. Please set one with: `!timezone [timezone]`\n"
	}

	// Calculate the difference between the two timezones
	if racer1.Timezone.Valid && racer2.Timezone.Valid {
		loc1, _ := time.LoadLocation(racer1.Timezone.String)
		loc2, _ := time.LoadLocation(racer2.Timezone.String)
		_, offset1 := time.Now().In(loc1).Zone()
		_, offset2 := time.Now().In(loc2).Zone()
		if offset1 == offset2 {
			msg += "You both are in **the same timezone**. Great!\n"
		} else {
			difference := math.Abs(float64(offset1 - offset2))
			hours := difference / 3600
			msg += "You are **" + floatToString(hours) + " hours** away from each other.\n"
		}
	}
	msg += "\n"

	// Find out if the racers have set their stream URL
	if racer1.StreamURL.Valid {
		msg += racer1NameToUse + " has a stream of: <" + racer1.StreamURL.String + ">\n"
	} else {
		msg += racer1NameToUse + ", your stream is **not currently set**. Please set one with: `!stream [url]`\n"
	}
	if racer2.StreamURL.Valid {
		msg += racer2NameToUse + " has a stream of: <" + racer2.StreamURL.String + ">\n"
	} else {
		msg += racer2NameToUse + ", your stream is **not currently set**. Please set one with: `!stream [url]`\n"
	}
	msg += "\n"

	// Give the welcome message
	if race.DatetimeScheduled.Valid {
		var racerThatNeedsToConfirmTime *User
		if race.ActiveRacer == 1 {
			racerThatNeedsToConfirmTime = racer2
		} else {
			racerThatNeedsToConfirmTime = racer1
		}
		msg += getRaceScheduledTime(race, racerThatNeedsToConfirmTime)
		msg += "We are currently waiting on " + racerThatNeedsToConfirmTime.Username + " to confirm that this time is good."
	} else {
		msg += "Please discuss the times that each of you are available to play this week.\n"
		msg += "You can suggest a time to your opponent with something like: `!time 6pm sat`\n"
		msg += "If they accept with `!timeok`, then the match will be officially scheduled."
	}

	discordSend(race.ChannelID, msg)
}

func printStatusScheduled(m *discordgo.MessageCreate, race *Race) {
	msg := getRaceScheduleMessage(race, race.Racer1) // Default to using the first racer's timezone
	discordSend(race.ChannelID, msg)
}
