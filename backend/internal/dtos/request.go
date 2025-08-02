package dtos

type (
	PaginatedReq struct {
		Page  int `form:"page" binding:"omitempty"`
		Limit int `form:"limit" binding:"omitempty"`
	}
)
