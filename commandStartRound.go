package main

import (
	"encoding/json"
	"errors"

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

	// Get the Discord guild members
	var discordMembers []*discordgo.Member
	if v, err := discordSession.GuildMembers(discordGuildID, "0", 1000); err != nil {
		msg := "Failed to get the Discord guild members: " + err.Error()
		log.Error(msg)
		discordSend(m.ChannelID, msg)
		return
	} else {
		discordMembers = v
	}

	// Get the Discord roles
	var discordRoles []*discordgo.Role
	if v, err := discordSession.GuildRoles(discordGuildID); err != nil {
		msg := "Failed to get the roles for the guild: " + err.Error()
		log.Error(msg)
		discordSend(m.ChannelID, msg)
		return
	} else {
		discordRoles = v
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

		var racer1DiscordID string
		var racer2DiscordID string
		if v1, v2, err := getDiscordIDsForMatch(tournament, discordMembers, discordRoles, player1Name, player2Name); err != nil {
			log.Error(err)
			discordSend(m.ChannelID, err.Error())
			return
		} else {
			racer1DiscordID = v1
			racer2DiscordID = v2
		}

		if dryRun {
			continue
		}

		var channelID string
		if v, err := createDiscordChannelForMatch(
			channelName,
			tournament.DiscordCategoryID,
			racer1DiscordID,
			racer2DiscordID,
		); err != nil {
			log.Error(err)
			discordSend(m.ChannelID, err.Error())
			return
		} else {
			channelID = v
		}

		buildsRemaining := make([]string, 0)
		for _, build := range builds {
			buildName := getBuildName(build)

			// The 0th element of the "builds.json" file is blank
			if buildName != "" {
				buildsRemaining = append(buildsRemaining, buildName)
			}
		}

		// Create the race in the database
		race := &Race{
			TournamentName:      tournament.Name,
			Racer1ChallongeID:   player1ID,
			Racer2ChallongeID:   player2ID,
			ChannelID:           channelID,
			ChannelName:         channelName,
			ChallongeURL:        tournament.ChallongeURL,
			ChallongeMatchID:    challongeMatchID,
			BracketRound:        round,
			State:               RaceStateInitial,
			CharactersRemaining: characters,
			BuildsRemaining:     buildsRemaining,
			Racer1Bans:          numBans,
			Racer2Bans:          numBans,
			Racer1Vetos:         numVetos,
			Racer2Vetos:         numVetos,
		}
		if err := modals.Races.Insert(racer1DiscordID, racer2DiscordID, race); err != nil {
			msg := "Failed to create the race in the database: " + err.Error()
			log.Error(msg)
			discordSend(m.ChannelID, msg)
			return
		}

		// We re-get the race in the database so that the racer fields are filled in properly
		if v, err := getRace(m.ChannelID); err != nil {
			msg := "Failed to get the race from the database: " + err.Error()
			log.Error(msg)
			discordSend(m.ChannelID, msg)
			return
		} else {
			race = v
		}

		// Send the introductory messages for the Discord channel
		announceStatus(m, race, true)

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

// Find the Discord ID of the two racers and add them to the database if they are not already
func getDiscordIDsForMatch(
	tournament Tournament,
	members []*discordgo.Member,
	roles []*discordgo.Role,
	player1Name string,
	player2Name string,
) (string, string, error) {
	if tournament.Ruleset == "team" {
		return getDiscordIDsForMatchTeam(members, roles, player1Name, player2Name)
	}

	return getDiscordIDsForMatch1v1(members, roles, player1Name, player2Name)
}

// Since this is a team match, we only need to find the team captain
func getDiscordIDsForMatchTeam(
	members []*discordgo.Member,
	roles []*discordgo.Role,
	team1CaptainName string,
	team2CaptainName string,
) (string, string, error) {
	var team1DiscordID string
	if v, err := getDiscordRoleIDByName(roles, team1CaptainName); err != nil {
		return "", "", err
	} else {
		team1DiscordID = v
	}

	var team2DiscordID string
	if v, err := getDiscordRoleIDByName(roles, team2CaptainName); err != nil {
		return "", "", err
	} else {
		team2DiscordID = v
	}

	team1CaptainDiscordUser := getDiscordTeamCaptain(members, team1DiscordID)
	if team1CaptainDiscordUser == nil {
		return "", "", errors.New("Failed to find \"" + team1CaptainName + "\" (the team captain) in the Discord server.")
	}

	team2CaptainDiscordUser := getDiscordTeamCaptain(members, team2DiscordID)
	if team2CaptainDiscordUser == nil {
		return "", "", errors.New("Failed to find \"" + team2CaptainName + "\" (the team captain) in the Discord server.")
	}

	if _, err := userGet(team1CaptainDiscordUser); err != nil {
		return "", "", errors.New("Failed to insert \"" + team1CaptainName + "\" (the team captain) into the database: " + err.Error())
	}

	if _, err := userGet(team2CaptainDiscordUser); err != nil {
		return "", "", errors.New("Failed to insert \"" + team2CaptainName + "\" (the team captain) into the database: " + err.Error())
	}

	return team1CaptainDiscordUser.ID, team2CaptainDiscordUser.ID, nil
}

func getDiscordIDsForMatch1v1(
	members []*discordgo.Member,
	roles []*discordgo.Role,
	player1Name string,
	player2Name string,
) (string, string, error) {
	discordUser1 := getDiscordUserByName(members, player1Name)
	if discordUser1 == nil {
		return "", "", errors.New("Failed to find \"" + player1Name + "\" in the Discord server.")
	}

	discordUser2 := getDiscordUserByName(members, player2Name)
	if discordUser2 == nil {
		return "", "", errors.New("Failed to find \"" + player2Name + "\" in the Discord server.")
	}

	var racer1 *User
	if v, err := userGet(discordUser1); err != nil {
		return "", "", errors.New("Failed to get \"" + player1Name + "\" from the database:" + err.Error())
	} else {
		racer1 = v
	}

	var racer2 *User
	if v, err := userGet(discordUser2); err != nil {
		return "", "", errors.New("Failed to get \"" + player1Name + "\" from the database:" + err.Error())
	} else {
		racer2 = v
	}

	return racer1.DiscordID, racer2.DiscordID, nil
}

func createDiscordChannelForMatch(
	channelName string,
	categoryID string,
	racer1DiscordID string,
	racer2DiscordID string,
) (string, error) {
	var channelID string
	if v, err := discordSession.GuildChannelCreate(discordGuildID, channelName, discordgo.ChannelTypeGuildText); err != nil {
		return "", errors.New("Failed to create the Discord channel of \"" + channelName + "\": " + err.Error())
	} else {
		channelID = v.ID
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

	permissions = append(permissions,
		&discordgo.PermissionOverwrite{
			ID:    racer1DiscordID,
			Type:  discordgo.PermissionOverwriteTypeMember,
			Allow: permissionsReadWrite,
		},
		&discordgo.PermissionOverwrite{
			ID:    racer2DiscordID,
			Type:  discordgo.PermissionOverwriteTypeMember,
			Allow: permissionsReadWrite,
		})
	if _, err := discordSession.ChannelEditComplex(channelID, &discordgo.ChannelEdit{
		PermissionOverwrites: permissions,
		ParentID:             categoryID,
	}); err != nil {
		return "", errors.New("Failed to edit the permissions for the new channel of \"" + channelName + "\": " + err.Error())
	}

	return channelID, nil
}
