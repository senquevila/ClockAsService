package data

import (
	"testing"
	"time"
	"github.com/google/uuid"
)

func TestAlarmCreation(t *testing.T) {
	id := uuid.New().String()
	alarm := Alarm{
		ID:          id,
		Name:        "Test Alarm",
		Description: "This is a test alarm.",
		Target:      time.Now().Add(1 * time.Hour),
		CreatedAt:   time.Now(),
	}
	if alarm.ID != id {
		t.Errorf("expected ID '%s', got %s", id, alarm.ID)
	}
	if alarm.Name != "Test Alarm" {
		t.Errorf("expected Name 'Test Alarm', got %s", alarm.Name)
	}
	if alarm.Description != "This is a test alarm." {
		t.Errorf("expected Description 'This is a test alarm.', got %s", alarm.Description)
	}
}
