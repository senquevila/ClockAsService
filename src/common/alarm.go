package common

import "time"

// Alarm represents a countdown to a target time
type Alarm struct {
	ID          string
	Name        string
	Description string
	Target      time.Time
	CreatedAt   time.Time
}
