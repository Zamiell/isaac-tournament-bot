package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"syscall"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	logging "github.com/op/go-logging"
)

var (
	projectPath string
	log         *logging.Logger
	modals      *Models
	builds      = make([]Build, 0)
)

func main() {
	// Initialize logging:
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

	// Get the project path:
	// https://stackoverflow.com/questions/18537257/
	if v, err := os.Executable(); err != nil {
		log.Fatal("Failed to get the path of the currently running executable:", err)
	} else {
		projectPath = filepath.Dir(v)
	}

	// Load the ".env" file which contains environment variables with secret values.
	if err := godotenv.Load(path.Join(projectPath, ".env")); err != nil {
		log.Fatal("Failed to load .env file:", err)
	}

	// Initialize the database model.
	if v, err := NewModels(); err != nil {
		log.Fatal("Failed to open the database:", err)
	} else {
		modals = v
	}
	defer modals.Close()

	// Initialize the other parts of the program.
	loadAllBuilds()
	discordInit()
	defer discordSession.Close()
	challongeInit()
	matchInit()
	languageInit()
	log.Info("The bot has successfully initialized.")

	// Wait here until CTRL-C or other term signal is received.
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}

func loadAllBuilds() {
	libPath := getLibraryPath()
	jsonFilePath := path.Join(libPath, "builds.json")
	var jsonFile []byte
	if v, err := ioutil.ReadFile(jsonFilePath); err != nil {
		log.Fatal("Failed to open \""+jsonFilePath+"\":", err)
	} else {
		jsonFile = v
	}

	if err := json.Unmarshal(jsonFile, &builds); err != nil {
		log.Fatal("Failed to unmarshal the builds:", err)
	}
}

func getLibraryPath() string {
	libPath := path.Join(projectPath, "lib", "node_modules", "isaac-racing-common", "src")
	if _, err := os.Stat(libPath); os.IsNotExist(err) {
		log.Fatal("The library path at \"" + libPath + "\" does not exist. Did you forget to run \"npm install\" in the \"lib\" subdirectory?")
	} else if err != nil {
		log.Fatal("Failed to check if the \""+libPath+"\" file exists:", err)
	}

	return libPath
}
