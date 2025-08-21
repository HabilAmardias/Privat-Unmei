package entity

import "time"

type (
	CourseOrder struct {
		ID               int
		StudentID        string
		CourseID         int
		Status           int
		TotalPrice       float64
		NumberOfSessions int
		AcceptedAt       *time.Time
		PaymentDue       *time.Time
		CreatedAt        time.Time
		UpdatedAt        time.Time
		DeletedAt        *time.Time
	}
)
