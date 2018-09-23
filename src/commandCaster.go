package main

import (
	"database/sql"

	"github.com/Zamiell/isaac-tournament-bot/src/models"
	"github.com/bwmarrin/discordgo"
)

func commandCaster(m *discordgo.MessageCreate, args []string) {
	// Check to see if this is a race channel (and get the race from the database)
	var race *models.Race
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

	// Check to see if someone is casting this match
	if len(race.Casts) == 0 {
		discordSend(m.ChannelID, "No-one has volunteered to cast this match yet.")
		return
	}

	// Display all of the casts for this match
	msg := ""
	for _, cast := range race.Casts {
		if cast.R1Permission && cast.R2Permission {
			msg += "`" + cast.Caster.Username + "` is approved to cast this match in " + languageMap[cast.Language] + " at: <" + cast.Caster.StreamURL.String + ">\n"
		} else {
			msg += "`" + cast.Caster.Username + "` has requested to cast this match in " + languageMap[cast.Language] + " at: <" + cast.Caster.StreamURL.String + ">\n"
			if !cast.R1Permission {
				msg += "`" + race.Racer1.Username + "` still needs to okay this with the `!casterok` command.\n"
			} else if !cast.R2Permission {
				msg += "`" + race.Racer2.Username + "` still needs to okay this with the `!casterok` command.\n"
			}
		}
	}
	discordSend(m.ChannelID, msg)
}
