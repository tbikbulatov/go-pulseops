package grpc

import (
	incidentv1 "github.com/tbikbulatov/go-pulseops/gen/incident/v1"
	"github.com/tbikbulatov/go-pulseops/internal/incident/domain"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func incidentToProto(incident domain.Incident) *incidentv1.Incident {
	return &incidentv1.Incident{
		Id:          incident.ID,
		Service:     incident.Service,
		Environment: incident.Environment,
		Severity:    incident.Severity.String(),
		Status:      incident.Status.String(),
		DedupKey:    incident.DedupKey,
		AlertCount:  int32(incident.AlertCount),
		FirstSeenAt: timestamppb.New(incident.FirstSeenAt),
		LastSeenAt:  timestamppb.New(incident.LastSeenAt),
		CreatedAt:   timestamppb.New(incident.CreatedAt),
		UpdatedAt:   timestamppb.New(incident.UpdatedAt),
		Version:     int32(incident.Version),
	}
}
