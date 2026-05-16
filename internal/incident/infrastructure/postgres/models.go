package postgres

import (
	"encoding/json"
	"time"

	incidentdomain "github.com/tbikbulatov/go-pulseops/internal/incident/domain"
	"github.com/tbikbulatov/go-pulseops/internal/shared/domain/valueobject"
)

type IncidentModel struct {
	ID          string `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Service     string
	Environment string
	Severity    valueobject.Severity
	Status      string
	DedupKey    string `gorm:"column:dedup_key"`
	AlertCount  int    `gorm:"column:alert_count"`
	FirstSeenAt time.Time
	LastSeenAt  time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (IncidentModel) TableName() string {
	return "incidents"
}

func NewIncidentModel(incident incidentdomain.Incident) IncidentModel {
	return IncidentModel{
		ID:          incident.ID,
		Service:     incident.Service,
		Environment: incident.Environment,
		Severity:    incident.Severity,
		Status:      incident.Status,
		DedupKey:    incident.DedupKey,
		AlertCount:  incident.AlertCount,
		FirstSeenAt: incident.FirstSeenAt,
		LastSeenAt:  incident.LastSeenAt,
		CreatedAt:   incident.CreatedAt,
		UpdatedAt:   incident.UpdatedAt,
	}
}

func (m IncidentModel) ToDomain() incidentdomain.Incident {
	return incidentdomain.Incident{
		ID:          m.ID,
		Service:     m.Service,
		Environment: m.Environment,
		Severity:    m.Severity,
		Status:      m.Status,
		DedupKey:    m.DedupKey,
		AlertCount:  m.AlertCount,
		FirstSeenAt: m.FirstSeenAt,
		LastSeenAt:  m.LastSeenAt,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

type IncidentEventModel struct {
	ID         string          `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	IncidentID string          `gorm:"type:uuid;column:incident_id"`
	Type       string          `gorm:"column:type"`
	Payload    json.RawMessage `gorm:"type:jsonb;column:payload"`
	CreatedAt  time.Time       `gorm:"column:created_at"`
}

func (IncidentEventModel) TableName() string {
	return "incident_events"
}

func NewIncidentEventModel(event incidentdomain.IncidentEvent) IncidentEventModel {
	return IncidentEventModel{
		ID:         event.ID,
		IncidentID: event.IncidentID,
		Type:       event.Type,
		Payload:    event.Payload,
		CreatedAt:  event.CreatedAt,
	}
}

func (m IncidentEventModel) ToDomain() incidentdomain.IncidentEvent {
	return incidentdomain.IncidentEvent{
		ID:         m.ID,
		IncidentID: m.IncidentID,
		Type:       m.Type,
		Payload:    m.Payload,
		CreatedAt:  m.CreatedAt,
	}
}

type ProcessedMessageModel struct {
	ID           string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	ConsumerName string    `gorm:"column:consumer_name"`
	MessageID    string    `gorm:"type:uuid;column:message_id"`
	ProcessedAt  time.Time `gorm:"column:processed_at"`
}

func (ProcessedMessageModel) TableName() string {
	return "processed_messages"
}
