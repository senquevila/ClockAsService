package common

import "time"

// Event represents a timer showing elapsed time since creation
type Event struct {
	ID          string
	Name        string
	Description string
	StartedAt   time.Time
}
