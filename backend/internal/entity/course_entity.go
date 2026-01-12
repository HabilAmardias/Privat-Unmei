package entity

import (
	"time"
)

type (
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
		Method           string
		Price            float64
		SessionDuration  int
		MaxSession       int
		TransactionCount int
		CreatedAt        time.Time
		UpdatedAt        time.Time
		DeletedAt        *time.Time
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
		MentorID        string
		Title           string
		Description     string
		Domicile        string
		Price           float64
		Method          string
		SessionDuration int
		MaxSession      int
		Topics          []CreateTopic
		Categories      []int
	}
	DeleteCourseParam struct {
		MentorID string
		CourseID int
	}
	MentorListCourseParam struct {
		PaginatedParam
		MentorID       string
		IsProtected    bool
		Search         *string
		CourseCategory *int
	}
	ListCourseParam struct {
		PaginatedParam
		Search         *string
		CourseCategory *int
		Method         *string
	}
	MentorListCourseQuery struct {
		ID              int
		Title           string
		Domicile        string
		Method          string
		Price           float64
		SessionDuration int
		MaxSession      int
	}
	CourseListQuery struct {
		MentorListCourseQuery
		MentorID           string
		MentorName         string
		MentorPublicID     string
		MentorProfileImage string
	}
	CourseDetailQuery struct {
		CourseListQuery
		Description string
	}
	CourseDetailParam struct {
		ID int
	}
	UpdateCourseQuery struct {
		Title            *string
		Description      *string
		Domicile         *string
		Method           *string
		Price            *float64
		SessionDuration  *int
		MaxSession       *int
		TransactionCount *int
	}
	UpdateCourseParam struct {
		MentorID string
		CourseID int
		UpdateCourseQuery
		CourseTopic      []CreateTopic
		CourseCategories []int
	}
)
