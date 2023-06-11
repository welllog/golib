package slices

// Diff compares slices s1 and s2, puts elements from s1 that do not exist in s2 into dst, and returns it.
func Diff[T comparable](dst, s1, s2 []T) []T {
	if len(s1) == 0 {
		return dst
	}

	m := make(map[T]struct{}, len(s2))
	for _, v := range s2 {
		m[v] = struct{}{}
	}

	for _, v := range s1 {
		if _, ok := m[v]; !ok {
			dst = append(dst, v)
		}
	}
	return dst
}

// DiffInPlaceFirst compares s1 and s2, moves elements from s1 that not in s2 to the front of s1,
// and returns this portion of s1. Please note that the order of elements in s1 will be altered.
func DiffInPlaceFirst[T comparable](s1, s2 []T) []T {
	if len(s1) == 0 || len(s2) == 0 {
		return s1
	}

	m := make(map[T]struct{}, len(s2))
	for _, v := range s2 {
		m[v] = struct{}{}
	}

	var remain int
	for i, v := range s1 {
		if _, ok := m[v]; !ok {
			s1[remain], s1[i] = s1[i], s1[remain]
			remain++
		}
	}
	return s1[:remain]
}

// Unique compares slice s, puts unique elements into dst, and returns it.
func Unique[T comparable](dst, s []T) []T {
	if len(s) == 0 {
		return dst
	}

	seen := make(map[T]struct{}, len(s))
	var uniqueCount int
	for _, v := range s {
		seen[v] = struct{}{}
		if uniqueCount < len(seen) {
			dst = append(dst, v)
			uniqueCount = len(seen)
		}
	}
	return dst
}

// UniqueInPlace compares slice s, moves unique elements to the front of s, and returns this portion of s.
// Please note that the order of elements in s will be altered.
func UniqueInPlace[T comparable](s []T) []T {
	if len(s) == 0 {
		return s
	}

	seen := make(map[T]struct{}, len(s))
	var remain, uniqueCount int
	for i, v := range s {
		seen[v] = struct{}{}
		if uniqueCount < len(seen) {
			s[remain], s[i] = s[i], s[remain]
			uniqueCount = len(seen)
			remain++
		}
	}
	return s[:remain]
}

func Equal[T comparable](s1, s2 []T) bool {
	if len(s1) != len(s2) {
		return false
	}

	if len(s1) == 0 {
		return true
	}

	s2 = s2[:len(s1)]
	for i, v := range s1 {
		if v != s2[i] {
			return false
		}
	}

	return true
}
