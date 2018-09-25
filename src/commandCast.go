package main

import (
	"database/sql"
	"strings"

	"github.com/Zamiell/isaac-tournament-bot/src/models"
	"github.com/bwmarrin/discordgo"
)

func commandCast(m *discordgo.MessageCreate, args []string) {
	if len(args) != 1 {
		commandCastPrint(m)
		return
	}
	language := strings.ToLower(args[0])

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

	// Create the user in the database if it does not already exist
	var user *models.User
	if v, err := userGet(m.Author); err != nil {
		msg := "Failed to get the user from the database: " + err.Error()
		log.Error(msg)
		discordSend(m.ChannelID, msg)
		return
	} else {
		user = v
	}

	// Check to see if they have a stream set
	if !user.StreamURL.Valid {
		discordSend(m.ChannelID, "You cannot volunteer to cast a match if you do not have a stream URL set. Please set one first with the `!stream` command.")
		return
	}

	// Check to see if this person is one of the two racers
	if m.Author.ID == race.Racer1.DiscordID || m.Author.ID == race.Racer2.DiscordID {
		discordSend(m.ChannelID, "You cannot cast a match that you are participating in.")
		return
	}

	// Check to see if this race has already been scheduled
	if race.State != "scheduled" {
		discordSend(m.ChannelID, "You cannot volunteer to cast a match until a time has been scheduled by both of the racers.")
		return
	}

	// Check to see if this is a valid language
	valid := false
	var languageFull string
	for k, v := range languageMap {
		if strings.ToLower(k) == language || strings.ToLower(v) == language {
			valid = true
			language = k
			languageFull = v
			break
		}
	}
	if !valid {
		msg := "That is not a valid language. Valid languages are:\n"
		for k, v := range languageMap {
			msg += k + " / " + v + ", "
		}
		msg = strings.TrimSuffix(msg, ", ")
		discordSend(m.ChannelID, msg)
		return
	}

	// Check to see if they are already casting this match
	for _, cast := range race.Casts {
		if cast.Caster.DiscordID == m.Author.ID {
			msg := "You have already volunteered to cast this match."
			discordSend(m.ChannelID, msg)
			return
		}
	}

	// Check to see if someone else is already casting this match in that language
	for _, cast := range race.Casts {
		if cast.Language == language {
			msg := "This match is already being casted in " + languageFull + " by `" + cast.Caster.Username + "`."
			discordSend(m.ChannelID, msg)
			return
		}
	}

	// Add them as a new caster
	if err := db.Casts.Insert(race.ChannelID, user.DiscordID, language); err != nil {
		msg := "Failed to insert the new cast in the database: " + err.Error()
		log.Error(msg)
		discordSend(m.ChannelID, msg)
		return
	}

	msg := user.Mention() + ", you are now registered as the " + languageFull + " caster for this match with the following stream: <" + user.StreamURL.String + ">\n"

	if race.Racer1.CasterAlwaysOk {
		if err := db.Casts.SetPermission(race.ChannelID, user.DiscordID, 1); err != nil {
			msg := "Failed to set the caster approval for racer 1 in the database: " + err.Error()
			log.Error(msg)
			discordSend(m.ChannelID, msg)
			return
		}
		msg += race.Racer1.Username + " has automatically approved all casters.\n"
	}
	if race.Racer2.CasterAlwaysOk {
		if err := db.Casts.SetPermission(race.ChannelID, user.DiscordID, 2); err != nil {
			msg := "Failed to set the caster approval for racer 2 in the database: " + err.Error()
			log.Error(msg)
			discordSend(m.ChannelID, msg)
			return
		}
		msg += race.Racer2.Username + " has automatically approved all casters.\n"
	}
	if !race.Racer1.CasterAlwaysOk || !race.Racer2.CasterAlwaysOk {
		if !race.Racer1.CasterAlwaysOk && !race.Racer2.CasterAlwaysOk {
			msg += "Both " + race.Racer1.Mention() + " and " + race.Racer2.Mention()
		} else if !race.Racer1.CasterAlwaysOk {
			msg += race.Racer1.Mention()
		} else if !race.Racer2.CasterAlwaysOk {
			msg += race.Racer2.Mention()
		}
		msg += " must agree to this with the `!casterok` command. If you do not agree, use the `!casternotok` command.\n"
		msg += "(You can also use the `!casteralwaysok` command to give blanket permission for everyone to cast.)"
	}
	discordSend(m.ChannelID, msg)
}

func commandCastPrint(m *discordgo.MessageCreate) {
	msg := "Volunteer to cast this match by doing: `!cast [language]`\n"
	msg += "e.g. `!cast en`"
	discordSend(m.ChannelID, msg)
}
