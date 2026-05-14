package main

import (
	"net/http"
	"os"

	"github.com/labstack/echo/v5"
	"github.com/tbikbulatov/go-pulseops/internal/platform/config"
	"github.com/tbikbulatov/go-pulseops/internal/platform/db"
	"github.com/tbikbulatov/go-pulseops/internal/platform/validation"
)

func main() {
	e := echo.New()
	e.Validator = validation.NewValidator()

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

	addr := ":" + cfg.App.Port
	e.Logger.Info("starting pulseops api", "addr", addr)

	if err := e.Start(addr); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}
