package query

import incidentdomain "github.com/tbikbulatov/go-pulseops/internal/incident/domain"

const (
	ListItemsDefaultLimit = 10
	ListItemsMaxLimit     = 100
)

type GetIncidentQuery struct {
	ID string
}

type ListIncidentsQuery struct {
	Status string
	Limit  int
	Offset int
}

type Filter struct {
	Status *incidentdomain.Status
	Limit  int
	Offset int
}
