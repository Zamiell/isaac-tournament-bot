package main

import (
	"github.com/bwmarrin/discordgo"
)

func commandRacingPlus(m *discordgo.MessageCreate, args []string) {
	discordSend(m.ChannelID, "Racing+ is a mod for The Binding of Isaac: Afterbirth+: https://isaacracing.net")
}
