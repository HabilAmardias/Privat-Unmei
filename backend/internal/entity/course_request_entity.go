package entity

import (
	"time"
)

type (
	CourseRequest struct {
		ID                  int
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
		ID         int
		StudentID  string
		CourseID   int
		TotalPrice float64
		Status     string
		Name       string
		Email      string
		CourseName string
	}
	MentorCourseRequestListParam struct {
		PaginatedParam
		MentorID string
		Status   *string
	}
	MentorCourseRequestDetailParam struct {
		CourseRequestID int
		MentorID        string
	}
	MentorCourseRequestDetailQuery struct {
		CourseRequestID     int
		CourseName          string
		StudentName         string
		StudentEmail        string
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
		ID          int
		StudentID   string
		CourseID    int
		TotalPrice  float64
		Status      string
		MentorName  string
		MentorEmail string
		CourseName  string
	}
	StudentCourseRequestListParam struct {
		PaginatedParam
		StudentID string
		Status    *string
		Search    *string
	}
	StudentCourseRequestDetailParam struct {
		CourseRequestID int
		StudentID       string
	}
	StudentCourseRequestDetailQuery struct {
		CourseRequestID     int
		CourseName          string
		MentorName          string
		MentorEmail         string
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
)
