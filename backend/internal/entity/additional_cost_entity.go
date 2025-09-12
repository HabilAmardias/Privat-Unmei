package entity

import "time"

type (
	AdditionalCost struct {
		ID        int
		Name      string
		Amount    float64
		CreatedAt time.Time
		UpdatedAt time.Time
		DeletedAt *time.Time
	}
)
