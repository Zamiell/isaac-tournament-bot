package models

import (
	"database/sql"
	"strconv"
	"time"

	"github.com/go-sql-driver/mysql"
)

type Races struct{}

// State 0 is freshly created before both players have confirmed a scheduled time
// State 1 is confirmed but before it starts
// State 2 is banning characters (5 minutes before starting)
// State 3 is banning items
// State 4 is in progress
// State 5 is completed
type Race struct {
	Racer1            Racer
	Racer1ID          int
	Racer2            Racer
	Racer2ID          int
	ChannelID         string
	BracketRound      string
	State             int
	DatetimeScheduled mysql.NullTime
	Caster            Racer
	CasterID          sql.NullInt64
	CasterP1          bool
	CasterP2          bool
	ActivePlayer      int
	Characters        string
	Builds            string
}

func (r *Race) Name() string {
	return "round-" + r.BracketRound + "-" + r.Racer1.Username + "-vs-" + r.Racer2.Username
}

func (*Races) Insert(racer1DiscordID string, racer2DiscordID string, race Race) error {
	var stmt *sql.Stmt
	if v, err := db.Prepare(`
		INSERT INTO races (
			racer1,
			racer2,
			channel_id,
			bracket_round,
			characters,
			builds
		) VALUES (
			(SELECT id FROM racers WHERE discord_id = ?),
			(SELECT id FROM racers WHERE discord_id = ?),
			?,
			?,
			?,
			?
		)
	`); err != nil {
		return err
	} else {
		stmt = v
	}
	defer stmt.Close()

	if _, err := stmt.Exec(
		racer1DiscordID,
		racer2DiscordID,
		race.ChannelID,
		race.BracketRound,
		race.Characters,
		race.Builds,
	); err != nil {
		return err
	}

	return nil
}

func (*Races) Delete(channelID string) error {
	var stmt *sql.Stmt
	if v, err := db.Prepare(`
		DELETE FROM races
		WHERE channel_id = ?
	`); err != nil {
		return err
	} else {
		stmt = v
	}
	defer stmt.Close()

	if _, err := stmt.Exec(channelID); err != nil {
		return err
	}

	return nil
}

func (*Races) Get(channelID string) (Race, error) {
	var race Race
	if err := db.QueryRow(`
		SELECT
			racer1,
			racer2,
			channel_id,
			bracket_round,
			state,
			datetime_scheduled,
			caster,
			caster_p1,
			caster_p2,
			active_player,
			characters,
			builds
		FROM races
		WHERE channel_id = ?
	`, channelID).Scan(
		&race.Racer1ID,
		&race.Racer2ID,
		&race.ChannelID,
		&race.BracketRound,
		&race.State,
		&race.DatetimeScheduled,
		&race.CasterID,
		&race.CasterP1,
		&race.CasterP2,
		&race.ActivePlayer,
		&race.Characters,
		&race.Builds,
	); err != nil {
		return race, err
	}

	return race, nil
}

func (*Races) SetState(channelID string, state int) error {
	var stmt *sql.Stmt
	if v, err := db.Prepare(`
		UPDATE races
		SET state = ?
		WHERE channel_id = ?
	`); err != nil {
		return err
	} else {
		stmt = v
	}
	defer stmt.Close()

	if _, err := stmt.Exec(state, channelID); err != nil {
		return err
	}

	return nil
}

func (*Races) SetDatetimeScheduled(channelID string, datetimeScheduled time.Time) error {
	var stmt *sql.Stmt
	if v, err := db.Prepare(`
		UPDATE races
		SET datetime_scheduled = ?
		WHERE channel_id = ?
	`); err != nil {
		return err
	} else {
		stmt = v
	}
	defer stmt.Close()

	if _, err := stmt.Exec(datetimeScheduled, channelID); err != nil {
		return err
	}

	return nil
}

func (*Races) UnsetDatetimeScheduled(channelID string) error {
	var stmt *sql.Stmt
	if v, err := db.Prepare(`
		UPDATE races
		SET datetime_scheduled = NULL
		WHERE channel_id = ?
	`); err != nil {
		return err
	} else {
		stmt = v
	}
	defer stmt.Close()

	if _, err := stmt.Exec(channelID); err != nil {
		return err
	}

	return nil
}

func (*Races) SetCaster(channelID string, casterID string) error {
	var stmt *sql.Stmt
	if v, err := db.Prepare(`
		UPDATE races
		SET caster = (SELECT id FROM racers WHERE discord_id = ?)
		WHERE channel_id = ?
	`); err != nil {
		return err
	} else {
		stmt = v
	}
	defer stmt.Close()

	if _, err := stmt.Exec(casterID, channelID); err != nil {
		return err
	}

	return nil
}

func (*Races) UnsetCaster(channelID string) error {
	var stmt *sql.Stmt
	if v, err := db.Prepare(`
		UPDATE races
		SET
			caster = NULL,
			caster_p1 = 0,
			caster_p2 = 0,
		WHERE channel_id = ?
	`); err != nil {
		return err
	} else {
		stmt = v
	}
	defer stmt.Close()

	if _, err := stmt.Exec(channelID); err != nil {
		return err
	}

	return nil
}

func (*Races) SetCasterApproval(channelID string, playerNum int) error {
	var stmt *sql.Stmt
	if v, err := db.Prepare(`
		UPDATE races
		SET caster_p` + strconv.Itoa(playerNum) + ` = 0
		WHERE channel_id = ?
	`); err != nil {
		return err
	} else {
		stmt = v
	}
	defer stmt.Close()

	if _, err := stmt.Exec(channelID); err != nil {
		return err
	}

	return nil
}
