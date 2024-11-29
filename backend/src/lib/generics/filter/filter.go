package filter

type Filter interface {
	IsComposite() bool
}

type comparasionOperator string

const (
	Equal              comparasionOperator = "eq"
	NotEqual           comparasionOperator = "ne"
	GreaterThan        comparasionOperator = "gt"
	GreaterThanOrEqual comparasionOperator = "ge"
	LessThan           comparasionOperator = "lt"
	LessThanOrEqual    comparasionOperator = "le"
	Like               comparasionOperator = "like"
	NotLike            comparasionOperator = "not_like"
	In                 comparasionOperator = "in"
	NotIn              comparasionOperator = "not_in"
	IsNull             comparasionOperator = "is_null"
	IsNotNull          comparasionOperator = "is_not_null"
)

type Leaf struct {
	Field      string
	Comparator comparasionOperator
	Value      interface{}
}

func (f Leaf) IsComposite() bool {
	return false
}

type LogicalOperator string

const (
	And LogicalOperator = "and"
	Or  LogicalOperator = "or"
	Not LogicalOperator = "not"
)

type Composite struct {
	Operator LogicalOperator
	Filters  []Filter
}

func (f Composite) IsComposite() bool {
	return true
}

func AndFilter(filters ...Filter) Composite {
	return Composite{Operator: And, Filters: filters}
}

func OrFilter(filters ...Filter) Composite {
	return Composite{Operator: Or, Filters: filters}
}

func NotFilter(filter Filter) Composite {
	return Composite{Operator: Not, Filters: []Filter{filter}}
}

func EqualFilter(field string, value interface{}) Leaf {
	return Leaf{Field: field, Comparator: Equal, Value: value}
}

func NotEqualFilter(field string, value interface{}) Leaf {
	return Leaf{Field: field, Comparator: NotEqual, Value: value}
}

func GreaterThanFilter(field string, value interface{}) Leaf {
	return Leaf{Field: field, Comparator: GreaterThan, Value: value}
}

func GreaterThanOrEqualFilter(field string, value interface{}) Leaf {
	return Leaf{Field: field, Comparator: GreaterThanOrEqual, Value: value}
}

func LessThanFilter(field string, value interface{}) Leaf {
	return Leaf{Field: field, Comparator: LessThan, Value: value}
}

func LessThanOrEqualFilter(field string, value interface{}) Leaf {
	return Leaf{Field: field, Comparator: LessThanOrEqual, Value: value}
}

func LikeFilter(field string, value interface{}) Leaf {
	return Leaf{Field: field, Comparator: Like, Value: value}
}

func NotLikeFilter(field string, value interface{}) Leaf {
	return Leaf{Field: field, Comparator: NotLike, Value: value}
}

func InFilter(field string, value interface{}) Leaf {
	return Leaf{Field: field, Comparator: In, Value: value}
}

func NotInFilter(field string, value interface{}) Leaf {
	return Leaf{Field: field, Comparator: NotIn, Value: value}
}

func IsNullFilter(field string) Leaf {
	return Leaf{Field: field, Comparator: IsNull, Value: nil}
}

func IsNotNullFilter(field string) Leaf {
	return Leaf{Field: field, Comparator: IsNotNull, Value: nil}
}

func NewCompositeFilter(operator LogicalOperator, filters ...Filter) Composite {
	return Composite{Operator: operator, Filters: filters}
}

func NewLeafFilter(field string, comparator comparasionOperator, value interface{}) Leaf {
	return Leaf{Field: field, Comparator: comparator, Value: value}
}
