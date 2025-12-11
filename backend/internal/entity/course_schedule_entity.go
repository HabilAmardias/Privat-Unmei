package entity

import "time"

type (
	CourseRequestSchedule struct {
		ID              int
		CourseRequestID string
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
	ConflictingSchedule struct {
		Date       time.Time
		StartTime  string
		EndTime    string
		ScheduleID int64
	}
)
