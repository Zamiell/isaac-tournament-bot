package main

import (
	"strconv"
)

var (
	buildExceptions = [][]string{
		{}, // #0 - n/a

		// -------------------
		// Treasure Room items
		// -------------------

		// #1 - Cricket's Body
		// Cain - The range down causes Cain to barely be able to hit anything
		// Samson - ???
		// Azazel - The Brimstone beam inherits the split shots but they are not very good
		// Tainted Azazel - The Brimstone beam inherits the split shots but they are not very good
		{"Cain", "Samson", "Azazel", "Tainted Azazel"},

		{}, // #2 - Cricket's Head

		// #3 - Dead Eye
		// Azazel - The Brimstone beam prevents it from working
		// Lilith - It is hard to be accurate when shots come from Incubus
		// Keeper - It is hard to be accurate with triple shot
		// Tainted Azazel - The Brimstone beam prevents it from working
		// Tainted Keeper - It is hard to be accurate with triple shot
		{"Azazel", "Lilith", "Keeper", "Tainted Azazel", "Tainted Keeper"},

		// #4 - Death's Touch
		// The Forgotten - The piercing shots do nothing for the bone club
		{"The Forgotten"},

		// #5 - Dr. Fetus
		{},

		// #6 - Ipecac
		// Azazel - The short-range brimstone causes self-damage
		{"Azazel"},

		{}, // #7 - Magic Mushroom
		{}, // #8 - Mom's Knife
		{}, // #9 - Polyphemus
		{}, // #10 - Proptosis
		{}, // #11 - Tech.5
		{}, // #12 - Tech X
		{}, // #13 - C Section

		// ----------------
		// Devil Room items
		// ----------------

		{}, // #14 - Brimstone
		{}, // #15 - Maw of the Void

		// ----------------
		// Angel Room items
		// ----------------

		// #16 - Crown of Light
		// Eve - Eve cannot use the razor with this start
		{"Eve"},

		{}, // #17 - Sacred Heart
		{}, // #18 - Spirit Sword
		{}, // #19 - Revelation

		// -----------------
		// Secret Room items
		// -----------------

		{}, // #20 - Epic Fetus

		// ------------
		// Custom items
		// ------------

		{}, // #21 - Sawblade

		// ------
		// Builds
		// ------

		{}, // #22 - 20/20 + The Inner Eye
		{}, // #23 - Chocolate Milk + Steven
		{}, // #24 - Godhead + Cupid's Arrow
		{}, // #25 - Haemolacria + The Sad Onion

		// #26 - Incubus + Incubus
		// The Forgotten - The extra Incubus with bone clubs are not very helpful
		{"The Forgotten"},

		// #27 - Monstro's Lung + The Sad Onion
		// Keeper - The charge rate is too low to be very useful
		// Tainted Keeper - The charge rate is too low to be very useful
		{"Keeper", "Tainted Keeper"},

		{}, // #28 - Technology + A Lump of Coal
		{}, // #29 - Twisted Pair + Twisted Pair
		{}, // #30 - Pointy Rib + Eve's Mascara

		// #31 - Fire Mind + Mysterious Liquid + 13 Luck
		// Azazel - The synergy is only useful with a tear build
		// The Forgotten - The synergy is only useful with a tear build
		// Tainted Azazel - The synergy is only useful with a tear build
		{"Azazel", "The Forgotten", "Tainted Azazel"},

		// #32 - Eye of the Occult + Loki's Horns + 15 Luck
		// Azazel - The synergy is only useful with a tear build
		// The Forgotten - The synergy is only useful with a tear build
		// Tainted Azazel - The synergy is only useful with a tear build
		{"Azazel", "The Forgotten", "Tainted Azazel"},

		{}, // #33 - Distant Admiration + Friend Zone + Forever Alone + BFFS!
	}
)

func buildsStart(race *Race, msg string) {
	race.State = "vetoBuilds"
	if err := modals.Races.SetState(race.ChannelID, race.State); err != nil {
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

func buildsRound(race *Race, msg string) {
	if race.NumVoted == 2 {
		// Both racers have voted, so get a new build
		race.NumVoted = 0
		if err := modals.Races.SetNumVoted(race.ChannelID, race.NumVoted); err != nil {
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

func getBuild(race *Race) string {
	// Get a random build
	randBuildNum := getRandomInt(0, len(race.BuildsRemaining)-1)
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
	if err := modals.Races.SetBuilds(race.ChannelID, race.Builds); err != nil {
		msg := "Failed to set the builds for race \"" + race.Name() + "\": " + err.Error()
		log.Error(msg)
		return msg
	}

	// Remove it from the available builds
	race.BuildsRemaining = deleteFromSlice(race.BuildsRemaining, randBuildNum)
	if err := modals.Races.SetBuildsRemaining(race.ChannelID, race.BuildsRemaining); err != nil {
		msg := "Failed to set the builds for race \"" + race.Name() + "\": " + err.Error()
		log.Error(msg)
		return msg
	}

	msg := "**Round " + strconv.Itoa(roundNum) + "**:\n"
	msg += "- Character: *" + character + "*\n"
	msg += "- Build: *" + randBuild + "*\n\n"
	return msg
}
