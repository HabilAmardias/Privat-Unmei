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
		YearsOfExperience int
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
	AddMentorPaymentInfo struct {
		PaymentMethodID int
		AccountNumber   string
	}
	MentorPaymentInfo struct {
		PaymentMethodID   int
		PaymentMethodName string
		AccountNumber     string
	}
	MentorProfileParam struct {
		ID string
	}
	MentorProfileQuery struct {
		ID                string
		Name              string
		PublicID          string
		Bio               string
		ProfileImage      string
		YearsOfExperience int
		Degree            string
		Major             string
		Campus            string
	}
	AddNewMentorParam struct {
		AdminID           string
		Name              string
		Email             string
		Password          string
		MentorPayments    []AddMentorPaymentInfo
		YearsOfExperience int
		Degree            string
		Major             string
		Campus            string
		MentorSchedules   []MentorSchedule
	}
	UpdateMentorParam struct {
		ID                string
		ProfileImage      multipart.File
		Name              *string
		Bio               *string
		YearsOfExperience *int
		Degree            *string
		Major             *string
		Campus            *string
		MentorPayments    []AddMentorPaymentInfo
		MentorSchedules   []MentorSchedule
	}
	UpdateMentorQuery struct {
		TotalRating       *float64
		RatingCount       *int
		YearsOfExperience *int
		Degree            *string
		Major             *string
		Campus            *string
	}
	DeleteMentorParam struct {
		ID      string
		AdminID string
	}
	ListMentorQuery struct {
		ID                string
		Name              string
		PublicID          string
		ProfileImage      string
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
	GetProfileMentorQuery struct {
		ProfileImage         string
		Name                 string
		Bio                  string
		YearsOfExperience    int
		Degree               string
		Major                string
		Campus               string
		MentorAvailabilities []MentorSchedule
		MentorPayments       []MentorPaymentInfo
	}
	AvailabilityResult struct {
		TotalRequested   int
		AvailableSlots   int
		UnavailableSlots []string
	}
	GetDOWAvailabilityParam struct {
		Role     int
		CourseID int
		UserID   string
	}
	GetMentorAvailabilityParam struct {
		MentorID string
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
