package main

import (
	"strconv"

	"github.com/bwmarrin/discordgo"
)

func commandRandBuild(m *discordgo.MessageCreate, args []string) {
	randBuildNum := getRandom(0, len(builds)-1)
	randBuild := builds[randBuildNum]
	msg := "Random build between 1 and " + strconv.Itoa(len(builds)) + ":\n"
	msg += "**" + strconv.Itoa(randBuildNum) + " - " + randBuild + "**"
	discordSend(m.ChannelID, msg)
}
