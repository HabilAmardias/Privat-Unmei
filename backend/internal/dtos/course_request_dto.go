package dtos

import "time"

type (
	PreferredSlot struct {
		Date      string   `json:"date" binding:"required"`
		StartTime TimeOnly `json:"start_time" binding:"required"`
	}
	CreateCourseRequestReq struct {
		PreferredSlots      []PreferredSlot `json:"preferred_slots" binding:"dive"`
		PaymentMethodID     int             `json:"payment_method_id" binding:"required"`
		NumberOfParticipant int             `json:"number_of_participant" binding:"required,gte=1"`
	}
	CreateCourseRequestRes struct {
		CourseRequestID int `json:"id"`
	}
	HandleCourseRequestReq struct {
		Accept *bool `json:"accept" binding:"required"`
	}
	PaymentDetailRes struct {
		CourseRequestID int        `json:"id"`
		StudentName     string     `json:"student_name"`
		MentorName      string     `json:"mentor_name"`
		CourseID        int        `json:"course_id"`
		CourseTitle     string     `json:"course_title"`
		PaymentMethod   string     `json:"payment_method"`
		AccountNumber   string     `json:"account_number"`
		Subtotal        float64    `json:"subtotal"`
		OperationalCost float64    `json:"operational_cost"`
		TotalCost       float64    `json:"total_cost"`
		ExpiredAt       *time.Time `json:"expired_at"`
	}
	MentorCourseRequestListReq struct {
		PaginatedReq
		Status *string `form:"status"`
	}
	MentorCourseRequestRes struct {
		ID         int     `json:"id"`
		StudentID  string  `json:"student_id"`
		CourseID   int     `json:"course_id"`
		TotalPrice float64 `json:"total_price"`
		Status     string  `json:"status"`
		Name       string  `json:"name"`
		Email      string  `json:"email"`
		CourseName string  `json:"course_name"`
	}
	MentorCourseRequestDetailRes struct {
		CourseRequestID     int                 `json:"course_request_id"`
		CourseName          string              `json:"course_name"`
		StudentID           string              `json:"student_id"`
		StudentName         string              `json:"student_name"`
		StudentEmail        string              `json:"student_email"`
		TotalPrice          float64             `json:"total_price"`
		Subtotal            float64             `json:"subtotal"`
		OperationalCost     float64             `json:"operational_cost"`
		NumberOfSessions    int                 `json:"number_of_sessions"`
		Status              string              `json:"status"`
		ExpiredAt           *time.Time          `json:"expired_at"`
		PaymentMethod       string              `json:"payment_method"`
		AccountNumber       string              `json:"account_number"`
		NumberOfParticipant int                 `json:"number_of_participant"`
		Schedules           []CourseScheduleRes `json:"schedules"`
	}
	StudentCourseRequestListReq struct {
		PaginatedReq
		Status *string `form:"status"`
		Search *string `form:"search"`
	}
	StudentCourseRequestRes struct {
		ID          int     `json:"id"`
		StudentID   string  `json:"student_id"`
		CourseID    int     `json:"course_id"`
		TotalPrice  float64 `json:"total_price"`
		Status      string  `json:"status"`
		MentorName  string  `json:"mentor_name"`
		MentorEmail string  `json:"mentor_email"`
		CourseName  string  `json:"course_name"`
	}
	StudentCourseRequestDetailRes struct {
		CourseRequestID     int                 `json:"course_request_id"`
		CourseName          string              `json:"course_name"`
		MentorName          string              `json:"mentor_name"`
		MentorEmail         string              `json:"mentor_email"`
		TotalPrice          float64             `json:"total_price"`
		Subtotal            float64             `json:"subtotal"`
		OperationalCost     float64             `json:"operational_cost"`
		PaymentMethodName   string              `json:"payment_method"`
		AccountNumber       string              `json:"account_number"`
		NumberOfSessions    int                 `json:"number_of_sessions"`
		Status              string              `json:"status"`
		ExpiredAt           *time.Time          `json:"expired_at"`
		NumberOfParticipant int                 `json:"number_of_participant"`
		Schedules           []CourseScheduleRes `json:"schedules"`
	}
	OperationalCostRes struct {
		Cost float64 `json:"operational_cost"`
	}
)
