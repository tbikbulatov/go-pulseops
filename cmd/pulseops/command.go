package main

import (
	"fmt"

	"github.com/tbikbulatov/go-pulseops/internal/platform/config"
	"github.com/tbikbulatov/go-pulseops/internal/platform/logger"
)

func run(args []string) error {
	cfg, err := config.NewConfig()
	if err != nil {
		return err
	}

	log, err := logger.New("pulseops", cfg.App.Env, cfg.App.LogLevel)
	if err != nil {
		return err
	}

	cmd := "api"
	if len(args) > 0 {
		cmd = args[0]
	}

	switch cmd {
	case "api":
		return runAPI(cfg, log)
	case "publish-alerts-once":
		return publishAlertsOutbox(cfg)
	case "incident-processor":
		return runIncidentProcessor(cfg)
	case "grpc":
		return runGRPC(cfg)
	case "kafka-ping":
		return runKafkaPing(cfg)
	default:
		return fmt.Errorf("unknown command %q", cmd)
	}
}
