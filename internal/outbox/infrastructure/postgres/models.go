package postgres

import (
	"encoding/json"
	"time"

	outboxdomain "github.com/tbikbulatov/go-pulseops/internal/outbox/domain"
)

type EventModel struct {
	ID            string          `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	EventID       string          `gorm:"type:uuid;column:event_id;default:gen_random_uuid()"`
	AggregateType string          `gorm:"column:aggregate_type"`
	AggregateID   string          `gorm:"type:uuid;column:aggregate_id"`
	EventType     string          `gorm:"column:event_type"`
	Payload       json.RawMessage `gorm:"type:jsonb;column:payload"`
	Status        string
	Attempts      int
	NextAttemptAt time.Time `gorm:"column:next_attempt_at"`
	CreatedAt     time.Time `gorm:"column:created_at"`
	PublishedAt   *time.Time
}

func (EventModel) TableName() string {
	return "outbox_events"
}

func (m EventModel) ToDomain() outboxdomain.Event {
	return outboxdomain.Event{
		ID:            m.ID,
		EventID:       m.EventID,
		AggregateType: m.AggregateType,
		AggregateID:   m.AggregateID,
		EventType:     m.EventType,
		Payload:       m.Payload,
		Status:        m.Status,
		Attempts:      m.Attempts,
		NextAttemptAt: m.NextAttemptAt,
		CreatedAt:     m.CreatedAt,
		PublishedAt:   m.PublishedAt,
	}
}
