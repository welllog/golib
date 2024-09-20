package heapz

type Slice[T any] struct {
	Values []T
	cmp    func(T, T) bool
}

// NewSlice returns a new heap with the given compare function.
func NewSlice[T any](cap int, cmp func(T, T) bool) Slice[T] {
	s := make([]T, 0, cap)
	return FromSlice(s, cmp)
}

// FromSlice creates a new heap from a slice.
func FromSlice[T any](s []T, cmp func(T, T) bool) Slice[T] {
	ss := Slice[T]{Values: s, cmp: cmp}
	build(ss.Values, ss.cmp, swap[T])
	return ss
}

// Push pushes the element x onto the heap.
func (s *Slice[T]) Push(x T) {
	s.Values = append(s.Values, x)
	up(s.Values, s.cmp, swap[T], len(s.Values)-1)
}

// Pop removes and returns the minimum element (according to compare function) from the heap.
func (s *Slice[T]) Pop() (T, bool) {
	var x, zero T

	n := len(s.Values)
	if n == 0 {
		return x, false
	} else if n == 1 {
		x = s.Values[0]
		s.Values[0] = zero
		s.Values = s.Values[:0]
		return x, true
	}

	n--
	s.Values[0], s.Values[n] = s.Values[n], s.Values[0]
	down(s.Values, s.cmp, swap[T], 0, n)
	x = s.Values[n]
	s.Values[n] = zero
	s.Values = s.Values[:n]
	return x, true
}

// Peek returns the minimum element (according to compare function) from the heap without removing it.
func (s *Slice[T]) Peek() (T, bool) {
	if len(s.Values) == 0 {
		var x T
		return x, false
	}

	return s.Values[0], true
}

// Len returns the number of elements in the heap.
func (s *Slice[T]) Len() int {
	return len(s.Values)
}

// Remove removes and returns the element at index i from the slice.
func (s *Slice[T]) Remove(i int) (T, bool) {
	var x, zero T
	if i < 0 || i >= len(s.Values) {
		return x, false
	}

	n := len(s.Values) - 1
	if n != i {
		s.Values[i], s.Values[n] = s.Values[n], s.Values[i]
		fix(s.Values, s.cmp, swap[T], i, n)
	}

	x = s.Values[n]
	s.Values[n] = zero
	s.Values = s.Values[:n]
	return x, true
}

// Fix re-establishes the heap ordering after the element at index i has changed its value.
func (s *Slice[T]) Fix(i int) {
	if i < 0 || i >= len(s.Values) {
		return
	}

	fix(s.Values, s.cmp, swap[T], i, len(s.Values))
}
