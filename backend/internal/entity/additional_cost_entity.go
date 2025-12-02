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
	GetOperationalCostParam struct {
		UserID string
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
	DeleteAdditionalCostParam struct {
		CostID  int
		AdminID string
	}
	GetAdditionalCostQuery struct {
		ID     int
		Name   string
		Amount float64
	}
	GetAllAdditionalCostParam struct {
		PaginatedParam
		AdminID string
	}
)
