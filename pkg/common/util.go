package common

import (
	"sort"
)

func GetLexicallySortedKeys[T any](stringMap map[string]T) []string {
	keys := make([]string, 0, len(stringMap))
	for k := range stringMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func SlicesContainsSameElements[T comparable](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}

	diff := make(map[T]int, len(a))
	for _, i := range a {
		// Count every occurrence
		diff[i]++
	}
	for _, j := range b {
		if _, exists := diff[j]; !exists {
			return false
		}
		diff[j]--
		if diff[j] == 0 {
			// Delete entry if we've found it the same amount of times
			delete(diff, j)
		}
	}
	return len(diff) == 0
}
