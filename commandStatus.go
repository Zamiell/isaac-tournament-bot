package main

import (
	"database/sql"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func commandStatus(m *discordgo.MessageCreate, args []string) {
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

	msg1 := "The current state of this race channel is: " + race.State
	discordSend(m.ChannelID, msg1)

	msg2 := "The current builds are: " + strings.Join(race.Builds, ", ")
	msg2 += " (" + strconv.Itoa(len(race.Builds)) + ")"
	discordSend(m.ChannelID, msg2)

	msg3 := "The builds left to pick are: " + strings.Join(race.BuildsRemaining, ", ")
	msg3 += " (" + strconv.Itoa(len(race.BuildsRemaining)) + ")"
	discordSend(m.ChannelID, msg3)
}
