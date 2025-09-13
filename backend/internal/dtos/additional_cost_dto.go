package dtos

type (
	CreateAdditionalCostReq struct {
		Name   string  `json:"name" binding:"required"`
		Amount float64 `json:"amount" binding:"required,gte=1"`
	}
	AdditionalCostIDRes struct {
		ID int `json:"id"`
	}
)
