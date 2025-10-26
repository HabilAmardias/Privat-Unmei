package entity

import "time"

type (
	PaymentMethod struct {
		ID        int
		Name      string
		CreatedAt time.Time
		UpdatedAt time.Time
		DeletedAt *time.Time
	}
	MentorPayment struct {
		PaymentMethodID int
		MentorID        string
		AccountNumber   string
		CreatedAt       time.Time
		UpdatedAt       time.Time
		DeletedAt       *time.Time
	}
	CreatePaymentMethodParam struct {
		AdminID    string
		MethodName string
	}
	DeletePaymentMethodParam struct {
		AdminID  string
		MethodID int
	}
	UpdatePaymentMethodParam struct {
		AdminID       string
		MethodID      int
		MethodNewName *string
	}
	GetAllPaymentMethodParam struct {
		SeekPaginatedParam
		Search *string
	}
	GetPaymentMethodQuery struct {
		ID            int
		Name          string
		AccountNumber string
	}
	GetMentorPaymentMethodParam struct {
		UserID   string
		MentorID string
	}
)
