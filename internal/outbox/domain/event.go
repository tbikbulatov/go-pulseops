package domain

import (
	"encoding/json"
	"time"
)

const (
	StatusPublished = "published"
	StatusPending   = "pending"
	StatusFailed    = "failed"
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

func (e *Event) MarkPublished(at time.Time) {
	e.Status = StatusPublished
	e.PublishedAt = &at
}
