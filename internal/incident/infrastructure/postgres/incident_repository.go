package postgres

import (
	"context"
	"errors"

	d "github.com/tbikbulatov/go-pulseops/internal/incident/domain"
	"github.com/tbikbulatov/go-pulseops/internal/platform/transaction"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
) (d.Incident, bool, error) {
	var model IncidentModel

	db := transaction.FromContext(ctx, r.db)
	err := db.
		Where("service = ? AND environment = ? AND dedup_key = ? AND status IN ?", service, environment, dedupKey, []d.Status{
			d.StatusOpen,
			d.StatusAcknowledged,
		}).
		Order("created_at ASC").
		First(&model).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return d.Incident{}, false, nil
	}
	if err != nil {
		return d.Incident{}, false, err
	}

	return model.ToDomain(), true, nil
}

func (r *IncidentRepository) FindByIDForUpdate(ctx context.Context, id string) (d.Incident, bool, error) {
	var model IncidentModel

	db := transaction.FromContext(ctx, r.db)
	err := db.
		Clauses(clause.Locking{Strength: clause.LockingStrengthUpdate}).
		Where("id = ?", id).
		First(&model).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return d.Incident{}, false, nil
	}
	if err != nil {
		return d.Incident{}, false, err
	}

	return model.ToDomain(), true, nil
}

func (r *IncidentRepository) Create(ctx context.Context, incident d.Incident) (d.Incident, error) {
	model := NewIncidentModel(incident)

	db := transaction.FromContext(ctx, r.db)
	if err := db.Create(&model).Error; err != nil {
		return d.Incident{}, err
	}

	return model.ToDomain(), nil
}

func (r *IncidentRepository) Save(ctx context.Context, incident d.Incident) error {
	model := NewIncidentModel(incident)

	db := transaction.FromContext(ctx, r.db)
	return db.Save(&model).Error
}
