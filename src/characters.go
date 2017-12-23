package main

import (
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

func charactersStart(race models.Race) {
	msg := race.Racer1.Mention() + " and " + race.Racer2.Mention() + " - the race is scheduled to start in 5 minutes.\n\n"
	if race.CasterID.Valid {
		msg += race.Caster.Mention() + ", you are scheduled to cast this match in 5 minutes at: <" + race.Caster.StreamURL.String + ">\n\n"
	}

	msg += "**Character Ban Phase**\n\n"
	msg += "- Each player gets to ban 3 characters.\n"
	msg += "- Use the `!ban` command to select a character.\n"
	msg += "  e.g. `!ban 3` (to ban the 3rd character in the list)\n\n"
	if race.ActivePlayer == 1 {
		msg += race.Racer1.Mention()
	} else if race.ActivePlayer == 2 {
		msg += race.Racer2.Mention()
	}
	msg += ", you start! (randomly decided)\n\n"
	msg += getRemaining(race, "characters")
	discordSend(race.ChannelID, msg)
}
