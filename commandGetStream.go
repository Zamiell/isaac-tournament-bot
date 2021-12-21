package main

import (
	"github.com/bwmarrin/discordgo"
)

func commandGetStream(m *discordgo.MessageCreate, args []string) {
	if len(args) != 1 {
		commandGetStreamPrint(m)
		return
	}

	// Get the Discord guild object
	var guild *discordgo.Guild
	if v, err := discord.Guild(discordGuildID); err != nil {
		msg := "Failed to get the Discord guild: " + err.Error()
		log.Error(msg)
		discordSend(m.ChannelID, msg)
		return
	} else {
		guild = v
	}

	// Find the discord ID of the user
	var discordUser *discordgo.User
	for _, member := range guild.Members {
		username := member.Nick
		if username == "" {
			username = member.User.Username
		}
		if username == args[0] {
			discordUser = member.User
			break
		}
	}
	if discordUser == nil {
		msg := "Failed to find \"" + args[0] + "\" in the Discord server."
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
