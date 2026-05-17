package main

import (
	"github.com/tbikbulatov/go-pulseops/internal/platform/config"
	"github.com/tbikbulatov/go-pulseops/internal/platform/kafka"
)

func runKafkaPing(cfg *config.Config) error {
	producer, err := kafka.NewProducer(cfg.Kafka)
	if err != nil {
		return err
	}
	defer producer.Close()

	return producer.Publish(
		kafka.TopicAlertReceived,
		"kafka-ping",
		[]byte(`{"type":"kafka.ping","message":"hello redpanda"}`),
	)
}
