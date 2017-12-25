package main

import (
	"github.com/bwmarrin/discordgo"
)

func commandDebug(m *discordgo.MessageCreate, args []string) {
	if !isAdmin(m) {
		return
	}
}
