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
)
