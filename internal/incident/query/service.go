package query

import (
	"context"

	incidentdomain "github.com/tbikbulatov/go-pulseops/internal/incident/domain"
	"github.com/tbikbulatov/go-pulseops/internal/platform/apperror"
)

type Service struct {
	reader Reader
}

func NewService(reader Reader) *Service {
	return &Service{reader: reader}
}

func (s *Service) GetIncident(ctx context.Context, q GetIncidentQuery) (incidentdomain.Incident, error) {
	if q.ID == "" {
		return incidentdomain.Incident{}, apperror.Wrap(apperror.CodeInvalidArgument, ErrInvalidQuery)
	}

	incident, found, err := s.reader.GetByID(ctx, q.ID)
	if err != nil {
		return incidentdomain.Incident{}, err
	}
	if !found {
		return incidentdomain.Incident{}, apperror.Wrap(apperror.CodeNotFound, ErrIncidentNotFound)
	}

	return incident, nil
}

func (s *Service) ListIncidents(ctx context.Context, q ListIncidentsQuery) ([]incidentdomain.Incident, error) {
	filter, err := normalizeFilter(q)
	if err != nil {
		return nil, err
	}

	return s.reader.List(ctx, filter)
}

func normalizeFilter(q ListIncidentsQuery) (Filter, error) {
	limit := q.Limit
	if limit == 0 {
		limit = ListItemsDefaultLimit
	}
	if limit < 0 || limit > ListItemsMaxLimit || q.Offset < 0 {
		return Filter{}, apperror.Wrap(apperror.CodeInvalidArgument, ErrInvalidQuery)
	}

	var status *incidentdomain.Status
	if q.Status != "" {
		parsedStatus, err := incidentdomain.NewStatus(q.Status)
		if err != nil {
			return Filter{}, apperror.Wrap(apperror.CodeInvalidArgument, err)
		}
		status = &parsedStatus
	}

	return Filter{
		Status: status,
		Limit:  limit,
		Offset: q.Offset,
	}, nil
}
