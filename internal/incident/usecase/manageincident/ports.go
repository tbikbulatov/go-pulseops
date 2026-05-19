package manageincident

import (
	"context"

	"github.com/tbikbulatov/go-pulseops/internal/incident/domain"
)

type IncidentRepository interface {
	FindByIDForUpdate(ctx context.Context, id string) (domain.Incident, bool, error)
	Save(ctx context.Context, incident domain.Incident) error
}

type IncidentEventRepository interface {
	Create(ctx context.Context, event domain.IncidentEvent) error
}

type TransactionManager interface {
	WithinTx(ctx context.Context, fn func(ctx context.Context) error) error
}
