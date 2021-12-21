package main

import (
	"github.com/bwmarrin/discordgo"
)

func commandEndRound(m *discordgo.MessageCreate, args []string) {
	if !isAdmin(m) {
		return
	}

	// Get all of the channels
	var channels []*discordgo.Channel
	if v, err := discordSession.GuildChannels(discordGuildID); err != nil {
		msg := "Failed to get the Discord server channels: " + err.Error()
		log.Error(msg)
		discordSend(m.ChannelID, msg)
		return
	} else {
		channels = v
	}

	deletedChannels := false
	for _, channel := range channels {
		raceChannel := false
		for _, tournament := range tournaments {
			if tournament.DiscordCategoryID == channel.ParentID {
				raceChannel = true
				break
			}
		}
		if !raceChannel {
			continue
		}

		// Delete it from the database
		deletedChannels = true
		if err := modals.Races.Delete(channel.ID); err != nil {
			msg := "Failed to delete the race from the database: " + err.Error()
			log.Error(msg)
			discordSend(m.ChannelID, msg)
			return
		}

		// Delete it from Discord
		if _, err := discordSession.ChannelDelete(channel.ID); err != nil {
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
