package entity

import (
	"mime/multipart"
	"time"
)

type (
	Student struct {
		ID        string
		CreatedAt time.Time
		UpdatedAt time.Time
		DeletedAt *time.Time
	}
	StudentRegisterParam struct {
		Name        string
		Email       string
		Password    string
		Bio         *string
		File        multipart.File
		ContentType string
		Status      string
	}
)
