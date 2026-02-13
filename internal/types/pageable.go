package types

import "math"

type Pageable[T any] struct {
	Page       int `json:"page"`
	PageSize   int `json:"page_size"`
	TotalPages int `json:"total_pages"`
	TotalItems int `json:"total_items"`
	Data       []T `json:"data"`
}

func NewPageable[T any](data []T, page, pageSize, totalItems int) Pageable[T] {
	totalPages := max(int(math.Ceil(float64(totalItems)/float64(pageSize))), 1)
	return Pageable[T]{
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
		TotalItems: totalItems,
		Data:       data,
	}
}
