package main

import (
	"strconv"

	"github.com/bwmarrin/discordgo"
)

func commandRandom(m *discordgo.MessageCreate, args []string) {
	if len(args) != 2 {
		commandRandomPrint(m)
		return
	}

	// Ensure that both arguments are numbers
	var min int
	if v, err := strconv.Atoi(args[0]); err != nil {
		discordSend(m.ChannelID, "\""+args[0]+"\" is not a number.")
		return
	} else {
		min = v
	}
	var max int
	if v, err := strconv.Atoi(args[1]); err != nil {
		discordSend(m.ChannelID, "\""+args[1]+"\" is not a number.")
		return
	} else {
		max = v
	}

	randNum := getRandom(min, max)
	msg := "Random number between " + strconv.Itoa(min) + " and " + strconv.Itoa(max) + ":\n"
	msg += "**" + strconv.Itoa(randNum) + "**"
	discordSend(m.ChannelID, msg)
}

func commandRandomPrint(m *discordgo.MessageCreate) {
	msg := "Get a random number with: `!random [min] [max]`\n"
	msg += "e.g. `!random 1 2`"
	discordSend(m.ChannelID, msg)
}
