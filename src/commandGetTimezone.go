package main

import (
	"github.com/Zamiell/isaac-tournament-bot/src/models"
	"github.com/bwmarrin/discordgo"
)

func commandGetTimezone(m *discordgo.MessageCreate, args []string) {
	if len(args) != 1 {
		commandGetTimezonePrint(m)
		return
	}

	// Get all of the users in the guild
	var guild *discordgo.Guild
	if v, err := discord.Guild(discordGuildID); err != nil {
		msg := "Failed to get the Discord guild: " + err.Error()
		log.Error(msg)
		discordSend(m.ChannelID, msg)
		return
	} else {
		guild = v
	}

	// Find the discord ID of the player
	var discordUser *discordgo.User
	for _, member := range guild.Members {
		username := member.Nick
		if username == "" {
			username = member.User.Username
		}
		if username == args[0] {
			discordUser = member.User
			break
		}
	}
	if discordUser == nil {
		msg := "Failed to find \"" + args[0] + "\" in the Discord server."
		log.Error(msg)
		discordSend(m.ChannelID, msg)
		return
	}

	var racer models.Racer
	if v, err := racerGet(discordUser); err != nil {
		msg := "Failed to get the racer from the database: " + err.Error()
		log.Error(msg)
		discordSend(m.ChannelID, msg)
		return
	} else {
		racer = v
	}

	msg := "The timezone for " + discordUser.Username + " is "
	if racer.Timezone.Valid {
		msg += "currently set to: **" + racer.Timezone.String + "**"
	} else {
		msg += "**not currently set**."
	}
	discordSend(m.ChannelID, msg)
}

func commandGetTimezonePrint(m *discordgo.MessageCreate) {
	msg := "Get another player's timezone with: `!gettimezone [username]`\n"
	msg += "For example: `!gettimezone Zamiel`"
	discordSend(m.ChannelID, msg)
}
