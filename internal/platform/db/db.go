package db

import (
	"context"
	"time"

	"github.com/tbikbulatov/go-pulseops/internal/platform/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewGorm(cfg config.PostgresConfig) (*gorm.DB, error) {
	gormDB, err := gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db, err := gormDB.DB()
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}

	return gormDB, nil
}
