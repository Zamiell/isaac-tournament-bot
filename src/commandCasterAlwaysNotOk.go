package main

import (
	"github.com/Zamiell/isaac-tournament-bot/src/models"
	"github.com/bwmarrin/discordgo"
)

func commandCasterAlwaysNotOk(m *discordgo.MessageCreate, args []string) {
	// Create the user in the database if it does not already exist
	var user *models.User
	if v, err := userGet(m.Author); err != nil {
		msg := "Failed to get the user from the database: " + err.Error()
		log.Error(msg)
		discordSend(m.ChannelID, msg)
		return
	} else {
		user = v
	}

	// Check to see if they have already enabled default caster approval
	if !user.CasterAlwaysOk {
		msg := "You have not yet enabled default caster approval. You can enable it with the `!casteralwaysok` command."
		discordSend(m.ChannelID, msg)
		return
	}

	// Set the new value
	if err := db.Users.SetCasterAlwaysOk(m.Author.ID, false); err != nil {
		msg := "Failed to update the default caster approval: " + err.Error()
		log.Error(msg)
		discordSend(m.ChannelID, msg)
		return
	}

	msg := "**" + user.Username + "** has disabled default caster approval.\n"
	msg += "(To enable this, use the `!casteralwaysok` command.)"
	discordSend(m.ChannelID, msg)
}
