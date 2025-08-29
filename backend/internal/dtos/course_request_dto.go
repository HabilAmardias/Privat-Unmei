package dtos

import "time"

type (
	PreferredSlot struct {
		Date      string   `json:"date" binding:"required"`
		StartTime TimeOnly `json:"start_time" binding:"required"`
	}
	CreateCourseRequstReq struct {
		PreferredSlots []PreferredSlot `json:"preferred_slots" binding:"dive"`
	}
	CreateCourseRequestRes struct {
		CourseRequestID int `json:"id"`
	}
	HandleCourseRequestReq struct {
		Accept *bool `json:"accept" binding:"required"`
	}
	PaymentDetailRes struct {
		CourseRequestID int       `json:"id"`
		MentorID        string    `json:"mentor_id"`
		MentorName      string    `json:"mentor_name"`
		GopayNumber     string    `json:"gopay_number"`
		CourseID        int       `json:"course_id"`
		CourseTitle     string    `json:"course_title"`
		Subtotal        float64   `json:"subtotal"`
		OperationalCost float64   `json:"operational_cost"`
		TotalCost       float64   `json:"total_cost"`
		ExpiredAt       time.Time `json:"expired_at"`
	}
)
