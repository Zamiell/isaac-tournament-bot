package main

import (
	"github.com/bwmarrin/discordgo"
)

func commandEndRound(m *discordgo.MessageCreate, args []string) {
	// Get all of the channels
	var channels []*discordgo.Channel
	if v, err := discord.GuildChannels(discordGuildID); err != nil {
		msg := "Failed to get the Discord server channels: " + err.Error()
		log.Error(msg)
		discordSend(m.ChannelID, msg)
		return
	} else {
		channels = v
	}

	deletedChannels := false
	for _, channel := range channels {
		if channel.ParentID != discordChannelCategoryID {
			continue
		}

		// Delete it from the database
		deletedChannels = true
		if err := db.Races.Delete(channel.ID); err != nil {
			msg := "Failed to delete the race from the database: " + err.Error()
			log.Error(msg)
			discordSend(m.ChannelID, msg)
			return
		}

		// Delete it from Discord
		if _, err := discord.ChannelDelete(channel.ID); err != nil {
			msg := "Failed to delete the \"" + channel.Name + "\" channel: " + err.Error()
			log.Error(msg)
			discordSend(m.ChannelID, msg)
			return
		}

		msg := "Deleted channel \"" + channel.Name + "\"."
		discordSend(m.ChannelID, msg)
		log.Info(msg)
	}

	if !deletedChannels {
		msg := "There were no channels to clean up."
		discordSend(m.ChannelID, msg)
		log.Info(msg)
	}
}
