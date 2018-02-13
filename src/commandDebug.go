package main

import (
	"github.com/bwmarrin/discordgo"
)

func commandDebug(m *discordgo.MessageCreate, args []string) {
	if !isAdmin(m) {
		return
	}

	log.Info("Tournaments:")
	for _, tournament := range tournaments { // This is a map indexed by "ChallongeURL"
		log.Info(tournament.Name + " - " + tournament.ChallongeURL + " - " + floatToString(tournament.ChallongeID) + " - " + tournament.Ruleset)
	}
}
