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
)
