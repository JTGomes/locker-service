package api

type Pagination struct {
	Limit  int
	Offset int
}

type Page[T any] struct {
	Data  []T   `json:"data"`
	Total int64 `json:"total"`
}

const (
	DefaultLimit = 20
	MaxLimit     = 100
)
