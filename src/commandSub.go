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

func getTimezone(timezone string) string {
	msg := timezone + " "
	msg += "(" + getTimezoneShort(timezone) + ")"
	return msg
}

func getTimezoneShort(timezone string) string {
	loc, _ := time.LoadLocation(timezone)
	shortName, _ := time.Now().In(loc).Zone()
	msg := shortName + " / " + getGMT(timezone)
	return msg
}

func getGMT(timezone string) string {
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
