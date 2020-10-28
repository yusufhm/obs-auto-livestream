package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/unrolled/secure"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// App contains the Router & DB and also contains methods.
type App struct {
	Router *mux.Router
	DB     *gorm.DB
}

var appConfig = struct {
	dbURL       string
	serverURL   string
	environment string
	isDev       func(string) bool
}{
	environment: "development",
	isDev: func(env string) bool {
		return env == "development"
	},
}

// Initialise sets up the DB and router.
func (a *App) Initialise() {
	log.Printf("%+v\n", appConfig)
	log.Printf("Using database '%v'", appConfig.dbURL)

	var err error
	a.DB, err = gorm.Open(sqlite.Open(appConfig.dbURL), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	secureMiddleware := secure.New(secure.Options{
		FrameDeny:     true,
		IsDevelopment: appConfig.isDev(appConfig.environment),
	})
	a.Router = mux.NewRouter()
	a.Router.Use(secureMiddleware.Handler)
	a.initialiseRoutes()
}

// Run calls the main server loop.
func (a *App) Run() {
	log.Printf("Starting server at 'http://%v'", appConfig.serverURL)
	log.Fatal(http.ListenAndServe(appConfig.serverURL, a.Router))
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func (a *App) getHome(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (a *App) getEvent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid event ID")
		return
	}

	ev := pageEvent{ID: uint(id)}
	if err := ev.getEvent(a.DB); err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			respondWithError(w, http.StatusNotFound, "Event not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, ev)
}

func (a *App) getEvents(w http.ResponseWriter, r *http.Request) {
	count, _ := strconv.Atoi(r.FormValue("count"))
	start, _ := strconv.Atoi(r.FormValue("start"))

	if count > 10 || count < 1 {
		count = 10
	}
	if start < 0 {
		start = 0
	}

	ev := pageEvent{}
	events, err := ev.getEvents(a.DB, start, count)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, events)
}

func (a *App) facebookWebhook(w http.ResponseWriter, r *http.Request) {
	var fbEv fbPageEvent
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&fbEv); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	var ev pageEvent
	if err := ev.createFromFbEvent(a.DB, &fbEv); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, ev)
}

func (a *App) deleteEvent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Event ID")
		return
	}

	ev := pageEvent{ID: uint(id)}
	if err := ev.deleteEvent(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (a *App) initialiseRoutes() {
	a.Router.HandleFunc("/", a.getHome).Methods("GET")
	a.Router.HandleFunc("/events", a.getEvents).Methods("GET")
	a.Router.HandleFunc("/facebook-webhook", a.facebookWebhook).Methods("POST")
	a.Router.HandleFunc("/event/{id:[0-9]+}", a.getEvent).Methods("GET")
	a.Router.HandleFunc("/event/{id:[0-9]+}", a.deleteEvent).Methods("DELETE")
}
