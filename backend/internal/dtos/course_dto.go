package dtos

type (
	CourseTopicReq struct {
		Title       string `json:"title" binding:"required"`
		Description string `json:"description" binding:"required"`
	}
	CourseTopicRes struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}
	CreateCourseReq struct {
		Title           string           `json:"title" binding:"required"`
		Description     string           `json:"description" binding:"required"`
		Domicile        string           `json:"domicile" binding:"required"`
		Price           float64          `json:"price" binding:"required,min=1"`
		Method          string           `json:"method" binding:"required"`
		SessionDuration int              `json:"session_duration_minutes" binding:"required,min=1"`
		MaxSession      int              `json:"max_total_session" binding:"required,min=1"`
		Topics          []CourseTopicReq `json:"course_topics" binding:"dive"`
		Categories      []int            `json:"course_categories"`
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
	GetCategoriesRes struct {
		CategoryID   int    `json:"category_id"`
		CategoryName string `json:"category_name"`
	}
	MentorListCourseRes struct {
		ID              int     `json:"id"`
		Title           string  `json:"title"`
		Domicile        string  `json:"domicile"`
		Method          string  `json:"method"`
		Price           float64 `json:"price"`
		SessionDuration int     `json:"session_duration_minutes"`
		MaxSession      int     `json:"max_total_session"`
	}
	CourseListRes struct {
		MentorListCourseRes
		MentorID    string `json:"mentor_id"`
		MentorName  string `json:"mentor_name"`
		MentorEmail string `json:"mentor_email"`
	}
	CourseDetailRes struct {
		CourseListRes
		Description      string             `json:"description"`
		Topics           []CourseTopicRes   `json:"topics"`
		CourseCategories []GetCategoriesRes `json:"course_categories"`
	}
	UpdateCourseReq struct {
		Title            *string          `json:"title"`
		Description      *string          `json:"description"`
		Domicile         *string          `json:"domicile"`
		Price            *float64         `json:"price" binding:"omitempty,min=1"`
		Method           *string          `json:"method"`
		SessionDuration  *int             `json:"session_duration_minutes" binding:"omitempty,min=1"`
		MaxSession       *int             `json:"max_total_session" binding:"omitempty,min=1,max=7"`
		CourseTopic      []CourseTopicReq `json:"course_topics"`
		CourseCategories []int            `json:"course_categories"`
	}
	UpdateCourseRes struct {
		ID int `json:"id"`
	}
)
