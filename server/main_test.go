package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var a App

func TestMain(m *testing.M) {
	appConfig.dbURL = "test.db"
	a.Initialise()

	ensureTableExists()
	code := m.Run()
	clearTable()
	os.Exit(code)
}

func ensureTableExists() {
	if err := a.DB.AutoMigrate(&pageEvent{}); err != nil {
		log.Fatal(err)
	}
}

func clearTable() {
	a.DB.Exec("DELETE FROM page_events")
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)
	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func addEvents(count int) {
	if count < 1 {
		count = 1
	}

	for i := 0; i < count; i++ {
		a.DB.Create(&pageEvent{
			PageID:          124124,
			LiveVideoID:     9798745,
			LiveVideoStatus: "preview",
		})
	}
}

func TestEmptyTable(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/events", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body != "[]" {
		t.Errorf("Expected an empty array. Got %s", body)
	}
}

func TestGetNonExistentEvent(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/event/11", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "Event not found" {
		t.Errorf("Expected the 'error' key of the response to be set to 'Event not found'. Got '%s'", m["error"])
	}
}

func TestFacebookWebhook(t *testing.T) {
	clearTable()

	var jsonStr = []byte(`{
		"object": "page",
		"entry": [
			{
				"id": "109768413992690",
				"time": 1602839070,
				"changes": [
					{
						"value": {
							"id": "220578652911665",
							"status": "vod"
						},
						"field": "live_videos"
					}
				]
			}
		]
	}`)
	req, _ := http.NewRequest("POST", "/facebook-webhook", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	var event pageEvent
	json.Unmarshal(response.Body.Bytes(), &event)

	if event.PageID != 109768413992690 {
		t.Errorf("Expected event page ID to be '109768413992690'. Got '%v'", event.PageID)
	}

	if event.LiveVideoID != 220578652911665 {
		t.Errorf("Expected event entry Live Video ID to be '220578652911665'. Got '%v'", event.LiveVideoID)
	}

	if event.LiveVideoStatus != "vod" {
		t.Errorf("Expected event status to be 'vod'. Got '%v'", event.LiveVideoStatus)
	}
}

func TestGetProduct(t *testing.T) {
	clearTable()
	addEvents(1)

	req, _ := http.NewRequest("GET", "/event/1", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestDeleteProduct(t *testing.T) {
	clearTable()
	addEvents(1)

	req, _ := http.NewRequest("GET", "/event/1", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("DELETE", "/event/1", nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("GET", "/event/1", nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)
}
