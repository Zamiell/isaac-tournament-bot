package main

import (
	"net/url"
	"strings"

	"github.com/Zamiell/isaac-tournament-bot/src/models"
	"github.com/bwmarrin/discordgo"
)

func commandStream(m *discordgo.MessageCreate, args []string) {
	if len(args) == 0 {
		commandStreamPrint(m)
		return
	}
	streamURL := args[0]

	// Create the user in the database if it does not already exist
	var racer models.Racer
	if v, err := racerGet(m.Author); err != nil {
		msg := "Failed to get the racer from the database: " + err.Error()
		log.Error(msg)
		discordSend(m.ChannelID, msg)
		return
	} else {
		racer = v
	}

	// Lower-case the URL
	streamURL = strings.ToLower(streamURL)

	// Fill in the URL in case they were lazy when typing it
	if strings.HasPrefix(streamURL, "http://") {
		streamURL = strings.Replace(streamURL, "http://", "https://", -1)
	}
	if strings.HasPrefix(streamURL, "https://twitch.tv/") {
		streamURL = strings.Replace(streamURL, "twitch.tv", "www.twitch.tv", -1)
	}
	if strings.HasPrefix(streamURL, "twitch.tv/") {
		streamURL = "https://www." + streamURL
	}
	if strings.HasPrefix(streamURL, "www.twitch.tv/") {
		streamURL = "https://" + streamURL
	}
	if strings.HasSuffix(streamURL, "/") {
		streamURL = strings.TrimSuffix(streamURL, "/")
	}

	// Validate that the stream is a full URL
	if _, err := url.ParseRequestURI(streamURL); err != nil {
		msg := "That is not a valid URL."
		discordSend(m.ChannelID, msg)
		return
	}

	// Set the new stream
	if err := db.Racers.SetStreamURL(m.Author.ID, streamURL); err != nil {
		msg := "Failed to update the stream: " + err.Error()
		log.Error(msg)
		discordSend(m.ChannelID, msg)
		return
	}

	msg := "The stream for **" + racer.Username + "** has been set to: <" + streamURL + ">"
	discordSend(m.ChannelID, msg)
}

func commandStreamPrint(m *discordgo.MessageCreate) {
	var racer models.Racer
	if v, err := racerGet(m.Author); err != nil {
		msg := "Failed to get the racer from the database: " + err.Error()
		log.Error(msg)
		discordSend(m.ChannelID, msg)
		return
	} else {
		racer = v
	}

	msg := m.Author.Mention() + ", your stream is "
	if racer.StreamURL.Valid {
		msg += "currently set to: <" + racer.StreamURL.String + ">\n\n"
	} else {
		msg += "**not currently set**.\n\n"
	}
	msg += "Set your stream with: `!stream [url]`\n"
	msg += "e.g. `!stream https://www.twitch.tv/zamiell`"
	discordSend(m.ChannelID, msg)
}
