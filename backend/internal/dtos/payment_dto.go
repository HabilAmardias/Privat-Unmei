package dtos

type (
	CreatePaymentMethodReq struct {
		Name string `json:"payment_method_name" binding:"required"`
	}
	CreatePaymentMethodRes struct {
		ID int `json:"payment_method_id"`
	}
)
