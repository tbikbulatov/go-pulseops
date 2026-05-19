package manageincident

import (
	"errors"

	"github.com/tbikbulatov/go-pulseops/internal/incident/domain"
	"github.com/tbikbulatov/go-pulseops/internal/platform/apperror"
)

func mapDomainError(err error) error {
	switch {
	case err == nil:
		return nil
	case errors.Is(err, domain.ErrInvalidStatus):
		return apperror.Wrap(apperror.CodeInvalidArgument, err)
	case errors.Is(err, domain.ErrIncidentResolved),
		errors.Is(err, domain.ErrIncidentAlreadyResolved),
		errors.Is(err, domain.ErrIncidentAlreadyAcknowledged):
		return apperror.Wrap(apperror.CodeFailedPrecondition, err)
	case errors.Is(err, domain.ErrUnexpectedIncidentStatus):
		return apperror.Wrap(apperror.CodeAborted, err)
	default:
		return err
	}
}
