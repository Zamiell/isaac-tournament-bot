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
		"Forgotten", // 15
		"Samael",    // 16
	}
)

func charactersBanStart(race *models.Race) {
	// Update the state
	race.State = "banningCharacters"
	if err := db.Races.SetState(race.ChannelID, race.State); err != nil {
		msg := "Failed to set the state for race \"" + race.Name() + "\": " + err.Error()
		log.Error(msg)
		discordSend(race.ChannelID, msg)
		return
	}
	log.Info("Race \"" + race.Name() + "\" is now in state: " + race.State)

	msg := matchBeginningAlert(race)

	msg += "**Character Ban Phase**\n\n"
	msg += "- Each racer gets to ban " + strconv.Itoa(numBans) + " characters.\n"
	msg += "- Use the `!ban` command to select a character.\n"
	msg += "  e.g. `!ban 3` (to ban the 3rd character in the list)\n\n"
	if race.ActiveRacer == 1 {
		msg += race.Racer1.Mention()
	} else if race.ActiveRacer == 2 {
		msg += race.Racer2.Mention()
	}
	msg += ", you start! (randomly decided)\n\n"

	msg += getBansRemaining(race, "characters")
	msg += getRemaining(race, "characters")
	discordSend(race.ChannelID, msg)
}

func charactersPickStart(race *models.Race, msg string) {
	// Set the state
	race.State = "pickingCharacters"
	if err := db.Races.SetState(race.ChannelID, race.State); err != nil {
		msg := "Failed to set the state for race \"" + race.Name() + "\": " + err.Error()
		log.Error(msg)
		discordSend(race.ChannelID, msg)
		return
	}
	log.Info("Race \"" + race.Name() + "\" is now in state: " + race.State)

	msg += "**Character Pick Phase**\n\n"
	msg += "- " + strconv.Itoa(tournaments[race.ChallongeURL].BestOf) + " characters need to be picked.\n"
	msg += "- Use the `!pick` command to select a character.\n"
	msg += "  e.g. `!pick 3` (to pick the 3rd character in the list)\n\n"
	if race.ActiveRacer == 1 {
		msg += race.Racer1.Mention()
	} else if race.ActiveRacer == 2 {
		msg += race.Racer2.Mention()
	}
	msg += ", you start!\n\n"

	msg += getPicksRemaining(race, "characters")
	msg += getRemaining(race, "characters")
	discordSend(race.ChannelID, msg)
}

func charactersVetoStart(race *models.Race) {
	// Update the state
	race.State = "vetoCharacters"
	if err := db.Races.SetState(race.ChannelID, race.State); err != nil {
		msg := "Failed to set the state for race \"" + race.Name() + "\": " + err.Error()
		log.Error(msg)
		discordSend(race.ChannelID, msg)
		return
	}
	log.Info("Race \"" + race.Name() + "\" is now in state: " + race.State)

	msg := matchBeginningAlert(race)

	msg += "**Character Veto Phase**\n\n"
	msg += "- " + strconv.Itoa(tournaments[race.ChallongeURL].BestOf) + " characters will randomly be chosen. Each racer will get one veto.\n"
	msg += "- Use the `!yes` and `!no` commands to answer the questions.\n\n"
	race.NumVoted = 2 // Set it to 2 so that it gives a new character
	charactersRound(race, msg)
}

func charactersRound(race *models.Race, msg string) {
	if race.NumVoted == 2 {
		// Both racers have voted, so get a new character
		race.NumVoted = 0
		if err := db.Races.SetNumVoted(race.ChannelID, race.NumVoted); err != nil {
			msg := "Failed to set the NumVoted for race \"" + race.Name() + "\": " + err.Error()
			log.Error(msg)
			discordSend(race.ChannelID, msg)
			return
		}

		if len(race.Characters) == tournaments[race.ChallongeURL].BestOf {
			charactersEnd(race, msg)
			return
		}

		msg += getCharacter(race)
	}

	if (race.Racer1Vetos == 0 && race.Racer2Vetos == 0) || // Both racers have used all of their vetos
		(race.ActiveRacer == 1 && race.Racer1Vetos == 0) || // It is racer 1's turn and they have already used their vetos
		(race.ActiveRacer == 2 && race.Racer2Vetos == 0) { // It is racer 2's turn and they have already used their vetos

		log.Info("Skipping racer " + strconv.Itoa(race.ActiveRacer) + "'s turn, since they do not have a veto.")

		race.NumVoted++
		if err := db.Races.SetNumVoted(race.ChannelID, race.NumVoted); err != nil {
			msg := "Failed to set the NumVoted for race \"" + race.Name() + "\": " + err.Error()
			log.Error(msg)
			discordSend(race.ChannelID, msg)
			return
		}

		incrementActiveRacer(race)
		charactersRound(race, msg)
		return
	}

	if race.ActiveRacer == 1 {
		msg += race.Racer1.Mention()
	} else if race.ActiveRacer == 2 {
		msg += race.Racer2.Mention()
	}
	msg += ", do you want to veto this character? Use `!yes` or `!no` to answer."
	discordSend(race.ChannelID, msg)
}

func charactersEnd(race *models.Race, msg string) {
	msg += "**Characters for the Match**\n\n"
	for i, character := range race.Characters {
		msg += strconv.Itoa(i+1) + ". " + character + "\n"
	}
	msg += "\n"

	// Reset the vetos
	race.Racer1Vetos = numVetos
	if err := db.Races.SetVetos(race.ChannelID, 1, numVetos); err != nil {
		msg := "Failed to set the vetos for \"" + race.Racer1.Username + "\" on race \"" + race.Name() + "\": " + err.Error()
		log.Error(msg)
		return
	}
	race.Racer2Vetos = numVetos
	if err := db.Races.SetVetos(race.ChannelID, 2, numVetos); err != nil {
		msg := "Failed to set the vetos for \"" + race.Racer2.Username + "\" on race \"" + race.Name() + "\": " + err.Error()
		log.Error(msg)
		return
	}

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

func getCharacter(race *models.Race) string {
	// Get a random character
	randCharacterNum := getRandom(0, len(race.CharactersRemaining)-1)
	randCharacter := race.CharactersRemaining[randCharacterNum]

	// Add it to the characters
	race.Characters = append(race.Characters, randCharacter)
	if err := db.Races.SetCharacters(race.ChannelID, race.Characters); err != nil {
		msg := "Failed to set the characters for race \"" + race.Name() + "\": " + err.Error()
		log.Error(msg)
		return msg
	}

	// Remove it from the available characters
	race.CharactersRemaining = deleteFromSlice(race.CharactersRemaining, randCharacterNum)
	if err := db.Races.SetCharactersRemaining(race.ChannelID, race.CharactersRemaining); err != nil {
		msg := "Failed to set the characters for race \"" + race.Name() + "\": " + err.Error()
		log.Error(msg)
		return msg
	}

	roundNum := len(race.Characters)
	msg := "**Round " + strconv.Itoa(roundNum) + "**:\n"
	msg += "- Character: *" + randCharacter + "*\n\n"
	return msg
}
