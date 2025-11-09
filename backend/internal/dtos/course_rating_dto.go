package dtos

import "time"

type (
	CreateRatingReq struct {
		Rating   int     `json:"rating" binding:"required,min=1,max=5"`
		Feedback *string `json:"feedback" binding:"omitempty,min=15"`
	}
	CreateRatingRes struct {
		RatingID int `json:"id"`
	}
	CourseRatingRes struct {
		ID          int       `json:"id"`
		CourseID    int       `json:"course_id"`
		StudentID   string    `json:"student_id"`
		StudentName string    `json:"name"`
		Rating      int       `json:"rating"`
		Feedback    *string   `json:"feedback"`
		CreatedAt   time.Time `json:"created_at"`
	}
	CourseRatingReq struct {
		PaginatedReq
	}
)
