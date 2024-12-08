package order

import (
	"testing"
)

func TestAscOrderBy(t *testing.T) {
	field := "name"
	order := AscOrderBy(field)
	if order.Field != field || order.Direction != Asc {
		t.Errorf("AscOrderBy(%s) = %v, want %v", field, order, OrderBy{Field: field, Direction: Asc})
	}
}

func TestDescOrderBy(t *testing.T) {
	field := "name"
	order := DescOrderBy(field)
	if order.Field != field || order.Direction != Desc {
		t.Errorf("DescOrderBy(%s) = %v, want %v", field, order, OrderBy{Field: field, Direction: Desc})
	}
}

func TestNewOrderBy(t *testing.T) {
	field := "name"
	direction := Asc
	order := NewOrderBy(field, direction)
	if order.Field != field || order.Direction != direction {
		t.Errorf("NewOrderBy(%s, %s) = %v, want %v", field, direction, order, OrderBy{Field: field, Direction: direction})
	}
}

func TestParse(t *testing.T) {
	tests := []struct {
		input    string
		expected []OrderBy
		hasError bool
	}{
		{"name:asc", []OrderBy{{Field: "name", Direction: Asc}}, false},
		{"name:asc,age:desc", []OrderBy{{Field: "name", Direction: Asc}, {Field: "age", Direction: Desc}}, false},
		{"invalid", nil, true},
	}

	for _, test := range tests {
		result, err := Parse(test.input)
		if (err != nil) != test.hasError {
			t.Errorf("Parse(%s) error = %v, wantErr %v", test.input, err, test.hasError)
			continue
		}
		if !test.hasError && !equal(result, test.expected) {
			t.Errorf("Parse(%s) = %v, want %v", test.input, result, test.expected)
		}
	}
}

func equal(a, b []OrderBy) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
