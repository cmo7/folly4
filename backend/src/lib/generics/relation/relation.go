// Package relation provides functionality to parse and handle relations.
//
// Relation represents a type alias for string to define relations.
//
// Parse takes a comma-separated string and returns a slice of Relation.
// Each element in the input string is split by commas and converted to a Relation type.
package relation

import "strings"

type Relation string

func Parse(s string) ([]Relation, error) {
	orderStrings := strings.Split(s, ",")
	relations := make([]Relation, len(orderStrings))
	for i, orderString := range orderStrings {
		relations[i] = Relation(orderString)
	}

	return relations, nil
}
