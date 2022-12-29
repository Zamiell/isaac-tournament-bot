package main

import (
	"database/sql"
	"strconv"
)

type Casts struct{}

type Cast struct {
	CasterID     int // The database ID of the user casting.
	Caster       *User
	R1Permission bool
	R2Permission bool
	Language     string
}

func (*Casts) Insert(channelID string, casterID string, language string) error {
	var stmt *sql.Stmt
	if v, err := db.Prepare(`
		INSERT INTO tournament_casts (
			race_id,
			caster,
			language
		) VALUES (
			(SELECT id FROM tournament_races WHERE channel_id = ?),
			(SELECT id FROM tournament_users WHERE discord_id = ?),
			?
		)
	`); err != nil {
		return err
	} else {
		stmt = v
	}
	defer stmt.Close()

	_, err := stmt.Exec(
		channelID,
		casterID,
		language,
	)
	return err
}

func (*Casts) Delete(channelID string, casterID string) error {
	var stmt *sql.Stmt
	if v, err := db.Prepare(`
		DELETE FROM tournament_casts
		WHERE
			race_id = (SELECT id FROM tournament_races WHERE channel_id = ?) AND
			caster = (SELECT id FROM tournament_users WHERE discord_id = ?)
	`); err != nil {
		return err
	} else {
		stmt = v
	}
	defer stmt.Close()

	if _, err := stmt.Exec(channelID, casterID); err != nil {
		return err
	}

	return nil
}

func (*Casts) GetAll(channelID string) ([]*Cast, error) {
	var rows *sql.Rows
	if v, err := db.Query(`
		SELECT
			caster,
			r1_permission,
			r2_permission,
			language
		FROM tournament_casts
		WHERE race_id = (SELECT id FROM tournament_races WHERE channel_id = ?)
	`, channelID); err != nil {
		return nil, err
	} else {
		rows = v
	}
	defer rows.Close()

	casts := make([]*Cast, 0)
	for rows.Next() {
		var cast Cast
		if err := rows.Scan(
			&cast.CasterID,
			&cast.R1Permission,
			&cast.R2Permission,
			&cast.Language,
		); err != nil {
			return nil, err
		}
		casts = append(casts, &cast)
	}

	return casts, nil
}

func (*Casts) SetPermission(channelID string, casterID string, racerNum int) error {
	var stmt *sql.Stmt
	if v, err := db.Prepare(`
		UPDATE tournament_casts
		SET r` + strconv.Itoa(racerNum) + `_permission = 1
		WHERE
			race_id = (SELECT id FROM tournament_races WHERE channel_id = ?) AND
			caster = (SELECT id FROM tournament_users WHERE discord_id = ?)
	`); err != nil {
		return err
	} else {
		stmt = v
	}
	defer stmt.Close()

	_, err := stmt.Exec(channelID, casterID)
	return err
}
