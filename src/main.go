package main

import (
	"os"
	"os/signal"
	"path"
	"syscall"

	"github.com/Zamiell/isaac-tournament-bot/src/models"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	logging "github.com/op/go-logging"
)

var (
	projectPath = path.Join(os.Getenv("GOPATH"), "src", "github.com", "Zamiell", "isaac-tournament-bot")
	log         *logging.Logger
	db          *models.Models
)

func main() {
	// Initialize logging
	// http://godoc.org/github.com/op/go-logging#Formatter
	log = logging.MustGetLogger("isaac-tournament-bot")
	loggingBackend := logging.NewLogBackend(os.Stdout, "", 0)
	logFormat := logging.MustStringFormatter( // https://golang.org/pkg/time/#Time.Format
		`%{time:Mon Jan 2 15:04:05 MST 2006} - %{level:.4s} - %{shortfile} - %{message}`,
	)
	loggingBackendFormatted := logging.NewBackendFormatter(loggingBackend, logFormat)
	logging.SetBackend(loggingBackendFormatted)

	// Welcome message
	log.Info("+--------------------------------+")
	log.Info("| Starting isaac-tournament-bot. |")
	log.Info("+--------------------------------+")

	// Load the ".env" file which contains environment variables with secret values
	if err := godotenv.Load(path.Join(projectPath, ".env")); err != nil {
		log.Fatal("Failed to load .env file:", err)
	}

	// Initialize the database model
	if v, err := models.Init(); err != nil {
		log.Fatal("Failed to open the database:", err)
	} else {
		db = v
	}
	defer db.Close()

	// Initialize the other parts of the program
	discordInit()
	defer discord.Close()
	challongeInit()
	matchInit()
	languageInit()
	log.Info("The bot has successfully initialized.")

	// Wait here until CTRL-C or other term signal is received
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}
