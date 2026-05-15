package http

import (
	uc "github.com/tbikbulatov/go-pulseops/internal/alert/usecase/ingestalert"
)

type IngestAlertRequest struct {
	IntegrationKey string
	ExternalID     string `json:"external_id" validate:"required,max=128"`
	Service        string `json:"service" validate:"required,max=128"`
	Environment    string `json:"environment" validate:"required,max=64"`
	Severity       string `json:"severity" validate:"required,oneof=info warning critical"`
	Name           string `json:"name" validate:"required,max=256"`
	Message        string `json:"message" validate:"required,max=2048"`
	DedupKey       string `json:"dedup_key" validate:"required,max=256"`
}

func (r *IngestAlertRequest) toCommand() uc.IngestAlertCommand {
	return uc.IngestAlertCommand{
		IntegrationKey: r.IntegrationKey,
		ExternalID:     r.ExternalID,
		Service:        r.Service,
		Environment:    r.Environment,
		Severity:       r.Severity,
		Name:           r.Name,
		Message:        r.Message,
		DedupKey:       r.DedupKey,
	}
}
