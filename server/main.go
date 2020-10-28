package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	// Check environment variables.
	if envDbURL := os.Getenv("OBS_AL_DB_URL"); envDbURL != "" {
		absDbURL, err := filepath.Abs(envDbURL)
		if err != nil {
			log.Fatal(err)
		}
		appConfig.dbURL = absDbURL
	}
	if envServerURL := os.Getenv("OBS_AL_SERVER_URL"); envServerURL != "" {
		appConfig.serverURL = envServerURL
	}
	// If environment is not specified, default is development.
	if environment := os.Getenv("OBS_AL_ENV"); environment != "" {
		appConfig.environment = environment
	}

	// Define flags.
	if appConfig.dbURL == "" {
		var flagDbURL string
		flag.StringVar(&flagDbURL, "db-url", "events.db", "Path to an SQLITE DB.")
		absDbURL, err := filepath.Abs(flagDbURL)
		if err != nil {
			log.Fatal(err)
		}
		appConfig.dbURL = absDbURL
	}
	if appConfig.serverURL == "" {
		flag.StringVar(&appConfig.serverURL, "server-url", "localhost:8010", "URL at which to serve.")
	}
	flag.Parse()

	a := App{}
	a.Initialise()
	a.Run()
}
