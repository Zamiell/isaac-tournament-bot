package main

import (
	"github.com/bwmarrin/discordgo"
)

func commandSetCasterNotOk(m *discordgo.MessageCreate, args []string) {
	if !isAdmin(m) {
		return
	}

	if len(args) == 0 {
		commandSetCasterNotOkPrint(m)
		return
	}

	// Get the Discord guild object
	var guild *discordgo.Guild
	if v, err := discordSession.Guild(discordGuildID); err != nil {
		msg := "Failed to get the Discord guild: " + err.Error()
		log.Error(msg)
		discordSend(m.ChannelID, msg)
		return
	} else {
		guild = v
	}

	// Find the Discord ID of the user
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

	m.Author = discordUser
	args = args[1:] // This will be an empty slice if there is nothing after the command
	commandCasterNotOk(m, args)
}

func commandSetCasterNotOkPrint(m *discordgo.MessageCreate) {
	msg := "Deny permission on behalf of a racer with: `!setcasternotok [username]`\n"
	msg += "e.g. `!setcasternotok Willy`"
	discordSend(m.ChannelID, msg)
}
