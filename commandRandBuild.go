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
	itemNames := getBuildItemNames(build)
	return strings.Join(itemNames, " + ")
}

func getBuildItemNames(build []IsaacItem) []string {
	itemNames := make([]string, 0)
	for _, item := range build {
		itemNames = append(itemNames, item.Name)
	}

	return itemNames
}
