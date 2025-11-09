package dtos

type (
	ListCourseCategoryReq struct {
		PaginatedReq
		Search *string `form:"search" binding:"omitempty"`
	}
	ListCourseCategoryRes struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}
	CreateCategoryReq struct {
		Name string `json:"name" binding:"required"`
	}
	UpdateCategoryReq struct {
		Name *string `json:"name"`
	}
	CategoryIDRes struct {
		ID int `json:"id"`
	}
)
