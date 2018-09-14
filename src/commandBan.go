package main

import (
	"database/sql"
	"strconv"

	"github.com/Zamiell/isaac-tournament-bot/src/models"
	"github.com/bwmarrin/discordgo"
)

func commandBan(m *discordgo.MessageCreate, args []string) {
	if len(args) == 0 {
		commandBanPrint(m)
		return
	}

	// Check to see if this is a race channel (and get the race from the database)
	var race models.Race
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
	var playerNum int
	if m.Author.ID == race.Racer1.DiscordID {
		playerNum = 1
	} else if m.Author.ID == race.Racer2.DiscordID {
		playerNum = 2
	} else {
		discordSend(m.ChannelID, "Only \""+race.Racer1.Username+"\" and \""+race.Racer2.Username+"\" can ban something.")
		return
	}

	// Check to see if this race is in the banning phase
	if race.State != "banningCharacters" &&
		race.State != "banningBuilds" {

		discordSend(m.ChannelID, "You can only ban something once the match has started.")
		return
	}

	// Check to see if it is their turn
	if race.ActivePlayer != playerNum {
		discordSend(m.ChannelID, "It is not your turn.")
		return
	}

	// Check to see if they have any bans left
	if (playerNum == 1 && race.Racer1Bans == 0) ||
		(playerNum == 2 && race.Racer2Bans == 0) {

		discordSend(m.ChannelID, "You do not have any bans left.")
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
	var thingsRemaining []string
	if race.State == "banningCharacters" {
		thingsRemaining = race.CharactersRemaining
	} else if race.State == "banningBuilds" {
		thingsRemaining = race.BuildsRemaining
	}
	if choice < 0 || choice >= len(thingsRemaining) {
		discordSend(m.ChannelID, "\""+args[0]+"\" is not a valid choice.")
		return
	}

	// Ban the item
	thing := thingsRemaining[choice]
	thingsRemaining = deleteFromSlice(thingsRemaining, choice)

	if race.State == "banningCharacters" {
		race.CharactersRemaining = thingsRemaining
		if err := db.Races.SetCharactersRemaining(race.ChannelID, race.CharactersRemaining); err != nil {
			msg := "Failed to set the characters remaining for race \"" + race.Name() + "\": " + err.Error()
			log.Error(msg)
			discordSend(m.ChannelID, msg)
			return
		}
	} else if race.State == "banningBuilds" {
		race.BuildsRemaining = thingsRemaining
		if err := db.Races.SetBuildsRemaining(race.ChannelID, race.BuildsRemaining); err != nil {
			msg := "Failed to set the builds remaining for race \"" + race.Name() + "\": " + err.Error()
			log.Error(msg)
			discordSend(m.ChannelID, msg)
			return
		}
	}

	// Decrement their bans
	var bansLeft int
	if playerNum == 1 {
		race.Racer1Bans--
		bansLeft = race.Racer1Bans
	} else if playerNum == 2 {
		race.Racer2Bans--
		bansLeft = race.Racer2Bans
	}
	if err := db.Races.SetBans(race.ChannelID, playerNum, bansLeft); err != nil {
		msg := "Failed to set the bans for racer " + strconv.Itoa(playerNum) + " on race \"" + race.Name() + "\": " + err.Error()
		log.Error(msg)
		discordSend(m.ChannelID, msg)
		return
	}

	incrementActivePlayer(&race)

	msg := m.Author.Mention() + " banned **" + thing + "**.\n"
	totalBansLeft := race.Racer1Bans + race.Racer2Bans
	if totalBansLeft > 0 {
		msg += getNext(race)
		msg += getBansRemaining(race, "characters")
		msg += getRemaining(race, "characters")
		discordSend(race.ChannelID, msg)
	} else {
		msg += "\n"
		charactersPickStart(race, msg)
	}
}

func commandBanPrint(m *discordgo.MessageCreate) {
	msg := "Ban something with: `!ban [number]`\n"
	msg += "e.g. `!ban 3`\n"
	discordSend(m.ChannelID, msg)
}
