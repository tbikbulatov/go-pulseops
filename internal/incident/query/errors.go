package query

import "errors"

var (
	ErrIncidentNotFound = errors.New("incident not found")
	ErrInvalidQuery     = errors.New("invalid query")
)
