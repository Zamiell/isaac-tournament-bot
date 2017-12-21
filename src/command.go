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
	commandHandlerMap["gettimezone"] = commandGetTimezone
	commandHandlerMap["stream"] = commandStream
	commandHandlerMap["getstream"] = commandGetStream

	// Match commands
	commandHandlerMap["time"] = commandTime
	commandHandlerMap["timeok"] = commandTimeOk
	commandHandlerMap["timedelete"] = commandTimeDelete
	commandHandlerMap["cast"] = commandCast
	commandHandlerMap["castcancel"] = commandCastCancel
	commandHandlerMap["caster"] = commandCaster
	commandHandlerMap["casterok"] = commandCasterOk
	commandHandlerMap["casternotok"] = commandCasterNotOk
	commandHandlerMap["startbans"] = commandStartBans
	commandHandlerMap["bansstart"] = commandStartBans

	// Admin-only commands
	commandHandlerMap["settimezone"] = commandSetTimezone
	commandHandlerMap["timezoneset"] = commandSetTimezone
	commandHandlerMap["startround"] = commandStartRound
	commandHandlerMap["roundstart"] = commandStartRound
	commandHandlerMap["beginround"] = commandStartRound
	commandHandlerMap["roundbegin"] = commandStartRound
	commandHandlerMap["endround"] = commandEndRound
	commandHandlerMap["roundend"] = commandEndRound
}

func commandHelpGetMsg() string {
	msg := "General commands (all channels):\n"
	msg += "```\n"
	msg += "Command                  Description\n"
	msg += "---------------------------------------------------------------------------------\n"
	msg += "!help                    Get a list of all of the commands\n"
	msg += "!r+                      Get info about the Racing+ mod\n"
	msg += "!timezone                Get your stored timezone\n"
	msg += "!timezone [timezone]     Set your stored timezone\n"
	msg += "!gettimezone [username]  Get the timezone of the specified person\n"
	msg += "!stream                  Get your stored stream URL\n"
	msg += "!stream [url]            Set your stored stream URL\n"
	msg += "!getstream [username]    Get the stream of the specified person\n"
	msg += "```\n"
	msg += "Match commands (in a match channel):\n"
	msg += "```\n"
	msg += "Command                  Description\n"
	msg += "---------------------------------------------------------------------------------\n"
	msg += "!time                    Get the currently scheduled time for the match"
	msg += "!time [date & time]      Suggest a time for the match to your opponent\n"
	msg += "!timeok                  Confirm that the suggested time is good\n"
	msg += "!timedelete              Delete the currently scheduled time\n"
	msg += "!cast                    Volunteer to be the caster to rebroadcast this match\n"
	msg += "!castcancel              Unvolunteer to be the caster"
	msg += "!caster                  Get the person who volunteered to cast\n"
	msg += "!casterok                Confirm that you are okay with the caster rebroadcasting\n"
	msg += "!casternotok             Undo your caster confirmation or deny the current caster\n"
	msg += "!startbans               Start choosing character and item bans\n"
	msg += "```\n"
	msg += "Admin-only commands:\n"
	msg += "```\n"
	msg += "Command               Description\n"
	msg += "---------------------------------------------------------------------------------\n"
	msg += "!settimezone             Set a player's timezone for them\n"
	msg += "!startround              Create channels for the current round of the tournament\n"
	msg += "!endround                Delete all of the channels for this round\n"
	msg += "```"

	return msg
}
