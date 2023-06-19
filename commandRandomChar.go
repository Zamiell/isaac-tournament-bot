package main

import (
	"strconv"

	"github.com/bwmarrin/discordgo"
)

func commandRandomChar(m *discordgo.MessageCreate, args []string) {
	randomCharacter, randomCharacterIndex := getRandomArrayElement(characters)
	msg := "Random character between 1 and " + strconv.Itoa(len(characters)) + ":\n"
	msg += "**" + strconv.Itoa(randomCharacterIndex+1) + " - " + randomCharacter + "**"
	discordSend(m.ChannelID, msg)
}
