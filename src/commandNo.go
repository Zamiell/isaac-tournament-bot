package main

import (
	"database/sql"

	"github.com/Zamiell/isaac-tournament-bot/src/models"
	"github.com/bwmarrin/discordgo"
)

func commandNo(m *discordgo.MessageCreate, args []string) {
	// Check to see if this is a race channel (and get the race from the database)
	var race models.Race
	if v, err := raceGet(m.ChannelID); err == sql.ErrNoRows {
		discordSend(m.ChannelID, "You can only use that command in a race channel.")
		return
	} else if err != nil {
		msg := "Failed to get the race from the database: " + err.Error()
		log.Error(msg)
		discordSend(m.ChannelID, msg)
		return
	} else {
		race = v
	}

	// Check to see if this person is one of the two racers
	var playerNum int
	if m.Author.ID == race.Racer1.DiscordID {
		playerNum = 1
	} else if m.Author.ID == race.Racer2.DiscordID {
		playerNum = 2
	} else {
		discordSend(m.ChannelID, "Only \""+race.Racer1.Username+"\" and \""+race.Racer2.Username+"\" can veto a build.")
		return
	}

	// Check to see if this race is in the item banning phase
	if race.State != "vetoBuilds" {
		discordSend(m.ChannelID, "You can only veto something once the match has started.")
		return
	}

	// Check to see if it is their turn
	if race.ActivePlayer != playerNum {
		discordSend(m.ChannelID, "It is not your turn.")
		return
	}

	// Set the number of people who have voted on this build
	race.NumVoted++
	if err := db.Races.SetNumVoted(race.ChannelID, race.NumVoted); err != nil {
		msg := "Failed to set the NumVoted for race \"" + race.Name() + "\": " + err.Error()
		log.Error(msg)
		discordSend(m.ChannelID, msg)
		return
	}

	incrementActivePlayer(&race)
	buildsRound(race, "")
	log.Info("Racer \"" + m.Author.Username + "\" declined to veto.")
}
