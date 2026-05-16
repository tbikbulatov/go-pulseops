package domain

import (
	"encoding/json"
	"time"
)

const (
	TypeIncidentCreated      = "incident.created"
	TypeIncidentUpdated      = "incident.updated"
	TypeIncidentAcknowledged = "incident.acknowledged"
	TypeIncidentResolved     = "incident.resolved"
	TypeAlertDeduplicated    = "alert.deduplicated"
)

type IncidentEvent struct {
	ID         string
	IncidentID string
	Type       string
	Payload    json.RawMessage
	CreatedAt  time.Time
}
