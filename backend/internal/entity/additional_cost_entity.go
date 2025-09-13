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
	CreateAdditionalCostParam struct {
		Name    string
		Amount  float64
		AdminID string
	}
	UpdateAdditonalCostParam struct {
		Amount  *float64
		CostID  int
		AdminID string
	}
)
