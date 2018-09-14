package models

import (
	"database/sql"
	"strconv"
	"time"

	"github.com/go-sql-driver/mysql"
)

type Races struct{}

type Race struct {
	TournamentName    string
	Racer1ID          int     // The "tournament_racers" database ID
	Racer1ChallongeID float64 // The "participant" ID; needed to automatically set the winner through the Challonge API
	Racer1            Racer
	Racer2ID          int     // The "tournament_racers" database ID
	Racer2ChallongeID float64 // The "participant" ID; needed to automatically set the winner through the Challonge API
	Racer2            Racer
	ChannelID         string // The Discord channel ID that was automatically created for this race
	ChallongeURL      string // The suffix of the Challonge URL for this tournament
	ChallongeMatchID  string
	BracketRound      string
	State             string
	/*
		State definitions:
		- "initial" is freshly created before both players have confirmed a scheduled time
		- "scheduled" is confirmed but before it starts
		- "banningCharacters" (triggered 5 minutes before starting)
		- "pickingCharacters"
		- "vetoBuilds"
		- "inProgress"
		- "completed" (after a score is reported)
	*/
	DatetimeScheduled   mysql.NullTime
	CasterID            sql.NullInt64
	Caster              Racer
	CasterP1            bool // Whether or not player 1 approves of the caster who volunteered
	CasterP2            bool // Whether or not player 2 approves of the caster who volunteered
	ActivePlayer        int
	CharactersRemaining []string
	Characters          []string
	BuildsRemaining     []string
	Builds              []string
	Racer1Bans          int
	Racer2Bans          int
	Racer1Vetos         int
	Racer2Vetos         int
	NumVoted            int
}

func (r *Race) Name() string {
	return r.Racer1.Username + "-vs-" + r.Racer2.Username
}

func (*Races) Insert(racer1DiscordID string, racer2DiscordID string, race Race) error {
	charactersRemaining := sliceToString(race.CharactersRemaining)
	buildsRemaining := sliceToString(race.BuildsRemaining)

	var stmt *sql.Stmt
	if v, err := db.Prepare(`
		INSERT INTO tournament_races (
			tournament_name,
			racer1,
			racer1_challonge_id,
			racer2,
			racer2_challonge_id,
			channel_id,
			challonge_url,
			challonge_match_id,
			bracket_round,
			state,
			characters_remaining,
			builds_remaining,
			racer1_bans,
			racer2_bans,
			racer1_vetos,
			racer2_vetos
		) VALUES (
			?,
			(SELECT id FROM tournament_racers WHERE discord_id = ?),
			?,
			(SELECT id FROM tournament_racers WHERE discord_id = ?),
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
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
		race.TournamentName,
		racer1DiscordID,
		race.Racer1ChallongeID,
		racer2DiscordID,
		race.Racer2ChallongeID,
		race.ChannelID,
		race.ChallongeURL,
		race.ChallongeMatchID,
		race.BracketRound,
		race.State,
		charactersRemaining,
		buildsRemaining,
		race.Racer1Bans,
		race.Racer2Bans,
		race.Racer1Vetos,
		race.Racer2Vetos,
	); err != nil {
		return err
	}

	return nil
}

func (*Races) Delete(channelID string) error {
	var stmt *sql.Stmt
	if v, err := db.Prepare(`
		DELETE FROM tournament_races
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
	var charactersRemaining, characters, buildsRemaining, builds string
	if err := db.QueryRow(`
		SELECT
			tournament_name,
			racer1,
			racer2,
			channel_id,
			challonge_url,
			challonge_match_id,
			bracket_round,
			state,
			datetime_scheduled,
			caster,
			caster_p1,
			caster_p2,
			active_player,
			characters_remaining,
			characters,
			builds_remaining,
			builds,
			racer1_bans,
			racer2_bans,
			racer1_vetos,
			racer2_vetos,
			num_voted
		FROM tournament_races
		WHERE channel_id = ?
	`, channelID).Scan(
		&race.TournamentName,
		&race.Racer1ID,
		&race.Racer2ID,
		&race.ChannelID,
		&race.ChallongeURL,
		&race.ChallongeMatchID,
		&race.BracketRound,
		&race.State,
		&race.DatetimeScheduled,
		&race.CasterID,
		&race.CasterP1,
		&race.CasterP2,
		&race.ActivePlayer,
		&charactersRemaining,
		&characters,
		&buildsRemaining,
		&builds,
		&race.Racer1Bans,
		&race.Racer2Bans,
		&race.Racer1Vetos,
		&race.Racer2Vetos,
		&race.NumVoted,
	); err != nil {
		return race, err
	}

	race.CharactersRemaining = stringToSlice(charactersRemaining)
	race.Characters = stringToSlice(characters)
	race.BuildsRemaining = stringToSlice(buildsRemaining)
	race.Builds = stringToSlice(builds)

	return race, nil
}

func (*Races) GetAllScheduled() ([]string, error) {
	var rows *sql.Rows
	if v, err := db.Query(`
		SELECT channel_id
		FROM tournament_races
		WHERE state = "scheduled"
	`); err != nil {
		return nil, err
	} else {
		rows = v
	}
	defer rows.Close()

	channelIDs := make([]string, 0)
	for rows.Next() {
		var channelID string
		if err := rows.Scan(
			&channelID,
		); err != nil {
			return nil, err
		}
		channelIDs = append(channelIDs, channelID)
	}

	return channelIDs, nil
}

func (*Races) GetNext() (string, error) {
	var channelID string
	if err := db.QueryRow(`
		SELECT channel_id
		FROM tournament_races
		WHERE
			state = "scheduled"
			AND datetime_scheduled > NOW()
		ORDER BY datetime_scheduled ASC
		LIMIT 1
	`).Scan(&channelID); err != nil {
		return channelID, err
	}

	return channelID, nil
}

func (*Races) SetState(channelID string, state string) error {
	var stmt *sql.Stmt
	if v, err := db.Prepare(`
		UPDATE tournament_races
		SET state = ?
		WHERE channel_id = ?
	`); err != nil {
		return err
	} else {
		stmt = v
	}
	defer stmt.Close()

	_, err := stmt.Exec(state, channelID)
	return err
}

func (*Races) SetDatetimeScheduled(channelID string, datetimeScheduled time.Time, activePlayer int) error {
	// activePlayer is the player who suggested the time
	var stmt *sql.Stmt
	if v, err := db.Prepare(`
		UPDATE tournament_races
		SET datetime_scheduled = ?, active_player = ?
		WHERE channel_id = ?
	`); err != nil {
		return err
	} else {
		stmt = v
	}
	defer stmt.Close()

	_, err := stmt.Exec(datetimeScheduled, activePlayer, channelID)
	return err
}

func (*Races) UnsetDatetimeScheduled(channelID string) error {
	var stmt *sql.Stmt
	if v, err := db.Prepare(`
		UPDATE tournament_races
		SET datetime_scheduled = NULL
		WHERE channel_id = ?
	`); err != nil {
		return err
	} else {
		stmt = v
	}
	defer stmt.Close()

	_, err := stmt.Exec(channelID)
	return err
}

func (*Races) SetCaster(channelID string, casterID string) error {
	var stmt *sql.Stmt
	if v, err := db.Prepare(`
		UPDATE tournament_races
		SET caster = (SELECT id FROM tournament_racers WHERE discord_id = ?)
		WHERE channel_id = ?
	`); err != nil {
		return err
	} else {
		stmt = v
	}
	defer stmt.Close()

	_, err := stmt.Exec(casterID, channelID)
	return err
}

func (*Races) UnsetCaster(channelID string) error {
	var stmt *sql.Stmt
	if v, err := db.Prepare(`
		UPDATE tournament_races
		SET
			caster = NULL,
			caster_p1 = 0,
			caster_p2 = 0
		WHERE channel_id = ?
	`); err != nil {
		return err
	} else {
		stmt = v
	}
	defer stmt.Close()

	_, err := stmt.Exec(channelID)
	return err
}

func (*Races) SetCasterApproval(channelID string, playerNum int) error {
	var stmt *sql.Stmt
	if v, err := db.Prepare(`
		UPDATE tournament_races
		SET caster_p` + strconv.Itoa(playerNum) + ` = 1
		WHERE channel_id = ?
	`); err != nil {
		return err
	} else {
		stmt = v
	}
	defer stmt.Close()

	_, err := stmt.Exec(channelID)
	return err
}

func (*Races) SetActivePlayer(channelID string, activePlayer int) error {
	var stmt *sql.Stmt
	if v, err := db.Prepare(`
		UPDATE tournament_races
		SET active_player = ?
		WHERE channel_id = ?
	`); err != nil {
		return err
	} else {
		stmt = v
	}
	defer stmt.Close()

	_, err := stmt.Exec(activePlayer, channelID)
	return err
}

func (*Races) SetCharactersRemaining(channelID string, characters []string) error {
	charactersString := sliceToString(characters)

	var stmt *sql.Stmt
	if v, err := db.Prepare(`
		UPDATE tournament_races
		SET characters_remaining = ?
		WHERE channel_id = ?
	`); err != nil {
		return err
	} else {
		stmt = v
	}
	defer stmt.Close()

	_, err := stmt.Exec(charactersString, channelID)
	return err
}

func (*Races) SetCharacters(channelID string, characters []string) error {
	charactersString := sliceToString(characters)

	var stmt *sql.Stmt
	if v, err := db.Prepare(`
		UPDATE tournament_races
		SET characters = ?
		WHERE channel_id = ?
	`); err != nil {
		return err
	} else {
		stmt = v
	}
	defer stmt.Close()

	_, err := stmt.Exec(charactersString, channelID)
	return err
}

func (*Races) SetBuildsRemaining(channelID string, builds []string) error {
	buildsString := sliceToString(builds)

	var stmt *sql.Stmt
	if v, err := db.Prepare(`
		UPDATE tournament_races
		SET builds_remaining = ?
		WHERE channel_id = ?
	`); err != nil {
		return err
	} else {
		stmt = v
	}
	defer stmt.Close()

	_, err := stmt.Exec(buildsString, channelID)
	return err
}

func (*Races) SetBuilds(channelID string, builds []string) error {
	buildsString := sliceToString(builds)

	var stmt *sql.Stmt
	if v, err := db.Prepare(`
		UPDATE tournament_races
		SET builds = ?
		WHERE channel_id = ?
	`); err != nil {
		return err
	} else {
		stmt = v
	}
	defer stmt.Close()

	_, err := stmt.Exec(buildsString, channelID)
	return err
}

func (*Races) SetBans(channelID string, playerNum int, bans int) error {
	var stmt *sql.Stmt
	if v, err := db.Prepare(`
		UPDATE tournament_races
		SET racer` + strconv.Itoa(playerNum) + `_bans = ?
		WHERE channel_id = ?
	`); err != nil {
		return err
	} else {
		stmt = v
	}
	defer stmt.Close()

	_, err := stmt.Exec(bans, channelID)
	return err
}

func (*Races) SetVetos(channelID string, playerNum int, vetos int) error {
	var stmt *sql.Stmt
	if v, err := db.Prepare(`
		UPDATE tournament_races
		SET racer` + strconv.Itoa(playerNum) + `_vetos = ?
		WHERE channel_id = ?
	`); err != nil {
		return err
	} else {
		stmt = v
	}
	defer stmt.Close()

	_, err := stmt.Exec(vetos, channelID)
	return err
}

func (*Races) SetNumVoted(channelID string, numVoted int) error {
	var stmt *sql.Stmt
	if v, err := db.Prepare(`
		UPDATE tournament_races
		SET num_voted = ?
		WHERE channel_id = ?
	`); err != nil {
		return err
	} else {
		stmt = v
	}
	defer stmt.Close()

	_, err := stmt.Exec(numVoted, channelID)
	return err
}

func (*Races) SetScore(channelID string, score string) error {
	var stmt *sql.Stmt
	if v, err := db.Prepare(`
		UPDATE tournament_races
		SET score = ?
		WHERE channel_id = ?
	`); err != nil {
		return err
	} else {
		stmt = v
	}
	defer stmt.Close()

	_, err := stmt.Exec(score, channelID)
	return err
}
