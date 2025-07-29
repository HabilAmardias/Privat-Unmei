package dtos

type (
	Response struct {
		Success bool        `json:"success"`
		Data    interface{} `json:"data"`
	}
	MessageResponse struct {
		Message string `json:"message"`
	}
	DetailsError struct {
		Title   string `json:"field"`
		Message string `json:"message"`
	}
)
