package entity

import "time"

type (
	Rbac struct {
		ID           int
		RoleID       int
		PermissionID int
		ResourceID   int
		CreatedAt    time.Time
		UpdatedAt    time.Time
		DeletedAt    *time.Time
	}
)
