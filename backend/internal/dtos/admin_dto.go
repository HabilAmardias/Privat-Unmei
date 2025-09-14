package dtos

type (
	AdminLoginReq struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,containsany=!@#?,min=8"`
	}
	AdminLoginRes struct {
		Token  string `json:"token"`
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
)
