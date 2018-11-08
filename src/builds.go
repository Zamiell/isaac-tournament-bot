package main

import (
	"strconv"

	"github.com/Zamiell/isaac-tournament-bot/src/models"
)

var (
	builds = []string{
		"20/20",                     // 1
		"Chocolate Milk",            // 2
		"Cricket's Body",            // 3
		"Cricket's Head",            // 4
		"Dead Eye",                  // 5
		"Death's Touch",             // 6
		"Dr. Fetus",                 // 7
		"Epic Fetus",                // 8
		"Ipecac",                    // 9
		"Judas' Shadow",             // 10
		"Lil' Brimstone",            // 11
		"Magic Mushroom",            // 12
		"Mom's Knife",               // 13
		"Monstro's Lung",            // 14
		"Polyphemus",                // 15
		"Proptosis",                 // 16
		"Sacrificial Dagger",        // 17
		"Tech.5",                    // 18
		"Tech X",                    // 19
		"Brimstone",                 // 20
		"Incubus",                   // 21
		"Maw of the Void",           // 22
		"Crown of Light",            // 23
		"Godhead",                   // 24
		"Sacred Heart",              // 25
		"Mutant Spider + Inner Eye", // 26
		"Technology + Coal",         // 27
		"Ludovico + Parasite",       // 28
		"Fire Mind + 13 luck",       // 29
		"Tech Zero + more",          // 30
		"Kamikaze! + Host Hat",      // 31
		"Mega Blast + more",         // 32
	}

	buildExceptions = [][]string{
		{"Samae"},                              // 1
		{"Samael"},                             // 2
		{"Cain", "Samson", "Azazel", "Samael"}, // 3
		{}, // 4
		{"Azazel", "Lilith", "Keeper"}, // 5
		{},         // 6
		{},         // 7
		{},         // 8
		{},         // 9
		{"Azazel"}, // 10
		{},         // 11
		{},         // 12
		{},         // 13
		{"Keeper"}, // 14
		{},         // 15
		{},         // 16
		{},         // 17
		{"Lilith"}, // 18
		{},         // 19
		{},         // 20
		{},         // 21
		{"Lilith"}, // 22
		{},         // 23
		{},         // 24
		{},         // 25
		{"Keeper"}, // 26
		{},         // 27
		{},         // 28
		{"Azazel"}, // 29
		{"Azazel"}, // 30
		{},         // 31
		{},         // 32
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
