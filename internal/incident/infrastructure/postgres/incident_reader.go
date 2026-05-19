package postgres

import (
	"context"
	"errors"

	incidentdomain "github.com/tbikbulatov/go-pulseops/internal/incident/domain"
	"github.com/tbikbulatov/go-pulseops/internal/incident/query"
	"github.com/tbikbulatov/go-pulseops/internal/platform/transaction"
	"gorm.io/gorm"
)

type IncidentReader struct {
	db *gorm.DB
}

func NewIncidentReader(db *gorm.DB) *IncidentReader {
	return &IncidentReader{db: db}
}

func (r *IncidentReader) GetByID(ctx context.Context, id string) (incidentdomain.Incident, bool, error) {
	var model IncidentModel

	db := transaction.FromContext(ctx, r.db)
	err := db.Where("id = ?", id).First(&model).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return incidentdomain.Incident{}, false, nil
	}
	if err != nil {
		return incidentdomain.Incident{}, false, err
	}

	return model.ToDomain(), true, nil
}

func (r *IncidentReader) List(ctx context.Context, filter query.Filter) ([]incidentdomain.Incident, error) {
	var models []IncidentModel

	db := transaction.FromContext(ctx, r.db).Order("updated_at DESC").Limit(filter.Limit).Offset(filter.Offset)
	if filter.Status != nil {
		db = db.Where("status = ?", *filter.Status)
	}

	if err := db.Find(&models).Error; err != nil {
		return nil, err
	}

	incidents := make([]incidentdomain.Incident, 0, len(models))
	for _, model := range models {
		incidents = append(incidents, model.ToDomain())
	}

	return incidents, nil
}
