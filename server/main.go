package main

import (
	"flag"
	"os"
)

func main() {
	var dbURL string
	var serverURL string

	// Check environment variables.
	if envDbURL := os.Getenv("OBS_AL_DB_URL"); envDbURL != "" {
		dbURL = envDbURL
	}
	if envServerURL := os.Getenv("OBS_AL_SERVER_URL"); envServerURL != "" {
		serverURL = envServerURL
	}

	// Define flags.
	if dbURL == "" {
		flag.StringVar(&dbURL, "db-url", "events.db", "Path to an SQLITE DB.")
	}
	if serverURL == "" {
		flag.StringVar(&serverURL, "server-url", "0.0.0.0:8010", "URL at which to serve.")
	}
	flag.Parse()

	a := App{}
	a.Initialise(&dbURL)
	a.Run(&serverURL)
}
