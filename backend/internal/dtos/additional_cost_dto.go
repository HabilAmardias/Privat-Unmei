package dtos

type (
	CreateAdditionalCostReq struct {
		Name   string  `json:"name" binding:"required"`
		Amount float64 `json:"amount" binding:"required,gte=1"`
	}
	UpdateAdditionalCostReq struct {
		Amount *float64 `json:"amount" binding:"required,gte=1"`
	}
	AdditionalCostIDRes struct {
		ID int `json:"id"`
	}
	GetAllAdditionalCostReq struct {
		PaginatedReq
	}
	GetAdditionalCostRes struct {
		ID     int     `json:"id"`
		Name   string  `json:"name"`
		Amount float64 `json:"amount"`
	}
)
