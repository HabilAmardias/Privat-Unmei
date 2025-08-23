package entity

import "time"

type (
	CourseRequestSchedule struct {
		ID              int
		CourseRequestID int
		SessionNumber   int
		ScheduledDate   time.Time
		StartTime       TimeOnly
		EndTime         TimeOnly
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
