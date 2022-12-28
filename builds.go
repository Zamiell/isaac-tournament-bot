package main

import (
	"strconv"
)

var (
	// This must be kept in sync with the build exceptions in "isaac-racing-server".
	buildExceptions = [][]string{
		// #0 - n/a
		{},

		// -------------------
		// Treasure Room items
		// -------------------

		// #1 - Cricket's Body
		{
			// The Brimstone beam inherits the split shots but they are not very good.
			"Azazel", // 7

			// The Brimstone beam inherits the split shots but they are not very good.
			"Tainted Azazel", // 28
		},

		{}, // #2 - Cricket's Head

		// #3 - Dead Eye
		{
			// The Brimstone beam prevents it from working.
			"Azazel", // 7

			// It is hard to be accurate when shots come from Incubus.
			"Lilith", // 13

			// It is hard to be accurate with triple shot.
			"Keeper", // 14

			// The Brimstone beam prevents it from working.
			"Tainted Azazel", // 28

			// It is hard to be accurate with quad shot.
			"Tainted Keeper", // 33
		},

		// #4 - Death's Touch
		{
			// The piercing shots do nothing for Azazel's Brimstone.
			"Azazel", // 7

			// The piercing shots do nothing for the bone club.
			"The Forgotten", // 16

			// The piercing shots do nothing for Tainted Azazel's Brimstone.
			"Tainted Azazel", // 28

			// The piercing shots do nothing for the bone club.
			"Tainted Forgotten", // 35
		},

		// #5 - Dr. Fetus
		{
			// Very annoying to use with the skeleton body.
			"Tainted Forgotten", // 35
		},

		// #6 - Ipecac
		{
			// The short-range brimstone causes self-damage.
			"Azazel", // 7

			// Can cause unavoidable damage if a clot is behind you or shoots at an obstacle near
			// you.
			"Tainted Eve", // 26
		},

		// #7 - Magic Mushroom
		{},

		// #8 - Mom's Knife
		{},

		// #9 - Polyphemus
		{},

		// #10 - Proptosis
		{},

		// #11 - Tech.5
		{},

		// #12 - Tech X
		{},

		// #13 - C Section
		{},

		// ----------------
		// Devil Room items
		// ----------------

		// #14 - Brimstone
		{
			// Gello fires the brim very slowly and auto-fire is not always accurate.
			"Tainted Lilith", // 32
		},

		// #15 - Maw of the Void
		{},

		// ----------------
		// Angel Room items
		// ----------------

		// #16 - Crown of Light
		{
			// Eve cannot use the razor with this start.
			"Eve", // 5

			// Crown is never full with the depleting hearts.
			"Tainted Magdalene", // 22

			// Crown can't be active with the clot mechanic.
			"Tainted Eve", // 26
		},

		// #17 - Sacred Heart
		{},

		// #18 - Spirit Sword
		{
			// Annoying because the sword goes to Incubus.
			"Lilith", // 13

			// No synergy with the bone club.
			"The Forgotten", // 16

			// No synergy with the bone club.
			"Tainted Forgotten", // 35
		},

		// #19 - Revelation
		{},

		// -----------------
		// Secret Room items
		// -----------------

		// #20 - Epic Fetus
		{
			// The target keeps moving and you can't control it, making it impossible to target
			// enemies.
			"Tainted Lilith", // 32
		},

		// ------------
		// Custom items
		// ------------

		// #21 - Sawblade
		{
			// Very complicated to play orbitals with her because she can't protect herself from
			// losing the deal with soul hearts.
			"Bethany", // 18

			// Impossible to play orbitals with Tainted Eve's clots, they will disappear very
			// quickly.
			"Tainted Eve", // 26

			// With his health mechanic, it is too dangerous to use orbitals.
			"Tainted Lost", // 31
		},

		// ------
		// Builds
		// ------

		// #22 - 20/20 + The Inner Eye
		{},

		// #23 - Chocolate Milk + Steven
		{},

		// #24 - Godhead + Cupid's Arrow
		{
			// Small damage up for a tears down, resulting in a loss of DPS overall.
			"Azazel", // 7

			// Does nothing with the bone club.
			"The Forgotten", // 16

			// Does nothing with the bone club.
			"Tainted Forgotten", // 35
		},

		// #25 - Haemolacria + The Sad Onion
		{},

		// #26 - Incubus + Incubus
		{},

		// #27 - Monstro's Lung + The Sad Onion
		{
			// Huge tears down, resulting in a loss of DPS overall.
			"Keeper", // 14

			// Tears down, worse than having no starter with the fetus.
			"Tainted Lilith", // 32

			// Huge tears down, resulting in a loss of DPS overall.
			"Tainted Keeper", // 33
		},

		// #28 - Technology + A Lump of Coal
		{},

		// #29 - Twisted Pair + Twisted Pair
		{},

		// #30 - Pointy Rib + Eve's Mascara
		{},

		// #31 - Fire Mind + Mysterious Liquid + 13 Luck.
		{
			// The synergy is only useful with a tear build.
			"Azazel", // 7

			// Luck does not apply to Incubus for some reason.
			"Lilith", // 13

			// The synergy is only useful with a tear build.
			"The Forgotten", // 16

			// The synergy is only useful with a tear build.
			"Tainted Azazel", // 28

			// Too dangerous to be synergistic.
			"Tainted Lost", // 31

			// The synergy is only useful with a tear build.
			"Tainted Forgotten", // 35
		},

		// #32 - Eye of the Occult + Loki's Horns + 15 Luck
		{
			// Homing brimstone is too powerful, resulting in a build with a low-skill requirement.
			"Azazel", // 7

			// It is only a damage up on the bone club.
			"The Forgotten", // 16

			// Homing brimstone is too powerful, resulting in a build with a low-skill requirement.
			"Tainted Azazel", // 28

			// It is only a damage up on the bone club.
			"Tainted Forgotten", // 35
		},

		// #33 - Distant Admiration + Friend Zone + Forever Alone + BFFS!
		{},
	}
)

func buildsBanStart(race *Race, msg string) {
	// Update the state
	race.State = RaceStateBanningBuilds
	if err := modals.Races.SetState(race.ChannelID, race.State); err != nil {
		msg := "Failed to set the state for race \"" + race.Name() + "\": " + err.Error()
		log.Error(msg)
		discordSend(race.ChannelID, msg)
		return
	}
	log.Info("Race \""+race.Name()+"\" is now in state:", race.State)

	// Initialize the number of bans
	race.Racer1Bans = numBans
	if err := modals.Races.SetBans(race.ChannelID, 1, race.Racer1Bans); err != nil {
		msg := "Failed to set the bans for racer 1 on race \"" + race.Name() + "\": " + err.Error()
		log.Error(msg)
		discordSend(race.ChannelID, msg)
		return
	}
	race.Racer2Bans = numBans
	if err := modals.Races.SetBans(race.ChannelID, 2, race.Racer2Bans); err != nil {
		msg := "Failed to set the bans for racer 2 on race \"" + race.Name() + "\": " + err.Error()
		log.Error(msg)
		discordSend(race.ChannelID, msg)
		return
	}

	msg += "**Build Ban Phase**\n\n"
	msg += "- Each racer gets to ban " + strconv.Itoa(numBans) + " builds.\n"
	msg += "- Use the `!ban` command to select a build.\n"
	msg += "  e.g. `!ban 3` (to ban the 3rd build in the list)\n\n"
	if race.ActiveRacer == 1 {
		msg += race.Racer1.Mention()
	} else if race.ActiveRacer == 2 {
		msg += race.Racer2.Mention()
	}
	msg += ", you start! (randomly decided)\n\n"

	msg += getBansRemaining(race)
	msg += getRemaining(race)
	discordSend(race.ChannelID, msg)
}

func buildsPickStart(race *Race, msg string) {
	// Set the state
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

	msg += getPicksRemaining(race)
	msg += getRemaining(race)
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

	msg += "**Build Ban Phase**\n\n"
	msg += "- " + strconv.Itoa(tournaments[race.ChallongeURL].BestOf) + " builds will randomly be chosen. Each racer will get one veto.\n"
	msg += "- Use the `!yes` and `!no` commands to answer the questions.\n\n"

	// The person who starts the vetos for the builds is the opposite of the person who got to start
	// the vetos for the character
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
			matchSetInProgressAndPrintSummary(race, msg)
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

func buildsEnd(race *Race, msg string) {
	// Unlike the "charactersEnd" function, we don't print the builds, since they will be displayed
	// in the summary.

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

	matchSetInProgressAndPrintSummary(race, msg)
}

func getBuild(race *Race) string {
	// Get a random build
	randBuildNum := getRandomInt(0, len(race.BuildsRemaining)-1)
	randBuild := race.BuildsRemaining[randBuildNum]

	// Check to see if the item synergizes
	roundNum := len(race.Builds) + 1
	character := race.Characters[roundNum-1]
	synergizes := true
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
