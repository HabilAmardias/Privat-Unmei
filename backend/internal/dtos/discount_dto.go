package dtos

type (
	CreateNewDiscountReq struct {
		NumberOfParticipant int     `json:"number_of_participant" binding:"required,gte=1"`
		Amount              float64 `json:"amount" binding:"required,gte=1"`
	}
	UpdateDiscountReq struct {
		Amount *float64 `json:"amount" binding:"omitempty,gte=1"`
	}
	DiscountIDRes struct {
		ID int `json:"id"`
	}
	GetAllDiscountReq struct {
		PaginatedReq
	}
	GetAllDiscountRes struct {
		ID                  int     `json:"id"`
		NumberOfParticipant int     `json:"number_of_participant"`
		Amount              float64 `json:"amount"`
	}
	GetDiscountRes struct {
		Amount float64 `json:"amount"`
	}
)
