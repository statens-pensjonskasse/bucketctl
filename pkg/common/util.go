package common

import (
	"sort"
	"strings"
)

func GetLexicallySortedKeys[T any](stringMap map[string]T) []string {
	keys := make([]string, 0, len(stringMap))
	for k := range stringMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func SlicesContainsSameElements[T comparable](aSlice, bSlice []T) bool {
	if len(aSlice) != len(bSlice) {
		return false
	}

	diff := make(map[T]int, len(aSlice))
	for _, aValue := range aSlice {
		// Count every occurrence
		diff[aValue]++
	}
	for _, bValue := range bSlice {
		if _, exists := diff[bValue]; !exists {
			return false
		}
		diff[bValue]--
		if diff[bValue] == 0 {
			// Delete entry if we've found it the same amount of times
			delete(diff, bValue)
		}
	}
	return len(diff) == 0
}

func SlicesContainsSameElementsIgnoringCase(aSlice, bSlice []string) bool {
	if len(aSlice) != len(bSlice) {
		return false
	}

	diff := make(map[string]int, len(aSlice))
	for _, aValue := range aSlice {
		// Count every occurrence
		diff[strings.ToLower(aValue)]++
	}
	for _, bValue := range bSlice {
		bLower := strings.ToLower(bValue)
		if _, exists := diff[bLower]; !exists {
			return false
		}
		diff[bLower]--
		if diff[bLower] == 0 {
			// Delete entry if we've found it the same amount of times
			delete(diff, bLower)
		}
	}
	return len(diff) == 0
}
