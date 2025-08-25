package entity

import "time"

type (
	CourseRequestSchedule struct {
		ID              int
		CourseRequestID int
		ScheduledDate   time.Time
		StartTime       string
		EndTime         string
		Status          string
		CreatedAt       time.Time
		UpdatedAt       time.Time
		DeletedAt       *time.Time
	}
	CreateRequestSchedule struct {
		Date      time.Time
		StartTime string
		EndTime   string
	}
)
