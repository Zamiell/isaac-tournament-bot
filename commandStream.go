package main

import (
	"net/url"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func commandStream(m *discordgo.MessageCreate, args []string) {
	if len(args) == 0 {
		commandStreamPrint(m)
		return
	}
	streamURL := args[0]

	// Create the user in the database if it does not already exist.
	var user *User
	if v, err := userGet(m.Author); err != nil {
		msg := "Failed to get the user from the database: " + err.Error()
		log.Error(msg)
		discordSend(m.ChannelID, msg)
		return
	} else {
		user = v
	}

	// Lower-case the URL.
	streamURL = strings.ToLower(streamURL)

	// Fill in the URL in case they were lazy when typing it.
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
	streamURL = strings.TrimSuffix(streamURL, "/")

	// Validate that the stream is a full URL.
	if _, err := url.ParseRequestURI(streamURL); err != nil {
		msg := "That is not a valid URL."
		discordSend(m.ChannelID, msg)
		return
	}

	// Set the new stream.
	if err := modals.Users.SetStreamURL(m.Author.ID, streamURL); err != nil {
		msg := "Failed to update the stream: " + err.Error()
		log.Error(msg)
		discordSend(m.ChannelID, msg)
		return
	}

	msg := "The stream for **" + user.Username + "** has been set to: <" + streamURL + ">"
	discordSend(m.ChannelID, msg)
}

func commandStreamPrint(m *discordgo.MessageCreate) {
	var user *User
	if v, err := userGet(m.Author); err != nil {
		msg := "Failed to get the user from the database: " + err.Error()
		log.Error(msg)
		discordSend(m.ChannelID, msg)
		return
	} else {
		user = v
	}

	msg := m.Author.Mention() + ", your stream is "
	if user.StreamURL.Valid {
		msg += "currently set to: <" + user.StreamURL.String + ">\n\n"
	} else {
		msg += "**not currently set**.\n\n"
	}
	msg += "Set your stream with: `!stream [url]`\n"
	msg += "e.g. `!stream https://www.twitch.tv/willy`"
	discordSend(m.ChannelID, msg)
}
