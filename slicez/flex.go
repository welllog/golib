package slicez

import "strconv"

type FlexSlice[T any] struct {
	values []T
	head   int
	tail   int
	len    int
}

func NewFlexSlice[T any](capacity int) FlexSlice[T] {
	var f FlexSlice[T]
	f.Init(capacity)
	return f
}

// Init initializes or clears the flex slice with the given capacity.
func (f *FlexSlice[T]) Init(capacity int) {
	if capacity < 0 {
		panic("slicez.FlexSlice Init: invalid capacity: " + strconv.Itoa(capacity))
	}

	if capacity == 0 {
		f.values = nil
	} else {
		f.values = make([]T, capacity)
	}
	f.head = 0
	f.tail = 0
	f.len = 0
}

// Len returns the number of elements in the flex slice.
func (f *FlexSlice[T]) Len() int {
	return f.len
}

// Cap returns the capacity of the flex slice.
func (f *FlexSlice[T]) Cap() int {
	return len(f.values)
}

// IsEmpty returns true if the flex slice is empty.
func (f *FlexSlice[T]) IsEmpty() bool {
	return f.len == 0
}

// IsFull returns true if the flex slice is full.
func (f *FlexSlice[T]) IsFull() bool {
	return f.len == len(f.values)
}

// Append adds elements to the end of the flex slice, growing it if necessary.
func (f *FlexSlice[T]) Append(values ...T) {
	f.tryGrow(len(values))

	for _, v := range values {
		f.values[f.tail] = v
		f.tail = f.backwardIndex(f.tail + 1)
	}
	f.len += len(values)
}

// Prepend adds elements to the front of the flex slice, growing it if necessary.
func (f *FlexSlice[T]) Prepend(values ...T) {
	f.tryGrow(len(values))

	for i := len(values) - 1; i >= 0; i-- {
		f.head = f.forwardIndex(f.head - 1)
		f.values[f.head] = values[i]
	}
	f.len += len(values)
}

// Pop removes and returns the last element of the flex slice.
func (f *FlexSlice[T]) Pop() (T, bool) {
	var zero T
	if f.len == 0 {
		return zero, false
	}

	f.tail = f.forwardIndex(f.tail - 1)
	value := f.values[f.tail]
	f.values[f.tail] = zero // Clear the value
	f.len--

	return value, true
}

// Shift removes and returns the first element of the flex slice.
func (f *FlexSlice[T]) Shift() (T, bool) {
	var zero T
	if f.len == 0 {
		return zero, false
	}

	value := f.values[f.head]
	f.values[f.head] = zero // Clear the value
	f.head = f.backwardIndex(f.head + 1)
	f.len--

	return value, true
}

// Front returns the first element of the flex slice without removing it.
func (f *FlexSlice[T]) Front() (T, bool) {
	var zero T
	if f.len == 0 {
		return zero, false
	}

	return f.values[f.head], true
}

// Back returns the last element of the flex slice without removing it.
func (f *FlexSlice[T]) Back() (T, bool) {
	var zero T
	if f.len == 0 {
		return zero, false
	}

	tailIndex := f.forwardIndex(f.tail - 1)
	return f.values[tailIndex], true
}

// Get returns the element at the given index in the deque.
func (f *FlexSlice[T]) Get(index int) (T, bool) {
	var zero T
	if f.withoutRange(index) {
		return zero, false
	}

	actualIndex := f.backwardIndex(f.head + index)
	return f.values[actualIndex], true
}

// Set sets the element at the given index in the flex slice to the specified value.
func (f *FlexSlice[T]) Set(index int, value T) bool {
	if f.withoutRange(index) {
		return false
	}

	actualIndex := f.backwardIndex(f.head + index)
	f.values[actualIndex] = value
	return true
}

// Remove removes and returns the element at the specified index in the flex slice.
func (f *FlexSlice[T]) Remove(index int) (T, bool) {
	var zero T
	if f.withoutRange(index) {
		return zero, false
	}

	if index == 0 {
		return f.Shift()
	}

	if index == f.len-1 {
		return f.Pop()
	}

	actualIndex := f.backwardIndex(f.head + index)
	value := f.values[actualIndex]

	// select move less expensive way to move elements
	if index < f.len/2 {
		// move front part
		fillIndex := actualIndex
		for i := index; i > 0; i-- {
			prevIndex := f.forwardIndex(fillIndex - 1)
			f.values[fillIndex] = f.values[prevIndex]
			fillIndex = prevIndex
		}

		f.values[f.head] = zero
		f.head = f.backwardIndex(f.head + 1)
	} else {
		// move back part
		fillIndex := actualIndex
		for i := index; i < (f.len - 1); i++ {
			nextIndex := f.backwardIndex(fillIndex + 1)
			f.values[fillIndex] = f.values[nextIndex]
			fillIndex = nextIndex
		}

		f.tail = f.forwardIndex(f.tail - 1)
		f.values[f.tail] = zero
	}
	f.len--

	return value, true
}

// InsertAt inserts elements at the specified index in the flex slice.
// If the index is out of bounds, it returns false.
func (f *FlexSlice[T]) InsertAt(index int, values ...T) bool {
	if index < 0 || index > f.len {
		return false
	}

	size := len(values)
	if size == 0 {
		return true
	}

	if index == 0 {
		f.Prepend(values...)
		return true
	}

	if index == f.len {
		f.Append(values...)
		return true
	}

	if f.len+size > len(f.values) {
		newCap := f.calCapacity(size)
		newValues := make([]T, newCap)
		actualIndex := index
		if f.len > 0 {
			actualIndex = f.backwardIndex(f.head + index)
		}

		if f.head < f.tail {
			n := copy(newValues, f.values[f.head:actualIndex])
			n += copy(newValues[n:], values)
			copy(newValues[n:], f.values[actualIndex:f.tail])
		} else {
			if actualIndex >= f.head {
				n := copy(newValues, f.values[f.head:actualIndex])
				n += copy(newValues[n:], values)
				n += copy(newValues[n:], f.values[actualIndex:])
				copy(newValues[n:], f.values[:f.tail])
			} else {
				n := copy(newValues, f.values[f.head:])
				n += copy(newValues[n:], f.values[:actualIndex])
				n += copy(newValues[n:], values)
				copy(newValues[n:], f.values[actualIndex:f.tail])
			}
		}

		f.values = newValues
		f.head = 0
		f.len += size
		f.tail = f.len
		return true
	}

	actualIndex := f.backwardIndex(f.head + index)
	if index < f.len/2 {
		// move front part
		migrateIndex := f.head
		for i := 0; i < index; i++ {
			dstIndex := f.forwardIndex(migrateIndex - size)
			f.values[dstIndex] = f.values[migrateIndex]
			migrateIndex = f.backwardIndex(migrateIndex + 1)
		}
		f.head = f.forwardIndex(f.head - size)
		actualIndex = f.forwardIndex(actualIndex - size)

	} else {
		// move back part
		migrateIndex := f.forwardIndex(f.tail - 1)
		for i := index; i < f.len; i++ {
			dstIndex := f.backwardIndex(migrateIndex + size)
			f.values[dstIndex] = f.values[migrateIndex]
			migrateIndex = f.forwardIndex(migrateIndex - 1)
		}
		f.tail = f.backwardIndex(f.tail + size)

	}

	for _, v := range values {
		f.values[actualIndex] = v
		actualIndex = f.backwardIndex(actualIndex + 1)
	}
	f.len += size

	return true
}

// Join appends all elements from another FlexSlice to the end of this FlexSlice.
func (f *FlexSlice[T]) Join(other FlexSlice[T]) {
	if other.len == 0 {
		return
	}

	if f.len+other.len > len(f.values) {
		f.grow(other.len)
	}

	for i := 0; i < other.len; i++ {
		f.values[f.tail] = other.values[(other.head+i)%len(other.values)]
		f.tail = f.backwardIndex(f.tail + 1)
	}
	f.len += other.len
}

// SubSlice returns a new FlexSlice containing elements from the specified range.
// If the end is negative, it will return all elements from start to the end of the flex slice.
// SubSlice reuses the original FlexSlice's values slice, so modify maybe affect the original slice.
func (f *FlexSlice[T]) SubSlice(start, end int) FlexSlice[T] {
	f2 := *f
	f2.tail = f2.head
	f2.len = 0

	if start >= f.len {
		return f2
	}

	if start < 0 {
		start = 0
	}

	if end < 0 || end > f.len {
		end = f.len
	}

	if start >= end {
		return f2
	}

	f2.head = f.backwardIndex(f2.head + start)
	newLen := end - start
	f2.len = newLen
	f2.tail = f.backwardIndex(f2.head + newLen)
	return f2
}

// ToSlice returns a slice containing all elements in the flex slice in order.
func (f *FlexSlice[T]) ToSlice() []T {
	result := make([]T, f.len)
	for i := 0; i < f.len; i++ {
		result[i] = f.values[(f.head+i)%len(f.values)]
	}

	return result
}

// Range iterates over the elements in the flex slice, calling the provided function for each element.
func (f *FlexSlice[T]) Range(fn func(index int, value T) bool) {
	for i := 0; i < f.len; i++ {
		if !fn(i, f.values[(f.head+i)%len(f.values)]) {
			break
		}
	}
}

// RevRange iterates over the elements in the flex slice in reverse order.
func (f *FlexSlice[T]) RevRange(fn func(index int, value T) bool) {
	for i := f.len - 1; i >= 0; i-- {
		if !fn(i, f.values[(f.head+i)%len(f.values)]) {
			break
		}
	}
}

// Shrink reduces the capacity of the flex slice if it is more than 4 times larger than the number of elements.
func (f *FlexSlice[T]) Shrink() bool {
	if len(f.values) <= 8 || f.len > len(f.values)/4 {
		return false
	}

	newCap := len(f.values) / 2
	if newCap < 8 {
		newCap = 8
	}

	newValues := make([]T, newCap)
	if f.len > 0 {
		if f.head < f.tail {
			copy(newValues, f.values[f.head:f.tail])
		} else {
			n := copy(newValues, f.values[f.head:])
			copy(newValues[n:], f.values[:f.tail])
		}
	}

	f.values = newValues
	f.head = 0
	f.tail = f.len
	return true
}

// Grow increases the capacity of the flex slice by the specified amount.
func (f *FlexSlice[T]) Grow(n uint) {
	f.tryGrow(int(n))
}

// Clear clears the flex slice, optionally shrinking its capacity.
func (f *FlexSlice[T]) Clear(shrink bool) {
	oldLen := f.len
	f.len = 0
	f.tail = f.head

	if shrink {
		if f.Shrink() {
			return
		}
	}

	var zero T
	for i := 0; i < oldLen; i++ {
		f.values[(f.head+i)%len(f.values)] = zero
	}
	f.head = 0
	f.tail = 0
}

func (f *FlexSlice[T]) tryGrow(n int) {
	if len(f.values) < f.len+n {
		f.grow(n)
	}
}

func (f *FlexSlice[T]) calCapacity(n int) int {
	minCap := f.len + n
	if len(f.values) == 0 {
		if minCap <= 8 {
			return 8
		}

		return minCap
	}

	oldCap := len(f.values)
	newCap := oldCap
	doubleCap := oldCap + oldCap

	if minCap > doubleCap {
		newCap = minCap
	} else {
		const threshold = 256
		if oldCap < threshold {
			newCap = doubleCap
		} else {
			for 0 < newCap && newCap < minCap {
				newCap += (newCap + 3*threshold) / 4
			}

			// Set new cap to the requested cap when
			// the new cap calculation overflowed.
			if newCap <= 0 {
				newCap = minCap
			}
		}
	}
	return newCap
}

func (f *FlexSlice[T]) grow(n int) {
	minCap := f.len + n
	if len(f.values) == 0 {
		if minCap <= 8 {
			f.values = make([]T, 8)
			return
		}

		f.values = make([]T, minCap)
		return
	}

	oldCap := len(f.values)
	newCap := oldCap
	doubleCap := oldCap + oldCap

	if minCap > doubleCap {
		newCap = minCap
	} else {
		const threshold = 256
		if oldCap < threshold {
			newCap = doubleCap
		} else {
			for 0 < newCap && newCap < minCap {
				newCap += (newCap + 3*threshold) / 4
			}

			// Set new cap to the requested cap when
			// the new cap calculation overflowed.
			if newCap <= 0 {
				newCap = minCap
			}
		}
	}

	newValues := make([]T, newCap)
	if f.head < f.tail {
		copy(newValues, f.values[f.head:f.tail])
	} else {
		n := copy(newValues, f.values[f.head:])
		copy(newValues[n:], f.values[:f.tail])
	}

	f.values = newValues
	f.head = 0
	f.tail = f.len
}

func (f *FlexSlice[T]) withoutRange(index int) bool {
	return index < 0 || index >= f.len
}

func (f *FlexSlice[T]) backwardIndex(index int) int {
	return index % len(f.values)
}

func (f *FlexSlice[T]) forwardIndex(index int) int {
	return (index + len(f.values)) % len(f.values)
}
