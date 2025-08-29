package dtos

type (
	RegisterStudentReq struct {
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,containsany=!@#?,min=8"`
		Bio      string `json:"bio" binding:"required"`
	}
	LoginStudentReq struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,containsany=!@#?,min=8"`
	}

	LoginStudentRes struct {
		Token string `json:"token"`
	}
	VerifyStudentReq struct {
		Token string `json:"token" binding:"required"`
	}
	SendResetTokenEmailReq struct {
		Email string `json:"email" binding:"required,email"`
	}
	ResetPasswordReq struct {
		NewPassword string `json:"new_password" binding:"required,containsany=!@#?,min=8"`
		Token       string `json:"token"`
	}
	ListStudentRes struct {
		ID           string `json:"id"`
		Name         string `json:"name"`
		Email        string `json:"email"`
		Bio          string `json:"bio"`
		ProfileImage string `json:"profile_image"`
		Status       string `json:"status"`
	}
	UpdateStudentReq struct {
		Name *string `form:"name"`
		Bio  *string `form:"bio"`
	}
	UpdateStudentRes struct {
		ID string `json:"id"`
	}
)
