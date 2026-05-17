package postgres

import (
	"context"

	incidentdomain "github.com/tbikbulatov/go-pulseops/internal/incident/domain"
	"github.com/tbikbulatov/go-pulseops/internal/platform/transaction"
	"gorm.io/gorm"
)

type IncidentEventRepository struct {
	db *gorm.DB
}

func NewIncidentEventRepository(db *gorm.DB) *IncidentEventRepository {
	return &IncidentEventRepository{db: db}
}

func (r *IncidentEventRepository) Create(ctx context.Context, event incidentdomain.IncidentEvent) error {
	model := NewIncidentEventModel(event)

	db := transaction.FromContext(ctx, r.db)
	return db.Create(&model).Error
}
