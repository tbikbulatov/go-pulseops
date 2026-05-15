package main

import (
	"context"
	"fmt"
	"time"

	"github.com/tbikbulatov/go-pulseops/internal/outbox/infrastructure/postgres"
	"github.com/tbikbulatov/go-pulseops/internal/platform/config"
	"github.com/tbikbulatov/go-pulseops/internal/platform/db"
	"github.com/tbikbulatov/go-pulseops/internal/platform/kafka"
)

func publishOutbox(cfg *config.Config) error {
	ctx := context.Background()

	producer, err := kafka.NewProducer(cfg.Kafka)
	if err != nil {
		return err
	}
	defer producer.Close()

	gorm, err := db.NewGorm(cfg.Postgres)
	if err != nil {
		return fmt.Errorf("init gorm: %w", err)
	}
	repo := postgres.NewOutboxRepository(gorm)

	events, err := repo.FindPending(ctx, 10)
	if err != nil {
		return fmt.Errorf("fetch pending outbox events: %w", err)
	}

	for _, event := range events {
		if err := producer.Publish(event.EventType, event.AggregateID, event.Payload); err != nil {
			return fmt.Errorf("publish outbox event %s: %w", event.ID, err)
		}

		event.MarkPublished(time.Now())

		if err := repo.Save(ctx, event); err != nil {
			return fmt.Errorf("save published outbox event %s: %w", event.ID, err)
		}
	}

	return nil
}
