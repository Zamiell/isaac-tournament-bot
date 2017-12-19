package main

import (
	"github.com/Zamiell/isaac-tournament-bot/src/models"
	"github.com/bwmarrin/discordgo"
)

// Get this racer from the database
// (and create an entry if it doesn't exist already)
func racerGet(u *discordgo.User) (models.Racer, error) {
	var racer models.Racer
	var exists bool
	if v, err := db.Racers.Exists(u.ID); err != nil {
		return racer, err
	} else {
		exists = v
	}

	if !exists {
		racer = models.Racer{
			DiscordID: u.ID,
			Username:  u.Username,
		}
		if err := db.Racers.Insert(racer); err != nil {
			return racer, err
		}
		return racer, nil
	}

	if racer, err := db.Racers.Get(u.ID); err != nil {
		return racer, err
	} else {
		return racer, nil
	}
}
