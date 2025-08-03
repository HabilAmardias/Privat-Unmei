package entity

type (
	PaginatedParam struct {
		Limit int
		Page  int
	}
	SeekPaginatedParam struct {
		Limit  int
		LastID int
	}
)
