package entity

import "time"

type (
	CourseRequest struct {
		ID               int
		StudentID        string
		CourseID         int
		Status           int
		TotalPrice       float64
		NumberOfSessions int
		AcceptedAt       *time.Time
		PaymentDue       *time.Time
		ExpiredAt        *time.Time
		CreatedAt        time.Time
		UpdatedAt        time.Time
		DeletedAt        *time.Time
	}
	PreferredSlot struct {
		Date      time.Time
		StartTime TimeOnly
	}
	CreateCourseRequestParam struct {
		CourseID       int
		StudentID      string
		PreferredSlots []PreferredSlot
	}
)
