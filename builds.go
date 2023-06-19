package main

import (
	"strconv"
)

func buildsBanStart(race *Race, msg string) {
	// Update the state.
	race.State = RaceStateBanningBuilds
	if err := modals.Races.SetState(race.ChannelID, race.State); err != nil {
		msg := "Failed to set the state for race \"" + race.Name() + "\": " + err.Error()
		log.Error(msg)
		discordSend(race.ChannelID, msg)
		return
	}
	log.Info("Race \""+race.Name()+"\" is now in state:", race.State)

	// Initialize the number of bans.
	race.Racer1Bans = numBuildBans
	if err := modals.Races.SetBans(race.ChannelID, 1, race.Racer1Bans); err != nil {
		msg := "Failed to set the bans for racer 1 on race \"" + race.Name() + "\": " + err.Error()
		log.Error(msg)
		discordSend(race.ChannelID, msg)
		return
	}
	race.Racer2Bans = numBuildBans
	if err := modals.Races.SetBans(race.ChannelID, 2, race.Racer2Bans); err != nil {
		msg := "Failed to set the bans for racer 2 on race \"" + race.Name() + "\": " + err.Error()
		log.Error(msg)
		discordSend(race.ChannelID, msg)
		return
	}

	msg += "**Build Ban Phase**\n\n"
	msg += "- Each racer gets to ban " + strconv.Itoa(numBuildBans) + " builds.\n"
	msg += "- Use the `!ban` command to select a build.\n"
	msg += "  e.g. `!ban 3` (to ban the 3rd build in the list)\n\n"
	if race.ActiveRacer == 1 {
		msg += race.Racer1.Mention()
	} else if race.ActiveRacer == 2 {
		msg += race.Racer2.Mention()
	}
	msg += ", you start! (randomly decided)\n\n"

	msg += getBansRemaining(race)
	msg += getRemainingThingsMsg(race)
	discordSend(race.ChannelID, msg)
}

func buildsPickStart(race *Race, msg string) {
	// Set the state.
	race.State = RaceStatePickingBuilds
	if err := modals.Races.SetState(race.ChannelID, race.State); err != nil {
		msg := "Failed to set the state for race \"" + race.Name() + "\": " + err.Error()
		log.Error(msg)
		discordSend(race.ChannelID, msg)
		return
	}
	log.Info("Race \""+race.Name()+"\" is now in state:", race.State)

	msg += "**Build Pick Phase**\n\n"
	msg += "- " + strconv.Itoa(tournaments[race.ChallongeURL].BestOf) + " builds need to be picked.\n"
	msg += "- Use the `!pick` command to select a build.\n"
	msg += "  e.g. `!pick 3` (to pick the 3rd build in the list)\n\n"
	if race.ActiveRacer == 1 {
		msg += race.Racer1.Mention()
	} else if race.ActiveRacer == 2 {
		msg += race.Racer2.Mention()
	}
	msg += ", you start!\n\n"

	msg += getPicksRemainingMsg(race)
	msg += getRemainingThingsMsg(race)
	discordSend(race.ChannelID, msg)
}

func buildsVetoStart(race *Race, msg string) {
	race.State = RaceStateVetoBuilds
	if err := modals.Races.SetState(race.ChannelID, race.State); err != nil {
		msg := "Failed to set the state for race \"" + race.Name() + "\": " + err.Error()
		log.Error(msg)
		discordSend(race.ChannelID, msg)
		return
	}
	log.Info("Race \""+race.Name()+"\" is now in state:", race.State)

	// Initialize the number of vetos.
	race.Racer1Vetos = numBuildVetos
	if err := modals.Races.SetVetos(race.ChannelID, 1, numBuildVetos); err != nil {
		msg := "Failed to set the vetos for \"" + race.Racer1.Username + "\" on race \"" + race.Name() + "\": " + err.Error()
		log.Error(msg)
		return
	}
	race.Racer2Vetos = numBuildVetos
	if err := modals.Races.SetVetos(race.ChannelID, 2, numBuildVetos); err != nil {
		msg := "Failed to set the vetos for \"" + race.Racer2.Username + "\" on race \"" + race.Name() + "\": " + err.Error()
		log.Error(msg)
		return
	}

	msg += "**Build Ban Phase**\n\n"
	msg += "- " + strconv.Itoa(tournaments[race.ChallongeURL].BestOf) + " builds will randomly be chosen. Each racer will get " + strconv.Itoa(numBuildVetos) + " veto"
	if numBuildVetos != 1 {
		msg += "s"
	}
	msg += ".\n"
	msg += "- Use the `!yes` and `!no` commands to answer the questions.\n\n"

	// The person who starts the vetos for the builds is the opposite of the person who got to
	// start the vetos for the character.
	newActiveRacer := race.FirstPicker + 1
	if newActiveRacer > 2 {
		newActiveRacer = 1
	}

	race.ActiveRacer = newActiveRacer
	if err := modals.Races.SetActiveRacer(race.ChannelID, race.ActiveRacer); err != nil {
		msg := "Failed to set the active racer for race \"" + race.Name() + "\": " + err.Error()
		log.Error(msg)
		discordSend(race.ChannelID, msg)
		return
	}

	race.NumVoted = 2 // Set it to 2 so that it gives a new build.
	buildsRound(race, msg)
}

func buildsRound(race *Race, msg string) {
	if race.NumVoted == 2 {
		// Both racers have voted, so get a new build.
		race.NumVoted = 0
		if err := modals.Races.SetNumVoted(race.ChannelID, race.NumVoted); err != nil {
			msg := "Failed to set the NumVoted for race \"" + race.Name() + "\": " + err.Error()
			log.Error(msg)
			discordSend(race.ChannelID, msg)
			return
		}

		if len(race.Builds) == tournaments[race.ChallongeURL].BestOf {
			matchSetInProgressAndPrintSummary(race, msg)
			return
		}

		msg += assignRandomBuild(race)
	}

	if (race.Racer1Vetos == 0 && race.Racer2Vetos == 0) || // Both racer have used all of their vetos.
		(race.ActiveRacer == 1 && race.Racer1Vetos == 0) || // It is racer 1's turn and they have already used their vetos.
		(race.ActiveRacer == 2 && race.Racer2Vetos == 0) { // It is racer 2's turn and they have already used their vetos.

		log.Info("Skipping racer " + strconv.Itoa(race.ActiveRacer) + "'s turn, since they do not have a veto.")

		race.NumVoted++
		if err := modals.Races.SetNumVoted(race.ChannelID, race.NumVoted); err != nil {
			msg := "Failed to set the NumVoted for race \"" + race.Name() + "\": " + err.Error()
			log.Error(msg)
			discordSend(race.ChannelID, msg)
			return
		}

		incrementActiveRacer(race)
		buildsRound(race, msg)
		return
	}

	if race.ActiveRacer == 1 {
		msg += race.Racer1.Mention()
	} else if race.ActiveRacer == 2 {
		msg += race.Racer2.Mention()
	}
	msg += ", do you want to veto this build? Use `!yes` or `!no` to answer."
	discordSend(race.ChannelID, msg)
}

func buildsEnd(race *Race, msg string) {
	// Unlike the "charactersEnd" function, we don't print the builds, since they will be displayed
	// in the summary.
	matchSetInProgressAndPrintSummary(race, msg)
}

func assignRandomBuild(race *Race) string {
	// Get a random build.
	randomBuildName, randomBuildIndex := getRandomArrayElement(race.BuildsRemaining)
	randomBuildNum := randomBuildIndex + 1

	// Check to see if the item synergizes.
	build := getBuildObjectFromBuildName(randomBuildName)
	roundNum := len(race.Builds) + 1
	characterName := race.Characters[roundNum-1]
	synergizes := true
	for _, bannedCharacter := range build.BannedCharacters {
		if bannedCharacter.Name == characterName {
			synergizes = false
			break
		}
	}
	if !synergizes {
		// Get a new random build.
		log.Info("The randomly selected build of \"" + randomBuildName + "\" does not synergize with \"" + characterName + "\". Trying again...")
		return assignRandomBuild(race)
	}

	// Add it to the builds.
	race.Builds = append(race.Builds, randomBuildName)
	if err := modals.Races.SetBuilds(race.ChannelID, race.Builds); err != nil {
		msg := "Failed to set the builds for race \"" + race.Name() + "\": " + err.Error()
		log.Error(msg)
		return msg
	}

	// Remove it from the available builds.
	race.BuildsRemaining = deleteFromSlice(race.BuildsRemaining, randomBuildNum)
	if err := modals.Races.SetBuildsRemaining(race.ChannelID, race.BuildsRemaining); err != nil {
		msg := "Failed to set the builds for race \"" + race.Name() + "\": " + err.Error()
		log.Error(msg)
		return msg
	}

	msg := "**Round " + strconv.Itoa(roundNum) + "**:\n"
	msg += "- Character: *" + characterName + "*\n"
	msg += "- Build: *" + randomBuildName + "*\n\n"
	return msg
}
