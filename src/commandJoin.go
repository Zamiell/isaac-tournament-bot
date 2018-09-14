package main

import (
	"github.com/bwmarrin/discordgo"
)

func commandJoin(m *discordgo.MessageCreate, args []string) {
	if !isAdmin(m) {
		return
	}

	clientID := "392184872921333770"
	msg := "<https://discordapp.com/oauth2/authorize?client_id=" + clientID + "&scope=bot&permissions=0>"
	discordSend(m.ChannelID, msg)
}
