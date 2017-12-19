package main

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

func commandRoundEnd(m *discordgo.MessageCreate, args []string) {
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
		if !strings.HasPrefix(channel.Name, "round-") {
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

		discordSend(m.ChannelID, "Deleted channel \""+channel.Name+"\".")
	}

	if !deletedChannels {
		discordSend(m.ChannelID, "I did not find any channels to clear up.")
	}
}
