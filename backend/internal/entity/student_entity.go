package entity

import (
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
)
