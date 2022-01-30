package main

import (
	"database/sql"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func commandCasterOk(m *discordgo.MessageCreate, args []string) {
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
	if m.Author.ID != race.Racer1.DiscordID && m.Author.ID != race.Racer2.DiscordID {
		discordSend(m.ChannelID, "You cannot give caster permission for a match that you are not participating in.")
		return
	}

	// Check to see if someone is casting this match
	if len(race.Casts) == 0 {
		discordSend(m.ChannelID, "No-one has volunteered to cast this match, so there is no need to give permission.")
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
	numNeedPermission := 0
	var cast *Cast
	for _, c := range race.Casts {
		if (racerNum == 1 && !c.R1Permission) ||
			(racerNum == 2 && !c.R2Permission) {

			numNeedPermission++
			cast = c
		}
	}
	if numNeedPermission == 0 {
		discordSend(m.ChannelID, "You have already given permission to all of the casters who have volunteered for this match.")
		return
	}

	// Get the corresponding cast
	// (there may be two or more casts for this match)
	if numNeedPermission >= 2 {
		// Check to see if they specified the caster's name that they are giving permission to
		// (they only need to do this if there are two or more casters that are awaiting permission)
		if len(args) != 1 {
			commandCasterOkPrint(m)
			return
		}

		for _, c := range race.Casts {
			if strings.EqualFold(c.Caster.Username, args[0]) {
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

	// Set permission
	if racerNum == 1 {
		cast.R1Permission = true
	} else if racerNum == 2 {
		cast.R2Permission = true
	}
	if err := modals.Casts.SetPermission(race.ChannelID, cast.Caster.DiscordID, racerNum); err != nil {
		msg := "Failed to set the caster permission in the database: " + err.Error()
		log.Error(msg)
		discordSend(m.ChannelID, msg)
		return
	}

	msg := "`" + racerName + "` has approved " + cast.Caster.Mention() + " as the caster for this match.\n"
	if cast.R1Permission && cast.R2Permission {
		msg += "Both racers have now approved this caster."
	} else if !cast.R1Permission {
		msg += race.Racer1.Mention() + " still needs to approve or disapprove this caster."
	} else if !cast.R2Permission {
		msg += race.Racer2.Mention() + " still needs to approve or disapprove this caster."
	}
	discordSend(m.ChannelID, msg)
}

func commandCasterOkPrint(m *discordgo.MessageCreate) {
	msg := "If there are two or more casters awaiting a response, then you need to specify the name of the caster by doing: `!casterok [username]`\n"
	msg += "e.g. `!casterok Willy`"
	discordSend(m.ChannelID, msg)
}
