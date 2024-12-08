package order

import (
	"fmt"
	"strings"
)

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

// Parse order by string format "field:direction,field:direction"
func Parse(s string) ([]OrderBy, error) {
	orderStrings := strings.Split(s, ",")
	orders := make([]OrderBy, len(orderStrings))
	for i, orderString := range orderStrings {
		orderBy, err := parseOrerBy(orderString)
		if err != nil {
			return nil, err
		}
		orders[i] = orderBy
	}

	return orders, nil
}

// OrderBy string format is "field:direction"
func parseOrerBy(s string) (OrderBy, error) {
	parts := strings.Split(s, ":")
	if len(parts) != 2 {
		return OrderBy{}, fmt.Errorf("invalid order by format: %s", s)
	}

	direction := OrderDirection(parts[1])

	return OrderBy{Field: parts[0], Direction: direction}, nil
}
