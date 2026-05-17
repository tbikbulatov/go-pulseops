package processalert

import (
	"time"

	"github.com/tbikbulatov/go-pulseops/internal/alert/domain"
	"github.com/tbikbulatov/go-pulseops/internal/shared/domain/valueobject"
)

type Command struct {
	MessageID     string
	AlertID       string
	IntegrationID string
	ExternalID    string
	Service       string
	Environment   string
	Severity      string
	Name          string
	Message       string
	DedupKey      string
}

func (c Command) ToAlert() (domain.Alert, error) {
	severity, err := valueobject.NewSeverity(c.Severity)
	if err != nil {
		return domain.Alert{}, err
	}

	return domain.Alert{
		ID:            c.AlertID,
		IntegrationID: c.IntegrationID,
		ExternalID:    c.ExternalID,
		Service:       c.Service,
		Environment:   c.Environment,
		Severity:      severity,
		Name:          c.Name,
		Message:       c.Message,
		DedupKey:      c.DedupKey,
		CreatedAt:     time.Now(),
	}, nil
}
