package entities

import (
	"math"
)

type Pagination struct {
	Page       int `json:"page"`
	PerPage    int `json:"page_size"`
	TotalPages int `json:"total_pages"`
	Total      int `json:"total"`
}

type PaginationQuery struct {
	Page    int
	PerPage int
}

func MakePaginationFromQueryWithTotal(total int, query PaginationQuery) *Pagination {
	page := 1
	if query.Page != 0 {
		page = query.Page
	}

	perPage := total
	if query.PerPage != 0 {
		perPage = query.PerPage
	}
	totalPages := 1
	if query.PerPage != 0 {
		totalPages = int(math.Ceil(float64(total) / float64(query.PerPage)))
	}
	return &Pagination{
		Page:       page,
		PerPage:    perPage,
		TotalPages: totalPages,
		Total:      total,
	}
}
