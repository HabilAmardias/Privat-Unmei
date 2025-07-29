package dtos

type (
	RegisterStudentReq struct {
		Name     string  `form:"name" binding:"required"`
		Email    string  `form:"email" binding:"required,email"`
		Password string  `form:"password" binding:"required"`
		Bio      *string `form:"bio"`
	}
)
