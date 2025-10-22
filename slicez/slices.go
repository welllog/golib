package slicez

import "github.com/welllog/golib/typez"

const (
	loopEnabled   = true
	loopThreshold = 4
)

// Diff compares slices s1 and s2, puts elements from s1 that do not exist in s2 into dst, and returns it.
func Diff[T comparable](dst, s1, s2 []T) []T {
	dst = dst[:0]
	if len(s1) == 0 {
		return dst
	}

	if len(s2) == 0 {
		return append(dst, s1...)
	}

	if loopEnabled && len(s2) <= loopThreshold {
		for _, v := range s1 {
			if Index(s2, v) < 0 {
				dst = append(dst, v)
			}
		}
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

	var remain int

	if loopEnabled && len(s2) <= loopThreshold {
		for i, v := range s1 {
			if Index(s2, v) < 0 {
				s1[remain], s1[i] = s1[i], s1[remain]
				remain++
			}
		}
		return s1[:remain]
	}

	m := make(map[T]struct{}, len(s2))
	for _, v := range s2 {
		m[v] = struct{}{}
	}

	for i, v := range s1 {
		if _, ok := m[v]; !ok {
			s1[remain], s1[i] = s1[i], s1[remain]
			remain++
		}
	}
	return s1[:remain]
}

// DiffSorted compares two sorted slices s1 and s2, puts elements from s1 that do not exist in s2 into dst, and returns it.
// It requires ascending sorted slices.
func DiffSorted[T typez.Ordered](dst, s1, s2 []T) []T {
	dst = dst[:0]
	if len(s1) == 0 {
		return dst
	}
	if len(s2) == 0 {
		return append(dst, s1...)
	}

	i, j := 0, 0
	for i < len(s1) && j < len(s2) {
		if s1[i] < s2[j] {
			dst = append(dst, s1[i])
			i++
		} else if s1[i] > s2[j] {
			j++
		} else {
			// Skip all equal elements
			current := s1[i]
			for i < len(s1) && s1[i] == current {
				i++
			}
			for j < len(s2) && s2[j] == current {
				j++
			}
		}
	}

	if i < len(s1) {
		dst = append(dst, s1[i:]...)
	}

	return dst
}

// DiffSortedInPlaceFirst compares two sorted slices s1 and s2, moves elements from s1 that do not exist in s2 to the front of s1,
// and returns this portion of s1. The returned slice maintains the sorted order.
// It requires ascending sorted slices.
//
// Note: This function only reorders elements in s1. All original elements are preserved,
// but the elements beyond the returned slice length are no longer sorted.
func DiffSortedInPlaceFirst[T typez.Ordered](s1, s2 []T) []T {
	if len(s1) == 0 {
		return s1[:0]
	}
	if len(s2) == 0 {
		return s1
	}

	remain := 0
	i, j := 0, 0
	for i < len(s1) && j < len(s2) {
		if s1[i] < s2[j] {
			// s1[remain] = s1[i]
			// swap to only change the order of elements in s1
			s1[remain], s1[i] = s1[i], s1[remain]
			remain++
			i++
		} else if s1[i] > s2[j] {
			j++
		} else {
			// Skip all equal elements
			current := s1[i]
			for i < len(s1) && s1[i] == current {
				i++
			}
			for j < len(s2) && s2[j] == current {
				j++
			}
		}
	}

	// if i < len(s1) {
	// 	remain += copy(s1[remain:], s1[i:])
	// }
	//
	// swap to only change the order of elements in s1
	for i < len(s1) {
		s1[remain], s1[i] = s1[i], s1[remain]
		remain++
		i++
	}

	return s1[:remain]
}

// Intersect compares slices s1 and s2, puts elements from s1 that are also present in s2 into dst, and returns it.
// This is not a mathematical intersection, as it will include duplicate elements from s1.
// For example, s1 = [1,2,2,3], s2 = [2,2,4], the result will be [2,2].
func Intersect[T comparable](dst, s1, s2 []T) []T {
	dst = dst[:0]
	if len(s1) == 0 || len(s2) == 0 {
		return dst
	}

	if loopEnabled && len(s2) <= loopThreshold {
		for _, v := range s1 {
			if Index(s2, v) >= 0 {
				dst = append(dst, v)
			}
		}
		return dst
	}

	m := make(map[T]struct{}, len(s2))
	for _, v := range s2 {
		m[v] = struct{}{}
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
// This is not a mathematical intersection, as it will include duplicate elements from s1.
func IntersectInPlaceFirst[T comparable](s1, s2 []T) []T {
	if len(s1) == 0 || len(s2) == 0 {
		return s1[:0]
	}

	var remain int

	if loopEnabled && len(s2) <= loopThreshold {
		for i, v := range s1 {
			if Index(s2, v) >= 0 {
				s1[remain], s1[i] = s1[i], s1[remain]
				remain++
			}
		}
		return s1[:remain]
	}

	m := make(map[T]struct{}, len(s2))
	for _, v := range s2 {
		m[v] = struct{}{}
	}

	for i, v := range s1 {
		if _, ok := m[v]; ok {
			s1[remain], s1[i] = s1[i], s1[remain]
			remain++
		}
	}
	return s1[:remain]
}

// IntersectSortedMultiset compares two sorted slices s1 and s2, and puts the multiset intersection into dst.
// Its behavior differs from Intersect. For example: s1 = [1,2,2,3], s2 = [2,2,2,4], the result will be [2,2].
// It requires ascending sorted slices.
func IntersectSortedMultiset[T typez.Ordered](dst, s1, s2 []T) []T {
	dst = dst[:0]
	if len(s1) == 0 || len(s2) == 0 {
		return dst
	}

	i, j := 0, 0
	for i < len(s1) && j < len(s2) {
		if s1[i] < s2[j] {
			i++
		} else if s1[i] > s2[j] {
			j++
		} else {
			// Count occurrences in both slices
			current := s1[i]
			count1, count2 := 0, 0
			for i < len(s1) && s1[i] == current {
				count1++
				i++
			}
			for j < len(s2) && s2[j] == current {
				count2++
				j++
			}
			// Add the minimum of the two counts to the destination
			intersectCount := count1
			if count2 < intersectCount {
				intersectCount = count2
			}
			for k := 0; k < intersectCount; k++ {
				dst = append(dst, current)
			}
		}
	}
	return dst
}

// IntersectSortedSet compares two sorted slices s1 and s2, puts elements that are present in both into dst, and returns it.
// Its behavior differs from Intersect. For example: s1 = [1,2,2,3], s2 = [2,2,2,4], the result will be [2].
// The result will not contain duplicate elements.
// It requires ascending sorted slices.
func IntersectSortedSet[T typez.Ordered](dst, s1, s2 []T) []T {
	dst = dst[:0]
	if len(s1) == 0 || len(s2) == 0 {
		return dst
	}

	i, j := 0, 0
	for i < len(s1) && j < len(s2) {
		if s1[i] < s2[j] {
			i++
		} else if s1[i] > s2[j] {
			j++
		} else {
			if len(dst) == 0 || dst[len(dst)-1] != s1[i] {
				dst = append(dst, s1[i])
			}
			i++
			j++
		}
	}
	return dst
}

// Unique compares slice s, puts unique elements into dst, and returns it.
func Unique[T comparable](dst, s []T) []T {
	dst = dst[:0]
	if len(s) == 0 {
		return dst
	}

	if loopEnabled && len(s) <= loopThreshold {
		for i, v := range s {
			if Index(s[:i], v) < 0 {
				dst = append(dst, v)
			}
		}
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

	var remain int

	if loopEnabled && len(s) <= loopThreshold {
		for i := range s {
			if Index(s[:remain], s[i]) < 0 {
				s[remain], s[i] = s[i], s[remain]
				remain++
			}
		}
		return s[:remain]
	}

	seen := make(map[T]struct{}, len(s))
	var uniqueCount int
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

// UniqueSorted puts unique elements from a sorted slice s into dst and returns it.
// It requires ascending sorted slices.
func UniqueSorted[T comparable](dst, s []T) []T {
	dst = dst[:0]
	if len(s) == 0 {
		return dst
	}

	dst = append(dst, s[0])
	for i := 1; i < len(s); i++ {
		if s[i] != s[i-1] {
			dst = append(dst, s[i])
		}
	}
	return dst
}

// UniqueInPlaceSorted moves unique elements in a sorted slice s to the front and returns this portion of s.
// It requires ascending sorted slices.
//
// WARNING: Unlike other InPlace functions that only reorder elements, this function
// overwrites duplicate elements in the original slice. Elements beyond the returned
// slice length are not preserved.
// If you need to preserve all original elements, consider using UniqueSorted instead.
func UniqueInPlaceSorted[T comparable](s []T) []T {
	if len(s) == 0 {
		return s[:0]
	}

	j := 1
	for i := 1; i < len(s); i++ {
		if s[i] != s[j-1] {
			s[j] = s[i]
			j++
		}
	}
	return s[:j]
}

// UniqueByKey through keyFn get the key of slice s, puts unique elements into dst, and returns it.
func UniqueByKey[T any, K comparable](dst, s []T, keyFn func(T) K) []T {
	dst = dst[:0]
	if len(s) == 0 {
		return dst
	}

	if loopEnabled && len(s) <= loopThreshold {
		for i, v := range s {
			key := keyFn(v)
			found := false
			for j := 0; j < i; j++ {
				if keyFn(s[j]) == key {
					found = true
					break
				}
			}
			if !found {
				dst = append(dst, v)
			}
		}
		return dst
	}

	seen := make(map[K]struct{}, len(s))
	var uniqueCount int
	for _, v := range s {
		k := keyFn(v)
		seen[k] = struct{}{}
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

	var remain int

	if loopEnabled && len(s) <= loopThreshold {
		for i := range s {
			key := keyFn(s[i])
			found := false
			for j := 0; j < remain; j++ {
				if keyFn(s[j]) == key {
					found = true
					break
				}
			}
			if !found {
				s[remain], s[i] = s[i], s[remain]
				remain++
			}
		}
		return s[:remain]
	}

	seen := make(map[K]struct{}, len(s))
	var uniqueCount int
	for i := range s {
		k := keyFn(s[i])
		seen[k] = struct{}{}
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

// Equal compares slices s1 and s2 element-wise and returns true if they have the same length
// and all corresponding elements are equal.
func Equal[T comparable](s1, s2 []T) bool {
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
// If end is negative, it will return all elements from start to the end of s.
func SubSlice[T any](s []T, start int, end int) []T {
	if start >= len(s) {
		return s[:0]
	} else if start < 0 {
		start = 0
	}

	if end < 0 || end > len(s) {
		end = len(s)
	}

	if start >= end {
		return s[:0]
	}

	return s[start:end]
}

// Contains returns true if v is present in s.
func Contains[T comparable](s []T, v T) bool {
	return Index(s, v) >= 0
}

// ContainsFunc returns true if there is an element in s that satisfies f(s[i]).
func ContainsFunc[T any](s []T, fn func(T) bool) bool {
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

	n := (len(s) + chunkSize - 1) / chunkSize
	chunks := make([][]T, 0, n)
	start := 0
	for start < len(s) {
		end := start + chunkSize
		if end > len(s) {
			end = len(s)
		}
		chunks = append(chunks, s[start:end])
		start = end
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

// Remove removes the element at the given index from s.
func Remove[T any](s []T, index int) ([]T, T, bool) {
	var zero T
	if index < 0 || index >= len(s) {
		return s, zero, false
	}

	last := len(s) - 1
	v := s[index]
	if index < last {
		copy(s[index:], s[index+1:])
	}
	s[last] = zero
	return s[:last], v, true
}

// BinarySearch performs a binary search on a sorted slice s for the value v.
func BinarySearch[T typez.Ordered](s []T, v T) (index int, found bool) {
	low, high := 0, len(s)-1
	if high < 0 {
		return 0, false
	}

	if s[low] > v {
		return 0, false
	}

	if s[high] < v {
		return len(s), false
	}

	for low <= high {
		mid := int(uint(low+high) >> 1)
		if s[mid] < v {
			low = mid + 1
		} else if s[mid] > v {
			high = mid - 1
		} else {
			return mid, true
		}
	}
	return low, false
}

// BinarySearchFunc performs a binary search on a sorted slice s using a custom comparison function fn for the value v.
func BinarySearchFunc[T, V any](s []T, v V, fn func(T, V) int) (index int, found bool) {
	low, high := 0, len(s)-1
	if high < 0 {
		return 0, false
	}

	if fn(s[low], v) > 0 {
		return 0, false
	}

	if fn(s[high], v) < 0 {
		return len(s), false
	}

	for low <= high {
		mid := int(uint(low+high) >> 1)
		cmp := fn(s[mid], v)
		if cmp < 0 {
			low = mid + 1
		} else if cmp > 0 {
			high = mid - 1
		} else {
			return mid, true
		}
	}
	return low, false
}
