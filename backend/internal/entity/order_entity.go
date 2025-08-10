package entity

import "time"

type (
	CourseOrder struct {
		ID           int
		StudentID    string
		CourseID     int
		Status       int
		Price        float64
		DurationDays int
		AcceptedAt   *time.Time
		PaymentDue   *time.Time
		StartDate    *time.Time
		EndDate      *time.Time
		CreatedAt    time.Time
		UpdatedAt    time.Time
		DeletedAt    *time.Time
	}
)
