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
		MentorID    string
		Title       string
		Description string
		Domicile    string
		MinPrice    float64
		MaxPrice    float64
		Method      string
		MinDuration int
		MaxDuration int
		Topics      []CreateTopic
		Categories  []int
	}
	DeleteCourseParam struct {
		MentorID string
		CourseID int
	}
	MentorListCourseParam struct {
		SeekPaginatedParam
		MentorID       string
		Search         *string
		CourseCategory *int
	}
	ListCourseParam struct {
		SeekPaginatedParam
		Search         *string
		CourseCategory *int
		Method         *string
	}
	MentorListCourseQuery struct {
		ID               int
		Title            string
		Domicile         string
		Method           string
		MinPrice         float64
		MaxPrice         float64
		MinDurationDays  int
		MaxDurationDays  int
		CourseCategories string
	}
	CourseListQuery struct {
		MentorListCourseQuery
		MentorID    string
		MentorName  string
		MentorEmail string
	}
	CourseDetailQuery struct {
		CourseListQuery
		Description string
		Topics      *[]CourseTopic
	}
	CourseDetailParam struct {
		ID int
	}
	UpdateCourseQuery struct {
		Title           *string
		Description     *string
		Domicile        *string
		MinPrice        *float64
		MaxPrice        *float64
		Method          *string
		MinDurationDays *int
		MaxDurationDays *int
	}
	UpdateCourseParam struct {
		MentorID string
		CourseID int
		UpdateCourseQuery
		CourseTopic      []CreateTopic
		CourseCategories []int
	}
)
