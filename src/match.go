package main

import (
	"os"
	"strconv"
	"time"

	"github.com/Zamiell/isaac-tournament-bot/src/models"
)

var (
	numBans  int
	numVetos int
)

func matchInit() {
	// Read the configuration from environment variables
	numBansString := os.Getenv("NUM_BANS")
	if len(numBansString) == 0 {
		log.Fatal("The \"NUM_BANS\" environment variable is blank. Set it in the \".env\" file.")
		return
	}

	if v, err := strconv.Atoi(numBansString); err != nil {
		log.Fatal("The \"NUM_BANS\" environment variable is not a number.")
		return
	} else {
		numBans = v
	}

	numVetosString := os.Getenv("NUM_VETOS")
	if len(numBansString) == 0 {
		log.Fatal("The \"NUM_VETOS\" environment variable is blank. Set it in the \".env\" file.")
		return
	}

	if v, err := strconv.Atoi(numVetosString); err != nil {
		log.Fatal("The \"NUM_VETOS\" environment variable is not a number.")
		return
	} else {
		numVetos = v
	}

	// Make sure the build exceptions match the builds
	if len(builds) != len(buildExceptions) {
		log.Fatal("The builds were updated without also modifying the build exceptions.")
		return
	}

	// Schedule Discord pings for when each scheduled match starts
	var channelIDs []string
	if v, err := db.Races.GetAllScheduled(); err != nil {
		log.Fatal("Failed to get the scheduled races: " + err.Error())
		return
	} else {
		channelIDs = v
	}
	for _, channelID := range channelIDs {
		var race *models.Race
		if v, err := raceGet(channelID); err != nil {
			log.Fatal("Failed to get the race from the database: " + err.Error())
			return
		} else {
			race = v
		}

		go matchStart(race)
	}
}

func matchStart(race *models.Race) {
	// Sleep until the match starts
	origStartTime := race.DatetimeScheduled.Time
	sleepDuration := race.DatetimeScheduled.Time.Sub(time.Now().UTC())
	if sleepDuration < 5*time.Minute {
		sleepDuration = 0
	} else {
		sleepDuration -= 5 * time.Minute
	}
	time.Sleep(sleepDuration)

	// Re-get the race from the database
	if v, err := raceGet(race.ChannelID); err != nil {
		msg := "Failed to re-get the race from the database: " + err.Error()
		log.Error(msg)
		discordSend(race.ChannelID, msg)
		return
	} else {
		race = v
	}

	// Check to see if the race has been rescheduled
	if origStartTime != race.DatetimeScheduled.Time {
		return
	}

	// Check to see if this match has started already
	if race.State != "scheduled" {
		log.Info("Reached the \"matchStart\" function when the state was " + race.State + ". Doing nothing.")
		return
	}

	// Randomly decide who starts
	race.ActiveRacer = getRandom(1, 2)
	if err := db.Races.SetActiveRacer(race.ChannelID, race.ActiveRacer); err != nil {
		msg := "Failed to set the active racer for race \"" + race.Name() + "\": " + err.Error()
		log.Error(msg)
		discordSend(race.ChannelID, msg)
		return
	}

	// Announce that the match is starting in the general channel
	msg := "------------------------------------------\n"
	msg += "A race is scheduled to begin in 5 minutes:\n\n"
	msg += matchGetDescription(race)
	discordSend(discordGeneralChannelID, msg)

	// TODO Dynamically handle the kind of tournament
	// charactersBanStart(race)
	charactersVetoStart(race)
}

func matchBeginningAlert(race *models.Race) string {
	// Alert the racers that the race is about to start
	msg := race.Racer1.Mention() + " and " + race.Racer2.Mention() + " - the race is scheduled to start in 5 minutes.\n\n"

	// Alert the casters that the race is about to start
	for _, cast := range race.Casts {
		msg += cast.Caster.Mention() + ", you are scheduled to cast this match in " + languageMap[cast.Language] + " in 5 minutes at: <" + cast.Caster.StreamURL.String + ">\n\n"
	}

	return msg
}

func matchGetDescription(race *models.Race) string {
	msg := "```\n" // This is necessary because underscores in usernames can mess up the formatting
	msg += race.TournamentName + "\n"
	msg += race.Name() + "\n"
	msg += "```\n"
	atLeastOneCaster := false
	for _, cast := range race.Casts {
		if cast.R1Permission && cast.R2Permission {
			atLeastOneCaster = true
			msg += "`" + cast.Caster.Username + "` has volunteered to cast the match at:\n"
			msg += "<" + cast.Caster.StreamURL.String + ">\n"
		}
	}
	if !atLeastOneCaster {
		msg += "No-one has volunteered to cast this match. You can watch both racers here:\n"
		msg += "<https://kadgar.net/live/" + race.Racer1.Username + "/" + race.Racer2.Username + ">"
	}
	return msg
}

func matchEnd(race *models.Race, msg string) {
	race.State = "inProgress"
	if err := db.Races.SetState(race.ChannelID, race.State); err != nil {
		msg := "Failed to set the state for race \"" + race.Name() + "\": " + err.Error()
		log.Error(msg)
		discordSend(race.ChannelID, msg)
		return
	}

	msg += "```\n"
	msg += "+---------------+\n"
	msg += "| Match Summary |\n"
	msg += "+---------------+\n"
	msg += "```\n\n"

	msg += "**Racer 1: **" + race.Racer1.Mention() + " - <" + race.Racer1.StreamURL.String + ">\n"
	msg += "**Racer 2: **" + race.Racer2.Mention() + " - <" + race.Racer2.StreamURL.String + ">\n"
	for _, cast := range race.Casts {
		msg += "**" + languageMap[cast.Language] + " Caster:** " + cast.Caster.Mention() + " - <" + cast.Caster.StreamURL.String + ">\n"
	}
	msg += "\n"

	ruleset := tournaments[race.ChallongeURL].Ruleset
	for i := 0; i < tournaments[race.ChallongeURL].BestOf; i++ {
		msg += "**Round " + strconv.Itoa(i+1) + "**:\n"
		msg += "- Character: *" + race.Characters[i] + "*\n"
		if ruleset == "seeded" {
			msg += "- Build: *" + race.Builds[i] + "*\n"
		}
		msg += "\n"
	}
	msg += "If I made a mistake, you can use `!randchar` "
	if ruleset == "seeded" {
		msg += "or `!randbuild` "
	}
	msg += "to manually get random characters"
	if ruleset == "seeded" {
		msg += " and builds"
	}
	msg += ".\n"
	msg += "When the race is over, please use the `!score [score]` command to report the results.\n"
	msg += "e.g. `!score 3-2`\n\n"
	msg += "Good luck and have fun!"
	discordSend(race.ChannelID, msg)
}
