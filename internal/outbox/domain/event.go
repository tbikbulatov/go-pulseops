package domain

import (
	"encoding/json"
	"time"
)

type Event struct {
	ID            string
	EventID       string
	AggregateType string
	AggregateID   string
	EventType     string
	Payload       json.RawMessage
	Status        string
	Attempts      int
	NextAttemptAt time.Time
	CreatedAt     time.Time
	PublishedAt   *time.Time
}
