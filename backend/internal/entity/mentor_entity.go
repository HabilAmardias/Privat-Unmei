package entity

import (
	"mime/multipart"
	"time"
)

type (
	Mentor struct {
		ID                string
		TotalRating       float64
		RatingCount       int
		Resume            string
		YearsOfExperience int
		WhatsappNumber    string
		Degree            string
		Major             string
		Campus            string
		CreatedAt         time.Time
		UpdatedAt         time.Time
		DeletedAt         *time.Time
	}
	AddNewMentorParam struct {
		Name              string
		Email             string
		Bio               string
		Password          string
		ResumeFile        multipart.File
		YearsOfExperience int
		WhatsappNumber    string
		Degree            string
		Major             string
		Campus            string
	}
	UpdateMentorParam struct {
		ID string
		UpdateMentorQuery
	}
	UpdateMentorQuery struct {
		TotalRating       *float64
		RatingCount       *int
		Resume            *string
		YearsOfExperience *int
		WhatsappNumber    *string
		Degree            *string
		Major             *string
		Campus            *string
	}
	DeleteMentorParam struct {
		ID string
	}
)
