package kafka

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/tbikbulatov/go-pulseops/internal/incident/usecase/processalert"
)

type AlertReceivedEvent struct {
	EventID       string `json:"event_id"`
	AlertID       string `json:"alert_id"`
	IntegrationID string `json:"integration_id"`
	ExternalID    string `json:"external_id"`
	Service       string `json:"service"`
	Environment   string `json:"environment"`
	Severity      string `json:"severity"`
	Name          string `json:"name"`
	Message       string `json:"message"`
	DedupKey      string `json:"dedup_key"`
}

type AlertReceivedProcessor interface {
	Handle(ctx context.Context, cmd processalert.Command) error
}

type AlertReceivedHandler struct {
	processor AlertReceivedProcessor
}

func NewAlertReceivedHandler(processor AlertReceivedProcessor) AlertReceivedHandler {
	return AlertReceivedHandler{
		processor: processor,
	}
}

func (h *AlertReceivedHandler) Handle(ctx context.Context, msg *sarama.ConsumerMessage) error {
	var event AlertReceivedEvent
	if err := json.Unmarshal(msg.Value, &event); err != nil {
		return fmt.Errorf("decode alert received event: %w", err)
	}

	messageID, err := event.MessageID()
	if err != nil {
		return err
	}

	return h.processor.Handle(ctx, event.ToCommand(messageID))
}

func (e AlertReceivedEvent) MessageID() (string, error) {
	if e.EventID != "" {
		return e.EventID, nil
	}
	if e.AlertID != "" {
		return e.AlertID, nil
	}

	return "", fmt.Errorf("alert received event has no event_id or alert_id")
}

func (e AlertReceivedEvent) ToCommand(messageID string) processalert.Command {
	return processalert.Command{
		MessageID:     messageID,
		AlertID:       e.AlertID,
		IntegrationID: e.IntegrationID,
		ExternalID:    e.ExternalID,
		Service:       e.Service,
		Environment:   e.Environment,
		Severity:      e.Severity,
		Name:          e.Name,
		Message:       e.Message,
		DedupKey:      e.DedupKey,
	}
}
