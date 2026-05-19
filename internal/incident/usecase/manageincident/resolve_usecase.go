package manageincident

import (
	"context"
	"encoding/json"
	"time"

	d "github.com/tbikbulatov/go-pulseops/internal/incident/domain"
	"github.com/tbikbulatov/go-pulseops/internal/platform/apperror"
)

type ResolveUsecase struct {
	tm           TransactionManager
	incidentRepo IncidentRepository
	eventRepo    IncidentEventRepository
}

func NewResolveUsecase(
	tm TransactionManager,
	incidentRepo IncidentRepository,
	eventRepo IncidentEventRepository,
) *ResolveUsecase {
	return &ResolveUsecase{
		tm:           tm,
		incidentRepo: incidentRepo,
		eventRepo:    eventRepo,
	}
}

func (uc *ResolveUsecase) Handle(ctx context.Context, cmd ResolveCommand) error {
	if cmd.IncidentID == "" || cmd.Actor == "" {
		return apperror.Wrap(apperror.CodeInvalidArgument, ErrInvalidCommand)
	}

	return uc.tm.WithinTx(ctx, func(ctx context.Context) error {
		return uc.handle(ctx, cmd)
	})
}

func (uc *ResolveUsecase) handle(ctx context.Context, cmd ResolveCommand) error {
	incident, err := uc.loadForUpdate(ctx, cmd.IncidentID)
	if err != nil {
		return err
	}
	if err := ensureExpectedStatus(incident, cmd.ExpectedStatus); err != nil {
		return err
	}
	if err := incident.Resolve(time.Now()); err != nil {
		return mapDomainError(err)
	}
	if err := uc.incidentRepo.Save(ctx, incident); err != nil {
		return err
	}

	return uc.eventRepo.Create(ctx, newResolveIncidentEvent(cmd))
}

func (uc *ResolveUsecase) loadForUpdate(ctx context.Context, id string) (d.Incident, error) {
	incident, found, err := uc.incidentRepo.FindByIDForUpdate(ctx, id)
	if err != nil {
		return d.Incident{}, err
	}
	if !found {
		return d.Incident{}, apperror.Wrap(apperror.CodeNotFound, ErrIncidentNotFound)
	}

	return incident, nil
}

func newResolveIncidentEvent(cmd ResolveCommand) d.IncidentEvent {
	body, _ := json.Marshal(cmd.toMap())

	return d.IncidentEvent{
		IncidentID: cmd.IncidentID,
		Type:       d.TypeIncidentResolved,
		Payload:    body,
		CreatedAt:  time.Now(),
	}
}
