package query

import (
	"context"

	incidentdomain "github.com/tbikbulatov/go-pulseops/internal/incident/domain"
)

type Reader interface {
	GetByID(ctx context.Context, id string) (incidentdomain.Incident, bool, error)
	List(ctx context.Context, filter Filter) ([]incidentdomain.Incident, error)
}
