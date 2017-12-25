package main

import (
	"strconv"

	"github.com/bwmarrin/discordgo"
)

func commandRandChar(m *discordgo.MessageCreate, args []string) {
	randCharNum := getRandom(0, len(characters)-1)
	randChar := characters[randCharNum]
	msg := "Random character between 1 and " + strconv.Itoa(len(characters)) + ":\n"
	msg += "**" + strconv.Itoa(randCharNum) + " - " + randChar + "**"
	discordSend(m.ChannelID, msg)
}
