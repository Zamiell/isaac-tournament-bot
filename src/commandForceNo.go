package main

import (
	"github.com/Zamiell/isaac-tournament-bot/src/models"
	"github.com/bwmarrin/discordgo"
)

func commandForceNo(m *discordgo.MessageCreate, args []string) {
	if !isAdmin(m) {
		return
	}

	var race models.Race
	if v, err := raceGet(m.ChannelID); err != nil {
		msg := "Failed to get the race from the database: " + err.Error()
		log.Error(msg)
		discordSend(m.ChannelID, msg)
		return
	} else {
		race = v
	}

	// Find the discord ID of the active racer
	var activeRacerDiscordID string
	if race.ActivePlayer == 1 {
		activeRacerDiscordID = race.Racer1.DiscordID
	} else if race.ActivePlayer == 2 {
		activeRacerDiscordID = race.Racer2.DiscordID
	}

	// Get the Discord guild object
	var guild *discordgo.Guild
	if v, err := discord.Guild(discordGuildID); err != nil {
		msg := "Failed to get the Discord guild: " + err.Error()
		log.Error(msg)
		discordSend(m.ChannelID, msg)
		return
	} else {
		guild = v
	}

	var discordUser *discordgo.User
	for _, member := range guild.Members {
		if member.User.ID == activeRacerDiscordID {
			discordUser = member.User
			break
		}
	}
	if discordUser == nil {
		msg := "Failed to find the active player in the Discord server."
		log.Error(msg)
		discordSend(m.ChannelID, msg)
		return
	}

	m.Author = discordUser
	commandNo(m, args)
}
