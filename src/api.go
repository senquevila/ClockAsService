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
	createdRaw, err := alarmStore.Create(alarm)
	if err != nil {
		http.Error(w, "Failed to create alarm", http.StatusInternalServerError)
		return
	}
	created, ok := createdRaw.(datapkg.Alarm)
	if !ok {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(created)
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
	w.Header().Set("Content-Type", "application/json")
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
	createdRaw, err := eventStore.Create(event)
	if err != nil {
		http.Error(w, "Failed to create event", http.StatusInternalServerError)
		return
	}
	created, ok := createdRaw.(datapkg.Event)
	if !ok {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(created)
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
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"id": id, "elapsed": elapsed.Seconds()})
}

func listAlarmsHandler(w http.ResponseWriter, r *http.Request) {
	raws, err := alarmStore.List()
	if err != nil {
		http.Error(w, "Failed to list alarms", http.StatusInternalServerError)
		return
	}
	var alarms []datapkg.Alarm
	for _, raw := range raws {
		if a, ok := raw.(datapkg.Alarm); ok {
			alarms = append(alarms, a)
		}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(alarms)
}

func listEventsHandler(w http.ResponseWriter, r *http.Request) {
	raws, err := eventStore.List()
	if err != nil {
		http.Error(w, "Failed to list events", http.StatusInternalServerError)
		return
	}
	var events []datapkg.Event
	for _, raw := range raws {
		if e, ok := raw.(datapkg.Event); ok {
			events = append(events, e)
		}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(events)
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
	http.HandleFunc("/alarms/list", listAlarmsHandler)
	http.HandleFunc("/events/create", createEventHandler)
	http.HandleFunc("/events/elapsed", getEventElapsedHandler)
	http.HandleFunc("/events/list", listEventsHandler)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
