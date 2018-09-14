package main

import (
	"github.com/Zamiell/isaac-tournament-bot/src/models"
	"github.com/bwmarrin/discordgo"
)

func commandCasterAlwaysOk(m *discordgo.MessageCreate, args []string) {
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
	if err := db.Racers.SetCasterAlwaysOk(m.Author.ID, true); err != nil {
		msg := "Failed to update the default caster approval: " + err.Error()
		log.Error(msg)
		discordSend(m.ChannelID, msg)
		return
	}

	msg := "**" + racer.Username + "** has enabled default caster approval."
	msg += "To disable this, use the `!casteralwaysnotok` command."
	discordSend(m.ChannelID, msg)
}
