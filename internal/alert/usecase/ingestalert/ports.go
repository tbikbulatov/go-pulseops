package ingestalert

import (
	"context"

	alertdomain "github.com/tbikbulatov/go-pulseops/internal/alert/domain"
	outboxdomain "github.com/tbikbulatov/go-pulseops/internal/outbox/domain"
)

type IngestAlertUsecase interface {
	Handle(ctx context.Context, cmd IngestAlertCommand) (IngestAlertResult, error)
}

type IntegrationRepository interface {
	FindActiveByKey(ctx context.Context, key string) (alertdomain.Integration, error)
}

type AlertRepository interface {
	Create(ctx context.Context, alert alertdomain.Alert) (alertdomain.Alert, error)
}

type OutboxRepository interface {
	Create(ctx context.Context, event outboxdomain.Event) error
}

type TransactionManager interface {
	WithinTx(ctx context.Context, fn func(ctx context.Context) error) error
}
