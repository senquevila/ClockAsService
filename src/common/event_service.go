package common

import (
	"database/sql"
	"time"
	"github.com/google/uuid"
)

type EventStorage struct {
	DB *sql.DB
}

func (e *EventStorage) CreateTable() error {
	eventTable := `CREATE TABLE IF NOT EXISTS events (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		description TEXT NOT NULL,
		started_at INTEGER NOT NULL,
		created_at INTEGER NOT NULL
	);`
	_, err := e.DB.Exec(eventTable)
	return err
}

func (e *EventStorage) Create(data interface{}) error {
	event, ok := data.(Event)
	if !ok {
		return sql.ErrConnDone
	}

	_, err := e.DB.Exec(
		"INSERT OR REPLACE INTO events (id, name, description, started_at, created_at) VALUES (?, ?, ?, ?, ?, ?)",
		uuid.New().String(), event.Name, event.Description, event.StartedAt.Unix(), time.Now().Unix(),
	)
	return err
}

func (e *EventStorage) Remove(id string) error {
	_, err := e.DB.Exec("DELETE FROM events WHERE id = ?", id)
	return err
}

func (e *EventStorage) List() ([]interface{}, error) {
	rows, err := e.DB.Query("SELECT id, name, description, started_at FROM events")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var events []interface{}
	for rows.Next() {
		var event Event
		var startedUnix int64
		if err := rows.Scan(&event.ID, &event.Name, &event.Description, &startedUnix); err != nil {
			return nil, err
		}
		event.StartedAt = time.Unix(startedUnix, 0)
		events = append(events, event)
	}
	return events, nil
}

func (e *EventStorage) FindByID(id string) (interface{}, error) {
	row := e.DB.QueryRow("SELECT id, name, description, started_at FROM events WHERE id = ?", id)
	var event Event
	var startedUnix int64
	if err := row.Scan(&event.ID, &event.Name, &event.Description, &startedUnix); err != nil {
		return nil, err
	}
	event.StartedAt = time.Unix(startedUnix, 0)
	return event, nil
}
