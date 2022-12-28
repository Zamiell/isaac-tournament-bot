package main

import (
	"database/sql"
	"strconv"

	"github.com/bwmarrin/discordgo"
)

func commandBan(m *discordgo.MessageCreate, args []string) {
	if len(args) == 0 {
		commandBanPrint(m)
		return
	}

	// Check to see if this is a race channel (and get the race from the database)
	var race *Race
	if v, err := getRace(m.ChannelID); err == sql.ErrNoRows {
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

	// Check to see if this race is in the banning phase
	if race.State != RaceStateBanningCharacters &&
		race.State != RaceStateBanningBuilds {

		discordSend(m.ChannelID, "You can only ban something once the match has started.")
		return
	}

	// Check to see if it is their turn
	if race.ActiveRacer != racerNum {
		discordSend(m.ChannelID, "It is not your turn.")
		return
	}

	// Check to see if they have any bans left
	if (racerNum == 1 && race.Racer1Bans == 0) ||
		(racerNum == 2 && race.Racer2Bans == 0) {

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
	if race.State == RaceStateBanningCharacters {
		thingsRemaining = race.CharactersRemaining
	} else if race.State == RaceStateBanningBuilds {
		thingsRemaining = race.BuildsRemaining
	}
	if choice < 0 || choice >= len(thingsRemaining) {
		discordSend(m.ChannelID, "\""+args[0]+"\" is not a valid choice.")
		return
	}

	// Ban the item
	thing := thingsRemaining[choice]
	thingsRemaining = deleteFromSlice(thingsRemaining, choice)

	if race.State == RaceStateBanningCharacters {
		race.CharactersRemaining = thingsRemaining
		if err := modals.Races.SetCharactersRemaining(race.ChannelID, race.CharactersRemaining); err != nil {
			msg := "Failed to set the characters remaining for race \"" + race.Name() + "\": " + err.Error()
			log.Error(msg)
			discordSend(m.ChannelID, msg)
			return
		}
	} else if race.State == RaceStateBanningBuilds {
		race.BuildsRemaining = thingsRemaining
		if err := modals.Races.SetBuildsRemaining(race.ChannelID, race.BuildsRemaining); err != nil {
			msg := "Failed to set the builds remaining for race \"" + race.Name() + "\": " + err.Error()
			log.Error(msg)
			discordSend(m.ChannelID, msg)
			return
		}
	}

	// Decrement their bans
	var bansLeft int
	if racerNum == 1 {
		race.Racer1Bans--
		bansLeft = race.Racer1Bans
	} else if racerNum == 2 {
		race.Racer2Bans--
		bansLeft = race.Racer2Bans
	}
	if err := modals.Races.SetBans(race.ChannelID, racerNum, bansLeft); err != nil {
		msg := "Failed to set the bans for racer " + strconv.Itoa(racerNum) + " on race \"" + race.Name() + "\": " + err.Error()
		log.Error(msg)
		discordSend(m.ChannelID, msg)
		return
	}

	incrementActiveRacer(race)

	msg := m.Author.Mention() + " banned **" + thing + "**.\n"
	totalBansLeft := race.Racer1Bans + race.Racer2Bans
	if totalBansLeft > 0 {
		msg += getNext(race)
		msg += getBansRemaining(race)
		msg += getRemaining(race)
		discordSend(race.ChannelID, msg)
	} else {
		msg += "\n"
		if race.State == RaceStateBanningCharacters {
			charactersPickStart(race, msg)
		} else if race.State == RaceStateBanningBuilds {
			buildsPickStart(race, msg)
		} else {
			msg += "The race state was invalid."
			discordSend(race.ChannelID, msg)
		}
	}
}

func commandBanPrint(m *discordgo.MessageCreate) {
	msg := "Ban something with: `!ban [number]`\n"
	msg += "e.g. `!ban 3`"
	discordSend(m.ChannelID, msg)
}
