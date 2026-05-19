package main

import (
	"context"
	"fmt"
	"net"
	"os/signal"
	"syscall"

	incidentv1 "github.com/tbikbulatov/go-pulseops/gen/incident/v1"
	incidentgrpc "github.com/tbikbulatov/go-pulseops/internal/incident/delivery/grpc"
	incidentpg "github.com/tbikbulatov/go-pulseops/internal/incident/infrastructure/postgres"
	"github.com/tbikbulatov/go-pulseops/internal/incident/query"
	"github.com/tbikbulatov/go-pulseops/internal/incident/usecase/manageincident"
	"github.com/tbikbulatov/go-pulseops/internal/platform/config"
	"github.com/tbikbulatov/go-pulseops/internal/platform/db"
	"github.com/tbikbulatov/go-pulseops/internal/platform/transaction"
	"google.golang.org/grpc"
)

func runGRPC(cfg *config.Config) error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	addr := ":" + cfg.App.GRPCPort
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("listen grpc on %s: %w", addr, err)
	}

	gormDB, err := db.NewGorm(cfg.Postgres)
	if err != nil {
		return fmt.Errorf("init gorm: %w", err)
	}

	incidentReader := incidentpg.NewIncidentReader(gormDB)
	querySrv := query.NewService(incidentReader)
	incidentQuerySrv := incidentgrpc.NewIncidentQueryServer(querySrv)

	tm := transaction.NewGormManager(gormDB)
	incidentRepo := incidentpg.NewIncidentRepository(gormDB)
	eventRepo := incidentpg.NewIncidentEventRepository(gormDB)
	ackUC := manageincident.NewAcknowledgeUsecase(tm, incidentRepo, eventRepo)
	resolveUC := manageincident.NewResolveUsecase(tm, incidentRepo, eventRepo)
	incidentCommandSrv := incidentgrpc.NewIncidentCommandService(ackUC, resolveUC)

	srv := grpc.NewServer()
	incidentv1.RegisterIncidentQueryServiceServer(srv, incidentQuerySrv)
	incidentv1.RegisterIncidentCommandServiceServer(srv, incidentCommandSrv)

	errCh := make(chan error, 1)
	go func() {
		errCh <- srv.Serve(lis)
	}()

	select {
	case <-ctx.Done():
		srv.GracefulStop()
		return nil
	case err := <-errCh:
		if err != nil {
			return fmt.Errorf("serve grpc: %w", err)
		}
		return nil
	}
}
