package main

import (
	"log"
	"os"

	"github.com/yusufhm/obs-auto-livestream/common/facebook"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file", err)
	}

	// Ensure the Facebook APP ID is available.
	facebook.AppID = os.Getenv("OBS_AUTO_LIVESTREAM_FACEBOOK_APP_ID")
	if facebook.AppID == "" {
		log.Fatal("Facebook App ID not found")
	}

	if len(os.Args) > 1 && os.Args[1] == "login" {
		facebook.Login()
	}
}
