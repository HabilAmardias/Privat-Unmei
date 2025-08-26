package entity

import "time"

type (
	CourseRequest struct {
		ID               int
		StudentID        string
		CourseID         int
		Status           string
		TotalPrice       float64
		NumberOfSessions int
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
	HandleCourseRequestParam struct {
		MentorID        string
		CourseRequestID int
		Accept          bool
	}
	ConfirmPaymentParam struct {
		MentorID        string
		CourseRequestID int
	}
)
