package main

import (
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func commandRandBuild(m *discordgo.MessageCreate, args []string) {
	numBuilds := len(builds) - 1
	randomBuildNum := getRandomInt(1, numBuilds) // The 0th build is blank
	randomBuild := builds[randomBuildNum]
	randomBuildName := getBuildName(randomBuild)
	msg := "Random build between 1 and " + strconv.Itoa(numBuilds) + ":\n"
	msg += "**" + strconv.Itoa(randomBuildNum) + " - " + randomBuildName + "**"
	discordSend(m.ChannelID, msg)
}

func getBuildName(build []IsaacItem) string {
	seperator := " + "
	msg := ""
	for _, item := range build {
		msg += item.Name
		msg += seperator
	}
	msg = strings.TrimSuffix(msg, seperator)

	return msg
}
