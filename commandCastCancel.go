package main

import (
	"database/sql"

	"github.com/bwmarrin/discordgo"
)

func commandCastCancel(m *discordgo.MessageCreate, args []string) {
	// Check to see if this is a race channel (and get the race from the database)
	var race *Race
	if v, err := getRace(m.ChannelID); err == sql.ErrNoRows {
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

	// Check to see if there are any casters registered
	if len(race.Casts) == 0 {
		discordSend(m.ChannelID, "No-one has volunteered to cast this match, so you do not need to cancel anything.")
		return
	}

	// Check to see if this person is one of the two racers
	if m.Author.ID == race.Racer1.DiscordID || m.Author.ID == race.Racer2.DiscordID {
		discordSend(m.ChannelID, "If you don't want someone to cast your match, use the `!casternotok` command.")
		return
	}

	// Check to see if this person is registered as a caster
	username := ""
	for _, cast := range race.Casts {
		if cast.Caster.DiscordID == m.Author.ID {
			username = cast.Caster.Username
			break
		}
	}
	if username == "" {
		discordSend(m.ChannelID, "You are not marked as casting this match, so there is no need to cancel anything.")
		return
	}

	// Delete the cast from the database
	if err := modals.Casts.Delete(race.ChannelID, m.Author.ID); err != nil {
		msg := "Failed to delete the cast from the database: " + err.Error()
		log.Error(msg)
		discordSend(m.ChannelID, msg)
		return
	}

	msg := "`" + username + "` has been removed as a caster for this match."
	discordSend(m.ChannelID, msg)
}
