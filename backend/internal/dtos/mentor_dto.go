package dtos

type (
	AddNewMentorReq struct {
		Name              string `form:"name" binding:"required"`
		Email             string `form:"email" binding:"required,email"`
		Bio               string `form:"bio" binding:"required"`
		Password          string `form:"password" binding:"required,containsany=!@#?,min=8"`
		YearsOfExperience int    `form:"years_of_experience" binding:"required,min=0"`
		WhatsappNumber    string `form:"whatsapp_number" binding:"required"`
		Degree            string `form:"degree" binding:"required"`
		Major             string `form:"major" binding:"required"`
		Campus            string `form:"campus" binding:"required"`
	}
	GeneratePasswordRes struct {
		Password string `json:"password"`
	}
	UpdateMentorForAdminReq struct {
		WhatsappNumber    *string `json:"whatsapp_number"`
		YearsOfExperience *int    `json:"years_of_experience"`
	}
	UpdateMentorForAdminRes struct {
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
		WhatsappNumber    string `json:"whatsapp_number"`
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
)
