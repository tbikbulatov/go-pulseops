package main

import (
	"fmt"

	"github.com/tbikbulatov/go-pulseops/internal/platform/config"
)

func run(args []string) error {
	cfg, err := config.NewConfig()
	if err != nil {
		return err
	}

	cmd := "api"
	if len(args) > 0 {
		cmd = args[0]
	}

	switch cmd {
	case "api":
		return runAPI(cfg)
	case "publish-once":
		return publishOutbox(cfg)
	case "kafka-ping":
		return runKafkaPing(cfg)
	default:
		return fmt.Errorf("unknown command %q", cmd)
	}
}
