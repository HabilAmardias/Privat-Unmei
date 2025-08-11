package entity

import "time"

type (
	CourseRating struct {
		ID        int
		CourseID  int
		StudentID string
		Rating    int
		Feedback  string
		CreatedAt time.Time
		UpdatedAt time.Time
		DeletedAt *time.Time
	}
	CreateRatingParam struct {
		StudentID string
		CourseID  int
		Rating    int
		Feedback  *string
	}
)
