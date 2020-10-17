package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/getlantern/systray"
	"github.com/getlantern/systray/example/icon"
	"github.com/joho/godotenv"
	"github.com/yusufhm/obs-auto-livestream/common/facebook"
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

	systray.Run(onReady, nil)
}

func onReady() {
	systray.SetIcon(icon.Data)
	systray.SetTitle("OBS Auto Livestream")
	systray.SetTooltip("OBS Auto Livestream")
	mLogin := systray.AddMenuItem("Login", "Login to Facebook")
	mQuit := systray.AddMenuItem("Quit", "Quit the whole app")

	go func() {
		for {
			select {
			case <-mLogin.ClickedCh:
				openLogin()
			case <-mQuit.ClickedCh:
				systray.Quit()
			}
		}
	}()
}

func openLogin() {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)
	fmt.Println(exPath)

	cmd := exec.Command(exPath+"/gui", "login")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	if err := json.NewDecoder(stdout).Decode(&facebook.UserAccessToken); err != nil {
		log.Fatal(err)
	}
	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}
	fmt.Println(facebook.GetPages())
}
