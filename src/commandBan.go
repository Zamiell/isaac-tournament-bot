package main

import (
	"strconv"
	"strings"

	"github.com/Zamiell/isaac-tournament-bot/src/models"
	"github.com/bwmarrin/discordgo"
)

func commandBan(m *discordgo.MessageCreate, args []string) {
	if len(args) == 0 {
		commandBanPrint(m)
		return
	}

	var race models.Race
	if v, err := raceGet(m.ChannelID); err != nil {
		msg := "Failed to get the race from the database: " + err.Error()
		log.Error(msg)
		discordSend(m.ChannelID, msg)
		return
	} else {
		race = v
	}

	// Check to see if this person is one of the two racers
	if m.Author.ID != race.Racer1.DiscordID && m.Author.ID != race.Racer2.DiscordID {
		discordSend(m.ChannelID, "Only \""+race.Racer1.Username+"\" and \""+race.Racer2.Username+"\" can ban something.")
		return
	}

	// Check to see if this race is in the banning phase
	if race.State != 2 {
		discordSend(m.ChannelID, "You can only ban something once the match has started.")
		return
	}

	// Check to see if it is their turn
	if (race.ActivePlayer == 1 && m.Author.ID != race.Racer1.DiscordID) ||
		(race.ActivePlayer == 2 && m.Author.ID != race.Racer2.DiscordID) {

		discordSend(m.ChannelID, "It is not your turn.")
		return
	}

	// Check to see if this is a valid number
	var banNum int
	if v, err := strconv.Atoi(args[0]); err != nil {
		discordSend(m.ChannelID, "\""+args[0]+"\" is not a number.")
		return
	} else {
		banNum = v
	}

	// Account for the fact that the array is 0 indexed and the choices presented to the user begin at 1
	banNum -= 1

	// Check to see if this is a valid index
	var thingsString string
	var thingsFull []string
	if race.State == 2 {
		thingsString = race.Characters
		thingsFull = characters
	} else if race.State == 3 {
		thingsString = race.Builds
		thingsFull = builds
	}
	things := strings.Split(thingsString, ",")
	if banNum < 0 || banNum >= len(things) {
		discordSend(m.ChannelID, "\""+args[0]+"\" is not a valid choice.")
		return
	}

	// Ban the item
	thing := things[banNum]
	things = deleteFromSlice(things, banNum)
	thingsString = strings.Join(things, ",")
	if race.State == 2 {
		race.Characters = thingsString
		if err := db.Races.SetCharacters(race.ChannelID, race.Characters); err != nil {
			msg := "Failed to set the characters for race \"" + race.Name() + "\": " + err.Error()
			log.Error(msg)
			discordSend(m.ChannelID, msg)
			return
		}
	} else if race.State == 3 {
		race.Builds = thingsString
		if err := db.Races.SetBuilds(race.ChannelID, race.Builds); err != nil {
			msg := "Failed to set the builds for race \"" + race.Name() + "\": " + err.Error()
			log.Error(msg)
			discordSend(m.ChannelID, msg)
			return
		}
	}

	// Increment the active player
	race.ActivePlayer++
	if race.ActivePlayer > 2 {
		race.ActivePlayer = 1
	}
	if err := db.Races.SetActivePlayer(race.ChannelID, race.ActivePlayer); err != nil {
		msg := "Failed to set the active player for race \"" + race.Name() + "\": " + err.Error()
		log.Error(msg)
		discordSend(race.ChannelID, msg)
		return
	}

	msg := m.Author.Username + " banned **" + thing + "**.\n"
	bansLeft := len(things) - len(thingsFull) + bansNum
	if bansLeft > 0 {
		msg += getNext(race)
		msg += getRemaining(race, "characters")
	} else {
		msg += "\n**Build Ban Phase**\n\n"
		msg += "- 5 builds will randomly be chosen. Each player will get one veto.\n"
		msg += "- Use the `!yes` and `!no` commands to answer the questions.\n\n"
		msg += ""
	}
	discordSend(race.ChannelID, msg)
}

func commandBanPrint(m *discordgo.MessageCreate) {
	msg := "Ban something with: `!ban [number]`\n"
	msg += "For example: `!ban 3`\n"
	discordSend(m.ChannelID, msg)
}
