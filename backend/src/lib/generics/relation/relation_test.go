package relation

import (
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		input    string
		expected []Relation
	}{
		{"a,b,c", []Relation{"a", "b", "c"}},
		{"", []Relation{""}},
		{"one,two,three", []Relation{"one", "two", "three"}},
		{"single", []Relation{"single"}},
	}

	for _, test := range tests {
		result, err := Parse(test.input)
		if err != nil {
			t.Errorf("Parse(%q) returned an error: %v", test.input, err)
		}
		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("Parse(%q) = %v, want %v", test.input, result, test.expected)
		}
	}
}
