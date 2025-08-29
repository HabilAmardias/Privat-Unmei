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
		Bio      string
		Status   string
	}
	StudentLoginParam struct {
		Email    string
		Password string
	}
	ResetPasswordParam struct {
		NewPassword string
		Token       string
	}
	ListStudentQuery struct {
		ID           string
		Name         string
		Email        string
		Bio          string
		ProfileImage string
		Status       string
	}
	ListStudentParam struct {
		PaginatedParam
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
)
