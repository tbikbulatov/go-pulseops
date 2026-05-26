package main

import (
	"log/slog"
	"net/http"

	echoprometheus "github.com/labstack/echo-prometheus"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"github.com/prometheus/client_golang/prometheus"
	alertsHttp "github.com/tbikbulatov/go-pulseops/internal/alert/delivery/http"
	alertPostgres "github.com/tbikbulatov/go-pulseops/internal/alert/infrastructure/postgres"
	"github.com/tbikbulatov/go-pulseops/internal/alert/usecase/ingestalert"
	outboxPostgres "github.com/tbikbulatov/go-pulseops/internal/outbox/infrastructure/postgres"
	"github.com/tbikbulatov/go-pulseops/internal/platform/config"
	"github.com/tbikbulatov/go-pulseops/internal/platform/db"
	"github.com/tbikbulatov/go-pulseops/internal/platform/metrics"
	"github.com/tbikbulatov/go-pulseops/internal/platform/transaction"
	"github.com/tbikbulatov/go-pulseops/internal/platform/validation"
	"gorm.io/gorm"
)

func runAPI(cfg *config.Config, logger *slog.Logger) error {
	e := echo.New()
	validator := validation.NewValidator()
	e.Validator = validator

	gorm, err := db.NewGorm(cfg.Postgres)
	if err != nil {
		logger.Error("failed to init gorm", "error", err)
		return err
	}

	promReg := prometheus.NewRegistry()
	appMetrics := metrics.New()
	appMetrics.Register(promReg)

	setupAPIMiddlewares(e, logger, promReg)
	setupSystemRoutes(e, gorm, promReg)
	setupAlertRoutes(e, gorm, validator, appMetrics.Alert)

	addr := ":" + cfg.App.Port
	logger.Info("starting pulseops api", "addr", addr)

	if err := e.Start(addr); err != nil {
		logger.Error("failed to start server", "error", err)
		return err
	}

	return nil
}

func setupAPIMiddlewares(e *echo.Echo, logger *slog.Logger, reg *prometheus.Registry) {
	e.Use(middleware.RequestID())
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogMethod:    true,
		LogURI:       true,
		LogStatus:    true,
		LogLatency:   true,
		LogRemoteIP:  true,
		LogRequestID: true,
		HandleError:  true,
		Skipper: func(c *echo.Context) bool {
			return c.Path() == "/healthz" || c.Path() == "/metrics"
		},
		LogValuesFunc: func(c *echo.Context, v middleware.RequestLoggerValues) error {
			attrs := []slog.Attr{
				slog.String("method", v.Method),
				slog.String("uri", v.URI),
				slog.String("route", c.Path()),
				slog.Int("status", v.Status),
				slog.Int("duration_ms", int(v.Latency.Milliseconds())),
				slog.String("remote_ip", v.RemoteIP),
				slog.String("request_id", v.RequestID),
			}
			if v.Error == nil {
				logger.LogAttrs(c.Request().Context(), slog.LevelInfo, "request", attrs...)
			} else {
				logger.LogAttrs(c.Request().Context(), slog.LevelError, "request_error",
					append(attrs, slog.String("err", v.Error.Error()))...,
				)
			}
			return nil
		},
	}))
	e.Use(echoprometheus.NewMiddlewareWithConfig(echoprometheus.MiddlewareConfig{
		Subsystem:                 "pulseops",
		Registerer:                reg,
		DoNotUseRequestPathFor404: true,
	}))
}

func setupSystemRoutes(e *echo.Echo, gorm *gorm.DB, reg *prometheus.Registry) {
	e.GET("/metrics", echoprometheus.NewHandlerWithConfig(echoprometheus.HandlerConfig{Gatherer: reg}))

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
}

func setupAlertRoutes(e *echo.Echo, gorm *gorm.DB, v *validation.Validator, m metrics.AlertMetrics) {
	txManager := transaction.NewGormManager(gorm)
	alertRepo := alertPostgres.NewAlertRepository(gorm)
	integrationRepo := alertPostgres.NewIntegrationRepository(gorm)
	outboxRepo := outboxPostgres.NewOutboxRepository(gorm)
	alertUsecase := ingestalert.NewUsecase(txManager, alertRepo, integrationRepo, outboxRepo)
	alertHandler := alertsHttp.NewAlertHandler(v, alertUsecase, m)

	e.POST("/v1/integrations/:integration_key/alerts", alertHandler.IngestAlert)
}
