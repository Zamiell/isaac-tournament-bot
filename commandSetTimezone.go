package main

import (
	"github.com/bwmarrin/discordgo"
)

func commandSetTimezone(m *discordgo.MessageCreate, args []string) {
	if !isAdmin(m) {
		return
	}

	if len(args) != 2 {
		commandSetTimezonePrint(m)
		return
	}
	username := args[0]

	// Get the Discord guild members
	var members []*discordgo.Member
	if v, err := discordSession.GuildMembers(discordGuildID, "0", 1000); err != nil {
		msg := "Failed to get the Discord guild members: " + err.Error()
		log.Error(msg)
		discordSend(m.ChannelID, msg)
		return
	} else {
		members = v
	}

	discordUser := getDiscordUserByName(members, username)
	if discordUser == nil {
		msg := "Failed to find \"" + username + "\" in the Discord server."
		log.Error(msg)
		discordSend(m.ChannelID, msg)
		return
	}

	m.Author = discordUser
	args = args[1:] // This will be an empty slice if there is nothing after the command
	commandTimezone(m, args)
}

func commandSetTimezonePrint(m *discordgo.MessageCreate) {
	msg := "Set another user's timezone with: `!settimezone [username] [timezone]`\n"
	msg += "e.g. `!settimezone Willy America/New_York`\n"
	msg += "The submitted timezone has to exactly match the TZ column of the following page:\n"
	msg += "<https://en.wikipedia.org/wiki/List_of_tz_database_time_zones>"
	discordSend(m.ChannelID, msg)
}
