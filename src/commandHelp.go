package main

import (
	"github.com/bwmarrin/discordgo"
)

func commandHelp(m *discordgo.MessageCreate, args []string) {
	discordSend(m.ChannelID, commandHelpGetMsg())
}
