package dtos

type (
	AdminLoginReq struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,containsany=!@#?,min=8"`
	}
	AdminLoginRes struct {
		Status string `json:"status"`
	}
	AdminVerificationReq struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,containsany=!@#?,min=8"`
	}
	AdminUpdatePasswordReq struct {
		Password string `json:"password" binding:"required,containsany=!@#?,min=8"`
	}
	GetStudentListReq struct {
		PaginatedReq
	}
	AdminIDRes struct {
		ID string `json:"id"`
	}
	AdminProfileRes struct {
		Name         string `json:"name"`
		Email        string `json:"email"`
		Bio          string `json:"bio"`
		ProfileImage string `json:"profile_image"`
		Status       string `json:"status"`
	}
)
