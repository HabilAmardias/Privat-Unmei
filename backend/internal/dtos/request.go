package dtos

type (
	PaginatedReq struct {
		Page  int `form:"page" binding:"omitempty"`
		Limit int `form:"limit" binding:"omitempty"`
	}
	SeekPaginatedReq struct {
		Limit  int  `form:"limit"`
		LastID *int `form:"last_id"`
	}
)
