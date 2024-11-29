package order

type OrderDirection string

const (
	Asc  OrderDirection = "asc"
	Desc OrderDirection = "desc"
)

type OrderBy struct {
	Field     string
	Direction OrderDirection
}

func AscOrderBy(field string) OrderBy {
	return OrderBy{Field: field, Direction: Asc}
}

func DescOrderBy(field string) OrderBy {
	return OrderBy{Field: field, Direction: Desc}
}

func NewOrderBy(field string, direction OrderDirection) OrderBy {
	return OrderBy{Field: field, Direction: direction}
}
