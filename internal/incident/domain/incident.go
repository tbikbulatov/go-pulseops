package domain

import (
	"time"

	alertdomain "github.com/tbikbulatov/go-pulseops/internal/alert/domain"
	"github.com/tbikbulatov/go-pulseops/internal/shared/domain/valueobject"
)

type Incident struct {
	ID          string
	Service     string
	Environment string
	Severity    valueobject.Severity
	Status      Status
	DedupKey    string
	AlertCount  int
	Version     int
	FirstSeenAt time.Time
	LastSeenAt  time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewFromAlert(alert alertdomain.Alert) Incident {
	now := alert.CreatedAt
	if now.IsZero() {
		now = time.Now()
	}

	return Incident{
		Service:     alert.Service,
		Environment: alert.Environment,
		Severity:    alert.Severity,
		Status:      StatusOpen,
		DedupKey:    alert.DedupKey,
		AlertCount:  1,
		Version:     1,
		FirstSeenAt: now,
		LastSeenAt:  now,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func (i *Incident) ApplyAlert(alert alertdomain.Alert) error {
	if i.Status == StatusResolved {
		return ErrIncidentResolved
	}

	seenAt := alert.CreatedAt
	if seenAt.IsZero() {
		seenAt = time.Now()
	}

	i.AlertCount++
	i.Version++
	i.LastSeenAt = seenAt
	i.UpdatedAt = seenAt

	if alert.Severity.HigherThan(i.Severity) {
		i.Severity = alert.Severity
	}

	return nil
}

func (i *Incident) Acknowledge(at time.Time) error {
	if i.Status == StatusResolved {
		return ErrIncidentResolved
	}
	if i.Status == StatusAcknowledged {
		return ErrIncidentAlreadyAcknowledged
	}
	if at.IsZero() {
		at = time.Now()
	}

	i.Status = StatusAcknowledged
	i.Version++
	i.UpdatedAt = at
	return nil
}

func (i *Incident) Resolve(at time.Time) error {
	if i.Status == StatusResolved {
		return ErrIncidentAlreadyResolved
	}
	if at.IsZero() {
		at = time.Now()
	}

	i.Status = StatusResolved
	i.Version++
	i.UpdatedAt = at
	return nil
}
