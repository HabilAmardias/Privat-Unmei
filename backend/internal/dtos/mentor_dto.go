package dtos

type (
	TimeOnly struct {
		Hour   int `json:"hour" binding:"omitempty,gte=0,lt=24"`
		Minute int `json:"minute" binding:"omitempty,gte=0,lt=60"`
		Second int `json:"second" binding:"omitempty,gte=0,lt=60"`
	}
	MentorAvailabilityReq struct {
		DayOfWeek int      `json:"day_of_week" binding:"required,gte=1,lte=7"`
		StartTime TimeOnly `json:"start_time" binding:"required"`
		EndTime   TimeOnly `json:"end_time" binding:"required"`
	}
	MentorAvailabilityRes struct {
		DayOfWeek int      `json:"day_of_week"`
		StartTime TimeOnly `json:"start_time"`
		EndTime   TimeOnly `json:"end_time"`
	}
	AddNewMentorReq struct {
		Name              string   `form:"name" binding:"required"`
		Email             string   `form:"email" binding:"required,email"`
		Bio               string   `form:"bio" binding:"required"`
		Password          string   `form:"password" binding:"required,containsany=!@#?,min=8"`
		YearsOfExperience int      `form:"years_of_experience" binding:"required,gte=0"`
		GopayNumber       string   `form:"gopay_number" binding:"required"`
		Degree            string   `form:"degree" binding:"required"`
		Major             string   `form:"major" binding:"required"`
		Campus            string   `form:"campus" binding:"required"`
		MentorSchedules   []string `form:"mentor_availability"`
	}
	GeneratePasswordRes struct {
		Password string `json:"password"`
	}
	UpdateMentorForAdminReq struct {
		GopayNumber       *string `json:"gopay_number"`
		YearsOfExperience *int    `json:"years_of_experience"`
	}
	UpdateMentorForAdminRes struct {
		ID string `json:"id"`
	}
	UpdateMentorReq struct {
		Name              *string  `form:"name"`
		Bio               *string  `form:"bio"`
		YearsOfExperience *int     `form:"years_of_experience" binding:"omitempty,gte=0"`
		GopayNumber       *string  `form:"gopay_number"`
		Degree            *string  `form:"degree"`
		Major             *string  `form:"major"`
		Campus            *string  `form:"campus"`
		MentorSchedules   []string `form:"mentor_availability"`
	}
	UpdateMentorRes struct {
		ID string `json:"id"`
	}
	ListMentorReq struct {
		PaginatedReq
		Search               *string `form:"search"`
		SortYearOfExperience *bool   `form:"sort_year_of_experience"`
	}
	ListMentorRes struct {
		ID                string `json:"id"`
		Name              string `json:"name"`
		Email             string `json:"email"`
		GopayNumber       string `json:"gopay_number"`
		YearsOfExperience int    `json:"years_of_experience"`
	}
	LoginMentorReq struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,containsany=!@#?,min=8"`
	}
	LoginMentorRes struct {
		Token string `json:"token"`
	}
	MentorChangePasswordReq struct {
		NewPassword string `json:"password" binding:"required,containsany=!@#?,min=8"`
	}
	GetProfileMentorRes struct {
		ResumeFile           string                  `json:"resume_file"`
		ProfileImage         string                  `json:"profile_image"`
		Name                 string                  `json:"name"`
		Bio                  string                  `json:"bio"`
		YearsOfExperience    int                     `json:"years_of_experience"`
		GopayNumber          string                  `json:"gopay_number"`
		Degree               string                  `json:"degree"`
		Major                string                  `json:"major"`
		Campus               string                  `json:"campus"`
		MentorAvailabilities []MentorAvailabilityRes `json:"mentor_availability"`
	}
	GetMentorProfileForStudentRes struct {
		MentorID                string                  `json:"id"`
		MentorName              string                  `json:"name"`
		MentorEmail             string                  `json:"email"`
		MentorBio               string                  `json:"bio"`
		MentorProfileImage      string                  `json:"profile_image"`
		MentorAverageRating     float64                 `json:"rating"`
		MentorResume            string                  `json:"resume"`
		MentorYearsOfExperience int                     `json:"years_of_experience"`
		MentorDegree            string                  `json:"degree"`
		MentorMajor             string                  `json:"major"`
		MentorCampus            string                  `json:"campus"`
		MentorAvailabilities    []MentorAvailabilityRes `json:"mentor_availability"`
	}
)
