package postgres

import (
	"time"

	alertdomain "github.com/tbikbulatov/go-pulseops/internal/alert/domain"
)

type IntegrationModel struct {
	ID        string `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Key       string `gorm:"column:key"`
	Name      string
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (IntegrationModel) TableName() string {
	return "integrations"
}

type AlertModel struct {
	ID            string `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	IntegrationID string `gorm:"type:uuid;column:integration_id"`
	ExternalID    string `gorm:"column:external_id"`
	Service       string
	Environment   string
	Severity      string
	Name          string
	Message       string
	DedupKey      string `gorm:"column:dedup_key"`
	CreatedAt     time.Time
}

func (AlertModel) TableName() string {
	return "alerts"
}

func (m IntegrationModel) ToDomain() alertdomain.Integration {
	return alertdomain.Integration{
		ID:        m.ID,
		Key:       m.Key,
		Name:      m.Name,
		Status:    m.Status,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func (m AlertModel) ToDomain() alertdomain.Alert {
	return alertdomain.Alert{
		ID:            m.ID,
		IntegrationID: m.IntegrationID,
		ExternalID:    m.ExternalID,
		Service:       m.Service,
		Environment:   m.Environment,
		Severity:      m.Severity,
		Name:          m.Name,
		Message:       m.Message,
		DedupKey:      m.DedupKey,
		CreatedAt:     m.CreatedAt,
	}
}
