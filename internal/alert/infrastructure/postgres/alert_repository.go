package postgres

import (
	"context"

	alertdomain "github.com/tbikbulatov/go-pulseops/internal/alert/domain"
	"github.com/tbikbulatov/go-pulseops/internal/platform/transaction"
	"gorm.io/gorm"
)

type AlertRepository struct {
	db *gorm.DB
}

func NewAlertRepository(db *gorm.DB) *AlertRepository {
	return &AlertRepository{db: db}
}

func (r *AlertRepository) Create(ctx context.Context, alert alertdomain.Alert) (alertdomain.Alert, error) {
	model := AlertModel{
		ID:            alert.ID,
		IntegrationID: alert.IntegrationID,
		ExternalID:    alert.ExternalID,
		Service:       alert.Service,
		Environment:   alert.Environment,
		Severity:      alert.Severity,
		Name:          alert.Name,
		Message:       alert.Message,
		DedupKey:      alert.DedupKey,
		CreatedAt:     alert.CreatedAt,
	}

	db := transaction.FromContext(ctx, r.db)
	if err := db.Create(&model).Error; err != nil {
		return alertdomain.Alert{}, err
	}

	return model.ToDomain(), nil
}
