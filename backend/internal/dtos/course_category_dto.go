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
	CreateCategoryReq struct {
		Name string `json:"name" binding:"required"`
	}
	CreateCategoryRes struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}
)
