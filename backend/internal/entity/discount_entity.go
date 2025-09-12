package entity

import "time"

type (
	Discount struct {
		ID                  int
		NumberOfParticipant int
		Amount              float64
		CreatedAt           time.Time
		UpdatedAt           time.Time
		DeletedAt           *time.Time
	}
)
