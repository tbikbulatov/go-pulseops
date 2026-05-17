package processalert

import (
	"context"

	"github.com/tbikbulatov/go-pulseops/internal/incident/domain"
)

type MessageDeduplicator interface {
	TryStartProcessing(ctx context.Context, consumerName string, messageID string) (bool, error)
}

type IncidentRepository interface {
	FindActiveByServiceEnvDedupKey(ctx context.Context, service, environment, dedupKey string) (domain.Incident, bool, error)
	Create(ctx context.Context, incident domain.Incident) (domain.Incident, error)
	Save(ctx context.Context, incident domain.Incident) error
}

type IncidentEventRepository interface {
	Create(ctx context.Context, event domain.IncidentEvent) error
}

type TransactionManager interface {
	WithinTx(ctx context.Context, fn func(ctx context.Context) error) error
}
