package ingestalert

import (
	"context"
	"encoding/json"
	"time"

	alertdomain "github.com/tbikbulatov/go-pulseops/internal/alert/domain"
	outboxdomain "github.com/tbikbulatov/go-pulseops/internal/outbox/domain"
	"github.com/tbikbulatov/go-pulseops/internal/shared/domain/valueobject"
)

type Usecase struct {
	tm         TransactionManager
	alertRepo  AlertRepository
	intgRepo   IntegrationRepository
	outboxRepo OutboxRepository
}

func NewUsecase(
	tm TransactionManager,
	alertRepo AlertRepository,
	intgRepo IntegrationRepository,
	outboxRepo OutboxRepository,
) *Usecase {
	return &Usecase{
		tm:         tm,
		alertRepo:  alertRepo,
		intgRepo:   intgRepo,
		outboxRepo: outboxRepo,
	}
}

func (uc *Usecase) Handle(ctx context.Context, cmd IngestAlertCommand) (IngestAlertResult, error) {
	alert, err := uc.populateFromCmd(cmd)
	if err != nil {
		return IngestAlertResult{}, err
	}

	if err := uc.tm.WithinTx(ctx, func(ctx context.Context) error {
		intg, err := uc.intgRepo.FindActiveByKey(ctx, cmd.IntegrationKey)
		if err != nil {
			return err
		}

		alert.IntegrationID = intg.ID
		alert, err = uc.alertRepo.Create(ctx, alert)
		if err != nil {
			return err
		}

		event := uc.createEvent(alert)
		if err := uc.outboxRepo.Create(ctx, event); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return IngestAlertResult{}, err
	}

	return IngestAlertResult{AlertID: alert.ID, Status: "accepted"}, nil
}

func (uc Usecase) populateFromCmd(cmd IngestAlertCommand) (alertdomain.Alert, error) {
	severity, err := valueobject.NewSeverity(cmd.Severity)
	if err != nil {
		return alertdomain.Alert{}, err
	}

	return alertdomain.Alert{
		ExternalID:  cmd.ExternalID,
		Service:     cmd.Service,
		Environment: cmd.Environment,
		Severity:    severity,
		Name:        cmd.Name,
		Message:     cmd.Message,
		DedupKey:    cmd.DedupKey,
		CreatedAt:   time.Now(),
	}, nil
}

func (uc Usecase) createEvent(alert alertdomain.Alert) outboxdomain.Event {
	payload, _ := json.Marshal(map[string]any{
		"alert_id":       alert.ID,
		"integration_id": alert.IntegrationID,
		"external_id":    alert.ExternalID,
		"service":        alert.Service,
		"environment":    alert.Environment,
		"severity":       alert.Severity,
		"name":           alert.Name,
		"message":        alert.Message,
		"dedup_key":      alert.DedupKey,
	})

	return outboxdomain.Event{
		AggregateType: "alert",
		AggregateID:   alert.ID,
		EventType:     "alert.received",
		Payload:       payload,
		Status:        "pending",
		Attempts:      0,
		NextAttemptAt: time.Now(),
		CreatedAt:     time.Now(),
	}
}
