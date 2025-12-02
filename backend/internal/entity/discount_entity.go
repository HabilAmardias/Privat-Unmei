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
	GetDiscountParam struct {
		UserID      string
		Participant int
	}
	CreateNewDiscountParam struct {
		AdminID             string
		NumberOfParticipant int
		Amount              float64
	}
	UpdateDiscountParam struct {
		AdminID    string
		DiscountID int
		Amount     *float64
	}
	DeleteDiscountParam struct {
		AdminID    string
		DiscountID int
	}
	GetAllDiscountParam struct {
		PaginatedParam
		AdminID string
	}
	GetDiscountQuery struct {
		ID                  int
		NumberOfParticipant int
		Amount              float64
	}
)
