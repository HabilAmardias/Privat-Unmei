package entity

import "time"

type (
	Payment struct {
		ID                int
		CourseRequestID   string
		SubTotal          float64
		OperationalCost   float64
		TotalPrice        float64
		PaymentMethodName string
		AccountNumber     string
		CreatedAt         time.Time
		UpdatedAt         time.Time
		DeletedAt         *time.Time
	}
	PaymentMethod struct {
		ID        int
		Name      string
		CreatedAt time.Time
		UpdatedAt time.Time
		DeletedAt *time.Time
	}
	MentorPayment struct {
		PaymentMethodID   int
		PaymentMethodName string
		MentorID          string
		AccountNumber     string
		CreatedAt         time.Time
		UpdatedAt         time.Time
		DeletedAt         *time.Time
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
		PaginatedParam
		Search *string
	}
	GetMentorPaymentMethodQuery struct {
		ID            int
		Name          string
		AccountNumber string
	}
	GetPaymentMethodQuery struct {
		ID   int
		Name string
	}
	GetMentorPaymentMethodParam struct {
		UserID   string
		MentorID string
	}
	GetMyPaymentMethodParam struct {
		MentorID string
	}
)
