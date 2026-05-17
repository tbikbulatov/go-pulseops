package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"

	"github.com/IBM/sarama"
	incidentkafka "github.com/tbikbulatov/go-pulseops/internal/incident/delivery/kafka"
	incidentpostgres "github.com/tbikbulatov/go-pulseops/internal/incident/infrastructure/postgres"
	"github.com/tbikbulatov/go-pulseops/internal/incident/usecase/processalert"
	"github.com/tbikbulatov/go-pulseops/internal/platform/config"
	"github.com/tbikbulatov/go-pulseops/internal/platform/db"
	platformkafka "github.com/tbikbulatov/go-pulseops/internal/platform/kafka"
	"github.com/tbikbulatov/go-pulseops/internal/platform/transaction"
	"gorm.io/gorm"
)

func runIncidentProcessor(cfg *config.Config) error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	saramaCfg := sarama.NewConfig()
	saramaCfg.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRange()
	saramaCfg.Consumer.Offsets.Initial = sarama.OffsetOldest

	group, err := sarama.NewConsumerGroup(cfg.Kafka.BrokerList(), cfg.Kafka.IncidentProcessorGroup, saramaCfg)
	if err != nil {
		return fmt.Errorf("create incident processor consumer group: %w", err)
	}
	defer group.Close()

	gormDB, err := db.NewGorm(cfg.Postgres)
	if err != nil {
		return fmt.Errorf("init gorm: %w", err)
	}

	handler := newIncidentProcessorHandler(gormDB)
	topics := []string{platformkafka.TopicAlertReceived}

	for ctx.Err() == nil {
		if err := group.Consume(ctx, topics, handler); err != nil {
			return fmt.Errorf("consume incident processor messages: %w", err)
		}
	}

	return nil
}

type incidentProcessorHandler struct {
	router *incidentkafka.Router
}

func newIncidentProcessorHandler(gormDB *gorm.DB) *incidentProcessorHandler {
	tm := transaction.NewGormManager(gormDB)
	deduplicator := incidentpostgres.NewMessageDeduplicator(gormDB)
	incidentRepo := incidentpostgres.NewIncidentRepository(gormDB)
	eventRepo := incidentpostgres.NewIncidentEventRepository(gormDB)
	processAlertUsecase := processalert.NewUsecase(tm, deduplicator, incidentRepo, eventRepo)

	router := incidentkafka.NewRouter(
		incidentkafka.NewAlertReceivedHandler(processAlertUsecase),
		incidentkafka.IncidentResolvedHandler{},
	)

	return &incidentProcessorHandler{
		router: router,
	}
}

func (incidentProcessorHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (incidentProcessorHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h incidentProcessorHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		if err := h.router.Handle(session.Context(), message); err != nil {
			return err
		}

		session.MarkMessage(message, "")
	}

	return nil
}
