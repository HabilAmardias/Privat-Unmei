package dtos

import "time"

type (
	PreferredSlotReq struct {
		Date      time.Time `json:"date" binding:"required"`
		StartTime TimeOnly  `json:"start_time" binding:"required"`
	}
	CreateCourseScheduleReq struct {
		PreferredSlots []PreferredSlotReq `json:"preferred_slots" binding:"dive"`
	}
)
