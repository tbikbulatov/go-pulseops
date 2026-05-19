package domain

import (
	"errors"
	"fmt"
)

type Status string

const (
	StatusOpen         Status = "open"
	StatusAcknowledged Status = "acknowledged"
	StatusResolved     Status = "resolved"
)

var ErrInvalidStatus = errors.New("invalid incident status")

func NewStatus(value string) (Status, error) {
	status := Status(value)
	if !isValidStatus(status) {
		return "", ErrInvalidStatus
	}

	return status, nil
}

func (s Status) String() string {
	return string(s)
}

func (s Status) Value() (string, error) {
	if !isValidStatus(s) {
		return "", ErrInvalidStatus
	}

	return s.String(), nil
}

func (s *Status) Scan(value any) error {
	switch v := value.(type) {
	case string:
		return s.scanString(v)
	case []byte:
		return s.scanString(string(v))
	default:
		return fmt.Errorf("scan incident status: unsupported value %T", value)
	}
}

func isValidStatus(s Status) bool {
	switch s {
	case StatusOpen, StatusAcknowledged, StatusResolved:
		return true
	default:
		return false
	}
}

func (s *Status) scanString(value string) error {
	status, err := NewStatus(value)
	if err != nil {
		return err
	}

	*s = status
	return nil
}
