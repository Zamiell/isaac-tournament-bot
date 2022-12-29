package main

import (
	"github.com/bwmarrin/discordgo"
)

func commandCheckRound(m *discordgo.MessageCreate, args []string) {
	if !isAdmin(m) {
		return
	}

	// Go through all of the tournaments.
	for _, tournament := range tournaments {
		startRound(m, tournament, true)
	}
}
