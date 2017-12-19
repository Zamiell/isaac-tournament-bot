package main

import (
	"github.com/Zamiell/isaac-tournament-bot/src/models"
)

// Get this race from the database
func raceGet(channelID string) (models.Race, error) {
	var race models.Race
	if v, err := db.Races.Get(channelID); err != nil {
		return race, err
	} else {
		race = v
	}

	// Now we have to fill in the "Racer1", "Racer2", and "Caster" fields
	if v, err := db.Racers.GetID(race.Racer1ID); err != nil {
		return race, err
	} else {
		race.Racer1 = v
	}

	if v, err := db.Racers.GetID(race.Racer2ID); err != nil {
		return race, err
	} else {
		race.Racer2 = v
	}

	if race.CasterID.Valid {
		if v, err := db.Racers.GetID(int(race.CasterID.Int64)); err != nil {
			return race, err
		} else {
			race.Caster = v
		}
	}

	return race, nil
}
