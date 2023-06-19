package main

import (
	"strconv"

	"github.com/bwmarrin/discordgo"
)

func commandRandomBuild(m *discordgo.MessageCreate, args []string) {
	randomBuild, randomBuildIndex := getRandomArrayElement(builds)
	msg := "Random build between 1 and " + strconv.Itoa(len(builds)) + ":\n"
	msg += "**" + strconv.Itoa(randomBuildIndex+1) + " - " + randomBuild.Name + "**"
	discordSend(m.ChannelID, msg)
}
