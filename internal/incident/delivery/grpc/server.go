package grpc

import (
	"context"

	incidentv1 "github.com/tbikbulatov/go-pulseops/gen/incident/v1"
	"github.com/tbikbulatov/go-pulseops/internal/incident/domain"
	"github.com/tbikbulatov/go-pulseops/internal/incident/query"
	"github.com/tbikbulatov/go-pulseops/internal/incident/usecase/manageincident"
)

type QueryService interface {
	GetIncident(ctx context.Context, q query.GetIncidentQuery) (domain.Incident, error)
	ListIncidents(ctx context.Context, q query.ListIncidentsQuery) ([]domain.Incident, error)
}

type IncidentQueryServer struct {
	incidentv1.UnimplementedIncidentQueryServiceServer
	service QueryService
}

func NewIncidentQueryServer(service QueryService) *IncidentQueryServer {
	return &IncidentQueryServer{service: service}
}

func (s *IncidentQueryServer) GetIncident(
	ctx context.Context,
	req *incidentv1.GetIncidentRequest,
) (*incidentv1.GetIncidentResponse, error) {
	incident, err := s.service.GetIncident(ctx, query.GetIncidentQuery{ID: req.GetId()})
	if err != nil {
		return &incidentv1.GetIncidentResponse{}, mapError(err)
	}

	return &incidentv1.GetIncidentResponse{Incident: incidentToProto(incident)}, nil
}

func (s *IncidentQueryServer) ListIncidents(
	ctx context.Context,
	req *incidentv1.ListIncidentsRequest,
) (*incidentv1.ListIncidentsResponse, error) {
	incidents, err := s.service.ListIncidents(ctx, query.ListIncidentsQuery{
		Status: req.GetStatus(),
		Limit:  int(req.GetLimit()),
		Offset: int(req.GetOffset()),
	})
	if err != nil {
		return nil, mapError(err)
	}

	items := make([]*incidentv1.Incident, 0, len(incidents))
	for _, incident := range incidents {
		items = append(items, incidentToProto(incident))
	}

	return &incidentv1.ListIncidentsResponse{Incidents: items}, nil
}

type AcknowledgeUsecase interface {
	Handle(ctx context.Context, c manageincident.AcknowledgeCommand) error
}

type ResolveUsecase interface {
	Handle(ctx context.Context, c manageincident.ResolveCommand) error
}

type IncidentCommandService struct {
	incidentv1.UnimplementedIncidentCommandServiceServer
	acknowledgeUC AcknowledgeUsecase
	resolveUC     ResolveUsecase
}

func NewIncidentCommandService(
	acknowledgeUC AcknowledgeUsecase,
	resolveUC ResolveUsecase,
) *IncidentCommandService {
	return &IncidentCommandService{
		acknowledgeUC: acknowledgeUC,
		resolveUC:     resolveUC,
	}
}

func (s *IncidentCommandService) AcknowledgeIncident(
	ctx context.Context,
	req *incidentv1.AcknowledgeIncidentRequest,
) (*incidentv1.AcknowledgeIncidentResponse, error) {
	err := s.acknowledgeUC.Handle(ctx, manageincident.AcknowledgeCommand{
		IncidentID:     req.GetIncidentId(),
		CommandID:      req.GetCommandId(),
		Actor:          req.GetActor(),
		ExpectedStatus: req.GetExpectedStatus(),
	})
	if err != nil {
		return nil, mapError(err)
	}

	return &incidentv1.AcknowledgeIncidentResponse{Success: true}, nil
}

func (s *IncidentCommandService) ResolveIncident(
	ctx context.Context,
	req *incidentv1.ResolveIncidentRequest,
) (*incidentv1.ResolveIncidentResponse, error) {
	err := s.resolveUC.Handle(ctx, manageincident.ResolveCommand{
		IncidentID:     req.GetIncidentId(),
		CommandID:      req.GetCommandId(),
		Actor:          req.GetActor(),
		Reason:         req.GetReason(),
		ExpectedStatus: req.GetExpectedStatus(),
	})
	if err != nil {
		return nil, mapError(err)
	}

	return &incidentv1.ResolveIncidentResponse{Success: true}, nil
}
