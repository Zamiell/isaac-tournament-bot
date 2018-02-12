package main

import (
	"errors"

	"github.com/Zamiell/isaac-tournament-bot/src/models"
	"github.com/bwmarrin/discordgo"
)

// Get this racer from the database
// (and create an entry if it doesn't exist already)
func racerGet(u *discordgo.User) (models.Racer, error) {
	var racer models.Racer

	// Get the Discord guild object
	var guild *discordgo.Guild
	if v, err := discord.Guild(discordGuildID); err != nil {
		return racer, err
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
		return racer, errors.New("Failed to find \"" + u.Username + "\" in the Discord server.")
	}

	// See if this racer exists in the database already
	var exists bool
	if v, err := db.Racers.Exists(u.ID); err != nil {
		return racer, err
	} else {
		exists = v
	}

	// This Discord ID already exists in the database, so return it
	if exists {
		if v, err := db.Racers.Get(u.ID); err != nil {
			return racer, err
		} else {
			racer = v
		}

		if racer.Username == username {
			// Their username in the database matches the Discord nickname
			return racer, nil
		}

		// Their Discord nickname has changed since they were added to the database,
		// so we need to update it
		racer.Username = username
		if err := db.Racers.SetUsername(u.ID, username); err != nil {
			return racer, err
		} else {
			return racer, nil
		}
	}

	// This Discord ID does not exist in the database, so create it
	racer = models.Racer{
		DiscordID: u.ID,
		Username:  u.Username,
	}
	if err := db.Racers.Insert(racer); err != nil {
		return racer, err
	}
	return racer, nil
}
