package main

import (
	"strconv"

	"github.com/Zamiell/isaac-tournament-bot/src/models"
	"github.com/bwmarrin/discordgo"
)

func commandScore(m *discordgo.MessageCreate, args []string) {
	if len(args) == 0 {
		commandScorePrint(m)
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
	var playerNum int
	if m.Author.ID == race.Racer1.DiscordID {
		playerNum = 1
	} else if m.Author.ID == race.Racer2.DiscordID {
		playerNum = 2
	} else {
		discordSend(m.ChannelID, "Only \""+race.Racer1.Username+"\" and \""+race.Racer2.Username+"\" can veto a build.")
		return
	}

	// Check to see if this race is in progress
	if race.State != "inProgress" {
		discordSend(m.ChannelID, "You can only report the score once you have finished picking characters and builds.")
		return
	}

	// Check to see if this score was reported in the correct format
	score := args[0]
	scoreValid := true
	if len(score) != 3 {
		scoreValid = false
	}
	if score[1] != '-' {
		scoreValid = false
	}
	var p1Wins, p2Wins int
	if v, err := strconv.Atoi(string(score[0])); err != nil {
		scoreValid = false
	} else if playerNum == 1 {
		p1Wins = v
	} else if playerNum == 2 {
		p2Wins = v
	}
	if v, err := strconv.Atoi(string(score[2])); err != nil {
		scoreValid = false
	} else if playerNum == 2 {
		p1Wins = v
	} else if playerNum == 1 {
		p2Wins = v
	}
	if !scoreValid {
		msg := "You must report the score in the following format: `!score #-#`\n"
		msg += "e.g. `!score 3-2`"
		discordSend(m.ChannelID, msg)
		return
	}

	// Make sure the score is in the right order
	// (player 1 must be first)
	score = strconv.Itoa(p1Wins) + "-" + strconv.Itoa(p2Wins)

	// Update the match on Challonge
	// https://api.challonge.com/v1/documents/matches/update
	challongeTournamentID := floatToString(tournaments[race.TournamentName].ChallongeID)
	apiURL := "https://api.challonge.com/v1/tournaments/" + challongeTournamentID + "/matches/" + race.ChallongeID + ".json?"
	apiURL += "api_key=" + challongeAPIKey + "&match[scores_csv]=" + score
	if _, err := challongeGetJSON("PUT", apiURL, nil); err != nil {
		msg := "Failed to get the tournament from Challonge: " + err.Error()
		log.Error(msg)
		discordSend(race.ChannelID, msg)
		return
	}

	msg := "Score successfully submitted."
	discordSend(m.ChannelID, msg)
	log.Info("User \"" + m.Author.Username + "\" submitted a score of \"" + score + "\" for race \"" + race.Name() + "\".")
}

func commandScorePrint(m *discordgo.MessageCreate) {
	msg := "Report the score of the match with: `!score [score]`\n"
	msg += "e.g. `!score 3-2`"
	discordSend(m.ChannelID, msg)
}
