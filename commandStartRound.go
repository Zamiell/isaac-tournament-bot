package main

import (
	"encoding/json"
	"math"
	"time"

	"github.com/bwmarrin/discordgo"
)

func commandStartRound(m *discordgo.MessageCreate, args []string) {
	if !isAdmin(m) {
		return
	}

	// Go through all of the tournaments
	for _, tournament := range tournaments {
		startRound(m, tournament, false)
	}
}

func startRound(m *discordgo.MessageCreate, tournament Tournament, dryRun bool) {
	// Get the tournament from Challonge
	apiURL := "https://api.challonge.com/v1/tournaments/" + floatToString(tournament.ChallongeID) + ".json?"
	apiURL += "api_key=" + challongeAPIKey + "&include_participants=1&include_matches=1"
	var raw []byte
	if v, err := challongeGetJSON("GET", apiURL, nil); err != nil {
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
	jsonTournament := vMap["tournament"].(map[string]interface{})

	// Get the Discord roles
	var roles []*discordgo.Role
	if v, err := discordSession.GuildRoles(discordGuildID); err != nil {
		log.Fatal("Failed to get the roles for the guild: " + err.Error())
		return
	} else {
		roles = v
	}

	// Get all of the open matches
	foundMatches := false
	foundPlayers := make([]string, 0)
	var round string
	for _, v := range jsonTournament["matches"].([]interface{}) {
		vMap := v.(map[string]interface{})
		match := vMap["match"].(map[string]interface{})
		if match["state"] != "open" {
			continue
		}

		// Local variables
		foundMatches = true
		player1ID := match["player1_id"].(float64)
		player2ID := match["player2_id"].(float64)
		player1Name := challongeGetParticipantName(jsonTournament, player1ID)
		player2Name := challongeGetParticipantName(jsonTournament, player2ID)
		challongeMatchID := floatToString(match["id"].(float64))
		channelName := player1Name + "-vs-" + player2Name

		// Check to see if we have already created a channel for either of these players
		playerAlreadyHasChannel := false
		for _, p := range foundPlayers {
			if p == player1Name || p == player2Name {
				playerAlreadyHasChannel = true
				break
			}
		}
		if playerAlreadyHasChannel {
			log.Info("Skipping match \"" + channelName + "\" since one of the players already has a channel open.")
			continue
		}
		foundPlayers = append(foundPlayers, player1Name)
		foundPlayers = append(foundPlayers, player2Name)
		round = floatToString(match["round"].(float64))

		// Get the Discord guild object
		var guild *discordgo.Guild
		if v, err := discordSession.Guild(discordGuildID); err != nil {
			msg := "Failed to get the Discord guild: " + err.Error()
			log.Error(msg)
			discordSend(m.ChannelID, msg)
			return
		} else {
			guild = v
		}

		var members []*discordgo.Member
		if v, err := discordSession.GuildMembers(discordGuildID, "0", 1000); err != nil {
			msg := "Failed to get the Discord guild members: " + err.Error()
			log.Error(msg)
			discordSend(m.ChannelID, msg)
			return
		} else {
			members = v
		}

		// Find the Discord ID of the two racers and add them to the database if they are not already
		var racer1 *User
		var racer2 *User
		var team1DiscordID string
		var team2DiscordID string
		var discordUser1 *discordgo.User
		var discordUser2 *discordgo.User
		if tournament.Ruleset == "team" {
			// This is a team match, so we only need to find the team captain
			for _, role := range roles {
				if role.Name == player1Name {
					team1DiscordID = role.ID
				}
			}
			for _, role := range roles {
				if role.Name == player2Name {
					team2DiscordID = role.ID
				}
			}
			for _, member := range members {
				if stringInSlice(discordTeamCaptainRoleID, member.Roles) && stringInSlice(team1DiscordID, member.Roles) {
					discordUser1 = member.User
					break
				}
			}

			if discordUser1 == nil {
				msg := "Failed to find \"" + player1Name + "\" (the team captain) in the Discord server."
				log.Error(msg)
				discordSend(m.ChannelID, msg)
				return
			}
			for _, member := range members {
				if stringInSlice(discordTeamCaptainRoleID, member.Roles) && stringInSlice(team2DiscordID, member.Roles) {
					discordUser2 = member.User
					break
				}
			}
			if discordUser2 == nil {
				msg := "Failed to find \"" + player2Name + "\" (the team captain) in the Discord server."
				log.Error(msg)
				discordSend(m.ChannelID, msg)
				return
			}
			if v, err := userGet(discordUser1); err != nil {
				msg := "Failed to get the user from the database: " + err.Error()
				log.Error(msg)
				discordSend(m.ChannelID, msg)
				return
			} else {
				racer1 = v
			}
			if v, err := userGet(discordUser2); err != nil {
				msg := "Failed to get the user from the database: " + err.Error()
				log.Error(msg)
				discordSend(m.ChannelID, msg)
				return
			} else {
				racer2 = v
			}
		} else {
			// This is a 1v1 match
			discordUser1 = discordGetUserFromGuild(guild, player1Name)
			if discordUser1 == nil {
				msg := "Failed to find \"" + player1Name + "\" in the Discord server."
				log.Error(msg)
				discordSend(m.ChannelID, msg)
				return
			}
			discordUser2 = discordGetUserFromGuild(guild, player2Name)
			if discordUser2 == nil {
				msg := "Failed to find \"" + player2Name + "\" in this Discord server."
				log.Error(msg)
				discordSend(m.ChannelID, msg)
				return
			}
			if v, err := userGet(discordUser1); err != nil {
				msg := "Failed to get the user from the database: " + err.Error()
				log.Error(msg)
				discordSend(m.ChannelID, msg)
				return
			} else {
				racer1 = v
			}
			if v, err := userGet(discordUser2); err != nil {
				msg := "Failed to get the user from the database: " + err.Error()
				log.Error(msg)
				discordSend(m.ChannelID, msg)
				return
			} else {
				racer2 = v
			}
		}

		if dryRun {
			continue
		}

		// Create a channel for this match
		var channelID string
		if v, err := discordSession.GuildChannelCreate(discordGuildID, channelName, discordgo.ChannelTypeGuildText); err != nil {
			msg := "Failed to create the Discord channel of \"" + channelName + "\": " + err.Error()
			log.Error(msg)
			discordSend(m.ChannelID, msg)
			return
		} else {
			channelID = v.ID
		}

		buildsRemaining := make([]string, 0)
		for _, build := range builds {
			buildName := getBuildName(build)
			buildsRemaining = append(buildsRemaining, buildName)
		}

		// Create the race in the database
		race := Race{
			TournamentName:      tournament.Name,
			Racer1ChallongeID:   player1ID,
			Racer2ChallongeID:   player2ID,
			ChannelID:           channelID,
			ChannelName:         channelName,
			ChallongeURL:        tournament.ChallongeURL,
			ChallongeMatchID:    challongeMatchID,
			BracketRound:        round,
			State:               "initial",
			CharactersRemaining: characters,
			BuildsRemaining:     buildsRemaining,
			Racer1Bans:          numBans,
			Racer2Bans:          numBans,
			Racer1Vetos:         numVetos,
			Racer2Vetos:         numVetos,
		}
		if err := modals.Races.Insert(racer1.DiscordID, racer2.DiscordID, race); err != nil {
			msg := "Failed to create the race in the database: " + err.Error()
			log.Error(msg)
			discordSend(m.ChannelID, msg)
			return
		}

		// Put the channel in the correct category and give access to the two racers
		// (channels in this category have "Read Text Channels & See Voice Channels" disabled for everyone except for admins/casters/bots)
		permissionsReadWrite := int64(discordgo.PermissionViewChannel |
			discordgo.PermissionSendMessages |
			discordgo.PermissionEmbedLinks |
			discordgo.PermissionAttachFiles |
			discordgo.PermissionReadMessageHistory)
		var permissions = make([]*discordgo.PermissionOverwrite, 0)
		permissions = append(permissions,
			&discordgo.PermissionOverwrite{
				ID:   discordEveryoneRoleID,
				Type: discordgo.PermissionOverwriteTypeRole,
				Deny: permissionsReadWrite,
			},

			// Allow bots to see + talk in this channel
			&discordgo.PermissionOverwrite{
				ID:    discordBotRoleID,
				Type:  discordgo.PermissionOverwriteTypeRole,
				Allow: permissionsReadWrite,
			},

			// Allow all casters to see + talk in this channel
			&discordgo.PermissionOverwrite{
				ID:    discordCasterRoleID,
				Type:  discordgo.PermissionOverwriteTypeRole,
				Allow: permissionsReadWrite,
			})
		if tournament.Ruleset == "team" {
			permissions = append(permissions,
				&discordgo.PermissionOverwrite{
					ID:    team1DiscordID,
					Type:  discordgo.PermissionOverwriteTypeRole,
					Allow: permissionsReadWrite,
				},
				&discordgo.PermissionOverwrite{
					ID:    team2DiscordID,
					Type:  discordgo.PermissionOverwriteTypeRole,
					Allow: permissionsReadWrite,
				})
		} else {
			permissions = append(permissions,
				&discordgo.PermissionOverwrite{
					ID:    racer1.DiscordID,
					Type:  discordgo.PermissionOverwriteTypeMember,
					Allow: permissionsReadWrite,
				},
				&discordgo.PermissionOverwrite{
					ID:    racer2.DiscordID,
					Type:  discordgo.PermissionOverwriteTypeMember,
					Allow: permissionsReadWrite,
				})
		}
		if _, err := discordSession.ChannelEditComplex(channelID, &discordgo.ChannelEdit{
			PermissionOverwrites: permissions,
			ParentID:             tournament.DiscordCategoryID,
		}); err != nil {
			msg := "Failed to edit the permissions for the new channel: " + err.Error()
			log.Error(msg)
			discordSend(m.ChannelID, msg)
			return
		}

		// Find out if the racers have set their timezone
		msg := ""
		if racer1.Timezone.Valid {
			msg += discordUser1.Mention() + " has a timezone of: " + getTimezone(racer1.Timezone.String) + "\n"
		} else {
			msg += discordUser1.Mention() + ", your timezone is **not currently set**. Please set one with: `!timezone [timezone]`\n"
		}
		if racer2.Timezone.Valid {
			msg += discordUser2.Mention() + " has a timezone of: " + getTimezone(racer2.Timezone.String) + "\n"
		} else {
			msg += discordUser2.Mention() + ", your timezone is **not currently set**. Please set one with: `!timezone [timezone]`\n"
		}

		// Calculate the difference between the two timezones
		if racer1.Timezone.Valid && racer2.Timezone.Valid {
			loc1, _ := time.LoadLocation(racer1.Timezone.String)
			loc2, _ := time.LoadLocation(racer2.Timezone.String)
			_, offset1 := time.Now().In(loc1).Zone()
			_, offset2 := time.Now().In(loc2).Zone()
			if offset1 == offset2 {
				msg += "You both are in **the same timezone**. Great!\n"
			} else {
				difference := math.Abs(float64(offset1 - offset2))
				hours := difference / 3600
				msg += "You are **" + floatToString(hours) + " hours** away from each other.\n"
			}
		}
		msg += "\n"

		// Find out if the racers have set their stream URL
		if racer1.StreamURL.Valid {
			msg += discordUser1.Mention() + " has a stream of: <" + racer1.StreamURL.String + ">\n"
		} else {
			msg += discordUser1.Mention() + ", your stream is **not currently set**. Please set one with: `!stream [url]`\n"
		}
		if racer2.StreamURL.Valid {
			msg += discordUser2.Mention() + " has a stream of: <" + racer2.StreamURL.String + ">\n"
		} else {
			msg += discordUser2.Mention() + ", your stream is **not currently set**. Please set one with: `!stream [url]`\n"
		}
		msg += "\n"

		// Give the welcome message
		msg += "Please discuss the times that each of you are available to play this week.\n"
		if tournament.Ruleset == "team" {
			msg += discordUser1.Mention() + " and " + discordUser2.Mention() + " are the team captains; I will only listen to them.\n"
		}
		msg += "You can use suggest a time to your opponent with something like: `!time 6pm sat`\n"
		msg += "If they accept with `!timeok`, then the match will be officially scheduled."
		discordSend(channelID, msg)

		log.Info("Started race: " + channelName)
	}

	if dryRun {
		msg := "Tournament \"" + tournament.Name + "\" looks good."
		discordSend(m.ChannelID, msg)
		log.Info(msg)
		return
	}

	// Rename the channel category
	categoryName := "Round " + round + " - " + tournament.Ruleset
	if _, err := discordSession.ChannelEdit(tournament.DiscordCategoryID, categoryName); err != nil {
		msg := "Failed to rename the channel category: " + err.Error()
		log.Error(msg)
		discordSend(m.ChannelID, msg)
		return
	}

	if foundMatches {
		msg := "Round " + round + " channels created for tournament \"" + tournament.Name + "\"."
		discordSend(m.ChannelID, msg)
		log.Info(msg)
	} else {
		msg := "There are no open matches on the Challonge bracket for tournament \"" + tournament.Name + "\"."
		discordSend(m.ChannelID, msg)
		log.Info(msg)
	}
}
