package main

import (
	"database/sql"
	"strings"

	"github.com/Zamiell/isaac-tournament-bot/src/models"
	"github.com/bwmarrin/discordgo"
)

func commandCasterNotOk(m *discordgo.MessageCreate, args []string) {
	// Check to see if this is a race channel (and get the race from the database)
	var race *models.Race
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
	if m.Author.ID != race.Racer1.DiscordID && m.Author.ID != race.Racer2.DiscordID {
		discordSend(m.ChannelID, "You cannot deny caster permission for a match that you are not participating in.")
		return
	}

	// Check to see if someone is casting this match
	if len(race.Casts) == 0 {
		discordSend(m.ChannelID, "No-one has volunteered to cast this match, so there is no need to deny permission.")
		return
	}

	// Find out whether they are racer 1 or racer 2
	racerNum := 1
	racerName := race.Racer1.Username
	if m.Author.ID == race.Racer2.DiscordID {
		racerNum = 2
		racerName = race.Racer2.Username
	}

	// Check to see if they have already given permission to everyone who has volunteered to cast
	numPermission := 0
	var cast *models.Cast
	for _, c := range race.Casts {
		if (racerNum == 1 && c.R1Permission) ||
			(racerNum == 2 && c.R2Permission) {

			numPermission++
			cast = c
		}
	}
	if numPermission == 0 {
		discordSend(m.ChannelID, "You have not yet give permission to any of the casters who have volunteered for this match.")
		return
	}

	// Get the corresponding cast
	// (there may be two or more casts for this match)
	if numPermission >= 2 {
		// Check to see if they specified the caster's name that they are denying permission to
		// (they only need to do this if there are two or more casters that are awaiting permission)
		if len(args) != 1 {
			commandCasterNotOkPrint(m)
			return
		}

		for _, c := range race.Casts {
			if strings.ToLower(c.Caster.Username) == strings.ToLower(args[0]) {
				cast = c
				break
			}
		}
		if cast == nil {
			msg := "`" + args[0] + "` has not volunteered to cast this match. Did you make a typo?"
			discordSend(m.ChannelID, msg)
			return
		}
	}

	// Delete this cast from the database
	if err := db.Casts.Delete(race.ChannelID, cast.Caster.DiscordID); err != nil {
		msg := "Failed to delete the cast from the database: " + err.Error()
		log.Error(msg)
		discordSend(m.ChannelID, msg)
		return
	}

	msg := "`" + racerName + "` has denied permission for " + cast.Caster.Mention() + " to rebroadcast the race. They have been removed as a registered caster for this match."
	discordSend(m.ChannelID, msg)
}

func commandCasterNotOkPrint(m *discordgo.MessageCreate) {
	msg := "Disapprove a caster by doing: `!casternotok`\n"
	msg += "If there are two or more casters that you need to disapprove, then you need to specify the name of the caster by doing: `!casternotok [username]\n"
	msg += "e.g. `!casternotok Willy`"
	discordSend(m.ChannelID, msg)
}
