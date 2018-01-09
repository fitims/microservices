package tools

import (
	"fmt"
)

// MapStringFunc is a function signature used for mapping string array to other string array
type MapStringFunc func(i string) string

// MapStrings will transform the provided array of string into another array of string
func MapStrings(mapFunc MapStringFunc, strings []string) []string {
	r := make([]string, len(strings))

	for i, s := range strings {
		r[i] = mapFunc(s)
	}
	return r
}

// SingleQuotes will surround string to single quotes
func SingleQuotes(s string) string {
	return fmt.Sprintf("'%s'", s)
}
