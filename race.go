package main

import (
	"database/sql"
)

type Race struct {
	TournamentName      string
	Racer1ID            int     // The "tournament_users" database ID
	Racer1ChallongeID   float64 // The "participant" ID; needed to automatically set the winner through the Challonge API
	Racer1              *User
	Racer2ID            int     // The "tournament_users" database ID
	Racer2ChallongeID   float64 // The "participant" ID; needed to automatically set the winner through the Challonge API
	Racer2              *User
	ChannelID           string // The Discord channel ID that was automatically created for this race
	ChannelName         string // The Discord channel name that was automatically created for this race
	ChallongeURL        string // The suffix of the Challonge URL for this tournament
	ChallongeMatchID    string
	BracketRound        string
	State               RaceState
	DatetimeScheduled   sql.NullTime
	ActiveRacer         int
	CharactersRemaining []string
	Characters          []string
	BuildsRemaining     []string
	Builds              []string
	Racer1Bans          int
	Racer2Bans          int
	Racer1Vetos         int
	Racer2Vetos         int
	NumVoted            int
	Casts               []*Cast
}

func (r *Race) Name() string {
	return r.Racer1.Username + "-vs-" + r.Racer2.Username
}

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
