package setz

// Set is a set of T.
type Set[T comparable] map[T]struct{}

// Add adds v to s.
func (s Set[T]) Add(v T) (first bool) {
	if _, ok := s[v]; ok {
		return false
	}
	s[v] = struct{}{}
	return true
}

// AddAll adds all values in vs to the set.
func (s Set[T]) AddAll(vs ...T) {
	for _, v := range vs {
		s[v] = struct{}{}
	}
}

// Has reports whether v is in s.
func (s Set[T]) Has(v T) bool {
	_, ok := s[v]
	return ok
}

// Contains reports whether v is in s.
func (s Set[T]) Contains(v T) bool {
	_, ok := s[v]
	return ok
}

// Delete deletes v from s.
func (s Set[T]) Delete(v T) (exists bool) {
	if _, ok := s[v]; !ok {
		return false
	}
	delete(s, v)
	return true
}

// Len returns the length of s.
func (s Set[T]) Len() int {
	return len(s)
}

// Values returns all values in s.
func (s Set[T]) Values(dst []T) []T {
	for k := range s {
		dst = append(dst, k)
	}
	return dst
}

// Range calls fn sequentially for each value present in the set.
func (s Set[T]) Range(fn func(T)) {
	for k := range s {
		fn(k)
	}
}

// Clear clears s.
func (s Set[T]) Clear() {
	for k := range s {
		delete(s, k)
	}
}

// Filter deletes all values in s that do not satisfy fn.
func (s Set[T]) Filter(fn func(T) bool) {
	for k := range s {
		if !fn(k) {
			delete(s, k)
		}
	}
}

// Merge adds all values in other to s.
func (s Set[T]) Merge(other Set[T]) {
	for k := range other {
		s[k] = struct{}{}
	}
}

// Diff deletes all values in s that are also in the other.
func (s Set[T]) Diff(other Set[T]) {
	for k := range other {
		delete(s, k)
	}
}

// Intersect deletes all values in s that are not in the other.
func (s Set[T]) Intersect(other Set[T]) {
	for k := range s {
		if _, ok := other[k]; !ok {
			delete(s, k)
		}
	}
}

// DiffWithSlice deletes all values in s that are also in the other.
func (s Set[T]) DiffWithSlice(other []T) {
	for _, v := range other {
		delete(s, v)
	}
}

// IntersectWithSlice deletes all values in s that are not in the other.
func (s Set[T]) IntersectWithSlice(other []T) {
	tmp := make(Set[T], len(other))
	tmp.AddAll(other...)

	for k := range s {
		if _, ok := tmp[k]; !ok {
			delete(s, k)
		}
	}
}

// AddMapKeysToSet adds all keys in m to s.
func AddMapKeysToSet[T comparable, V any](s Set[T], m map[T]V) {
	for k := range m {
		s[k] = struct{}{}
	}
}

// AddMapValuesToSet adds all values in m to s.
func AddMapValuesToSet[T comparable, K comparable](s Set[T], m map[K]T) {
	for _, v := range m {
		s[v] = struct{}{}
	}
}

// FromSlice returns a new set containing all values in src.
func FromSlice[T comparable](src []T) Set[T] {
	dst := make(Set[T], len(src))
	for _, v := range src {
		dst[v] = struct{}{}
	}

	return dst
}

// FromMapKeys returns a new set containing all keys in src.
func FromMapKeys[T comparable, V any](src map[T]V) Set[T] {
	dst := make(Set[T], len(src))
	for k := range src {
		dst[k] = struct{}{}
	}

	return dst
}

// FromMapValues returns a new set containing all values in src.
func FromMapValues[T comparable, K comparable](m map[K]T) Set[T] {
	dst := make(Set[T], len(m))
	for _, v := range m {
		dst[v] = struct{}{}
	}

	return dst
}
