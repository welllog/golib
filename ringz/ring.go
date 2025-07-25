package ringz

import "strconv"

type Ring[T any] struct {
	values []T
	head   int
	tail   int
	cap    int
}

// New returns a new ring with the given capacity.
func New[T any](cap int) Ring[T] {
	var r Ring[T]
	r.Init(cap)
	return r
}

// Init initializes or clears the ring.
func (r *Ring[T]) Init(cap int) {
	if cap <= 0 {
		panic("ringz.Ring Init: invalid capacity: " + strconv.Itoa(cap))
	}

	r.values = make([]T, cap)
	r.head = -1
	r.tail = -1
	r.cap = cap
}

// IsEmpty returns true if the ring is empty.
func (r *Ring[T]) IsEmpty() bool {
	return r.head == -1
}

// IsFull returns true if the ring is full.
func (r *Ring[T]) IsFull() bool {
	return (r.tail+1)%r.cap == r.head
}

// Push pushes the value to queue tail.
func (r *Ring[T]) Push(value T) bool {
	if r.IsFull() {
		return false
	}

	if r.IsEmpty() {
		r.head = 0
	}

	r.tail = (r.tail + 1) % r.cap
	r.values[r.tail] = value
	return true
}

// Pop removes and returns the value from queue head.
func (r *Ring[T]) Pop() (T, bool) {
	var zero T
	if r.IsEmpty() {
		return zero, false
	}

	value := r.values[r.head]
	r.values[r.head] = zero
	if r.head == r.tail {
		r.head = -1
		r.tail = -1
	} else {
		r.head = (r.head + 1) % r.cap
	}
	return value, true
}

// Peek returns the value from queue head without removing it.
func (r *Ring[T]) Peek() (T, bool) {
	if r.IsEmpty() {
		var zero T
		return zero, false
	}

	return r.values[r.head], true
}

// PushWithGrow pushes the value to queue tail and expands the ring if it is full.
func (r *Ring[T]) PushWithGrow(value T) {
	if r.IsFull() {
		r.Recap(r.cap * 2)
	}

	r.Push(value)
}

// Len returns the number of elements in the ring.
func (r *Ring[T]) Len() int {
	if r.IsEmpty() {
		return 0
	}

	if r.head <= r.tail {
		return r.tail - r.head + 1
	}

	return r.cap - r.head + r.tail + 1
}

// Cap returns the capacity of the ring.
func (r *Ring[T]) Cap() int {
	return r.cap
}

// Recap changes the capacity of the ring.
func (r *Ring[T]) Recap(cap int) bool {
	if cap <= 0 || cap == r.cap {
		return false
	}

	l := r.Len()
	if cap < l {
		return false
	}

	newValues := make([]T, cap)
	if r.IsEmpty() {
		r.values = newValues
		r.cap = cap
		r.head = -1
		r.tail = -1
		return true
	}

	if r.head <= r.tail {
		copy(newValues, r.values[r.head:r.tail+1])
	} else {
		n := copy(newValues, r.values[r.head:])
		copy(newValues[n:], r.values[:r.tail+1])
	}

	r.head = 0
	r.tail = l - 1
	r.values = newValues
	r.cap = cap
	return true
}
