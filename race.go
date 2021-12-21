package main

// Get this race from the database
func raceGet(channelID string) (*Race, error) {
	var race *Race
	if v, err := modals.Races.Get(channelID); err != nil {
		return race, err
	} else {
		race = v
	}

	// We also to fill in the "Racer1" and "Racer2" fields
	if v, err := modals.Users.GetFromUserID(race.Racer1ID); err != nil {
		return race, err
	} else {
		race.Racer1 = v
	}
	if v, err := modals.Users.GetFromUserID(race.Racer2ID); err != nil {
		return race, err
	} else {
		race.Racer2 = v
	}

	// We also have to fill in the "Casts" field
	if v, err := modals.Casts.GetAll(race.ChannelID); err != nil {
		return race, err
	} else {
		race.Casts = v
	}

	// We also have to fill in the "Caster" field(s), if any
	for _, cast := range race.Casts {
		if v, err := modals.Users.GetFromUserID(cast.CasterID); err != nil {
			return race, err
		} else {
			cast.Caster = v
		}
	}

	return race, nil
}
