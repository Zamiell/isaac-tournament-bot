package main

import (
	"encoding/json"
	"math"
	"strings"

	"github.com/Zamiell/isaac-tournament-bot/src/models"
	"github.com/bwmarrin/discordgo"
)

func commandRoundStart(m *discordgo.MessageCreate, args []string) {
	// Get the tournament from Challonge
	apiURL := "https://api.challonge.com/v1/tournaments/" + floatToString(challongeTournamentID) + ".json?"
	apiURL += "api_key=" + challongeAPIKey + "&include_participants=1&include_matches=1"
	var raw []byte
	if v, err := challongeGetJSON(apiURL); err != nil {
		msg := "Failed to get the tournament from Challonge: " + err.Error()
		log.Error(msg)
		discordSend(m.ChannelID, msg)
		return
	} else {
		raw = v
	}

	vMap := make(map[string]interface{})
	if err := json.Unmarshal(raw, &vMap); err != nil {
		msg := "Failed to unmarshal the Challonge JSON: " + err.Error()
		log.Error(msg)
		discordSend(m.ChannelID, msg)
		return
	}
	tournament := vMap["tournament"].(map[string]interface{})

	// Get all of the open matches
	foundMatches := false
	var round string
	for _, v := range tournament["matches"].([]interface{}) {
		vMap := v.(map[string]interface{})
		match := vMap["match"].(map[string]interface{})
		if match["state"] != "open" {
			continue
		}

		// Local variables
		foundMatches = true
		player1Name := challongeGetParticipantName(tournament, match["player1_id"].(float64))
		player2Name := challongeGetParticipantName(tournament, match["player2_id"].(float64))
		round = floatToString(match["round"].(float64))
		channelName := "round-" + round + "-" + player1Name + "-vs-" + player2Name
		log.Info("Doing actions for:", channelName)

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

		// Find the discord ID of the two players and add them to the database if they are not already
		var discord1 *discordgo.User
		for _, member := range guild.Members {
			username := member.Nick
			if username == "" {
				username = member.User.Username
			}
			if username == player1Name {
				discord1 = member.User
				break
			}
		}
		if discord1 == nil {
			msg := "Failed to find \"" + player1Name + "\" in the Discord server."
			log.Error(msg)
			discordSend(m.ChannelID, msg)
			return
		}

		var discord2 *discordgo.User
		for _, member := range guild.Members {
			username := member.Nick
			if username == "" {
				username = member.User.Username
			}
			if username == player2Name {
				discord2 = member.User
				break
			}
		}
		if discord2 == nil {
			msg := "Failed to find \"" + player2Name + "\" in this Discord server."
			log.Error(msg)
			discordSend(m.ChannelID, msg)
			return
		}

		var racer1 models.Racer
		if v, err := racerGet(discord1); err != nil {
			msg := "Failed to get the racer from the database: " + err.Error()
			log.Error(msg)
			discordSend(m.ChannelID, msg)
			return
		} else {
			racer1 = v
		}

		var racer2 models.Racer
		if v, err := racerGet(discord2); err != nil {
			msg := "Failed to get the racer from the database: " + err.Error()
			log.Error(msg)
			discordSend(m.ChannelID, msg)
			return
		} else {
			racer2 = v
		}

		// Create a channel for this match
		var channelID string
		if v, err := discord.GuildChannelCreate(discordGuildID, channelName, "text"); err != nil {
			msg := "Failed to create the Discord channel of \"" + channelName + "\": " + err.Error()
			log.Error(msg)
			discordSend(m.ChannelID, msg)
			return
		} else {
			channelID = v.ID
		}

		// Create the race in the database
		race := models.Race{
			ChannelID:    channelID,
			BracketRound: round,
			Characters:   strings.Join(characters, ","),
			Builds:       strings.Join(builds, ","),
		}
		log.Debug("-------------")
		log.Debug(discord2.ID)
		log.Debug(discord2.Username)
		log.Debug("-------------")
		if err := db.Races.Insert(racer1.DiscordID, racer2.DiscordID, race); err != nil {
			msg := "Failed to create the race in the database: " + err.Error()
			log.Error(msg)
			discordSend(m.ChannelID, msg)
			return
		}

		// Find out if the players have set their timezone
		msg := ""
		if racer1.Timezone.Valid {
			msg += discord1.Mention() + " has a timezone of: " + getTimezone(racer1.Timezone.Int64) + "\n"
		} else {
			msg += discord1.Mention() + ", your timezone is not currently set. Please set one with: `!timezone [timezone]`\n"
		}
		if racer2.Timezone.Valid {
			msg += discord2.Mention() + " has a timezone of: " + getTimezone(racer2.Timezone.Int64) + "\n"
		} else {
			msg += discord2.Mention() + ", your timezone is not currently set. Please set one with: `!timezone [timezone]`\n"
		}
		if racer1.Timezone.Valid && racer2.Timezone.Valid {
			if racer1.Timezone.Int64 == racer2.Timezone.Int64 {
				msg += "You both are in **the same timezone**. Great!\n"
			} else {
				difference := math.Abs(float64(racer1.Timezone.Int64 - racer2.Timezone.Int64))
				msg += "You are **" + floatToString(difference) + " hours** away from each other.\n"
			}
		}
		msg += "\n"

		// Find out if the players have set their stream URL
		if racer1.StreamURL.Valid {
			msg += discord1.Mention() + " has a stream of: <" + racer1.StreamURL.String + ">\n"
		} else {
			msg += discord1.Mention() + ", your stream is not currently set. Please set one with: `!stream [url]`\n"
		}
		if racer2.Timezone.Valid {
			msg += discord2.Mention() + " has a stream of: <" + racer1.StreamURL.String + ">\n"
		} else {
			msg += discord2.Mention() + ", your stream is not currently set. Please set one with: `!stream [url]`\n"
		}
		msg += "\n"

		// Give the welcome message
		msg += "Please communicate with your opponent the times that you are available to play over the course of the week.\n"
		msg += "You can use the bot to suggest a time to your opponent: `!schedule 02/06/2018 22:00`\n"
		msg += "If they accept with `!confirm`, then the match will be officially scheduled."
		discordSend(channelID, msg)
	}

	if foundMatches {
		discordSend(m.ChannelID, "Round "+round+" channels created.")
	} else {
		discordSend(m.ChannelID, "I was not able to find any open matches on the Challonge bracket.")
	}
}
