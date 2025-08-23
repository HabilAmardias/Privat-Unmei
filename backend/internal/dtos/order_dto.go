package dtos

type (
	PreferredSlot struct {
		Date      string   `json:"date" binding:"required"`
		StartTime TimeOnly `json:"start_time" binding:"required"`
	}
	CreateOrderReq struct {
		CourseID       int             `json:"course_id" binding:"required"`
		PreferredSlots []PreferredSlot `json:"preferred_slots" binding:"dive"`
	}
	CreateOrderRes struct {
		CourseRequestID int64 `json:"id"`
	}
)
