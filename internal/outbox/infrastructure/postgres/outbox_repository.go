package postgres

import (
	"context"

	outboxdomain "github.com/tbikbulatov/go-pulseops/internal/outbox/domain"
	"github.com/tbikbulatov/go-pulseops/internal/platform/transaction"
	"gorm.io/gorm"
)

type OutboxRepository struct {
	db *gorm.DB
}

func NewOutboxRepository(db *gorm.DB) *OutboxRepository {
	return &OutboxRepository{db: db}
}

func (r *OutboxRepository) Create(ctx context.Context, event outboxdomain.Event) error {
	model := EventModel{
		ID:            event.ID,
		EventID:       event.EventID,
		AggregateType: event.AggregateType,
		AggregateID:   event.AggregateID,
		EventType:     event.EventType,
		Payload:       event.Payload,
		Status:        event.Status,
		Attempts:      event.Attempts,
		NextAttemptAt: event.NextAttemptAt,
		CreatedAt:     event.CreatedAt,
		PublishedAt:   event.PublishedAt,
	}

	db := transaction.FromContext(ctx, r.db)
	return db.Create(&model).Error
}

func (r *OutboxRepository) FindPending(ctx context.Context, limit int) ([]outboxdomain.Event, error) {
	var models []EventModel

	db := transaction.FromContext(ctx, r.db)
	err := db.Where("status = ? AND next_attempt_at <= now()", outboxdomain.StatusPending).
		Order("created_at ASC").
		Limit(limit).
		Find(&models).
		Error

	if err != nil {
		return nil, err
	}

	events := make([]outboxdomain.Event, 0, len(models))
	for _, model := range models {
		events = append(events, model.ToDomain())
	}

	return events, nil
}

func (r *OutboxRepository) Save(ctx context.Context, event outboxdomain.Event) error {
	model := EventModel{
		ID:            event.ID,
		EventID:       event.EventID,
		AggregateType: event.AggregateType,
		AggregateID:   event.AggregateID,
		EventType:     event.EventType,
		Payload:       event.Payload,
		Status:        event.Status,
		Attempts:      event.Attempts,
		NextAttemptAt: event.NextAttemptAt,
		CreatedAt:     event.CreatedAt,
		PublishedAt:   event.PublishedAt,
	}

	db := transaction.FromContext(ctx, r.db)
	return db.Save(&model).Error
}
