package main

import (
	"net/http"
	"os"

	"github.com/labstack/echo/v5"
	alertsHttp "github.com/tbikbulatov/go-pulseops/internal/alert/delivery/http"
	"github.com/tbikbulatov/go-pulseops/internal/alert/usecase/ingestalert"
	"github.com/tbikbulatov/go-pulseops/internal/platform/config"
	"github.com/tbikbulatov/go-pulseops/internal/platform/db"
	"github.com/tbikbulatov/go-pulseops/internal/platform/validation"
)

func main() {
	e := echo.New()
	validator := validation.NewValidator()
	e.Validator = validator

	cfg, err := config.NewConfig()
	if err != nil {
		e.Logger.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	gorm, err := db.NewGorm(cfg.Postgres)
	if err != nil {
		e.Logger.Error("failed to init gorm", "error", err)
		os.Exit(1)
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

	alertUsecase := &ingestalert.Usecase{}
	alertHandler := alertsHttp.NewAlertHandler(validator, alertUsecase)
	e.POST("/v1/integrations/:integration_key/alerts", alertHandler.IngestAlert)

	addr := ":" + cfg.App.Port
	e.Logger.Info("starting pulseops api", "addr", addr)

	if err := e.Start(addr); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}
