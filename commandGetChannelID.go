package main

import (
	"github.com/bwmarrin/discordgo"
)

func commandGetChannelID(m *discordgo.MessageCreate, args []string) {
	if !isAdmin(m) {
		return
	}

	if len(args) == 0 {
		commandGetChannelIDPrint(m)
		return
	}

	// Get the Discord guild object.
	var guild *discordgo.Guild
	if v, err := discordSession.Guild(discordGuildID); err != nil {
		msg := "Failed to get the Discord guild: " + err.Error()
		log.Error(msg)
		discordSend(m.ChannelID, msg)
		return
	} else {
		guild = v
	}

	// Go through all of the channels.
	for _, channel := range guild.Channels {
		if channel.Name == args[0] {
			msg := "Found channel \"" + args[0] + "\": **" + channel.ID + "**"
			discordSend(m.ChannelID, msg)
			return
		}
	}
}

func commandGetChannelIDPrint(m *discordgo.MessageCreate) {
	msg := "Get the ID of a Discord channel with: `!getchannelid [name]`\n"
	msg += "e.g. `!getchannelid general`"
	discordSend(m.ChannelID, msg)
}
