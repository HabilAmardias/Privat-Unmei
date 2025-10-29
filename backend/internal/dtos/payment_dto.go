package dtos

type (
	CreatePaymentMethodReq struct {
		Name string `json:"payment_method_name" binding:"required"`
	}
	CreatePaymentMethodRes struct {
		ID int `json:"payment_method_id"`
	}
	DeletePaymentMethodRes struct {
		ID int `json:"payment_method_id"`
	}
	UpdatePaymentMethodReq struct {
		MethodNewName *string `json:"payment_method_name"`
	}
	UpdatePaymentMethodRes struct {
		ID int `json:"payment_method_id"`
	}
	GetAllPaymentMethodReq struct {
		SeekPaginatedReq
		Search *string `form:"search"`
	}
	GetPaymentMethodRes struct {
		ID   int    `json:"payment_method_id"`
		Name string `json:"payment_method_name"`
	}
	GetMentorPaymentMethodRes struct {
		ID            int    `json:"payment_method_id"`
		Name          string `json:"payment_method_name"`
		AccountNumber string `json:"account_number"`
	}
)
