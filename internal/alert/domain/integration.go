package domain

import "time"

type Integration struct {
	ID        string
	Key       string
	Name      string
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}
