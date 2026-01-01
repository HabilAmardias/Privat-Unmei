package entity

import (
	"mime/multipart"
	"time"
)

type (
	Student struct {
		ID          string
		VerifyToken *string
		ResetToken  *string
		CreatedAt   time.Time
		UpdatedAt   time.Time
		DeletedAt   *time.Time
	}
	StudentRegisterParam struct {
		Name     string
		Email    string
		Password string
		Status   string
	}
	StudentLoginParam struct {
		Email    string
		Password string
	}
	ResetPasswordParam struct {
		NewPassword string
		ID          string
		Token       string
	}
	VerifyStudentParam struct {
		Token string
		ID    string
	}
	ListStudentQuery struct {
		ID           string
		Name         string
		PublicID     string
		Bio          string
		ProfileImage string
		Status       string
	}
	ListStudentParam struct {
		PaginatedParam
		AdminID string
	}
	UpdateStudentParam struct {
		ID           string
		Name         *string
		Bio          *string
		ProfileImage multipart.File
	}
	StudentChangePasswordParam struct {
		ID          string
		NewPassword string
	}
	StudentProfileParam struct {
		ID string
	}
	StudentProfileQuery struct {
		ID           string
		Name         string
		PublicID     string
		Bio          string
		ProfileImage string
		Status       string
	}
	RefreshTokenParam struct {
		UserID string
		Role   int
	}
)
