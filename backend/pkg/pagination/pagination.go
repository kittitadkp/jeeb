package pagination

import (
	"net/http"
	"strconv"
)

const (
	defaultPage  = 1
	defaultLimit = 20
	maxLimit     = 100
)

type Params struct {
	Page  int    `json:"page"`
	Limit int    `json:"limit"`
	Sort  string `json:"sort"`
}

type Meta struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

func FromRequest(r *http.Request) Params {
	q := r.URL.Query()

	page, _ := strconv.Atoi(q.Get("page"))
	if page < 1 {
		page = defaultPage
	}

	limit, _ := strconv.Atoi(q.Get("limit"))
	if limit < 1 || limit > maxLimit {
		limit = defaultLimit
	}

	sort := q.Get("sort")
	if sort == "" {
		sort = "-created_at"
	}

	return Params{Page: page, Limit: limit, Sort: sort}
}

func NewMeta(page, limit int, total int64) *Meta {
	totalPages := int(total) / limit
	if int(total)%limit != 0 {
		totalPages++
	}
	return &Meta{Page: page, Limit: limit, Total: total, TotalPages: totalPages}
}

func (p Params) Skip() int64 {
	return int64((p.Page - 1) * p.Limit)
}
