package dsz

// Set[T comparable] is a set of T.
type Set[T comparable] map[T]struct{}

// Add adds v to s.
func (s Set[T]) Add(v T) (first bool) {
	if _, ok := s[v]; ok {
		return false
	}
	s[v] = struct{}{}
	return true
}

// MultiAdd adds vs to s.
func (s Set[T]) MultiAdd(vs ...T) {
	for _, v := range vs {
		s[v] = struct{}{}
	}
}

// Has reports whether v is in s.
func (s Set[T]) Has(v T) bool {
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

// Diff deletes all values in s that are also in other.
func (s Set[T]) Diff(other Set[T]) {
	for k := range other {
		delete(s, k)
	}
}

// Intersect deletes all values in s that are not in other.
func (s Set[T]) Intersect(other Set[T]) {
	for k := range s {
		if _, ok := other[k]; !ok {
			delete(s, k)
		}
	}
}

// DiffWithSlice deletes all values in s that are also in other.
func (s Set[T]) DiffWithSlice(other []T) {
	for _, v := range other {
		delete(s, v)
	}
}

// IntersectWithSlice deletes all values in s that are not in other.
func (s Set[T]) IntersectWithSlice(other []T) {
	tmp := make(Set[T], len(other))
	tmp.MultiAdd(other...)

	for k := range s {
		if _, ok := tmp[k]; !ok {
			delete(s, k)
		}
	}
}

// SetFromSlice fills dst with src
func SetFromSlice[T comparable](dst Set[T], src []T) {
	for _, v := range src {
		dst[v] = struct{}{}
	}
}

// SetFromMapKeys fills dst with keys from src
func SetFromMapKeys[T comparable, V any](dst Set[T], src map[T]V) {
	for k := range src {
		dst[k] = struct{}{}
	}
}

// SetFromMapValues fills dst with values from src
func SetFromMapValues[T comparable, K comparable](dst Set[T], m map[K]T) {
	for _, v := range m {
		dst[v] = struct{}{}
	}
}
