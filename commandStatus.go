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

	printStatus(m, race)
}

func printStatus(m *discordgo.MessageCreate, race *Race) {

}

func printStatusPre() {
	// Find out if the racers have set their timezone
	msg := ""
	if racer1.Timezone.Valid {
		msg += discordUser1.Mention() + " has a timezone of: " + getTimezone(racer1.Timezone.String) + "\n"
	} else {
		msg += discordUser1.Mention() + ", your timezone is **not currently set**. Please set one with: `!timezone [timezone]`\n"
	}
	if racer2.Timezone.Valid {
		msg += discordUser2.Mention() + " has a timezone of: " + getTimezone(racer2.Timezone.String) + "\n"
	} else {
		msg += discordUser2.Mention() + ", your timezone is **not currently set**. Please set one with: `!timezone [timezone]`\n"
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
		msg += discordUser1.Mention() + " has a stream of: <" + racer1.StreamURL.String + ">\n"
	} else {
		msg += discordUser1.Mention() + ", your stream is **not currently set**. Please set one with: `!stream [url]`\n"
	}
	if racer2.StreamURL.Valid {
		msg += discordUser2.Mention() + " has a stream of: <" + racer2.StreamURL.String + ">\n"
	} else {
		msg += discordUser2.Mention() + ", your stream is **not currently set**. Please set one with: `!stream [url]`\n"
	}
	msg += "\n"

	// Give the welcome message
	msg += "Please discuss the times that each of you are available to play this week.\n"
	if tournament.Ruleset == "team" {
		msg += discordUser1.Mention() + " and " + discordUser2.Mention() + " are the team captains; I will only listen to them.\n"
	}
	msg += "You can suggest a time to your opponent with something like: `!time 6pm sat`\n"
	msg += "If they accept with `!timeok`, then the match will be officially scheduled."
	discordSend(channelID, msg)

}
