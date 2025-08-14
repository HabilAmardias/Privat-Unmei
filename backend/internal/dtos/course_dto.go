package dtos

type (
	TimeOnly struct {
		Hour   int `json:"hour" binding:"omitempty,min=0,max=24"`
		Minute int `json:"minute" binding:"omitempty,min=0,max=60"`
		Second int `json:"second" binding:"omitempty,min=0,max=60"`
	}
	CourseAvailabilityReq struct {
		DayOfWeek string   `json:"day_of_week" binding:"required"`
		StartTime TimeOnly `json:"start_time" binding:"required"`
		EndTime   TimeOnly `json:"end_time" binding:"required"`
	}
	CourseAvailabilityRes struct {
		DayOfWeek string   `json:"day_of_week"`
		StartTime TimeOnly `json:"start_time"`
		EndTime   TimeOnly `json:"end_time"`
	}
	CourseTopicReq struct {
		Title       string `json:"title" binding:"required"`
		Description string `json:"description" binding:"required"`
	}
	CourseTopicRes struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}
	CreateCourseReq struct {
		Title              string                  `json:"title" binding:"required"`
		Description        string                  `json:"description" binding:"required"`
		Domicile           string                  `json:"domicile" binding:"required"`
		MinPrice           float64                 `json:"min_price" binding:"required,min=1"`
		MaxPrice           float64                 `json:"max_price" binding:"required,min=1"`
		Method             string                  `json:"method" binding:"required"`
		MinDuration        int                     `json:"min_duration_days" binding:"required,min=1"`
		MaxDuration        int                     `json:"max_duration_days" binding:"required,min=1"`
		CourseAvailability []CourseAvailabilityReq `json:"available_schedules" binding:"dive"`
		Topics             []CourseTopicReq        `json:"course_topics" binding:"dive"`
		Categories         []int                   `json:"course_categories"`
	}
	CreateCourseRes struct {
		ID int `json:"id"`
	}
	DeleteCourseRes struct {
		ID int `json:"id"`
	}
	MentorListCourseReq struct {
		SeekPaginatedReq
		Search         *string `form:"search"`
		CourseCategory *int    `form:"course_category"`
	}
	ListCourseReq struct {
		SeekPaginatedReq
		Search         *string `form:"search"`
		CourseCategory *int    `form:"course_category"`
		Method         *string `form:"method"`
	}
	MentorListCourseRes struct {
		ID               int      `json:"id"`
		Title            string   `json:"title"`
		Domicile         string   `json:"domicile"`
		Method           string   `json:"method"`
		MinPrice         float64  `json:"min_price"`
		MaxPrice         float64  `json:"max_price"`
		MinDurationDays  int      `json:"min_duration_days"`
		MaxDurationDays  int      `json:"max_duration_days"`
		CourseCategories []string `json:"course_categories"`
	}
	CourseListRes struct {
		MentorListCourseRes
		MentorID    string `json:"mentor_id"`
		MentorName  string `json:"mentor_name"`
		MentorEmail string `json:"mentor_email"`
	}
	CourseDetailRes struct {
		CourseListRes
		Description  string                  `json:"description"`
		Topics       []CourseTopicRes        `json:"topics"`
		Availability []CourseAvailabilityRes `json:"course_availability"`
	}
	UpdateCourseReq struct {
		Title            *string                 `json:"title"`
		Description      *string                 `json:"description"`
		Domicile         *string                 `json:"domicile"`
		MinPrice         *float64                `json:"min_price" binding:"omitempty,min=1"`
		MaxPrice         *float64                `json:"max_price" binding:"omitempty,min=1"`
		Method           *string                 `json:"method"`
		MinDurationDays  *int                    `json:"min_duration_days" binding:"omitempty,min=1"`
		MaxDurationDays  *int                    `json:"max_duration_days" binding:"omitempty,min=1"`
		CourseSchedule   []CourseAvailabilityReq `json:"course_availabilities"`
		CourseTopic      []CourseTopicReq        `json:"course_topics"`
		CourseCategories []int                   `json:"course_categories"`
	}
	UpdateCourseRes struct {
		ID int `json:"id"`
	}
)
