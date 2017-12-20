package main

import (
	"github.com/bwmarrin/discordgo"
)

var (
	commandHandlerMap = make(map[string]func(*discordgo.MessageCreate, []string))
)

func commandInit() {
	// General commands
	commandHandlerMap["help"] = commandHelp
	commandHandlerMap["r+"] = commandRacingPlus
	commandHandlerMap["racing+"] = commandRacingPlus
	commandHandlerMap["racingplus"] = commandRacingPlus
	commandHandlerMap["timezone"] = commandTimezone
	commandHandlerMap["stream"] = commandStream
	commandHandlerMap["schedule"] = commandSchedule
	commandHandlerMap["confirm"] = commandConfirm
	commandHandlerMap["reschedule"] = commandReschedule

	// Admin-only commands
	commandHandlerMap["startround"] = commandStartRound
	commandHandlerMap["beginround"] = commandStartRound
	commandHandlerMap["endround"] = commandEndRound
}

func commandHelpGetMsg() string {
	msg := "General commands:\n"
	msg += "```\n"
	msg += "Command                    Description\n"
	msg += "----------------------------------------------------------------------------------\n"
	msg += "!help                      Get a list of all of the commands\n"
	msg += "!r+                        Get info about the Racing+ mod\n"
	msg += "!timezone [timezone]       Set your timezone\n"
	msg += "!stream [url]              Set your stream URL\n"
	msg += "!schedule [date and time]  Suggest a time for the match to your opponent\n"
	msg += "!confirm                   Confirm that the suggested time is good\n"
	msg += "!reschedule                Delete the currently scheduled time\n"
	msg += "!startbans                 Start choosing character and item bans\n"
	msg += "```\n"
	msg += "Admin-only commands:\n"
	msg += "```\n"
	msg += "Command                    Description\n"
	msg += "----------------------------------------------------------------------------------\n"
	msg += "!startround                Create channels for the current round of the tournament\n"
	msg += "!endround                  Delete all of the channels for this round\n"
	msg += "```"

	return msg
}
