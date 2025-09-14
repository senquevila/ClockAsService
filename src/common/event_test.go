package common

import (
	"testing"
	"github.com/google/uuid"
)

func TestEventCreation(t *testing.T) {
	id := uuid.New().String()
	event := Event{
		ID:          id,
		Name:        "Test Event",
		Description: "This is a test event.",
	}
	if event.ID != id {
		t.Errorf("expected ID '%s', got %s", id, event.ID)
	}
	if event.Name != "Test Event" {
		t.Errorf("expected Name 'Test Event', got %s", event.Name)
	}
	if event.Description != "This is a test event." {
		t.Errorf("expected Description 'This is a test event.', got %s", event.Description)
	}
}
