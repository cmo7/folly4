package filter

import (
	"testing"
)

func TestParseLeaf(t *testing.T) {
	tests := []struct {
		input    string
		expected Leaf
		hasError bool
	}{
		{"name:eq:John", Leaf{Field: "name", Comparator: ComparatorEqual, Value: "John"}, false},
		{"age:gt:30", Leaf{Field: "age", Comparator: ComparatorGreaterThan, Value: "30"}, false},
		{"invalid:format", Leaf{}, true},
	}

	for _, test := range tests {
		result, err := ParseLeaf(test.input)
		if (err != nil) != test.hasError {
			t.Errorf("ParseLeaf(%s) error = %v, expected error = %v", test.input, err, test.hasError)
		}
		if result != test.expected {
			t.Errorf("ParseLeaf(%s) = %v, expected %v", test.input, result, test.expected)
		}
	}
}

func TestParseComposite(t *testing.T) {
	tests := []struct {
		input    string
		expected Composite
		hasError bool
	}{
		{"and(name:eq:John,age:gt:30)", Composite{Operator: LogicalAnd, Filters: []Filter{
			Leaf{Field: "name", Comparator: ComparatorEqual, Value: "John"},
			Leaf{Field: "age", Comparator: ComparatorGreaterThan, Value: "30"},
		}}, false},
		{"or(name:eq:John,age:gt:30)", Composite{Operator: LogicalOr, Filters: []Filter{
			Leaf{Field: "name", Comparator: ComparatorEqual, Value: "John"},
			Leaf{Field: "age", Comparator: ComparatorGreaterThan, Value: "30"},
		}}, false},
		{"not(name:eq:John)", Composite{Operator: LogicalNot, Filters: []Filter{
			Leaf{Field: "name", Comparator: ComparatorEqual, Value: "John"},
		}}, false},
		{"invalid(operator)", Composite{}, true},
	}

	for _, test := range tests {
		result, err := ParseComposite(test.input)
		if (err != nil) != test.hasError {
			t.Errorf("ParseComposite(%s) error = %v, expected error = %v", test.input, err, test.hasError)
		}
		if !test.hasError && !compareComposites(result, test.expected) {
			t.Errorf("ParseComposite(%s) = %v, expected %v", test.input, result, test.expected)
		}
	}
}

func TestParse(t *testing.T) {
	tests := []struct {
		input    string
		expected Filter
		hasError bool
	}{
		{"name:eq:John", Leaf{Field: "name", Comparator: ComparatorEqual, Value: "John"}, false},
		{"and(name:eq:John,age:gt:30)", Composite{Operator: LogicalAnd, Filters: []Filter{
			Leaf{Field: "name", Comparator: ComparatorEqual, Value: "John"},
			Leaf{Field: "age", Comparator: ComparatorGreaterThan, Value: "30"},
		}}, false},
		{"invalid:format", nil, true},
	}

	for _, test := range tests {
		result, err := Parse(test.input)
		if (err != nil) != test.hasError {
			t.Errorf("Parse(%s) error = %v, expected error = %v", test.input, err, test.hasError)
		}
		if !test.hasError && !compareFilters(result, test.expected) {
			t.Errorf("Parse(%s) = %v, expected %v", test.input, result, test.expected)
		}
	}
}

func TestAnd(t *testing.T) {
	filters := []Filter{
		Leaf{Field: "name", Comparator: ComparatorEqual, Value: "John"},
		Leaf{Field: "age", Comparator: ComparatorGreaterThan, Value: "30"},
	}
	expected := Composite{Operator: LogicalAnd, Filters: filters}
	result := And(filters...)
	if !compareComposites(result, expected) {
		t.Errorf("And(%v) = %v, expected %v", filters, result, expected)
	}
}

func TestOr(t *testing.T) {
	filters := []Filter{
		Leaf{Field: "name", Comparator: ComparatorEqual, Value: "John"},
		Leaf{Field: "age", Comparator: ComparatorGreaterThan, Value: "30"},
	}
	expected := Composite{Operator: LogicalOr, Filters: filters}
	result := Or(filters...)
	if !compareComposites(result, expected) {
		t.Errorf("Or(%v) = %v, expected %v", filters, result, expected)
	}
}

func TestNot(t *testing.T) {
	filter := Leaf{Field: "name", Comparator: ComparatorEqual, Value: "John"}
	expected := Composite{Operator: LogicalNot, Filters: []Filter{filter}}
	result := Not(filter)
	if !compareComposites(result, expected) {
		t.Errorf("Not(%v) = %v, expected %v", filter, result, expected)
	}
}

func TestLeafComparators(t *testing.T) {
	tests := []struct {
		name     string
		function func(string, interface{}) Leaf
		field    string
		value    interface{}
		expected Leaf
	}{
		{"Equal", Equal, "name", "John", Leaf{Field: "name", Comparator: ComparatorEqual, Value: "John"}},
		{"NotEqual", NotEqual, "name", "John", Leaf{Field: "name", Comparator: ComparatorNotEqual, Value: "John"}},
		{"GreaterThan", GreaterThan, "age", 30, Leaf{Field: "age", Comparator: ComparatorGreaterThan, Value: 30}},
		{"GreaterThanOrEqual", GreaterThanOrEqual, "age", 30, Leaf{Field: "age", Comparator: ComparatorGreaterThanOrEqual, Value: 30}},
		{"LessThan", LessThan, "age", 30, Leaf{Field: "age", Comparator: ComparatorLessThan, Value: 30}},
		{"LessThanOrEqual", LessThanOrEqual, "age", 30, Leaf{Field: "age", Comparator: ComparatorLessThanOrEqual, Value: 30}},
		{"Like", Like, "name", "Jo%", Leaf{Field: "name", Comparator: ComparatorLike, Value: "Jo%"}},
		{"NotLike", NotLike, "name", "Jo%", Leaf{Field: "name", Comparator: ComparatorNotLike, Value: "Jo%"}},
		{"In", In, "age", "20,30", Leaf{Field: "age", Comparator: ComparatorIn, Value: "20,30"}},
		{"NotIn", NotIn, "age", "20,30", Leaf{Field: "age", Comparator: ComparatorNotIn, Value: "20,30"}},
		{"IsNull", func(field string, _ interface{}) Leaf { return IsNull(field) }, "name", nil, Leaf{Field: "name", Comparator: ComparatorIsNull, Value: nil}},
		{"IsNotNull", func(field string, _ interface{}) Leaf { return IsNotNull(field) }, "name", nil, Leaf{Field: "name", Comparator: ComparatorIsNotNull, Value: nil}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := test.function(test.field, test.value)
			if result != test.expected {
				t.Errorf("%s(%s, %v) = %v, expected %v", test.name, test.field, test.value, result, test.expected)
			}
		})
	}
}

func compareFilters(a, b Filter) bool {
	if a.IsComposite() != b.IsComposite() {
		return false
	}
	if a.IsComposite() {
		return compareComposites(a.(Composite), b.(Composite))
	}
	return a.(Leaf) == b.(Leaf)
}

func compareComposites(a, b Composite) bool {
	if a.Operator != b.Operator {
		return false
	}
	if len(a.Filters) != len(b.Filters) {
		return false
	}
	for i := range a.Filters {
		if !compareFilters(a.Filters[i], b.Filters[i]) {
			return false
		}
	}
	return true
}
