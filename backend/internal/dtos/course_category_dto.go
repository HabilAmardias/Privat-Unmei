package dtos

type (
	ListCourseCategoryReq struct {
		SeekPaginatedReq
		Search *string `form:"search" binding:"omitempty"`
	}
	ListCourseCategoryRes struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}
)
