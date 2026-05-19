package manageincident

import "errors"

var (
	ErrIncidentNotFound    = errors.New("incident not found")
	ErrInvalidCommand      = errors.New("invalid incident command")
	ErrWrongExpectedStatus = errors.New("wrong expected status")
)
