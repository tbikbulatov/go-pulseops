package postgres

import (
	"context"
	"errors"

	incidentdomain "github.com/tbikbulatov/go-pulseops/internal/incident/domain"
	"github.com/tbikbulatov/go-pulseops/internal/platform/transaction"
	"gorm.io/gorm"
)

type IncidentRepository struct {
	db *gorm.DB
}

func NewIncidentRepository(db *gorm.DB) *IncidentRepository {
	return &IncidentRepository{db: db}
}

func (r *IncidentRepository) FindActiveByServiceEnvDedupKey(
	ctx context.Context,
	service string,
	environment string,
	dedupKey string,
) (incidentdomain.Incident, bool, error) {
	var model IncidentModel

	db := transaction.FromContext(ctx, r.db)
	err := db.
		Where("service = ? AND environment = ? AND dedup_key = ? AND status IN ?", service, environment, dedupKey, []string{
			incidentdomain.StatusOpen,
			incidentdomain.StatusAcknowledged,
		}).
		Order("created_at ASC").
		First(&model).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return incidentdomain.Incident{}, false, nil
	}
	if err != nil {
		return incidentdomain.Incident{}, false, err
	}

	return model.ToDomain(), true, nil
}

func (r *IncidentRepository) Create(ctx context.Context, incident incidentdomain.Incident) (incidentdomain.Incident, error) {
	model := NewIncidentModel(incident)

	db := transaction.FromContext(ctx, r.db)
	if err := db.Create(&model).Error; err != nil {
		return incidentdomain.Incident{}, err
	}

	return model.ToDomain(), nil
}

func (r *IncidentRepository) Save(ctx context.Context, incident incidentdomain.Incident) error {
	model := NewIncidentModel(incident)

	db := transaction.FromContext(ctx, r.db)
	return db.Save(&model).Error
}
