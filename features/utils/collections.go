package utils

import (
	"cmp"
	"sort"
)

// ------------------------------------------------------------------------------------------------
// Maps
// ------------------------------------------------------------------------------------------------

// GetSortedKeys returns a sorted slice of keys from a map.
// K must be a comparable type, which is a constraint satisfied by all types that can be map keys.
func GetSortedKeys[K cmp.Ordered, V any](m map[K]V) (keys []K) {
	for k := range m {
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})

	return
}

// GetSortedKeys returns a sorted slice of values from a map.
// K must be a comparable type, which is a constraint satisfied by all types that can be map keys.
func GetSortedValues[K cmp.Ordered, V any](m map[K]V) (values []V) {
	for _, key := range GetSortedKeys[K, V](m) {
		values = append(values, m[key])
	}

	return
}

// GetOneMapValue randomly returns a value from the map
func GetOneMapValue[K cmp.Ordered, V any](m map[K]V) (value V) {
	for _, val := range m {
		value = val
		return
	}

	return
}

// GetFirstMapValue returns the value corresponding to the first key, having sorted the keys beforehand
func GetFirstMapValue[K cmp.Ordered, V any](m map[K]V) (value V) {
	for _, key := range GetSortedKeys[K, V](m) {
		value = m[key]
		return
	}

	return
}

// InSlice returns true if the slice s contains the given element el
func InSlice[V comparable](s []V, el V) bool {
	for _, val := range s {
		if val == el {
			return true
		}
	}

	return false
}
