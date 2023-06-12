package slicesz

// Diff compares slices s1 and s2, puts elements from s1 that do not exist in s2 into dst, and returns it.
func Diff[T comparable](dst, s1, s2 []T) []T {
	if len(s1) == 0 {
		return dst
	}

	if len(s2) == 0 {
		return append(dst, s1...)
	}

	m := make(map[T]struct{}, len(s2))
	for i := range s2 {
		m[s2[i]] = struct{}{}
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
	for i := range s2 {
		m[s2[i]] = struct{}{}
	}

	var remain int
	for i := range s1 {
		if _, ok := m[s1[i]]; !ok {
			s1[remain], s1[i] = s1[i], s1[remain]
			remain++
		}
	}
	return s1[:remain]
}

// Intersect compares slices s1 and s2, puts elements from s1 that are also present in s2 into dst, and returns it.
func Intersect[T comparable](dst, s1, s2 []T) []T {
	if len(s1) == 0 || len(s2) == 0 {
		return dst
	}

	m := make(map[T]struct{}, len(s2))
	for i := range s2 {
		m[s2[i]] = struct{}{}
	}

	for _, v := range s1 {
		if _, ok := m[v]; ok {
			dst = append(dst, v)
		}
	}
	return dst
}

// IntersectInPlaceFirst compares s1 and s2, moves elements from s1 that are also present in s2 to the front of s1,
// and returns this portion of s1. Please note that the order of elements in s1 will be altered.
func IntersectInPlaceFirst[T comparable](s1, s2 []T) []T {
	if len(s1) == 0 || len(s2) == 0 {
		return s1[:0]
	}

	m := make(map[T]struct{}, len(s2))
	for i := range s2 {
		m[s2[i]] = struct{}{}
	}

	var remain int
	for i := range s1 {
		if _, ok := m[s1[i]]; ok {
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
	for i := range s {
		seen[s[i]] = struct{}{}
		if uniqueCount < len(seen) {
			s[remain], s[i] = s[i], s[remain]
			uniqueCount = len(seen)
			remain++
		}
	}
	return s[:remain]
}

// Filter puts elements from s that satisfy predicate into dst, and returns it.
func Filter[T any](dst, s []T, predicate func(T) bool) []T {
	for _, v := range s {
		if predicate(v) {
			dst = append(dst, v)
		}
	}
	return dst
}

// FilterInPlace moves elements that satisfy predicate to the front of s,
// and returns this portion of s. Please note that the order of elements in s will be altered.
func FilterInPlace[T any](s []T, predicate func(T) bool) []T {
	var remain int
	for i := range s {
		if predicate(s[i]) {
			s[remain], s[i] = s[i], s[remain]
			remain++
		}
	}
	return s[:remain]
}

// Equal compares slices s1 and s2, and returns true if they are equal.
// If s1 or s2 is nil, it returns true.
func Equal[T comparable](s1, s2 []T) bool {
	if len(s1) != len(s2) {
		return false
	}

	s2 = s2[:len(s1)]
	for i := range s1 {
		if s1[i] != s2[i] {
			return false
		}
	}

	return true
}

// Index returns the index of the first instance of v in s, or -1 if v is not present in s.
func Index[T comparable](s []T, v T) int {
	for i := range s {
		if v == s[i] {
			return i
		}
	}
	return -1
}

// Contains returns true if v is present in s.
func Contains[T comparable](s []T, v T) bool {
	return Index(s, v) >= 0
}

// Chunk splits slice s into chunks of size chunkSize, and returns the result.
func Chunk[T any](s []T, chunkSize int) [][]T {
	if chunkSize < 1 || len(s) <= chunkSize {
		return [][]T{s}
	}

	n := len(s) / chunkSize
	chunks := make([][]T, 0, n+1)
	var start, end int
	for i := 0; i < n; i++ {
		end = start + chunkSize
		chunks = append(chunks, s[start:end])
		start = end
	}
	if len(s) > start {
		chunks = append(chunks, s[start:])
	}
	return chunks
}

// ChunkProcess splits slice s into chunks of size chunkSize, and calls process on each chunk.
func ChunkProcess[T any](s []T, chunkSize int, process func([]T) error) error {
	if chunkSize < 1 || len(s) <= chunkSize {
		return process(s)
	}

	n := len(s) / chunkSize
	var start, end int
	for i := 0; i < n; i++ {
		end = start + chunkSize
		if err := process(s[start:end]); err != nil {
			return err
		}
		start = end
	}
	if len(s) > start {
		return process(s[start:])
	}
	return nil
}

// Copy copies length elements from s starting at position start.
// If length is negative, it will copy all elements from start to the end of s.
func Copy[T any](s []T, start, length int) []T {
	l := len(s)
	if l == 0 || start >= l || length == 0 {
		return nil
	}

	if start < 0 {
		start = 0
	}

	max := l - start
	if length < 0 || length > max {
		length = max
	}

	return append([]T(nil), s[start:start+length]...)
}
