package main

import (
	"errors"

	"github.com/bwmarrin/discordgo"
)

// Get this user from the database
// (and create an entry if it does not exist already)
func userGet(u *discordgo.User) (*User, error) {
	var user *User

	// Get the Discord guild object
	var guild *discordgo.Guild
	if v, err := discordSession.Guild(discordGuildID); err != nil {
		return user, err
	} else {
		guild = v
	}

	// Get their custom nickname for the Discord server, if any
	username := ""
	for _, member := range guild.Members {
		if member.User.ID != u.ID {
			continue
		}

		username = member.Nick
		if username == "" {
			username = member.User.Username
		}
	}
	if username == "" {
		return user, errors.New("Failed to find \"" + u.Username + "\" in the Discord server.")
	}

	// See if this user exists in the database already
	var exists bool
	if v, err := modals.Users.Exists(u.ID); err != nil {
		return user, err
	} else {
		exists = v
	}

	// This Discord ID already exists in the database, so return it
	if exists {
		if v, err := modals.Users.GetFromDiscordID(u.ID); err != nil {
			return user, err
		} else {
			user = v
		}

		if user.Username == username {
			// Their username in the database matches the Discord nickname
			return user, nil
		}

		// Their Discord nickname has changed since they were added to the database,
		// so we need to update it
		user.Username = username
		if err := modals.Users.SetUsername(u.ID, username); err != nil {
			return user, err
		} else {
			return user, nil
		}
	}

	// This Discord ID does not exist in the database, so create it
	user = &User{
		DiscordID: u.ID,
		Username:  u.Username,
	}
	err := modals.Users.Insert(user)
	return user, err
}
