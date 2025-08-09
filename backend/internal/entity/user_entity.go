package entity

import "time"

type (
	User struct {
		ID           string
		Name         string
		Email        string
		Password     string
		Bio          string
		ProfileImage string
		Status       string
		CreatedAt    time.Time
		UpdatedAt    time.Time
		DeletedAt    *time.Time
	}
	UpdateUserQuery struct {
		Name         *string
		Password     *string
		Bio          *string
		ProfileImage *string
		Status       *string
	}
)
