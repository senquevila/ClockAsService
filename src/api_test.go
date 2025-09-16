package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	datapkg "ClockAsService/src/data"
	"ClockAsService/src/services"

	_ "github.com/mattn/go-sqlite3"
)

func setupHandlersForTest(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to open db: %v", err)
	}
	// ensure db closed at end of test
	t.Cleanup(func() { db.Close() })
	// caller test will not run concurrently so global stores are ok
	if err := db.Ping(); err != nil {
		t.Fatalf("db ping failed: %v", err)
	}

	alarmStore = &services.AlarmStorage{DB: db}
	eventStore = &services.EventStorage{DB: db}
	if err := alarmStore.CreateTable(); err != nil {
		t.Fatalf("CreateTable alarm failed: %v", err)
	}
	if err := eventStore.CreateTable(); err != nil {
		t.Fatalf("CreateTable event failed: %v", err)
	}
}

func TestCreateAlarm_RejectsPastTarget(t *testing.T) {
	setupHandlersForTest(t)

	payload := map[string]interface{}{
		"name":        "past",
		"description": "past",
		"target":      time.Now().Add(-time.Hour).Format(time.RFC3339),
	}
	raw, _ := json.Marshal(payload)
	req := httptest.NewRequest("POST", "/alarms/create", bytes.NewReader(raw))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	createAlarmHandler(w, req)
	resp := w.Result()
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
	var body map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode body: %v", err)
	}
	if _, ok := body["error"]; !ok {
		t.Fatalf("expected error field in response")
	}
}

func TestGetAlarmCountdown_ClampsToZero(t *testing.T) {
	setupHandlersForTest(t)

	// create an alarm directly in storage with target in the past
	alarm := datapkg.Alarm{
		Name:        "past",
		Description: "past",
		Target:      time.Now().Add(-time.Hour),
	}
	createdRaw, err := alarmStore.Create(alarm)
	if err != nil {
		t.Fatalf("failed to create alarm in storage: %v", err)
	}
	created := createdRaw.(datapkg.Alarm)

	req := httptest.NewRequest("GET", "/alarms/countdown?id="+created.ID, nil)
	w := httptest.NewRecorder()

	getAlarmCountdownHandler(w, req)
	resp := w.Result()
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	var body map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode body: %v", err)
	}
	if body["countdown"].(float64) != 0 {
		t.Fatalf("expected countdown 0, got %v", body["countdown"])
	}
}
