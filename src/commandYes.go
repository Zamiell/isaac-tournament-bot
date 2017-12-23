package main

import (
	"github.com/Zamiell/isaac-tournament-bot/src/models"
	"github.com/bwmarrin/discordgo"
)

func commandYes(m *discordgo.MessageCreate, args []string) {
	var race models.Race
	if v, err := raceGet(m.ChannelID); err != nil {
		msg := "Failed to get the race from the database: " + err.Error()
		log.Error(msg)
		discordSend(m.ChannelID, msg)
		return
	} else {
		race = v
	}

	// Check to see if this person is one of the two racers
	if m.Author.ID != race.Racer1.DiscordID && m.Author.ID != race.Racer2.DiscordID {
		discordSend(m.ChannelID, "Only \""+race.Racer1.Username+"\" and \""+race.Racer2.Username+"\" can veto a build.")
		return
	}

	// Check to see if this race is in the item banning phase
	if race.State != 3 {
		discordSend(m.ChannelID, "You can only veto something once the match has started.")
		return
	}

	// Check to see if it is their turn
	if (race.ActivePlayer == 1 && m.Author.ID != race.Racer2.DiscordID) ||
		(race.ActivePlayer == 2 && m.Author.ID != race.Racer1.DiscordID) {

		discordSend(m.ChannelID, "It is not your turn.")
		return
	}
}
