package valueobject

import (
	"database/sql/driver"
	"errors"
	"fmt"
)

type Severity string

const (
	SeverityInfo     Severity = "info"
	SeverityWarning  Severity = "warning"
	SeverityCritical Severity = "critical"
)

var ErrInvalidSeverity = errors.New("invalid severity")

func NewSeverity(value string) (Severity, error) {
	severity := Severity(value)
	if !isValidSeverity(severity) {
		return "", ErrInvalidSeverity
	}

	return severity, nil
}

func (s Severity) String() string {
	return string(s)
}

func (s Severity) Value() (driver.Value, error) {
	if !isValidSeverity(s) {
		return nil, ErrInvalidSeverity
	}

	return s.String(), nil
}

func (s *Severity) Scan(value any) error {
	switch v := value.(type) {
	case string:
		return s.scanString(v)
	case []byte:
		return s.scanString(string(v))
	default:
		return fmt.Errorf("scan severity: unsupported value %T", value)
	}
}

func isValidSeverity(s Severity) bool {
	switch s {
	case SeverityInfo, SeverityWarning, SeverityCritical:
		return true
	default:
		return false
	}
}

func (s Severity) Rank() int {
	switch s {
	case SeverityCritical:
		return 3
	case SeverityWarning:
		return 2
	case SeverityInfo:
		return 1
	default:
		return 0
	}
}

func (s Severity) HigherThan(other Severity) bool {
	return s.Rank() > other.Rank()
}

func (s *Severity) scanString(value string) error {
	severity, err := NewSeverity(value)
	if err != nil {
		return err
	}

	*s = severity
	return nil
}
