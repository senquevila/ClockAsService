package services

import (
	"database/sql"
	"testing"
	"time"

	datapkg "ClockAsService/src/data"

	_ "github.com/mattn/go-sqlite3"
)

func setupEventStorage(t *testing.T) *EventStorage {
	t.Helper()
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to open in-memory db: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	s := &EventStorage{DB: db}
	if err := s.CreateTable(); err != nil {
		t.Fatalf("failed to create events table: %v", err)
	}
	return s
}

func TestEventStorage_CreateListFindRemove(t *testing.T) {
	s := setupEventStorage(t)

	original := datapkg.Event{
		Name:        "Test Event Service",
		Description: "service test",
		StartedAt:   time.Now().Add(-30 * time.Minute),
	}

	if err := s.Create(original); err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	list, err := s.List()
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("expected 1 event after create, got %d", len(list))
	}

	got, ok := list[0].(datapkg.Event)
	if !ok {
		t.Fatalf("expected datapkg.Event from List, got %T", list[0])
	}

	if got.Name != original.Name {
		t.Errorf("expected Name %q, got %q", original.Name, got.Name)
	}
	if got.Description != original.Description {
		t.Errorf("expected Description %q, got %q", original.Description, got.Description)
	}
	if got.StartedAt.Unix() != original.StartedAt.Unix() {
		t.Errorf("expected StartedAt Unix %d, got %d", original.StartedAt.Unix(), got.StartedAt.Unix())
	}
	if got.ID == "" {
		t.Errorf("expected generated ID, got empty string")
	}

	// FindByID
	foundRaw, err := s.FindByID(got.ID)
	if err != nil {
		t.Fatalf("FindByID failed: %v", err)
	}
	found, ok := foundRaw.(datapkg.Event)
	if !ok {
		t.Fatalf("expected datapkg.Event from FindByID, got %T", foundRaw)
	}
	if found.ID != got.ID {
		t.Errorf("expected ID %s, got %s", got.ID, found.ID)
	}

	// Remove
	if err := s.Remove(got.ID); err != nil {
		t.Fatalf("Remove failed: %v", err)
	}
	listAfter, err := s.List()
	if err != nil {
		t.Fatalf("List after remove failed: %v", err)
	}
	if len(listAfter) != 0 {
		t.Fatalf("expected 0 events after remove, got %d", len(listAfter))
	}
}

func TestEventStorage_Create_WrongType(t *testing.T) {
	s := setupEventStorage(t)
	if err := s.Create(12345); err == nil {
		t.Fatalf("expected error when creating with wrong type, got nil")
	}
}
