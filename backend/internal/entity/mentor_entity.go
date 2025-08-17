package entity

import (
	"mime/multipart"
	"time"
)

type (
	TimeOnly struct {
		Hour   int
		Minute int
		Second int
	}
	MentorSchedule struct {
		DayOfWeek int
		StartTime TimeOnly
		EndTime   TimeOnly
	}
	Mentor struct {
		ID                string
		TotalRating       float64
		RatingCount       int
		Resume            string
		YearsOfExperience int
		GopayNumber       string
		Degree            string
		Major             string
		Campus            string
		CreatedAt         time.Time
		UpdatedAt         time.Time
		DeletedAt         *time.Time
	}
	MentorAvailability struct {
		ID        int
		MentorID  string
		DayOfWeek int
		StartTime TimeOnly
		EndTime   TimeOnly
		CreatedAt time.Time
		UpdatedAt time.Time
		DeletedAt *time.Time
	}
	AddNewMentorParam struct {
		Name              string
		Email             string
		Bio               string
		Password          string
		ResumeFile        multipart.File
		YearsOfExperience int
		GopayNumber       string
		Degree            string
		Major             string
		Campus            string
		MentorSchedules   []MentorSchedule
	}
	UpdateMentorParam struct {
		ID                string
		Resume            multipart.File
		ProfileImage      multipart.File
		Name              *string
		Password          *string
		Bio               *string
		YearsOfExperience *int
		GopayNumber       *string
		Degree            *string
		Major             *string
		Campus            *string
		MentorSchedules   []MentorSchedule
	}
	UpdateMentorQuery struct {
		TotalRating       *float64
		RatingCount       *int
		Resume            *string
		YearsOfExperience *int
		GopayNumber       *string
		Degree            *string
		Major             *string
		Campus            *string
	}
	DeleteMentorParam struct {
		ID string
	}
	ListMentorQuery struct {
		ID                string
		Name              string
		Email             string
		GopayNumber       string
		YearsOfExperience int
	}
	ListMentorParam struct {
		PaginatedParam
		Search               *string
		SortYearOfExperience *bool
	}
	LoginMentorParam struct {
		Email    string
		Password string
	}
	MentorChangePasswordParam struct {
		ID          string
		NewPassword string
	}
)
