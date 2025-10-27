package entity

import "time"

type (
	Admin struct {
		ID        string
		CreatedAt time.Time
		UpdatedAt time.Time
		DeletedAt *time.Time
	}
	AdminLoginParam struct {
		Email    string
		Password string
	}
	AdminVerificationParam struct {
		AdminID  string
		Email    string
		Password string
	}
	AdminUpdatePasswordParam struct {
		AdminID  string
		Password string
	}
	AdminProfileParam struct {
		AdminID string
	}
	AdminProfileQuery struct {
		Name         string
		Email        string
		Bio          string
		ProfileImage string
		Status       string
	}
)
