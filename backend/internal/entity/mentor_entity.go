package entity

import (
	"fmt"
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

func (to *TimeOnly) ToString() string {
	var hour string
	var minute string
	var second string

	hour = fmt.Sprintf("%d", to.Hour)
	if to.Hour < 10 {
		hour = fmt.Sprintf("0%d", to.Hour)
	}

	minute = fmt.Sprintf("%d", to.Minute)
	if to.Minute < 10 {
		minute = fmt.Sprintf("0%d", to.Minute)
	}

	second = fmt.Sprintf("%d", to.Second)
	if to.Second < 10 {
		second = fmt.Sprintf("0%d", to.Second)
	}

	return fmt.Sprintf("%s:%s:%s", hour, minute, second)
}
