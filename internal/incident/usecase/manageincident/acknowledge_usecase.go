package manageincident

import (
	"context"
	"encoding/json"
	"time"

	d "github.com/tbikbulatov/go-pulseops/internal/incident/domain"
	"github.com/tbikbulatov/go-pulseops/internal/platform/apperror"
)

type AcknowledgeUsecase struct {
	tm           TransactionManager
	incidentRepo IncidentRepository
	eventRepo    IncidentEventRepository
}

func NewAcknowledgeUsecase(
	tm TransactionManager,
	incidentRepo IncidentRepository,
	eventRepo IncidentEventRepository,
) *AcknowledgeUsecase {
	return &AcknowledgeUsecase{
		tm:           tm,
		incidentRepo: incidentRepo,
		eventRepo:    eventRepo,
	}
}

func (uc *AcknowledgeUsecase) Handle(ctx context.Context, cmd AcknowledgeCommand) error {
	if cmd.IncidentID == "" || cmd.Actor == "" {
		return apperror.Wrap(apperror.CodeInvalidArgument, ErrInvalidCommand)
	}

	return uc.tm.WithinTx(ctx, func(ctx context.Context) error {
		return uc.handle(ctx, cmd)
	})
}

func (uc *AcknowledgeUsecase) handle(ctx context.Context, cmd AcknowledgeCommand) error {
	incident, err := uc.loadForUpdate(ctx, cmd.IncidentID)
	if err != nil {
		return err
	}
	if err := ensureExpectedStatus(incident, cmd.ExpectedStatus); err != nil {
		return err
	}
	if err := incident.Acknowledge(time.Now()); err != nil {
		return mapDomainError(err)
	}
	if err := uc.incidentRepo.Save(ctx, incident); err != nil {
		return err
	}

	return uc.eventRepo.Create(ctx, newAckIncidentEvent(cmd))
}

func (uc *AcknowledgeUsecase) loadForUpdate(ctx context.Context, id string) (d.Incident, error) {
	incident, found, err := uc.incidentRepo.FindByIDForUpdate(ctx, id)
	if err != nil {
		return d.Incident{}, err
	}
	if !found {
		return d.Incident{}, apperror.Wrap(apperror.CodeNotFound, ErrIncidentNotFound)
	}

	return incident, nil
}

func ensureExpectedStatus(incident d.Incident, expected string) error {
	if expected == "" {
		return nil
	}

	status, err := d.NewStatus(expected)
	if err != nil {
		return apperror.Wrap(apperror.CodeInvalidArgument, err)
	}
	if incident.Status != status {
		return apperror.Wrap(apperror.CodeAborted, ErrWrongExpectedStatus)
	}

	return nil
}

func newAckIncidentEvent(cmd AcknowledgeCommand) d.IncidentEvent {
	body, _ := json.Marshal(cmd.toMap())

	return d.IncidentEvent{
		IncidentID: cmd.IncidentID,
		Type:       d.TypeIncidentAcknowledged,
		Payload:    body,
		CreatedAt:  time.Now(),
	}
}
