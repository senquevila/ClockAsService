package services

import (
	datapkg "ClockAsService/src/data"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type AlarmStorage struct {
	DB *sql.DB
}

func (a *AlarmStorage) CreateTable() error {
	alarmTable := `CREATE TABLE IF NOT EXISTS alarms (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		description TEXT NOT NULL,
		target INTEGER NOT NULL,
		created_at INTEGER NOT NULL
	);`
	_, err := a.DB.Exec(alarmTable)
	return err
}

func (a *AlarmStorage) Create(raw interface{}) (interface{}, error) {
	alarm, ok := raw.(datapkg.Alarm)
	if !ok {
		return nil, sql.ErrConnDone
	}
	id := uuid.New().String()
	created := time.Now()
	_, err := a.DB.Exec(
		"INSERT OR REPLACE INTO alarms (id, name, description, target, created_at) VALUES (?, ?, ?, ?, ?)",
		id, alarm.Name, alarm.Description, alarm.Target.Unix(), created.Unix(),
	)
	if err != nil {
		return nil, err
	}
	alarm.ID = id
	alarm.CreatedAt = created
	return alarm, nil
}

func (a *AlarmStorage) Remove(id string) error {
	_, err := a.DB.Exec("DELETE FROM alarms WHERE id = ?", id)
	return err
}

func (a *AlarmStorage) List() ([]interface{}, error) {
	rows, err := a.DB.Query("SELECT id, name, description, target, created_at FROM alarms")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var alarms []interface{}
	for rows.Next() {
		var alarm datapkg.Alarm
		var targetUnix, createdUnix int64
		if err := rows.Scan(&alarm.ID, &alarm.Name, &alarm.Description, &targetUnix, &createdUnix); err != nil {
			return nil, err
		}
		alarm.Target = time.Unix(targetUnix, 0)
		alarm.CreatedAt = time.Unix(createdUnix, 0)
		alarms = append(alarms, alarm)
	}
	return alarms, nil
}

func (a *AlarmStorage) FindByID(id string) (interface{}, error) {
	row := a.DB.QueryRow("SELECT id, name, description, target, created_at FROM alarms WHERE id = ?", id)
	var alarm datapkg.Alarm
	var targetUnix, createdUnix int64
	if err := row.Scan(&alarm.ID, &alarm.Name, &alarm.Description, &targetUnix, &createdUnix); err != nil {
		return nil, err
	}
	alarm.Target = time.Unix(targetUnix, 0)
	alarm.CreatedAt = time.Unix(createdUnix, 0)
	return alarm, nil
}
