// Package filter provides functionality for creating and parsing filters
// that can be used to query data. Filters can be either leaf filters or
// composite filters.
//
// Leaf filters represent simple conditions on a single field, such as
// "name equals John" or "age greater than 30".
//
// Composite filters represent logical combinations of other filters,
// such as "name equals John AND age greater than 30".
//
// The package defines various comparison operators for leaf filters,
// such as equality, inequality, greater than, less than, etc. It also
// defines logical operators for composite filters, such as AND, OR, and NOT.
//
// The package provides functions for creating both leaf and composite
// filters, as well as a Parse function for parsing filter strings into
// Filter objects. The filter strings can represent either leaf filters
// or composite filters, and composite filters can be nested.
package filter

import (
	"fmt"
	"strings"
)

type Filter interface {
	IsComposite() bool
	ToString() string
}

type comparisonOperator string

const (
	ComparatorEqual              comparisonOperator = "eq"
	ComparatorNotEqual           comparisonOperator = "ne"
	ComparatorGreaterThan        comparisonOperator = "gt"
	ComparatorGreaterThanOrEqual comparisonOperator = "ge"
	ComparatorLessThan           comparisonOperator = "lt"
	ComparatorLessThanOrEqual    comparisonOperator = "le"
	ComparatorLike               comparisonOperator = "like"
	ComparatorNotLike            comparisonOperator = "not_like"
	ComparatorIn                 comparisonOperator = "in"
	ComparatorNotIn              comparisonOperator = "not_in"
	ComparatorIsNull             comparisonOperator = "is_null"
	ComparatorIsNotNull          comparisonOperator = "is_not_null"
)

type Leaf struct {
	Field      string
	Comparator comparisonOperator
	Value      interface{}
}

func (f Leaf) IsComposite() bool {
	return false
}

func (f Leaf) ToString() string {
	return fmt.Sprintf("%s:%s:%v", f.Field, f.Comparator, f.Value)
}

type LogicalOperator string

const (
	LogicalAnd LogicalOperator = "and"
	LogicalOr  LogicalOperator = "or"
	LogicalNot LogicalOperator = "not"
)

type Composite struct {
	Operator LogicalOperator
	Filters  []Filter
}

func (f Composite) IsComposite() bool {
	return true
}

func (f Composite) ToString() string {
	filterStrings := make([]string, 0, len(f.Filters))
	for _, filter := range f.Filters {
		filterStrings = append(filterStrings, filter.ToString())
	}
	return fmt.Sprintf("%s(%s)", f.Operator, strings.Join(filterStrings, ","))
}

func And(filters ...Filter) Composite {
	return Composite{Operator: LogicalAnd, Filters: filters}
}

func Or(filters ...Filter) Composite {
	return Composite{Operator: LogicalOr, Filters: filters}
}

func Not(filter Filter) Composite {
	return Composite{Operator: LogicalNot, Filters: []Filter{filter}}
}

func Equal(field string, value interface{}) Leaf {
	return Leaf{Field: field, Comparator: ComparatorEqual, Value: value}
}

func NotEqual(field string, value interface{}) Leaf {
	return Leaf{Field: field, Comparator: ComparatorNotEqual, Value: value}
}

func GreaterThan(field string, value interface{}) Leaf {
	return Leaf{Field: field, Comparator: ComparatorGreaterThan, Value: value}
}

func GreaterThanOrEqual(field string, value interface{}) Leaf {
	return Leaf{Field: field, Comparator: ComparatorGreaterThanOrEqual, Value: value}
}

func LessThan(field string, value interface{}) Leaf {
	return Leaf{Field: field, Comparator: ComparatorLessThan, Value: value}
}

func LessThanOrEqual(field string, value interface{}) Leaf {
	return Leaf{Field: field, Comparator: ComparatorLessThanOrEqual, Value: value}
}

func Like(field string, value interface{}) Leaf {
	return Leaf{Field: field, Comparator: ComparatorLike, Value: value}
}

func NotLike(field string, value interface{}) Leaf {
	return Leaf{Field: field, Comparator: ComparatorNotLike, Value: value}
}

func In(field string, value interface{}) Leaf {
	return Leaf{Field: field, Comparator: ComparatorIn, Value: value}
}

func NotIn(field string, value interface{}) Leaf {
	return Leaf{Field: field, Comparator: ComparatorNotIn, Value: value}
}

func IsNull(field string) Leaf {
	return Leaf{Field: field, Comparator: ComparatorIsNull, Value: nil}
}

func IsNotNull(field string) Leaf {
	return Leaf{Field: field, Comparator: ComparatorIsNotNull, Value: nil}
}

func NewComposite(operator LogicalOperator, filters ...Filter) Composite {
	return Composite{Operator: operator, Filters: filters}
}

func NewLeaf(field string, comparator comparisonOperator, value interface{}) Leaf {
	return Leaf{Field: field, Comparator: comparator, Value: value}
}

// Parse takes a filter string and returns a Filter object or an error.
// The filter string can represent either a leaf filter or a composite filter.
//
// Leaf Filter Format: FieldName:Comparator:Value
// Example: name:eq:John
//
// Composite Filter Format: Operator(Filter1, Filter2, ...)
// Example: and(name:eq:John,age:gt:30)
//
// Composite filters can be nested and can contain both leaf filters and other composite filters.
// Example: and(name:eq:John,or(age:gt:30,age:lt:20))
//
// If the filter string starts with a logical operator (and, or, not) followed by a parenthesis,
// it is considered a composite filter and will be parsed by ParseComposite.
// Otherwise, it is considered a leaf filter and will be parsed by ParseLeaf.
//
// Parameters:
//
//	s - The filter string to be parsed.
//
// Returns:
//
//	Filter - The parsed Filter object.
//	error - An error if the filter string is invalid.
func Parse(s string) (Filter, error) {
	// Leaf Filter Format: FieldName:Comparator:Value
	// Composite Filter Format: Operator(Filter1, Filter2, ...)
	// Example: and(name:eq:John,age:gt:30)

	// The filter list is enclosed in parentheses
	// The filter list can contain both leaf filters and composite filters
	// Filters in the filter list are separated by commas, but they can also be nested
	// Example: and(name:eq:John,or(age:gt:30,age:lt:20))

	// If the filter starts with an logical operator and parenthesis, it is a composite filter
	if strings.HasPrefix(s, "and(") || strings.HasPrefix(s, "or(") || strings.HasPrefix(s, "not(") {
		return ParseComposite(s)
	}

	// Otherwise, it is a leaf filter
	return ParseLeaf(s)

}

func ParseLeaf(s string) (Leaf, error) {
	parts := strings.Split(s, ":")
	if len(parts) != 3 {
		return Leaf{}, fmt.Errorf("invalid leaf filter format: %s", s)
	}

	comparator := comparisonOperator(parts[1])
	return Leaf{Field: parts[0], Comparator: comparator, Value: parts[2]}, nil
}

func ParseComposite(s string) (Composite, error) {
	operator := strings.Split(s, "(")[0]
	if operator != "and" && operator != "or" && operator != "not" {
		return Composite{}, fmt.Errorf("invalid composite filter operator: %s", operator)
	}

	// Remove the operator and the last parenthesis
	s = s[len(operator)+1 : len(s)-1]

	filterList, err := ParseFilterList(s)
	if err != nil {
		return Composite{}, err
	}
	return Composite{Operator: LogicalOperator(operator), Filters: filterList}, nil
}

func ParseFilterList(s string) ([]Filter, error) {
	// Split the filter list by commas
	parts := strings.Split(s, ",")
	filters := make([]Filter, 0, len(parts))
	for _, part := range parts {
		filter, err := Parse(part)
		if err != nil {
			return nil, err
		}
		filters = append(filters, filter)
	}
	return filters, nil
}
