package postgres

import (
	"context"

	alertdomain "github.com/tbikbulatov/go-pulseops/internal/alert/domain"
	"github.com/tbikbulatov/go-pulseops/internal/platform/transaction"
	"gorm.io/gorm"
)

type IntegrationRepository struct {
	db *gorm.DB
}

func NewIntegrationRepository(db *gorm.DB) *IntegrationRepository {
	return &IntegrationRepository{db: db}
}

func (r *IntegrationRepository) FindActiveByKey(ctx context.Context, key string) (alertdomain.Integration, error) {
	var model IntegrationModel

	db := transaction.FromContext(ctx, r.db)
	if err := db.Where("key = ? AND status = ?", key, "active").First(&model).Error; err != nil {
		return alertdomain.Integration{}, err
	}

	return model.ToDomain(), nil
}
