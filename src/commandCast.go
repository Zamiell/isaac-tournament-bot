package main

import (
	"database/sql"

	"github.com/Zamiell/isaac-tournament-bot/src/models"
	"github.com/bwmarrin/discordgo"
)

func commandCast(m *discordgo.MessageCreate, args []string) {
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

	// Create the user in the database if it does not already exist
	var caster models.Racer
	if v, err := racerGet(m.Author); err != nil {
		msg := "Failed to get the racer from the database: " + err.Error()
		log.Error(msg)
		discordSend(m.ChannelID, msg)
		return
	} else {
		caster = v
	}

	// Check to see if they have a stream set
	if !caster.StreamURL.Valid {
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

	// Check to see if someone is already casting this match
	if race.CasterID.Valid {
		discordSend(m.ChannelID, race.Caster.Username+" has already volunteered to cast this match.")
		return
	}

	// Set them as the new caster
	if err := db.Races.SetCaster(m.ChannelID, m.Author.ID); err != nil {
		msg := "Failed to set the new caster in the database: " + err.Error()
		log.Error(msg)
		discordSend(m.ChannelID, msg)
		return
	}

	msg := caster.Mention() + ", you are now registered as the caster for this match at the following stream: <" + caster.StreamURL.String + ">\n"

	if race.Racer1.CasterAlwaysOk {
		if err := db.Races.SetCasterApproval(race.Racer1.DiscordID, 1); err != nil {
			msg := "Failed to set the caster approval in the database: " + err.Error()
			log.Error(msg)
			discordSend(m.ChannelID, msg)
			return
		}
		msg += race.Racer1.Username + " has automatically approved all casters.\n"
	}
	if race.Racer2.CasterAlwaysOk {
		if err := db.Races.SetCasterApproval(race.Racer2.DiscordID, 2); err != nil {
			msg := "Failed to set the caster approval in the database: " + err.Error()
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
