package entity

type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	PerPage    int         `json:"per_page"`
	TotalPages int         `json:"total_pages"`
}

type Meta struct {
	Page       int   `json:"page"`
	PerPage    int   `json:"per_page"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

type PaginatedResult struct {
	Data interface{} `json:"data"`
	Meta Meta        `json:"meta"`
}

func NewPaginatedResult(data interface{}, page, perPage int, total int64) *PaginatedResult {
	totalPages := (int(total) + perPage - 1) / perPage
	return &PaginatedResult{
		Data: data,
		Meta: Meta{
			Page:       page,
			PerPage:    perPage,
			Total:      total,
			TotalPages: totalPages,
		},
	}
}
