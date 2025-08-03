package entity

import "time"

type (
	CourseCategory struct {
		ID        int
		Name      string
		CreatedAt time.Time
		UpdatedAt time.Time
		DeletedAt *time.Time
	}
	CourseCategoryAssignment struct {
		CourseID   int
		CategoryID int
		CreatedAt  time.Time
		UpdatedAt  time.Time
		DeletedAt  *time.Time
	}
	ListCourseCategoryParam struct {
		SeekPaginatedParam
		Search *string
	}
	ListCourseCategoryQuery struct {
		ID   int
		Name string
	}
	CreateCategoryQuery struct {
		ID   int
		Name string
	}
	CreateCategoryParam struct {
		Name string
	}
	UpdateCategoryParam struct {
		ID   int
		Name *string
	}
)
