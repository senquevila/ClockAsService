package services

import (
	"database/sql"
	"testing"
	"time"

	datapkg "ClockAsService/src/data"

	_ "github.com/mattn/go-sqlite3"
)

func setupAlarmStorage(t *testing.T) *AlarmStorage {
	t.Helper()
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to open in-memory db: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	s := &AlarmStorage{DB: db}
	if err := s.CreateTable(); err != nil {
		t.Fatalf("failed to create alarms table: %v", err)
	}
	return s
}

func TestAlarmStorage_CreateListFindRemove(t *testing.T) {
	s := setupAlarmStorage(t)

	original := datapkg.Alarm{
		Name:        "Test Alarm Service",
		Description: "service test",
		Target:      time.Now().Add(2 * time.Hour),
	}

	if err := s.Create(original); err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	list, err := s.List()
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("expected 1 alarm after create, got %d", len(list))
	}

	got, ok := list[0].(datapkg.Alarm)
	if !ok {
		t.Fatalf("expected datapkg.Alarm from List, got %T", list[0])
	}

	if got.Name != original.Name {
		t.Errorf("expected Name %q, got %q", original.Name, got.Name)
	}
	if got.Description != original.Description {
		t.Errorf("expected Description %q, got %q", original.Description, got.Description)
	}
	if got.Target.Unix() != original.Target.Unix() {
		t.Errorf("expected Target Unix %d, got %d", original.Target.Unix(), got.Target.Unix())
	}
	if got.ID == "" {
		t.Errorf("expected generated ID, got empty string")
	}

	// FindByID
	foundRaw, err := s.FindByID(got.ID)
	if err != nil {
		t.Fatalf("FindByID failed: %v", err)
	}
	found, ok := foundRaw.(datapkg.Alarm)
	if !ok {
		t.Fatalf("expected datapkg.Alarm from FindByID, got %T", foundRaw)
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
		t.Fatalf("expected 0 alarms after remove, got %d", len(listAfter))
	}
}

func TestAlarmStorage_Create_WrongType(t *testing.T) {
	s := setupAlarmStorage(t)
	if err := s.Create("not an alarm"); err == nil {
		t.Fatalf("expected error when creating with wrong type, got nil")
	}
}
