package domain

import "errors"

var (
	ErrIncidentResolved            = errors.New("incident is resolved")
	ErrIncidentAlreadyResolved     = errors.New("incident is already resolved")
	ErrIncidentAlreadyAcknowledged = errors.New("incident already acknowledged")
	ErrUnexpectedIncidentStatus    = errors.New("unexpected incident status")
)
