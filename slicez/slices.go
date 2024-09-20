package slicez

// Diff compares slices s1 and s2, puts elements from s1 that do not exist in s2 into dst, and returns it.
func Diff[T comparable](dst, s1, s2 []T) []T {
	dst = dst[:0]
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
	dst = dst[:0]
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
	dst = dst[:0]
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

// UniqueByKey through keyFn get the key of slice s, puts unique elements into dst, and returns it.
func UniqueByKey[T any, K comparable](dst, s []T, keyFn func(T) K) []T {
	dst = dst[:0]
	if len(s) == 0 {
		return dst
	}

	seen := make(map[K]struct{}, len(s))
	var uniqueCount int
	for _, v := range s {
		seen[keyFn(v)] = struct{}{}
		if uniqueCount < len(seen) {
			dst = append(dst, v)
			uniqueCount = len(seen)
		}
	}
	return dst
}

// UniqueByKeyInPlace through keyFn get the key of slice s, moves unique elements to the front of s, and returns this portion of s.
func UniqueByKeyInPlace[T any, K comparable](s []T, keyFn func(T) K) []T {
	if len(s) == 0 {
		return s
	}

	seen := make(map[K]struct{}, len(s))
	var remain, uniqueCount int
	for i := range s {
		seen[keyFn(s[i])] = struct{}{}
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
	dst = dst[:0]
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

// IndexFunc returns the first index i satisfying f(s[i]), or -1 if there is no such index.
func IndexFunc[T any](s []T, fn func(T) bool) int {
	for i := range s {
		if fn(s[i]) {
			return i
		}
	}
	return -1
}

// SubSlice returns a slice of s from start to end.
func SubSlice[T any](s []T, start int, end int) []T {
	if start > len(s) {
		return nil
	} else if start < 0 {
		start = 0
	}

	if end < 0 || end > len(s) {
		end = len(s)
	}

	if start >= end {
		return nil
	}

	return s[start:end]
}

// Contains returns true if v is present in s.
func Contains[T comparable](s []T, v T) bool {
	return Index(s, v) >= 0
}

// ContainsFunc returns true if there is an element in s that satisfies f(s[i]).
func ContainsFunc[T comparable](s []T, fn func(T) bool) bool {
	return IndexFunc(s, fn) >= 0
}

// Chunk splits slice s into chunks of size chunkSize, and returns the result.
func Chunk[T any](s []T, chunkSize int) [][]T {
	if len(s) == 0 {
		return nil
	}

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
	if len(s) == 0 {
		return nil
	}

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

	maxn := l - start
	if length < 0 || length > maxn {
		length = maxn
	}

	return append([]T(nil), s[start:start+length]...)
}

// Values returns a new slice containing the values returned by applying fn to each element of s.
func Values[T, V any](fn func(T) V, ss ...[]T) []V {
	var n int
	for _, s := range ss {
		n += len(s)
	}

	ret := make([]V, n)
	n = 0
	for _, s := range ss {
		for _, v := range s {
			ret[n] = fn(v)
			n++
		}
	}
	return ret
}
