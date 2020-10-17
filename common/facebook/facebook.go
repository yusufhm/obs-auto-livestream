package facebook

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"

	fb "github.com/huandu/facebook/v2"
	"github.com/webview/webview"
)

type accessToken struct {
	Token                    string
	DataAccessExpirationTime int
}

var UserAccessToken accessToken
var AppID string

func Login() {
	loginURL := fmt.Sprintf(
		"https://www.facebook.com/v8.0/dialog/oauth?client_id=%v&redirect_uri=%v&state=%v&response_type=token",
		AppID,
		url.QueryEscape("https://www.facebook.com/connect/login_success.html"),
		url.QueryEscape("\"{st=state123abc,ds=123456789}\""),
	)

	w := webview.New(true)
	defer w.Destroy()
	w.SetTitle("Facebook Login")

	w.Bind("pageUrl", func(pageUrl string) {
		u, err := url.Parse(pageUrl)
		if err != nil {
			log.Fatal(err)
		}

		fragments, _ := url.ParseQuery(u.Fragment)
		if u.Hostname() == "www.facebook.com" &&
			u.Path == "/connect/login_success.html" &&
			fragments["access_token"] != nil {
			UserAccessToken.Token = fragments["access_token"][0]
			daet, err := strconv.Atoi(fragments["data_access_expiration_time"][0])
			if err != nil {
				log.Fatal(err)
			}
			UserAccessToken.DataAccessExpirationTime = daet
			uatJSON, err := json.Marshal(UserAccessToken)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(string(uatJSON))
			w.Terminate()
		}
	})

	w.Init("pageUrl(window.location.href)")

	w.SetSize(800, 600, webview.HintNone)
	w.Navigate(loginURL)
	w.Run()
}

func GetUserPermissions() {
	res, err := fb.Get("/me/permissions", fb.Params{
		"access_token": UserAccessToken.Token,
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	jsonBytes, _ := json.MarshalIndent(res["data"], "", "    ")
	fmt.Println("User permissions:", string(jsonBytes))
}

func GetPages() []int {
	res, err := fb.Get("/me/accounts", fb.Params{
		"access_token": UserAccessToken.Token,
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	var pageIDs []int
	for _, page := range res["data"].([]interface{}) {
		pageID, err := strconv.Atoi(page.(map[string]interface{})["id"].(string))
		if err != nil {
			log.Fatal(err)
		}
		pageIDs = append(pageIDs, pageID)
	}
	return pageIDs
}

func GetLiveVideo(videoID string) {
	res, err := fb.Get(videoID, fb.Params{
		"access_token": UserAccessToken.Token,
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	jsonBytes, _ := json.MarshalIndent(res["data"], "", "    ")
	fmt.Println("Live video:", string(jsonBytes))
}
