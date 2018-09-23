package models

import (
	"database/sql"
)

type Users struct{}

type User struct {
	DiscordID string
	Username  string
	// Matches the TZ column of this page:
	// https://en.wikipedia.org/wiki/List_of_tz_database_time_zones
	Timezone       sql.NullString
	StreamURL      sql.NullString
	CasterAlwaysOk bool
}

func (u *User) Mention() string {
	return "<@" + u.DiscordID + ">"
}

func (*Users) Exists(discordID string) (bool, error) {
	var id int
	if err := db.QueryRow(`
		SELECT id
		FROM tournament_users
		WHERE discord_id = ?
	`, discordID).Scan(&id); err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return true, nil
}

func (*Users) Insert(user *User) error {
	var stmt *sql.Stmt
	if v, err := db.Prepare(`
		INSERT INTO tournament_users (
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

	_, err := stmt.Exec(
		user.DiscordID,
		user.Username,
	)
	return err
}

func (*Users) GetFromDiscordID(discordID string) (*User, error) {
	var user *User
	err := db.QueryRow(`
		SELECT
			discord_id,
			username,
			timezone,
			stream_url,
			caster_always_ok
		FROM tournament_users
		WHERE discord_id = ?
	`, discordID).Scan(
		&user.DiscordID,
		&user.Username,
		&user.Timezone,
		&user.StreamURL,
		&user.CasterAlwaysOk,
	)
	return user, err
}

func (*Users) GetFromUserID(userID int) (*User, error) {
	var user *User
	err := db.QueryRow(`
		SELECT
			discord_id,
			username,
			timezone,
			stream_url,
			caster_always_ok
		FROM tournament_users
		WHERE id = ?
	`, userID).Scan(
		&user.DiscordID,
		&user.Username,
		&user.Timezone,
		&user.StreamURL,
		&user.CasterAlwaysOk,
	)
	return user, err
}

func (*Users) SetUsername(discordID string, username string) error {
	var stmt *sql.Stmt
	if v, err := db.Prepare(`
		UPDATE tournament_users
		SET username = ?
		WHERE discord_id = ?
	`); err != nil {
		return err
	} else {
		stmt = v
	}
	defer stmt.Close()

	_, err := stmt.Exec(username, discordID)
	return err
}

func (*Users) SetTimezone(discordID string, timezone string) error {
	var stmt *sql.Stmt
	if v, err := db.Prepare(`
		UPDATE tournament_users
		SET timezone = ?
		WHERE discord_id = ?
	`); err != nil {
		return err
	} else {
		stmt = v
	}
	defer stmt.Close()

	_, err := stmt.Exec(timezone, discordID)
	return err
}

func (*Users) SetStreamURL(discordID string, streamURL string) error {
	var stmt *sql.Stmt
	if v, err := db.Prepare(`
		UPDATE tournament_users
		SET stream_url = ?
		WHERE discord_id = ?
	`); err != nil {
		return err
	} else {
		stmt = v
	}
	defer stmt.Close()

	_, err := stmt.Exec(streamURL, discordID)
	return err
}

func (*Users) SetCasterAlwaysOk(discordID string, ok bool) error {
	var stmt *sql.Stmt
	if v, err := db.Prepare(`
		UPDATE tournament_users
		SET caster_always_ok = ?
		WHERE discord_id = ?
	`); err != nil {
		return err
	} else {
		stmt = v
	}
	defer stmt.Close()

	_, err := stmt.Exec(ok, discordID)
	return err
}
