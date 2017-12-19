package models

import (
	"database/sql"
)

type Racers struct{}

type Racer struct {
	DiscordID string
	Username  string
	Timezone  sql.NullInt64 // Expressed in hours from GMT; e.g. -5 represents GMT-5
	StreamURL sql.NullString
}

func (r *Racer) Mention() string {
	return "<@" + r.DiscordID + ">"
}

func (*Racers) Exists(discordID string) (bool, error) {
	var id int
	if err := db.QueryRow(`
		SELECT id
		FROM racers
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
		INSERT INTO racers (
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
		FROM racers
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
		FROM racers
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

func (*Racers) SetTimeZone(discordID string, timezone int) error {
	var stmt *sql.Stmt
	if v, err := db.Prepare(`
		UPDATE racers
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
		UPDATE racers
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
