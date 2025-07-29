package entity

import "time"

type (
	User struct {
		ID           string
		Name         string
		Email        string
		Password     string
		Bio          *string
		ProfileImage *string
		CreatedAt    time.Time
		UpdatedAt    time.Time
		DeletedAt    *time.Time
	}
)
