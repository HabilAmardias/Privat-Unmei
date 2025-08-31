package entity

import (
	"time"
)

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
		SeekPaginatedParam
		MentorID string
		Status   *string
	}
	MentorCourseRequestDetailParam struct {
		CourseRequestID int
		MentorID        string
	}
	MentorCourseRequestDetailQuery struct {
		CourseRequestID  int
		CourseName       string
		StudentName      string
		StudentEmail     string
		TotalPrice       float64
		Subtotal         float64
		OperationalCost  float64
		NumberOfSessions int
		Status           string
		ExpiredAt        *time.Time
		Schedules        []CourseRequestSchedule
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
		SeekPaginatedParam
		StudentID string
		Status    *string
		Search    *string
	}
	StudentCourseRequestDetailParam struct {
		CourseRequestID int
		StudentID       string
	}
	StudentCourseRequestDetailQuery struct {
		CourseRequestID  int
		CourseName       string
		MentorName       string
		MentorEmail      string
		TotalPrice       float64
		Subtotal         float64
		OperationalCost  float64
		NumberOfSessions int
		Status           string
		ExpiredAt        *time.Time
		Schedules        []CourseRequestSchedule
	}
)
