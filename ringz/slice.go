package ringz

type Slice[T any] struct {
	values []T
	head int
	tail int
	cap int
}

// NewSlice returns a new slice with the given capacity.
func NewSlice[T any](cap int) Slice[T] {
	var s Slice[T]
	s.Init(cap)
	return s
}

// Init initializes or clears the slice s.
func (s *Slice[T]) Init(cap int) {
	s.values = make([]T, cap)
	s.head = -1
	s.tail = -1
	s.cap = cap
}

// IsEmpty returns true if the slice is empty.
func (s *Slice[T]) IsEmpty() bool {
	return s.head == -1
}

// IsFull returns true if the slice is full.
func (s *Slice[T]) IsFull() bool {
	return (s.tail + 1) % s.cap == s.head
}

// Push pushes the value to queue tail.
func (s *Slice[T]) Push(value T) bool {
	if s.IsFull() {
		return false
	}

	if s.IsEmpty() {
		s.head = 0
	}

	s.tail = (s.tail + 1) % s.cap
	s.values[s.tail] = value
	return true
}

// Pop removes and returns the value from queue head.
func (s *Slice[T]) Pop() (T, bool) {
	var zero T
	if s.IsEmpty() {
		return zero, false
	}

	value := s.values[s.head]
	s.values[s.head] = zero
	if s.head == s.tail {
		s.head = -1
		s.tail = -1
	} else {
		s.head = (s.head + 1) % s.cap
	}
	return value, true
}

// Peek returns the value from queue head without removing it.
func (s *Slice[T]) Peek() (T, bool) {
	if s.IsEmpty() {
		var zero T
		return zero, false
	}

	return s.values[s.head], true
}

// Len returns the number of elements in the slice.
func (s *Slice[T]) Len() int {
	if s.IsEmpty() {
        return 0
    }

    if s.head <= s.tail {
        return s.tail - s.head + 1
    }

    return s.cap - s.head + s.tail + 1
}

// Recap changes the capacity of the slice.
func (s *Slice[T]) Recap(cap int) bool {
	if cap <= 0 || cap == s.cap {
		return false
	}

	l := s.Len()
	if cap < l {
		return false
	}

	newValues := make([]T, cap)
	if s.IsEmpty() {
		s.values = newValues
		s.cap = cap
		s.head = -1
		s.tail = -1
		return true
	}

	if s.head <= s.tail {
		copy(newValues, s.values[s.head:s.tail+1])
	} else {
		n := copy(newValues, s.values[s.head:])
        copy(newValues[n:], s.values[:s.tail+1])
	}

	s.head = 0
    s.tail = l - 1
    s.values = newValues
    s.cap = cap
	return true
}
