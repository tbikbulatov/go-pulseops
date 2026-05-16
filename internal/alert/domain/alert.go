package domain

import (
	"time"

	"github.com/tbikbulatov/go-pulseops/internal/shared/domain/valueobject"
)

type Alert struct {
	ID            string
	IntegrationID string
	ExternalID    string
	Service       string
	Environment   string
	Severity      valueobject.Severity
	Name          string
	Message       string
	DedupKey      string
	CreatedAt     time.Time
}
