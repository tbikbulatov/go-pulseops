package main

import (
	"net/http"

	"github.com/labstack/echo/v5"
	alertsHttp "github.com/tbikbulatov/go-pulseops/internal/alert/delivery/http"
	alertPostgres "github.com/tbikbulatov/go-pulseops/internal/alert/infrastructure/postgres"
	"github.com/tbikbulatov/go-pulseops/internal/alert/usecase/ingestalert"
	outboxPostgres "github.com/tbikbulatov/go-pulseops/internal/outbox/infrastructure/postgres"
	"github.com/tbikbulatov/go-pulseops/internal/platform/config"
	"github.com/tbikbulatov/go-pulseops/internal/platform/db"
	"github.com/tbikbulatov/go-pulseops/internal/platform/transaction"
	"github.com/tbikbulatov/go-pulseops/internal/platform/validation"
)

func runAPI(cfg *config.Config) error {
	e := echo.New()
	validator := validation.NewValidator()
	e.Validator = validator

	gorm, err := db.NewGorm(cfg.Postgres)
	if err != nil {
		e.Logger.Error("failed to init gorm", "error", err)
		return err
	}

	e.GET("/healthz", func(c *echo.Context) error {
		db, err := gorm.DB()
		if err != nil {
			return c.JSON(http.StatusOK, map[string]any{"status": "fail"})
		}
		if pingErr := db.Ping(); pingErr != nil {
			return c.JSON(http.StatusOK, map[string]any{"status": "fail"})
		}

		return c.JSON(http.StatusOK, map[string]any{"status": "ok"})
	})

	txManager := transaction.NewGormManager(gorm)
	alertRepo := alertPostgres.NewAlertRepository(gorm)
	integrationRepo := alertPostgres.NewIntegrationRepository(gorm)
	outboxRepo := outboxPostgres.NewOutboxRepository(gorm)
	alertUsecase := ingestalert.NewUsecase(txManager, alertRepo, integrationRepo, outboxRepo)
	alertHandler := alertsHttp.NewAlertHandler(validator, alertUsecase)
	e.POST("/v1/integrations/:integration_key/alerts", alertHandler.IngestAlert)

	addr := ":" + cfg.App.Port
	e.Logger.Info("starting pulseops api", "addr", addr)

	if err := e.Start(addr); err != nil {
		e.Logger.Error("failed to start server", "error", err)
		return err
	}

	return nil
}
