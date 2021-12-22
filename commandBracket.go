package main

import (
	"github.com/bwmarrin/discordgo"
)

func commandBracket(m *discordgo.MessageCreate, args []string) {
	msg := "The Challonge bracket(s):\n"
	for urlSuffix := range tournaments {
		msg += "<https://challonge.com/" + urlSuffix + ">\n"
	}
	discordSend(m.ChannelID, msg)
}
