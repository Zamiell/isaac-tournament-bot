package main

import (
	"strconv"

	"github.com/Zamiell/isaac-tournament-bot/src/models"
)

var (
	characters = []string{
		"Isaac",     // 1
		"Magdalene", // 2
		"Cain",      // 3
		"Judas",     // 4
		"Blue Baby", // 5
		"Eve",       // 6
		"Samson",    // 7
		"Azazel",    // 8
		"Lazarus",   // 9
		"Eden",      // 10
		"The Lost",  // 11
		"Lilith",    // 12
		"Keeper",    // 13
		"Apollyon",  // 14
		"Samael",    // 15
	}
)

func charactersBanStart(race models.Race) {
	// Alert the players that the race is about to start
	// (this cannot be in the "matchStart()" function because we need to have everything in one message, or it can get out of order)
	msg := race.Racer1.Mention() + " and " + race.Racer2.Mention() + " - the race is scheduled to start in 5 minutes.\n\n"
	if race.CasterID.Valid {
		msg += race.Caster.Mention() + ", you are scheduled to cast this match in 5 minutes at: <" + race.Caster.StreamURL.String + ">\n\n"
	}

	msg += "**Character Ban Phase**\n\n"
	msg += "- Each player gets to ban " + strconv.Itoa(numBans) + " characters.\n"
	msg += "- Use the `!ban` command to select a character.\n"
	msg += "  e.g. `!ban 3` (to ban the 3rd character in the list)\n\n"
	if race.ActivePlayer == 1 {
		msg += race.Racer1.Mention()
	} else if race.ActivePlayer == 2 {
		msg += race.Racer2.Mention()
	}
	msg += ", you start! (randomly decided)\n\n"

	msg += getBansRemaining(race, "characters")
	msg += getRemaining(race, "characters")
	discordSend(race.ChannelID, msg)
}

func charactersPickStart(race models.Race, msg string) {
	// Set the state
	race.State = "pickingCharacters"
	if err := db.Races.SetState(race.ChannelID, race.State); err != nil {
		msg := "Failed to set the state for race \"" + race.Name() + "\": " + err.Error()
		log.Error(msg)
		discordSend(race.ChannelID, msg)
		return
	}

	msg += "**Character Pick Phase**\n\n"
	msg += "- " + strconv.Itoa(tournaments[race.ChallongeURL].BestOf) + " characters need to be picked.\n"
	msg += "- Use the `!pick` command to select a character.\n"
	msg += "  e.g. `!pick 3` (to pick the 3rd character in the list)\n\n"
	if race.ActivePlayer == 1 {
		msg += race.Racer1.Mention()
	} else if race.ActivePlayer == 2 {
		msg += race.Racer2.Mention()
	}
	msg += ", you start!\n\n"

	msg += getPicksRemaining(race, "characters")
	msg += getRemaining(race, "characters")
	discordSend(race.ChannelID, msg)
}

func charactersEnd(race models.Race, msg string) {
	msg += "**Characters for the Match**\n\n"
	for i, character := range race.Characters {
		msg += strconv.Itoa(i+1) + ". " + character + "\n"
	}
	msg += "\n"

	ruleset := tournaments[race.ChallongeURL].Ruleset
	if ruleset == "seeded" {
		buildsStart(race, msg)
	} else if ruleset == "unseeded" || ruleset == "team" {
		matchEnd(race, msg)
	} else {
		msg += "Unknown tournament ruleset for tournament \"" + race.TournamentName + "\"."
		log.Error(msg)
		discordSend(race.ChannelID, msg)
	}
}
