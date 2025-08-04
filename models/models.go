package models

import "time"

type Todo struct {
	ID    int      `json:"id"`
	Text  string   `json:"text"`
	Done  bool     `json:"done"`
	Notes []string `json:"notes"`
}

type Config struct {
	LastInput  string `json:"last_input"`
	RunningPID int    `json:"running_pid"`
}

type LogEntry struct {
	Timestamp time.Time `json:"timestamp"`
	EventType string    `json:"event_type"`
	TodoID    int       `json:"todo_id"`
	TodoData  Todo      `json:"todo_data"`
}
