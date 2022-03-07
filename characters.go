package main

import (
	"strconv"
)

var (
	characters = []string{
		"Isaac",
		"Magdalene",
		"Cain",
		"Judas",
		"Blue Baby",
		"Eve",
		"Samson",
		"Azazel",
		"Lazarus",
		"Eden",
		"The Lost",
		"Lilith",
		"Keeper",
		"Apollyon",
		"Forgotten",
		"Bethany",
		"Jacob & Esau", // Meme character
		"Tainted Isaac",
		"Tainted Magdalene",
		"Tainted Judas",
		"Tainted Blue Baby",
		"Tainted Eve",
		"Tainted Samson",
		"Tainted Azazel",
		// "Tainted Lazarus", // Meme character
		// "Tainted Eden", // Meme character
		"Tainted Lost",
		"Tainted Lilith",
		"Tainted Keeper",
		"Tainted Apollyon",
		// "Tainted Forgotten", // Meme character
		"Tainted Bethany",
		"Tainted Jacob",
	}
)

func charactersBanStart(race *Race) {
	// Update the state
	race.State = RaceStateBanningCharacters
	if err := modals.Races.SetState(race.ChannelID, race.State); err != nil {
		msg := "Failed to set the state for race \"" + race.Name() + "\": " + err.Error()
		log.Error(msg)
		discordSend(race.ChannelID, msg)
		return
	}
	log.Info("Race \""+race.Name()+"\" is now in state:", race.State)

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

func charactersPickStart(race *Race, msg string) {
	// Set the state
	race.State = RaceStatePickingCharacters
	if err := modals.Races.SetState(race.ChannelID, race.State); err != nil {
		msg := "Failed to set the state for race \"" + race.Name() + "\": " + err.Error()
		log.Error(msg)
		discordSend(race.ChannelID, msg)
		return
	}
	log.Info("Race \""+race.Name()+"\" is now in state:", race.State)

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

func charactersVetoStart(race *Race) {
	// Update the state
	race.State = RaceStateVetoCharacters
	if err := modals.Races.SetState(race.ChannelID, race.State); err != nil {
		msg := "Failed to set the state for race \"" + race.Name() + "\": " + err.Error()
		log.Error(msg)
		discordSend(race.ChannelID, msg)
		return
	}
	log.Info("Race \""+race.Name()+"\" is now in state:", race.State)

	msg := matchBeginningAlert(race)

	msg += "**Character Veto Phase**\n\n"
	msg += "- " + strconv.Itoa(tournaments[race.ChallongeURL].BestOf) + " characters will randomly be chosen. Each racer will get one veto.\n"
	msg += "- Use the `!yes` and `!no` commands to answer the questions.\n\n"
	race.NumVoted = 2 // Set it to 2 so that it gives a new character
	charactersRound(race, msg)
}

func charactersRound(race *Race, msg string) {
	if race.NumVoted == 2 {
		// Both racers have voted, so get a new character
		race.NumVoted = 0
		if err := modals.Races.SetNumVoted(race.ChannelID, race.NumVoted); err != nil {
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
		if err := modals.Races.SetNumVoted(race.ChannelID, race.NumVoted); err != nil {
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

func charactersEnd(race *Race, msg string) {
	msg += "**Characters for the Match**\n\n"
	for i, character := range race.Characters {
		msg += strconv.Itoa(i+1) + ". " + character + "\n"
	}
	msg += "\n"

	// Reset the vetos
	race.Racer1Vetos = numVetos
	if err := modals.Races.SetVetos(race.ChannelID, 1, numVetos); err != nil {
		msg := "Failed to set the vetos for \"" + race.Racer1.Username + "\" on race \"" + race.Name() + "\": " + err.Error()
		log.Error(msg)
		return
	}
	race.Racer2Vetos = numVetos
	if err := modals.Races.SetVetos(race.ChannelID, 2, numVetos); err != nil {
		msg := "Failed to set the vetos for \"" + race.Racer2.Username + "\" on race \"" + race.Name() + "\": " + err.Error()
		log.Error(msg)
		return
	}

	ruleset := tournaments[race.ChallongeURL].Ruleset
	if ruleset == "seeded" {
		buildsVetoStart(race, msg)
	} else if ruleset == "unseeded" || ruleset == "team" {
		matchEnd(race, msg)
	} else {
		msg += "Unknown tournament ruleset for tournament \"" + race.TournamentName + "\"."
		log.Error(msg)
		discordSend(race.ChannelID, msg)
	}
}

func getCharacter(race *Race) string {
	// Get a random character
	randCharacterNum := getRandomInt(0, len(race.CharactersRemaining)-1)
	randCharacter := race.CharactersRemaining[randCharacterNum]

	// Add it to the characters
	race.Characters = append(race.Characters, randCharacter)
	if err := modals.Races.SetCharacters(race.ChannelID, race.Characters); err != nil {
		msg := "Failed to set the characters for race \"" + race.Name() + "\": " + err.Error()
		log.Error(msg)
		return msg
	}

	// Remove it from the available characters
	race.CharactersRemaining = deleteFromSlice(race.CharactersRemaining, randCharacterNum)
	if err := modals.Races.SetCharactersRemaining(race.ChannelID, race.CharactersRemaining); err != nil {
		msg := "Failed to set the characters for race \"" + race.Name() + "\": " + err.Error()
		log.Error(msg)
		return msg
	}

	roundNum := len(race.Characters)
	msg := "**Round " + strconv.Itoa(roundNum) + "**:\n"
	msg += "- Character: *" + randCharacter + "*\n\n"
	return msg
}
