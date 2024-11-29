package pagination

type Page[T any] struct {
	Content  []T   `json:"content"`
	Page     int   `json:"page"`
	Size     int   `json:"size"`
	Total    int64 `json:"total"`
	Filtered int64 `json:"filtered"`
}

func NewPage[T any](content []T, page int, size int, total int64, filtered int64) Page[T] {
	return Page[T]{Content: content, Page: page, Size: size, Total: total, Filtered: filtered}
}

type Pageable struct {
	Page int
	Size int
}

func NewPageable(page int, size int) Pageable {
	return Pageable{Page: page, Size: size}
}
