package dtos

type (
	CreateRatingReq struct {
		Rating   int     `json:"rating" binding:"required,min=1,max=5"`
		Feedback *string `json:"feedback" binding:"omitempty,min=15"`
	}
	CreateRatingRes struct {
		RatingID int `json:"id"`
	}
)
