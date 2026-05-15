package domain

import "time"

type Alert struct {
	ID            string
	IntegrationID string
	ExternalID    string
	Service       string
	Environment   string
	Severity      string
	Name          string
	Message       string
	DedupKey      string
	CreatedAt     time.Time
}
