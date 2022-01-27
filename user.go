package main

import (
	"database/sql"

	"github.com/bwmarrin/discordgo"
)

type User struct {
	DiscordID string
	Username  string
	// Matches the TZ column of this page:
	// https://en.wikipedia.org/wiki/List_of_tz_database_time_zones
	Timezone       sql.NullString
	StreamURL      sql.NullString
	CasterAlwaysOk bool
}

func (u *User) Mention() string {
	return "<@" + u.DiscordID + ">"
}

func (u *User) GetTimezone() string {
	if !u.Timezone.Valid {
		return "UTC"
	}

	return u.Timezone.String
}

// Get this user from the database
// (and create an entry if it does not exist already)
func userGet(u *discordgo.User) (*User, error) {
	var user *User

	// Get the Discord guild members
	var members []*discordgo.Member
	if v, err := discordSession.GuildMembers(discordGuildID, "0", 1000); err != nil {
		return user, err
	} else {
		members = v
	}

	var username string
	if v, err := getDiscordNameByID(members, u.ID); err != nil {
		return user, err
	} else {
		username = v
	}

	// See if this user exists in the database already
	var exists bool
	if v, err := modals.Users.Exists(u.ID); err != nil {
		return user, err
	} else {
		exists = v
	}

	// This Discord ID already exists in the database, so return it
	if exists {
		if v, err := modals.Users.GetFromDiscordID(u.ID); err != nil {
			return user, err
		} else {
			user = v
		}

		if user.Username == username {
			// Their username in the database matches the Discord nickname
			return user, nil
		}

		// Their Discord username/nickname has changed since they were added to the database,
		// so we need to update it
		user.Username = username
		if err := modals.Users.SetUsername(u.ID, username); err != nil {
			return user, err
		} else {
			return user, nil
		}
	}

	// This Discord ID does not exist in the database, so create it
	user = &User{
		DiscordID: u.ID,
		Username:  u.Username,
	}
	err := modals.Users.Insert(user)
	return user, err
}
