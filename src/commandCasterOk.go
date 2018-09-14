package main

import (
	"database/sql"

	"github.com/Zamiell/isaac-tournament-bot/src/models"
	"github.com/bwmarrin/discordgo"
)

func commandCasterOk(m *discordgo.MessageCreate, args []string) {
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
	if m.Author.ID != race.Racer1.DiscordID && m.Author.ID != race.Racer2.DiscordID {
		discordSend(m.ChannelID, "You cannot submit approval for a match that you are not participanting in.")
		return
	}

	// Check to see if someone is casting this match
	if !race.CasterID.Valid {
		discordSend(m.ChannelID, "No-one has volunteered to cast this match, so there is no need to submit approval.")
		return
	}

	// Find out whether they are player 1 or player 2
	playerNum := 1
	racerName := race.Racer1.Username
	if m.Author.ID == race.Racer2.DiscordID {
		playerNum = 2
		racerName = race.Racer2.Username
	}

	// Set approval
	if err := db.Races.SetCasterApproval(m.ChannelID, playerNum); err != nil {
		msg := "Failed to set the caster approval in the database: " + err.Error()
		log.Error(msg)
		discordSend(m.ChannelID, msg)
		return
	}

	msg := racerName + " has approved " + race.Caster.Mention() + " as the caster for this match.\n"
	if playerNum == 1 {
		if race.CasterP2 {
			msg += "Both racers have now approved this caster."
		} else {
			msg += race.Racer2.Mention() + " still needs to approve or disapprove this caster."
		}
	} else if playerNum == 2 {
		if race.CasterP1 {
			msg += "Both racers have now approved this caster."
		} else {
			msg += race.Racer1.Mention() + " still needs to approve or disapprove this caster."
		}
	}
	discordSend(m.ChannelID, msg)
}
