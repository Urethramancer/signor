package stringer

import (
	"sort"
	"strings"
)

// RemoveDuplicateStrings returns a sorted list with no duplicates. Case matters.
func RemoveDuplicateStrings(list []string) []string {
	m := map[string]bool{}
	if len(list) < 2 {
		return list
	}

	for _, x := range list {
		m[x] = true
	}

	list = []string{}
	for k := range m {
		list = append(list, k)
	}

	sort.Strings(list)
	return list
}

// RemoveDuplicateStringsAndTitle converts the names to lowercase and returns a sorted list with no duplicates and the first letter capitalised.
func RemoveDuplicateStringsAndTitle(list []string) []string {
	m := map[string]bool{}
	if len(list) < 2 {
		return list
	}

	for _, x := range list {
		x = strings.ToLower(x)
		x = strings.Title(x)
		m[x] = true
	}

	list = []string{}
	for k := range m {
		list = append(list, k)
	}

	sort.Strings(list)
	return list
}
