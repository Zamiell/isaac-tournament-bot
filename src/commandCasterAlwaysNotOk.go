package main

import (
	"github.com/Zamiell/isaac-tournament-bot/src/models"
	"github.com/bwmarrin/discordgo"
)

func commandCasterAlwaysNotOk(m *discordgo.MessageCreate, args []string) {
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

	// Set the new value
	if err := db.Racers.SetCasterAlwaysOk(m.Author.ID, false); err != nil {
		msg := "Failed to update the default caster approval: " + err.Error()
		log.Error(msg)
		discordSend(m.ChannelID, msg)
		return
	}

	msg := "**" + racer.Username + "** has disabled default caster approval.\n"
	msg += "To enable this, use the `!casteralwaysok` command."
	discordSend(m.ChannelID, msg)
}
