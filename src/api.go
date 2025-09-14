package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	datapkg "ClockAsService/src/data"
	"ClockAsService/src/services"

	_ "github.com/mattn/go-sqlite3"
)

// request shapes
type AlarmRequest struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Target      time.Time `json:"target"`
}

type EventRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

var alarmStore *services.AlarmStorage
var eventStore *services.EventStorage

func createAlarmHandler(w http.ResponseWriter, r *http.Request) {
	var req AlarmRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	alarm := datapkg.Alarm{
		Name:        req.Name,
		Description: req.Description,
		Target:      req.Target,
	}
	if err := alarmStore.Create(alarm); err != nil {
		http.Error(w, "Failed to create alarm", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func getAlarmCountdownHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	raw, err := alarmStore.FindByID(id)
	if err != nil {
		http.Error(w, "Alarm not found", http.StatusNotFound)
		return
	}
	alarm, ok := raw.(datapkg.Alarm)
	if !ok {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}
	countdown := time.Until(alarm.Target)
	json.NewEncoder(w).Encode(map[string]interface{}{"id": id, "countdown": countdown.Seconds()})
}

func createEventHandler(w http.ResponseWriter, r *http.Request) {
	var req EventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	event := datapkg.Event{
		Name:        req.Name,
		Description: req.Description,
		StartedAt:   time.Now(),
	}
	if err := eventStore.Create(event); err != nil {
		http.Error(w, "Failed to create event", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func getEventElapsedHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	raw, err := eventStore.FindByID(id)
	if err != nil {
		http.Error(w, "Event not found", http.StatusNotFound)
		return
	}
	event, ok := raw.(datapkg.Event)
	if !ok {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}
	elapsed := time.Since(event.StartedAt)
	json.NewEncoder(w).Encode(map[string]interface{}{"id": id, "elapsed": elapsed.Seconds()})
}

func main() {
	db, err := sql.Open("sqlite3", "clock.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	alarmStore = &services.AlarmStorage{DB: db}
	eventStore = &services.EventStorage{DB: db}

	if err := alarmStore.CreateTable(); err != nil {
		panic(err)
	}
	if err := eventStore.CreateTable(); err != nil {
		panic(err)
	}

	http.HandleFunc("/alarms/create", createAlarmHandler)
	http.HandleFunc("/alarms/countdown", getAlarmCountdownHandler)
	http.HandleFunc("/events/create", createEventHandler)
	http.HandleFunc("/events/elapsed", getEventElapsedHandler)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
