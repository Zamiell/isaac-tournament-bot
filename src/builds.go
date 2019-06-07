package main

import (
	"strconv"

	"github.com/Zamiell/isaac-tournament-bot/src/models"
)

var (
	builds = []string{
		"20/20",                            // 1
		"Cricket's Body",                   // 2
		"Cricket's Head",                   // 3
		"Dead Eye",                         // 4
		"Death's Touch",                    // 5
		"Dr. Fetus",                        // 6
		"Epic Fetus",                       // 7
		"Ipecac",                           // 8
		"Judas' Shadow",                    // 9
		"Lil' Brimstone",                   // 10
		"Magic Mushroom",                   // 11
		"Mom's Knife",                      // 12
		"Monstro's Lung",                   // 13
		"Polyphemus",                       // 14
		"Proptosis",                        // 15
		"Sacrificial Dagger",               // 16
		"Tech.5",                           // 17
		"Tech X",                           // 18
		"Brimstone",                        // 19
		"Incubus",                          // 20
		"Maw of the Void",                  // 21
		"Crown of Light",                   // 22
		"Godhead",                          // 23
		"Sacred Heart",                     // 24
		"Chocolate Milk + Steven",          // 25
		"Jacob's Ladder + There's Options", // 26
		"Mutant Spider + Inner Eye",        // 27
		"Technology + Coal",                // 28
		"Ludovico + Parasite",              // 29
		"Fire Mind + 13 luck",              // 30
		"Tech Zero + more",                 // 31
		"Kamikaze! + Host Hat",             // 32
		"Mega Blast + more",                // 33
	}

	buildExceptions = [][]string{
		{"Samael"},                             // 1
		{"Cain", "Samson", "Azazel", "Samael"}, // 2
		{},                                     // 3
		{"Azazel", "Lilith", "Keeper"},         // 4
		{"The Forgotten"},                      // 5
		{},                                     // 6
		{},                                     // 7
		{"Azazel"},                             // 8
		{"Azazel", "The Forgotten"},            // 9
		{"The Forgotten"},                      // 10
		{},                                     // 11
		{},                                     // 12
		{"Keeper"},                             // 13
		{},                                     // 14
		{},                                     // 15
		{},                                     // 16
		{"Lilith"},                             // 17
		{},                                     // 18
		{},                                     // 19
		{"The Forgotten"},                      // 20
		{"Lilith"},                             // 21
		{"Eve"},                                // 22
		{},                                     // 23
		{},                                     // 24
		{"Samael"},                             // 25
		{"Azazel"},                             // 26
		{"Azazel", "Keeper"},                   // 27
		{},                                     // 28
		{},                                     // 29
		{"Azazel", "The Forgotten"},            // 30
		{"Azazel"},                             // 31
		{},                                     // 32
		{},                                     // 33
	}
)

func buildsStart(race *models.Race, msg string) {
	race.State = "vetoBuilds"
	if err := db.Races.SetState(race.ChannelID, race.State); err != nil {
		msg := "Failed to set the state for race \"" + race.Name() + "\": " + err.Error()
		log.Error(msg)
		discordSend(race.ChannelID, msg)
		return
	}
	log.Info("Race \"" + race.Name() + "\" is now in state: " + race.State)

	msg += "**Build Ban Phase**\n\n"
	msg += "- " + strconv.Itoa(tournaments[race.ChallongeURL].BestOf) + " builds will randomly be chosen. Each racer will get one veto.\n"
	msg += "- Use the `!yes` and `!no` commands to answer the questions.\n\n"
	race.NumVoted = 2 // Set it to 2 so that it gives a new build
	buildsRound(race, msg)
}

func buildsRound(race *models.Race, msg string) {
	if race.NumVoted == 2 {
		// Both racers have voted, so get a new build
		race.NumVoted = 0
		if err := db.Races.SetNumVoted(race.ChannelID, race.NumVoted); err != nil {
			msg := "Failed to set the NumVoted for race \"" + race.Name() + "\": " + err.Error()
			log.Error(msg)
			discordSend(race.ChannelID, msg)
			return
		}

		if len(race.Builds) == tournaments[race.ChallongeURL].BestOf {
			matchEnd(race, msg)
			return
		}

		msg += getBuild(race)
	}

	if (race.Racer1Vetos == 0 && race.Racer2Vetos == 0) || // Both racer have used all of their vetos
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

func getBuild(race *models.Race) string {
	// Get a random build
	randBuildNum := getRandom(0, len(race.BuildsRemaining)-1)
	log.Debug("randBuildNum:", randBuildNum)
	randBuild := race.BuildsRemaining[randBuildNum]
	log.Debug("randBuild:", randBuild)

	// Check to see if the item synergizes
	roundNum := len(race.Builds) + 1
	character := race.Characters[roundNum-1]
	log.Debug("character:", character)
	synergizes := true
	log.Debug("buildExceptions[randBuildNum]:", buildExceptions[randBuildNum])
	for _, exceptedCharacter := range buildExceptions[randBuildNum] {
		log.Debug("exceptedCharacter:", exceptedCharacter)
		if exceptedCharacter == character {
			synergizes = false
			break
		}
	}
	if !synergizes {
		// Get a new random build
		log.Info("The randomly selected build of \"" + randBuild + "\" does not synergize with \"" + character + "\". Trying again...")
		return getBuild(race)
	}

	// Add it to the builds
	race.Builds = append(race.Builds, randBuild)
	if err := db.Races.SetBuilds(race.ChannelID, race.Builds); err != nil {
		msg := "Failed to set the builds for race \"" + race.Name() + "\": " + err.Error()
		log.Error(msg)
		return msg
	}

	// Remove it from the available builds
	race.BuildsRemaining = deleteFromSlice(race.BuildsRemaining, randBuildNum)
	if err := db.Races.SetBuildsRemaining(race.ChannelID, race.BuildsRemaining); err != nil {
		msg := "Failed to set the builds for race \"" + race.Name() + "\": " + err.Error()
		log.Error(msg)
		return msg
	}

	msg := "**Round " + strconv.Itoa(roundNum) + "**:\n"
	msg += "- Character: *" + character + "*\n"
	msg += "- Build: *" + randBuild + "*\n\n"
	return msg
}
