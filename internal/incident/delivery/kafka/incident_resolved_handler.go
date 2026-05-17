package kafka

import (
	"context"

	"github.com/IBM/sarama"
)

type IncidentResolvedHandler struct {
}

func (h *IncidentResolvedHandler) Handle(ctx context.Context, msg *sarama.ConsumerMessage) error {
	return nil
}
