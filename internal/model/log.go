package model

import "time"

// Log represents a development journal entry
type Log struct {
	ID            int64
	CreatedAt     time.Time
	TaskName      string
	EstimateHours float64
	ActualHours   float64
	AIMinutes     *int
	Problem       *string
	Solution      *string
	Learning      *string
}

// LogInput represents user input for creating a new log entry
type LogInput struct {
	TaskName      string
	EstimateHours float64
	ActualHours   float64
	AIMinutes     int
	Problem       string
	Solution      string
	Learning      string
}
