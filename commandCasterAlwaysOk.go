package main

import (
	"github.com/bwmarrin/discordgo"
)

func commandCasterAlwaysOk(m *discordgo.MessageCreate, args []string) {
	// Create the user in the database if it does not already exist
	var user *User
	if v, err := userGet(m.Author); err != nil {
		msg := "Failed to get the user from the database: " + err.Error()
		log.Error(msg)
		discordSend(m.ChannelID, msg)
		return
	} else {
		user = v
	}

	// Check to see if they have already enabled default caster approval
	if user.CasterAlwaysOk {
		msg := "You have already enabled default caster approval. You can disable it with the `!casteralwaysnotok` command."
		discordSend(m.ChannelID, msg)
		return
	}

	// Set the new value
	if err := modals.Users.SetCasterAlwaysOk(m.Author.ID, true); err != nil {
		msg := "Failed to update the default caster approval: " + err.Error()
		log.Error(msg)
		discordSend(m.ChannelID, msg)
		return
	}

	msg := "**" + user.Username + "** has enabled default caster approval.\n"
	msg += "(To disable this, use the `!casteralwaysnotok` command.)"
	discordSend(m.ChannelID, msg)
}
