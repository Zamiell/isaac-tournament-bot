package main

import (
	"github.com/Zamiell/isaac-tournament-bot/src/models"
	"github.com/bwmarrin/discordgo"
)

func commandTimeOk(m *discordgo.MessageCreate, args []string) {
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
		discordSend(m.ChannelID, "Only \""+race.Racer1.Username+"\" and \""+race.Racer2.Username+"\" can reschedule this match.")
		return
	}

	// Check to see if this race has already been scheduled
	if race.State != 0 {
		discordSend(m.ChannelID, "Both racers have already agreed to a time, so you cannot confirm.")
		return
	}

	// Check to see if a time has been suggested
	if !race.DatetimeScheduled.Valid {
		discordSend(m.ChannelID, "No-one has suggested a time for the match yet, so you cannot confirm it.")
		return
	}

	// Set the state to 1
	if err := db.Races.SetState(m.ChannelID, 1); err != nil {
		msg := "Failed to set the state: " + err.Error()
		log.Error(msg)
		discordSend(m.ChannelID, msg)
		return
	}

	msg := "The race time has been confirmed. Before the match begins, use the `!startbans` command to ban characters and items.\n"
	msg += "(To delete this time and start over, use the `!timedelete` command.)"
	discordSend(m.ChannelID, msg)
	log.Info("Race \"" + race.Name() + "\" confirmed; set to state 1.")
}
