package entity

import "time"

type (
	CourseRequest struct {
		ID               int
		StudentID        string
		CourseID         int
		Status           string
		SubTotal         float64
		OperationalCost  float64
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
	GetPaymentDetailParam struct {
		UserID          string
		CourseRequestID int
	}
	PaymentDetailQuery struct {
		CourseRequestID int
		MentorID        string
		MentorName      string
		GopayNumber     string
		CourseID        int
		CourseTitle     string
		Subtotal        float64
		OperationalCost float64
		TotalCost       float64
		ExpiredAt       time.Time
	}
)
