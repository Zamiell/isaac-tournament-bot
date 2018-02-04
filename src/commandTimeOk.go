package main

import (
	"database/sql"

	"github.com/Zamiell/isaac-tournament-bot/src/models"
	"github.com/bwmarrin/discordgo"
)

func commandTimeOk(m *discordgo.MessageCreate, args []string) {
	// Create the user in the database if it does not already exist
	var racer models.Racer
	if v, err := racerGet(m.Author); err != nil {
		msg := "Failed to get the racer from the database: " + err.Error()
		log.Error(msg)
		discordSend(m.ChannelID, msg)
		return
	} else {
		racer = v
	}

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
	var activePlayer int
	if m.Author.ID == race.Racer1.DiscordID {
		activePlayer = 1
	} else if m.Author.ID == race.Racer2.DiscordID {
		activePlayer = 2
	} else {
		discordSend(m.ChannelID, "Only \""+race.Racer1.Username+"\" and \""+race.Racer2.Username+"\" can confirm the time for this match.")
		return
	}

	// Check to see if this race has already been scheduled
	if race.State != "initial" {
		discordSend(m.ChannelID, "Both racers have already agreed to a time, so you cannot confirm.")
		return
	}

	// Check to see if a time has been suggested
	if !race.DatetimeScheduled.Valid {
		discordSend(m.ChannelID, "No-one has suggested a time for the match yet, so you cannot confirm it.")
		return
	}

	// Check to see if they were the one who suggested the time
	if activePlayer == race.ActivePlayer {
		discordSend(m.ChannelID, "The other racer needs to confirm the time, not you.")
		return
	}

	// Check to see if this person has a timezone specified
	if !racer.Timezone.Valid {
		discordSend(m.ChannelID, "You must specify a timezone with the `!timezone` command before you can confirm the time for the match.")
		return
	}

	// Check to see if this person has a stream specified
	if !racer.StreamURL.Valid {
		discordSend(m.ChannelID, "You must specify a stream URL with the `!stream` command before you can confirm the time for the match.")
		return
	}

	// Set the state
	race.State = "scheduled"
	if err := db.Races.SetState(m.ChannelID, race.State); err != nil {
		msg := "Failed to set the state for race \"" + race.Name() + "\": " + err.Error()
		log.Error(msg)
		discordSend(m.ChannelID, msg)
		return
	}

	msg := "The race time has been confirmed. I will notify you 5 minutes before the match begins.\n"
	msg += "(To delete this time and start over, use the `!timedelete` command.)"
	discordSend(m.ChannelID, msg)
	log.Info("Race \"" + race.Name() + "\" scheduled; set to state \"" + race.State + "\".")

	// Sleep until the match starts
	// (use a goroutine so that the rest of the program doesn't block)
	go matchStart(race)
}
