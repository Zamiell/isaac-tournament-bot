package main

import (
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
	humanize "github.com/dustin/go-humanize"
)

// Check to see if they are an administrator
func isAdmin(m *discordgo.MessageCreate) bool {
	var member *discordgo.Member
	if v, err := discordSession.GuildMember(discordGuildID, m.Author.ID); err != nil {
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
	return "`" + msg + "`"
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
	dateFormatString2 := "Monday, January "
	dateString := datetimeLocal.Format(dateFormatString2)
	dateString += humanize.Ordinal(datetimeLocal.Day())
	timeFormatString := "3:04 PM"
	timeString := datetimeLocal.Format(timeFormatString)

	msg := dateString + " @ " + timeString + " (" + getTimezoneShort(timezone) + ")"
	return msg
}

/*
	Match subroutines
*/

func getBansRemaining(race *Race) string {
	bansLeft := race.Racer1Bans + race.Racer2Bans
	msg := "**" + strconv.Itoa(bansLeft) + " ban"
	if bansLeft != 1 {
		msg += "s"
	}
	msg += " to go.**\n"
	return msg
}

func getPicksRemaining(race *Race) string {
	var things []string
	if race.State == RaceStateBanningCharacters || race.State == RaceStatePickingCharacters {
		things = race.Characters
	} else if race.State == RaceStateBanningBuilds || race.State == RaceStatePickingBuilds {
		things = race.Builds
	} else {
		log.Error("The \"getPicksRemaining\" function was called when the race state was invalid.")
		return "error"
	}

	picksLeft := tournaments[race.ChallongeURL].BestOf - len(things)
	msg := "**" + strconv.Itoa(picksLeft) + " pick"
	if picksLeft != 1 {
		msg += "s"
	}
	msg += " to go.**\n"
	return msg
}

func getRemaining(race *Race) string {
	var thing string
	var thingsRemaining []string
	if race.State == RaceStateBanningCharacters || race.State == RaceStatePickingCharacters {
		thing = "characters"
		thingsRemaining = race.CharactersRemaining
	} else if race.State == RaceStateBanningBuilds || race.State == RaceStatePickingBuilds {
		thing = "builds"
		thingsRemaining = race.BuildsRemaining
	} else {
		log.Error("The \"getRemaining\" function was called when the race state was invalid.")
		return "error"
	}

	// Build column 1
	lines := make([]string, 0)
	column1length := 0
	halfwayPoint := int(float64((len(thingsRemaining) - 1) / 2))
	for i := 0; i <= halfwayPoint; i++ {
		line := strconv.Itoa(i+1) + " - " + thingsRemaining[i]
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
	for i := halfwayPoint + 1; i < len(thingsRemaining); i++ {
		line := strconv.Itoa(i+1) + " - " + thingsRemaining[i]
		lines[lineCounter] += line
		lineCounter++
	}

	// Make the string
	msg := "Current " + thing + " remaining:\n\n"
	msg += "```\n"
	for _, line := range lines {
		msg += line + "\n"
	}
	msg += "```"
	return msg
}

func getNext(race *Race) string {
	var msg string
	if race.ActiveRacer == 1 {
		msg = race.Racer1.Mention()
	} else if race.ActiveRacer == 2 {
		msg = race.Racer2.Mention()
	}
	msg += ", you're next!\n\n"
	return msg
}

func incrementActiveRacer(race *Race) {
	// Increment the active racer
	race.ActiveRacer++
	if race.ActiveRacer > 2 {
		race.ActiveRacer = 1
	}
	if err := modals.Races.SetActiveRacer(race.ChannelID, race.ActiveRacer); err != nil {
		msg := "Failed to set the active racer for race \"" + race.Name() + "\": " + err.Error()
		log.Error(msg)
		discordSend(race.ChannelID, msg)
		return
	}
}
