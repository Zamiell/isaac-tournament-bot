package main

import (
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/Zamiell/isaac-tournament-bot/src/models"
	"github.com/bwmarrin/discordgo"
	humanize "github.com/dustin/go-humanize"
)

// Check to see if they are an administrator
func isAdmin(m *discordgo.MessageCreate) bool {
	var member *discordgo.Member
	if v, err := discord.GuildMember(discordGuildID, m.Author.ID); err != nil {
		msg := "Failed to get the presence for the user: " + err.Error()
		log.Error(msg)
		discordSend(m.ChannelID, msg)
		return false
	} else {
		member = v
	}
	isAdmin := false
	for _, role := range member.Roles {
		if role == discordAdminRoleID {
			isAdmin = true
			break
		}
	}
	if !isAdmin {
		discordSend(m.ChannelID, "Only admins can perform this command.")
		return false
	}

	return true
}

/*
	Scheduling subroutines
*/

func getTimezone(timezone string) string {
	msg := timezone + " "
	msg += "(" + getTimezoneShort(timezone) + " / " + getTimezoneGMT(timezone) + ")"
	return msg
}

func getTimezoneShort(timezone string) string {
	loc, _ := time.LoadLocation(timezone)
	shortName, _ := time.Now().In(loc).Zone()
	return shortName
}

func getTimezoneGMT(timezone string) string {
	loc, _ := time.LoadLocation(timezone)
	_, offset := time.Now().In(loc).Zone()
	gmt := offset / 3600
	msg := "GMT"
	if offset >= 0 {
		msg += "+"
	}
	msg += strconv.Itoa(gmt)
	return msg
}

func getDate(datetime time.Time, timezone string) string {
	loc, _ := time.LoadLocation(timezone)
	datetimeLocal := datetime.In(loc)
	dateFormatString := "Monday, January "
	dateString := datetimeLocal.Format(dateFormatString)
	dateString += humanize.Ordinal(datetimeLocal.Day())
	timeFormatString := "3:04 PM"
	timeString := datetimeLocal.Format(timeFormatString)

	msg := "**" + dateString + " @ " + timeString + " (" + getTimezoneShort(timezone) + ")**"
	return msg
}

/*
	Ban subroutines
*/

func getRemaining(race models.Race, thing string) string {
	// Characters and builds are stored in the array as a comma delimited string,
	// so we must convert it back into a slice
	var thingArray string
	var thingsFull []string
	if thing == "characters" {
		thingArray = race.Characters
		thingsFull = characters
	} else if thing == "builds" {
		thingArray = race.Builds
		thingsFull = builds
	} else {
		log.Error("The \"getRemaining\" function was passed an invalid phase.")
		return ""
	}
	things := strings.Split(thingArray, ",")

	// Build column 1
	lines := make([]string, 0)
	column1length := 0
	halfwayPoint := int(math.Floor(float64((len(things) - 1) / 2)))
	log.Info("halfwayPoint:", halfwayPoint)
	for i := 0; i <= halfwayPoint; i++ {
		line := strconv.Itoa(i+1) + " - " + things[i]
		lines = append(lines, line)
		if len(line) > column1length {
			column1length = len(line)
		}
	}

	// Add padding to column 1
	column1length += 6 // A minimum of 6 spaces in between columns
	for i := 0; i < len(lines); i++ {
		for len(lines[i]) < column1length {
			lines[i] += " "
		}
	}

	// Build column 2
	lineCounter := 0
	for i := halfwayPoint + 1; i < len(things); i++ {
		line := strconv.Itoa(i+1) + " - " + things[i]
		lines[lineCounter] += line
		lineCounter += 1
	}

	// Make the string
	bansLeft := len(things) - len(thingsFull) + bansNum
	msg := "**" + strconv.Itoa(bansLeft) + " ban"
	if bansLeft > 0 {
		msg += "s"
	}
	msg += " to go.**\n"
	msg += "Current " + thing + " remaining:\n\n"
	msg += "```\n"
	for _, line := range lines {
		msg += line + "\n"
	}
	msg += "```"
	return msg
}

func getNext(race models.Race) string {
	var msg string
	if race.ActivePlayer == 1 {
		msg = race.Racer1.Mention()
	} else if race.ActivePlayer == 2 {
		msg = race.Racer2.Mention()
	}
	msg += ", you're next!\n\n"
	return msg
}
