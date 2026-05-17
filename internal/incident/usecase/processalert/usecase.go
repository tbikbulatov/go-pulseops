package processalert

import (
	"context"
	"encoding/json"
	"time"

	incidentdomain "github.com/tbikbulatov/go-pulseops/internal/incident/domain"
)

const ConsumerName = "incident-alert-received"

type Usecase struct {
	tm           TransactionManager
	deduplicator MessageDeduplicator
	incidentRepo IncidentRepository
	eventRepo    IncidentEventRepository
}

func NewUsecase(
	tm TransactionManager,
	deduplicator MessageDeduplicator,
	incidentRepo IncidentRepository,
	eventRepo IncidentEventRepository,
) *Usecase {
	return &Usecase{
		tm:           tm,
		deduplicator: deduplicator,
		incidentRepo: incidentRepo,
		eventRepo:    eventRepo,
	}
}

func (uc *Usecase) Handle(ctx context.Context, cmd Command) error {
	return uc.tm.WithinTx(ctx, func(ctx context.Context) error {
		return uc.handle(ctx, cmd)
	})
}

func (uc *Usecase) handle(ctx context.Context, cmd Command) error {
	started, err := uc.deduplicator.TryStartProcessing(ctx, ConsumerName, cmd.MessageID)
	if err != nil {
		return err
	}
	if !started {
		return nil
	}

	alert, err := cmd.ToAlert()
	if err != nil {
		return err
	}

	incident, found, err := uc.incidentRepo.FindActiveByServiceEnvDedupKey(ctx, cmd.Service, cmd.Environment, cmd.DedupKey)
	if err != nil {
		return err
	}

	eventType := incidentdomain.TypeAlertDeduplicated
	if !found {
		incident = incidentdomain.NewFromAlert(alert)
		incident, err = uc.incidentRepo.Create(ctx, incident)
		if err != nil {
			return err
		}
		eventType = incidentdomain.TypeIncidentCreated
	} else {
		if err := incident.ApplyAlert(alert); err != nil {
			return err
		}
		if err := uc.incidentRepo.Save(ctx, incident); err != nil {
			return err
		}
	}

	event, err := uc.newIncidentEvent(incident.ID, eventType, cmd)
	if err != nil {
		return err
	}

	return uc.eventRepo.Create(ctx, event)
}

func (uc *Usecase) newIncidentEvent(incidentID string, eventType string, cmd Command) (incidentdomain.IncidentEvent, error) {
	payload, err := json.Marshal(map[string]any{
		"message_id":     cmd.MessageID,
		"alert_id":       cmd.AlertID,
		"integration_id": cmd.IntegrationID,
		"external_id":    cmd.ExternalID,
		"service":        cmd.Service,
		"environment":    cmd.Environment,
		"severity":       cmd.Severity,
		"name":           cmd.Name,
		"message":        cmd.Message,
		"dedup_key":      cmd.DedupKey,
	})
	if err != nil {
		return incidentdomain.IncidentEvent{}, err
	}

	return incidentdomain.IncidentEvent{
		IncidentID: incidentID,
		Type:       eventType,
		Payload:    payload,
		CreatedAt:  time.Now(),
	}, nil
}
