package main

import (
	"os"
	"strings"
	"sync"

	"github.com/bwmarrin/discordgo"
)

const (
	discordAdminRoleName       = "Admins"
	discordBotRoleName         = "Bots"
	discordCasterRoleName      = "Casters"
	discordTeamCaptainRoleName = "Team Captain"
	discordGeneralChannelName  = "general"
)

var (
	discord                  *discordgo.Session
	discordBotID             string
	discordGuildName         string
	discordGuildID           string
	discordAdminRoleID       string
	discordBotRoleID         string
	discordCasterRoleID      string
	discordEveryoneRoleID    string
	discordGeneralChannelID  string
	discordTeamCaptainRoleID string
	commandMutex             = new(sync.Mutex)
)

func discordInit() {
	// Read the OAuth secret from the environment variable
	discordToken := os.Getenv("DISCORD_TOKEN")
	if len(discordToken) == 0 {
		log.Fatal("The \"DISCORD_TOKEN\" environment variable is blank. Set it in the \".env\" file.")
		return
	}

	discordGuildName = os.Getenv("DISCORD_SERVER_NAME")
	if len(discordGuildName) == 0 {
		log.Fatal("The \"DISCORD_SERVER_NAME\" environment variable is blank. Set it in the \".env\" file.")
		return
	}

	// Bot accounts must be prefixed with "Bot"
	if d, err := discordgo.New("Bot " + discordToken); err != nil {
		log.Fatal("Failed to create a Discord session:", err)
		return
	} else {
		discord = d
	}

	// Register function handlers for various events
	discord.AddHandler(discordReady)
	discord.AddHandler(discordMessageCreate)

	// Register function handlers for all of the commands
	commandInit()

	// Open the websocket and begin listening
	if err := discord.Open(); err != nil {
		log.Fatal("Error opening Discord session: ", err)
		return
	}
}

/*
	Event handlers
*/

func discordReady(s *discordgo.Session, event *discordgo.Ready) {
	log.Info("Discord bot connected with username: " + event.User.Username)
	discordBotID = event.User.ID

	// Get the guild ID
	var guilds []*discordgo.UserGuild
	if v, err := s.UserGuilds(0, "", ""); err != nil {
		log.Fatal("Failed to get the Discord guilds:", err)
		return
	} else {
		guilds = v
	}

	foundGuild := false
	for _, guild := range guilds {
		if guild.Name == discordGuildName {
			foundGuild = true
			discordGuildID = guild.ID
			break
		}
	}
	if !foundGuild {
		log.Fatal("Failed to find the ID of the \"" + discordGuildName + "\" Discord server.")
	}

	// Get the ID of several roles
	var roles []*discordgo.Role
	if v, err := discord.GuildRoles(discordGuildID); err != nil {
		log.Fatal("Failed to get the roles for the guild: " + err.Error())
		return
	} else {
		roles = v
	}
	for _, role := range roles {
		if role.Name == discordAdminRoleName {
			discordAdminRoleID = role.ID
		} else if role.Name == discordBotRoleName {
			discordBotRoleID = role.ID
		} else if role.Name == discordCasterRoleName {
			discordCasterRoleID = role.ID
		} else if role.Name == discordTeamCaptainRoleName {
			discordTeamCaptainRoleID = role.ID
		} else if role.Name == "@everyone" {
			discordEveryoneRoleID = role.ID
		}
	}
	if discordAdminRoleID == "" {
		log.Fatal("Failed to find the role of \"" + discordAdminRoleName + "\".")
	} else if discordBotRoleID == "" {
		log.Fatal("Failed to find the role of \"" + discordBotRoleName + "\".")
	} else if discordCasterRoleID == "" {
		log.Fatal("Failed to find the role of \"" + discordCasterRoleName + "\".")
	} else if discordEveryoneRoleID == "" {
		log.Fatal("Failed to find the role of \"@everyone\".")
	}

	// Get the ID of the general channel
	var channels []*discordgo.Channel
	if v, err := discord.GuildChannels(discordGuildID); err != nil {
		log.Fatal("Failed to get the Discord server channels: " + err.Error())
	} else {
		channels = v
	}
	found := false
	for _, channel := range channels {
		if channel.Name == discordGeneralChannelName {
			found = true
			discordGeneralChannelID = channel.ID
			break
		}
	}
	if !found {
		log.Fatal("Failed to find the \"" + discordGeneralChannelName + "\" channel.")
	}
}

func discordMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Log the message
	var channelName string
	if v, err := discord.Channel(m.ChannelID); err != nil {
		log.Error("Failed to get the channel name for the channel ID of \""+m.ChannelID+"\":", err)
	} else {
		channelName = v.Name
	}
	log.Info("[#" + channelName + "] <" + m.Author.Username + "#" + m.Author.Discriminator + "> " + m.Content)

	// Ignore all messages created by the bot itself
	if m.Author.ID == discordBotID {
		return
	}

	// First, look for people mentioning the bot
	// (the second condition accounts for if we have a server nickname)
	if strings.Contains(m.Content, "<@"+discordBotID+">") || strings.Contains(m.Content, "<@!"+discordBotID+">") {
		discordSend(m.ChannelID, "ping me again\nI DARE YOU")
		return
	}

	// Second, look for exact greetings
	if strings.ToLower(m.Content) == "hello willy" {
		discordSend(m.ChannelID, "hello buttface, idgaf")
		return
	}

	// Commands for this bot will start with a "!", so we can ignore everything else
	args := strings.Split(m.Content, " ")
	command := args[0]
	args = args[1:] // This will be an empty slice if there is nothing after the command
	if !strings.HasPrefix(command, "!") {
		return
	}
	command = strings.TrimPrefix(command, "!")
	command = strings.ToLower(command) // Commands are case-insensitive

	// Check to see if there is a command handler for this command
	if _, ok := commandHandlerMap[command]; !ok {
		discordSend(m.ChannelID, "That is not a valid command.")
		return
	}

	commandMutex.Lock()
	commandHandlerMap[command](m, args)
	commandMutex.Unlock()
}

/*
	Miscellaneous functions
*/

func discordSend(channelID string, msg string) {
	if _, err := discord.ChannelMessageSend(channelID, msg); err != nil {
		log.Error("Failed to send \"" + msg + "\" to \"" + channelID + "\": " + err.Error())
		return
	}
}
