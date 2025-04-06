package main

// This file contains polyfill implementations for features available in Go 1.23+
// but not in earlier versions like Go 1.19.3, which is used in the current development environment.
// These implementations allow us to maintain compatibility with newer Go versions while
// still being able to build in the current environment.

// MapKeys returns the keys of the map m as a slice.
// This is a polyfill for maps.Keys in Go 1.22+.
func MapKeys[M ~map[K]V, K comparable, V any](m M) []K {
	r := make([]K, 0, len(m))
	for k := range m {
		r = append(r, k)
	}
	return r
}

// MapValues returns the values of the map m as a slice.
// This is a polyfill for maps.Values in Go 1.22+.
func MapValues[M ~map[K]V, K comparable, V any](m M) []V {
	r := make([]V, 0, len(m))
	for _, v := range m {
		r = append(r, v)
	}
	return r
}

// SliceContains reports whether v is present in s.
// This is a polyfill for slices.Contains in Go 1.21+.
func SliceContains[S ~[]E, E comparable](s S, v E) bool {
	for _, vs := range s {
		if vs == v {
			return true
		}
	}
	return false
}

// SliceEqual reports whether two slices are equal: the same length and all elements equal.
// This is a polyfill for slices.Equal in Go 1.21+.
func SliceEqual[S ~[]E, E comparable](s1, s2 S) bool {
	if len(s1) != len(s2) {
		return false
	}
	for i := range s1 {
		if s1[i] != s2[i] {
			return false
		}
	}
	return true
}

// SliceIndex returns the index of the first occurrence of v in s,
// or -1 if not present.
// This is a polyfill for slices.Index in Go 1.21+.
func SliceIndex[S ~[]E, E comparable](s S, v E) int {
	for i := range s {
		if s[i] == v {
			return i
		}
	}
	return -1
}

// CloneMap returns a copy of m.
// This is a polyfill for maps.Clone in Go 1.21+.
func CloneMap[M ~map[K]V, K comparable, V any](m M) M {
	result := make(M, len(m))
	for k, v := range m {
		result[k] = v
	}
	return result
}

// Min returns the smaller of a and b.
// This is a polyfill for min in Go 1.21+.
func Min[T Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}

// Max returns the larger of a and b.
// This is a polyfill for max in Go 1.21+.
func Max[T Ordered](a, b T) T {
	if a > b {
		return a
	}
	return b
}

// Ordered is a constraint that permits any ordered type: any type
// that supports the operators < <= >= >.
// This is defined in Go 1.18 with generics.
type Ordered interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64 |
		~string
}
