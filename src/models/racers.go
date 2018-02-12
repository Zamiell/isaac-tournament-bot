package models

import (
	"database/sql"
)

type Racers struct{}

type Racer struct {
	DiscordID string
	Username  string
	// Matches the TZ column of this page:
	// https://en.wikipedia.org/wiki/List_of_tz_database_time_zones
	Timezone  sql.NullString
	StreamURL sql.NullString
}

func (r *Racer) Mention() string {
	return "<@" + r.DiscordID + ">"
}

func (*Racers) Exists(discordID string) (bool, error) {
	var id int
	if err := db.QueryRow(`
		SELECT id
		FROM tournament_racers
		WHERE discord_id = ?
	`, discordID).Scan(&id); err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return true, nil
}

func (*Racers) Insert(racer Racer) error {
	var stmt *sql.Stmt
	if v, err := db.Prepare(`
		INSERT INTO tournament_racers (
			discord_id,
			username
		) VALUES (
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
		racer.DiscordID,
		racer.Username,
	); err != nil {
		return err
	}

	return nil
}

func (*Racers) Get(discordID string) (Racer, error) {
	var racer Racer
	if err := db.QueryRow(`
		SELECT
			discord_id,
			username,
			timezone,
			stream_url
		FROM tournament_racers
		WHERE discord_id = ?
	`, discordID).Scan(
		&racer.DiscordID,
		&racer.Username,
		&racer.Timezone,
		&racer.StreamURL,
	); err != nil {
		return racer, err
	}

	return racer, nil
}

func (*Racers) GetID(racerID int) (Racer, error) {
	var racer Racer
	if err := db.QueryRow(`
		SELECT
			discord_id,
			username,
			timezone,
			stream_url
		FROM tournament_racers
		WHERE id = ?
	`, racerID).Scan(
		&racer.DiscordID,
		&racer.Username,
		&racer.Timezone,
		&racer.StreamURL,
	); err != nil {
		return racer, err
	}

	return racer, nil
}

func (*Racers) SetUsername(discordID string, username string) error {
	var stmt *sql.Stmt
	if v, err := db.Prepare(`
		UPDATE tournament_racers
		SET username = ?
		WHERE discord_id = ?
	`); err != nil {
		return err
	} else {
		stmt = v
	}
	defer stmt.Close()

	if _, err := stmt.Exec(username, discordID); err != nil {
		return err
	}

	return nil
}

func (*Racers) SetTimezone(discordID string, timezone string) error {
	var stmt *sql.Stmt
	if v, err := db.Prepare(`
		UPDATE tournament_racers
		SET timezone = ?
		WHERE discord_id = ?
	`); err != nil {
		return err
	} else {
		stmt = v
	}
	defer stmt.Close()

	if _, err := stmt.Exec(timezone, discordID); err != nil {
		return err
	}

	return nil
}

func (*Racers) SetStreamURL(discordID string, streamURL string) error {
	var stmt *sql.Stmt
	if v, err := db.Prepare(`
		UPDATE tournament_racers
		SET stream_url = ?
		WHERE discord_id = ?
	`); err != nil {
		return err
	} else {
		stmt = v
	}
	defer stmt.Close()

	if _, err := stmt.Exec(streamURL, discordID); err != nil {
		return err
	}

	return nil
}
