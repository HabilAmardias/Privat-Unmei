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
	SortInfo struct {
		Name string `json:"name"`
		ASC  bool   `json:"asc"`
	}
	FilterInfo struct {
		Name  string `json:"name"`
		Value any    `json:"value"`
	}
	PaginatedInfo struct {
		Page     int          `json:"page"`
		Limit    int          `json:"limit"`
		TotalRow int64        `json:"total_row"`
		SortBy   []SortInfo   `json:"sort_by,omitempty"`
		FilterBy []FilterInfo `json:"filter_by,omitempty"`
	}
	SeekPaginatedInfo struct {
		LastID   int          `json:"last_id"`
		Limit    int          `json:"limit"`
		TotalRow int64        `json:"total_row"`
		SortBy   []SortInfo   `json:"sort_by,omitempty"`
		FilterBy []FilterInfo `json:"filter_by,omitempty"`
	}
	PaginatedResponse[T interface{}] struct {
		Entries  []T           `json:"entries"`
		PageInfo PaginatedInfo `json:"page_info"`
	}
	SeekPaginatedResponse[T interface{}] struct {
		Entries  []T               `json:"entries"`
		PageInfo SeekPaginatedInfo `json:"page_info"`
	}
)
