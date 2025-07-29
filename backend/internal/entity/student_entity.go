package entity

import "time"

type (
	Student struct {
		ID        string
		CreatedAt time.Time
		UpdatedAt time.Time
		DeletedAt *time.Time
	}
)
