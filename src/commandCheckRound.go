package main

import (
	"github.com/bwmarrin/discordgo"
)

func commandCheckRound(m *discordgo.MessageCreate, args []string) {
	if !isAdmin(m) {
		return
	}

	// Go through all of the tournaments
	for tournamentName, _ := range tournaments {
		startRound(m, tournamentName, true) // The third argument is "dryRun"
	}
}
