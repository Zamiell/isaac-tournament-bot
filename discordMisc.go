package main

import (
	"errors"

	"github.com/bwmarrin/discordgo"
)

func discordSend(channelID string, msg string) {
	if _, err := discordSession.ChannelMessageSend(channelID, msg); err != nil {
		log.Error("Failed to send \"" + msg + "\" to \"" + channelID + "\": " + err.Error())
		return
	}
}

// Other calls to "discordSession.GuildMembers" should be refactored here, but I don't have the
// heart to do this right now.
func getDiscordMembers() ([]*discordgo.Member, error) {
	members, err := discordSession.GuildMembers(discordGuildID, "0", 1000)

	if err != nil {
		return nil, errors.New("Failed to get the Discord guild members: " + err.Error())
	}

	return members, nil
}

func getDiscordNameByID(members []*discordgo.Member, id string) (string, error) {
	for _, member := range members {
		if member.User.ID == id {
			if member.Nick != "" {
				return member.Nick, nil
			}

			return member.User.Username, nil
		}
	}

	return "", errors.New("Failed to find a Discord member matching ID: " + id)
}

func getDiscordUserByID(members []*discordgo.Member, id string) *discordgo.User {
	for _, member := range members {
		if member.User.ID == id {
			return member.User
		}
	}

	return nil
}

func getDiscordUserByName(members []*discordgo.Member, name string) *discordgo.User {
	for _, member := range members {
		username := member.Nick
		if username == "" {
			username = member.User.Username
		}

		if username == name {
			return member.User
		}
	}

	return nil
}

func getDiscordRoleIDByName(roles []*discordgo.Role, name string) (string, error) {
	for _, role := range roles {
		if role.Name == name {
			return role.ID, nil
		}
	}

	return "", errors.New("Failed to find a Discord role matching name: " + name)
}

func getDiscordTeamCaptain(members []*discordgo.Member, teamRoleID string) *discordgo.User {
	for _, member := range members {
		if stringInSlice(discordTeamCaptainRoleID, member.Roles) && stringInSlice(teamRoleID, member.Roles) {
			return member.User
		}
	}

	return nil
}
