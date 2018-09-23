package main

import (
	"database/sql"

	"github.com/Zamiell/isaac-tournament-bot/src/models"
	"github.com/bwmarrin/discordgo"
)

func commandForceYes(m *discordgo.MessageCreate, args []string) {
	if !isAdmin(m) {
		return
	}

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

	// Find the discord ID of the active racer
	var activeRacerDiscordID string
	if race.ActiveRacer == 1 {
		activeRacerDiscordID = race.Racer1.DiscordID
	} else if race.ActiveRacer == 2 {
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
		msg := "Failed to find the active racer in the Discord server."
		log.Error(msg)
		discordSend(m.ChannelID, msg)
		return
	}

	m.Author = discordUser
	commandYes(m, args)
}
