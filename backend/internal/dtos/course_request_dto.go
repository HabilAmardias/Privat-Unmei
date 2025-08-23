package dtos

type (
	PreferredSlot struct {
		Date      string   `json:"date" binding:"required"`
		StartTime TimeOnly `json:"start_time" binding:"required"`
	}
	CreateCourseRequstReq struct {
		CourseID       int             `json:"course_id" binding:"required"`
		PreferredSlots []PreferredSlot `json:"preferred_slots" binding:"dive"`
	}
	CreateCourseRequestRes struct {
		CourseRequestID int `json:"id"`
	}
)
