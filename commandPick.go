package main

import (
	"database/sql"
	"strconv"

	"github.com/bwmarrin/discordgo"
)

func commandPick(m *discordgo.MessageCreate, args []string) {
	if len(args) == 0 {
		commandPickPrint(m)
		return
	}

	// Check to see if this is a race channel (and get the race from the database)
	var race *Race
	if v, err := raceGet(m.ChannelID); err == sql.ErrNoRows {
		discordSend(m.ChannelID, "You can only use that command in a race channel.")
		return
	} else if err != nil {
		msg := "Failed to get the race from the database: " + err.Error()
		log.Error(msg)
		discordSend(m.ChannelID, msg)
		return
	} else {
		race = v
	}

	// Check to see if this person is one of the two racers
	var racerNum int
	if m.Author.ID == race.Racer1.DiscordID {
		racerNum = 1
	} else if m.Author.ID == race.Racer2.DiscordID {
		racerNum = 2
	} else {
		discordSend(m.ChannelID, "Only \""+race.Racer1.Username+"\" and \""+race.Racer2.Username+"\" can ban something.")
		return
	}

	// Check to see if this race is in the picking phase
	if race.State != "pickingCharacters" &&
		race.State != "pickingBuilds" {

		discordSend(m.ChannelID, "You can only pick something once the banning phase has finished.")
		return
	}

	// Check to see if it is their turn
	if race.ActiveRacer != racerNum {
		discordSend(m.ChannelID, "It is not your turn.")
		return
	}

	// Check to see if this is a valid number
	var choice int
	if v, err := strconv.Atoi(args[0]); err != nil {
		discordSend(m.ChannelID, "\""+args[0]+"\" is not a number.")
		return
	} else {
		choice = v
	}

	// Account for the fact that the array is 0 indexed and the choices presented to the user begin at 1
	choice--

	// Check to see if this is a valid index
	var thingsRemaining, things []string
	if race.State == "pickingCharacters" {
		thingsRemaining = race.CharactersRemaining
		things = race.Characters
	} else if race.State == "pickingBuilds" {
		thingsRemaining = race.BuildsRemaining
		things = race.Builds
	}
	if choice < 0 || choice >= len(thingsRemaining) {
		discordSend(m.ChannelID, "\""+args[0]+"\" is not a valid choice.")
		return
	}

	// Pick the thing
	thing := thingsRemaining[choice]
	thingsRemaining = deleteFromSlice(thingsRemaining, choice)
	things = append(things, thing)

	if race.State == "pickingCharacters" {
		race.CharactersRemaining = thingsRemaining
		if err := modals.Races.SetCharactersRemaining(race.ChannelID, race.CharactersRemaining); err != nil {
			msg := "Failed to set the characters remaining for race \"" + race.Name() + "\": " + err.Error()
			log.Error(msg)
			discordSend(m.ChannelID, msg)
			return
		}

		race.Characters = things
		if err := modals.Races.SetCharacters(race.ChannelID, race.Characters); err != nil {
			msg := "Failed to set the characters for race \"" + race.Name() + "\": " + err.Error()
			log.Error(msg)
			discordSend(m.ChannelID, msg)
			return
		}
	} else if race.State == "pickingBuilds" {
		race.BuildsRemaining = thingsRemaining
		if err := modals.Races.SetBuildsRemaining(race.ChannelID, race.BuildsRemaining); err != nil {
			msg := "Failed to set the builds remaining for race \"" + race.Name() + "\": " + err.Error()
			log.Error(msg)
			discordSend(m.ChannelID, msg)
			return
		}

		race.Builds = things
		if err := modals.Races.SetBuilds(race.ChannelID, race.Characters); err != nil {
			msg := "Failed to set the builds for race \"" + race.Name() + "\": " + err.Error()
			log.Error(msg)
			discordSend(m.ChannelID, msg)
			return
		}
	}

	incrementActiveRacer(race)

	msg := m.Author.Mention() + " picked **" + thing + "**.\n"
	picksLeft := tournaments[race.ChallongeURL].BestOf - len(things)
	if picksLeft > 0 {
		msg += getNext(race)
		msg += getPicksRemaining(race, "characters")
		msg += getRemaining(race, "characters")
		discordSend(race.ChannelID, msg)
	} else {
		msg += "\n"
		charactersEnd(race, msg)
	}
}

func commandPickPrint(m *discordgo.MessageCreate) {
	msg := "Pick something with: `!pick [number]`\n"
	msg += "e.g. `!pick 3`"
	discordSend(m.ChannelID, msg)
}
