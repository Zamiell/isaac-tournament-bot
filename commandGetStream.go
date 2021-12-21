package main

import (
	"github.com/bwmarrin/discordgo"
)

func commandGetStream(m *discordgo.MessageCreate, args []string) {
	if len(args) != 1 {
		commandGetStreamPrint(m)
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

	var user *User
	if v, err := userGet(discordUser); err != nil {
		msg := "Failed to get the user from the database: " + err.Error()
		log.Error(msg)
		discordSend(m.ChannelID, msg)
		return
	} else {
		user = v
	}

	msg := "The stream for `" + user.Username + "` is "
	if user.StreamURL.Valid {
		msg += "currently set to:\n<" + user.StreamURL.String + ">"
	} else {
		msg += "**not currently set**."
	}
	discordSend(m.ChannelID, msg)
}

func commandGetStreamPrint(m *discordgo.MessageCreate) {
	msg := "Get another user's stream with: `!getstream [username]`\n"
	msg += "e.g. `!getstream Willy`"
	discordSend(m.ChannelID, msg)
}
