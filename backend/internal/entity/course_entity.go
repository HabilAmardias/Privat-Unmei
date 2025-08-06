package entity

import "time"

type (
	TimeOnly struct {
		Hour   int
		Minute int
		Second int
	}
	CreateSchedule struct {
		DayOfWeek string
		StartTime TimeOnly
		EndTime   TimeOnly
	}
	CreateTopic struct {
		Title       string
		Description string
	}
	Course struct {
		ID               int
		MentorID         string
		Title            string
		Description      string
		Domicile         string
		MinPrice         float64
		MaxPrice         float64
		MinDuration      int
		MaxDuration      int
		Method           string
		TransactionCount int
		CreatedAt        time.Time
		UpdatedAt        time.Time
		DeletedAt        *time.Time
	}
	CourseAvailability struct {
		ID        int
		CourseID  int
		DayOfWeek string
		StartTime TimeOnly
		EndTime   TimeOnly
		CreatedAt time.Time
		UpdatedAt time.Time
		DeletedAt *time.Time
	}
	CourseTopic struct {
		ID          int
		CourseID    int
		Title       string
		Description string
		CreatedAt   time.Time
		UpdatedAt   time.Time
		DeletedAt   *time.Time
	}
	CreateCourseParam struct {
		MentorID           string
		Title              string
		Description        string
		Domicile           string
		MinPrice           float64
		MaxPrice           float64
		Method             string
		MinDuration        int
		MaxDuration        int
		CourseAvailability []CreateSchedule
		Topics             []CreateTopic
	}
	DeleteCourseParam struct {
		MentorID string
		CourseID int
	}
)
