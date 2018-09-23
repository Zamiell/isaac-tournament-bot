package main

import (
	"github.com/Zamiell/isaac-tournament-bot/src/models"
)

// Get this race from the database
func raceGet(channelID string) (*models.Race, error) {
	var race *models.Race
	if v, err := db.Races.Get(channelID); err != nil {
		return race, err
	} else {
		race = v
	}

	// We also to fill in the "Racer1" and "Racer2" fields
	if v, err := db.Users.GetFromUserID(race.Racer1ID); err != nil {
		return race, err
	} else {
		race.Racer1 = v
	}
	if v, err := db.Users.GetFromUserID(race.Racer2ID); err != nil {
		return race, err
	} else {
		race.Racer2 = v
	}

	// We also have to fill in the "Casts" field
	if v, err := db.Casts.GetAll(race.ChannelID); err != nil {
		return race, err
	} else {
		race.Casts = v
	}

	// We also have to fill in the "Caster" field(s), if any
	for _, cast := range race.Casts {
		if v, err := db.Users.GetFromUserID(cast.CasterID); err != nil {
			return race, err
		} else {
			cast.Caster = v
		}
	}

	return race, nil
}
