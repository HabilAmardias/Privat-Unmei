package dtos

type (
	TimeOnly struct {
		Hour   int `json:"hour" binding:"omitempty,min=0,max=24"`
		Minute int `json:"minute" binding:"omitempty,min=0,max=60"`
		Second int `json:"second" binding:"omitempty,min=0,max=60"`
	}
	CreateSchedule struct {
		DayOfWeek string   `json:"day_of_week" binding:"required"`
		StartTime TimeOnly `json:"start_time" binding:"required"`
		EndTime   TimeOnly `json:"end_time" binding:"required"`
	}
	CreateTopic struct {
		Title       string `json:"title" binding:"required"`
		Description string `json:"description" binding:"required"`
	}
	CreateCourseReq struct {
		Title              string           `json:"title" binding:"required"`
		Description        string           `json:"description" binding:"required"`
		Domicile           string           `json:"domicile" binding:"required"`
		MinPrice           float64          `json:"min_price" binding:"required,min=1"`
		MaxPrice           float64          `json:"max_price" binding:"required,min=1"`
		Method             string           `json:"method" binding:"required"`
		MinDuration        int              `json:"min_duration_days" binding:"required,min=1"`
		MaxDuration        int              `json:"max_duration_days" binding:"required,min=1"`
		CourseAvailability []CreateSchedule `json:"available_schedules" binding:"dive"`
		Topics             []CreateTopic    `json:"course_topics" binding:"dive"`
	}
	CreateCourseRes struct {
		ID int `json:"id"`
	}
)
