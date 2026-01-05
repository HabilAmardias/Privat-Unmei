package entity

import (
	"time"
)

type (
	CourseRequest struct {
		ID                  string
		StudentID           string
		CourseID            int
		Status              string
		NumberOfParticipant int
		NumberOfSessions    int
		ExpiredAt           *time.Time
		CreatedAt           time.Time
		UpdatedAt           time.Time
		DeletedAt           *time.Time
	}
	PreferredSlot struct {
		Date      time.Time
		StartTime TimeOnly
	}
	CreateCourseRequestParam struct {
		CourseID            int
		StudentID           string
		NumberOfParticipant int
		PreferredSlots      []PreferredSlot
		PaymentMethodID     int
	}
	HandleCourseRequestParam struct {
		MentorID        string
		CourseRequestID string
		Accept          bool
	}
	ConfirmPaymentParam struct {
		MentorID        string
		CourseRequestID string
	}
	GetPaymentDetailParam struct {
		UserID          string
		CourseRequestID string
	}
	PaymentDetailQuery struct {
		CourseRequestID string
		StudentName     string
		MentorName      string
		CourseID        int
		CourseTitle     string
		PaymentMethod   string
		AccountNumber   string
		Subtotal        float64
		OperationalCost float64
		TotalCost       float64
		ExpiredAt       *time.Time
	}
	MentorCourseRequestQuery struct {
		ID         string
		StudentID  string
		CourseID   int
		TotalPrice float64
		Status     string
		Name       string
		PublicID   string
		CourseName string
	}
	MentorCourseRequestListParam struct {
		PaginatedParam
		MentorID string
		Status   *string
	}
	MentorCourseRequestDetailParam struct {
		CourseRequestID string
		MentorID        string
	}
	MentorCourseRequestDetailQuery struct {
		CourseRequestID     string
		CourseName          string
		StudentID           string
		StudentName         string
		StudentPublicID     string
		NumberOfParticipant int
		TotalPrice          float64
		Subtotal            float64
		OperationalCost     float64
		PaymentMethod       string
		AccountNumber       string
		NumberOfSessions    int
		Status              string
		ExpiredAt           *time.Time
		Schedules           []CourseRequestSchedule
	}
	StudentCourseRequestQuery struct {
		ID             string
		StudentID      string
		CourseID       int
		TotalPrice     float64
		Status         string
		MentorName     string
		MentorPublicID string
		CourseName     string
	}
	StudentCourseRequestListParam struct {
		PaginatedParam
		StudentID string
		Status    *string
		Search    *string
	}
	StudentCourseRequestDetailParam struct {
		CourseRequestID string
		StudentID       string
	}
	StudentCourseRequestDetailQuery struct {
		CourseRequestID     string
		CourseName          string
		CourseID            int
		MentorName          string
		MentorID            string
		MentorPublicID      string
		TotalPrice          float64
		Subtotal            float64
		OperationalCost     float64
		PaymentMethodName   string
		AccountNumber       string
		NumberOfSessions    int
		Status              string
		NumberOfParticipant int
		ExpiredAt           *time.Time
		Schedules           []CourseRequestSchedule
	}
	IsReviewedParam struct {
		CourseID  int
		StudentID string
	}
)
