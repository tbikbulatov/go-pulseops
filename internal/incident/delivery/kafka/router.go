package kafka

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"
	platformkafka "github.com/tbikbulatov/go-pulseops/internal/platform/kafka"
)

type Router struct {
	alertReceivedHandler    AlertReceivedHandler
	incidentResolvedHandler IncidentResolvedHandler
}

func NewRouter(
	alertReceivedHandler AlertReceivedHandler,
	incidentResolvedHandler IncidentResolvedHandler,
) *Router {
	return &Router{
		alertReceivedHandler:    alertReceivedHandler,
		incidentResolvedHandler: incidentResolvedHandler,
	}
}

func (r *Router) Handle(ctx context.Context, msg *sarama.ConsumerMessage) error {
	switch msg.Topic {
	case platformkafka.TopicAlertReceived:
		return r.alertReceivedHandler.Handle(ctx, msg)
	case platformkafka.TopicIncidentResolved:
		return r.incidentResolvedHandler.Handle(ctx, msg)
	default:
		return fmt.Errorf("unsupported topic: %s", msg.Topic)
	}
}
