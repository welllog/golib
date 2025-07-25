package slicez

import "strconv"

type FlexSlice1[T any] struct {
	values []T
	head   int
	tail   int
	len    int
}

func NewFlex1[T any](capacity int) FlexSlice1[T] {
	var f FlexSlice1[T]
	f.Init(capacity)
	return f
}

// Init initializes or clears the flex slice with the given capacity.
func (f *FlexSlice1[T]) Init(capacity int) {
	if capacity < 0 {
		panic("slicez.FlexSlice1 Init: invalid capacity: " + strconv.Itoa(capacity))
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
func (f *FlexSlice1[T]) Len() int {
	return f.len
}

// Cap returns the capacity of the flex slice.
func (f *FlexSlice1[T]) Cap() int {
	return len(f.values)
}

// IsEmpty returns true if the flex slice is empty.
func (f *FlexSlice1[T]) IsEmpty() bool {
	return f.len == 0
}

// IsFull returns true if the flex slice is full.
func (f *FlexSlice1[T]) IsFull() bool {
	return f.len == len(f.values)
}

// Append adds elements to the end of the flex slice, growing it if necessary.
func (f *FlexSlice1[T]) Append(values ...T) {
	if len(f.values) < f.len+len(values) {
		f.grow(len(values))
	}

	for _, v := range values {
		f.values[f.tail] = v
		f.tail = (f.tail + 1) % len(f.values)
	}
	f.len += len(values)
}

// Prepend adds elements to the front of the flex slice, growing it if necessary.
func (f *FlexSlice1[T]) Prepend(values ...T) {
	if len(f.values) < f.len+len(values) {
		f.grow(len(values))
	}

	for i := len(values) - 1; i >= 0; i-- {
		f.head = (f.head - 1 + len(f.values)) % len(f.values)
		f.values[f.head] = values[i]
	}
	f.len += len(values)
}

// Pop removes and returns the last element of the flex slice.
func (f *FlexSlice1[T]) Pop() (T, bool) {
	var zero T
	if f.len == 0 {
		return zero, false
	}

	f.tail = (f.tail - 1 + len(f.values)) % len(f.values)
	value := f.values[f.tail]
	f.values[f.tail] = zero // Clear the value
	f.len--

	return value, true
}

// Shift removes and returns the first element of the flex slice.
func (f *FlexSlice1[T]) Shift() (T, bool) {
	var zero T
	if f.len == 0 {
		return zero, false
	}

	value := f.values[f.head]
	f.values[f.head] = zero // Clear the value
	f.head = (f.head + 1) % len(f.values)
	f.len--

	return value, true
}

// Get returns the element at the given index in the deque.
func (f *FlexSlice1[T]) Get(index int) (T, bool) {
	var zero T
	if f.withoutRange(index) {
		return zero, false
	}

	actualIndex := (f.head + index) % len(f.values)
	return f.values[actualIndex], true
}

// Front returns the first element of the flex slice without removing it.
func (f *FlexSlice1[T]) Front() (T, bool) {
	var zero T
	if f.len == 0 {
		return zero, false
	}

	return f.values[f.head], true
}

// Back returns the last element of the flex slice without removing it.
func (f *FlexSlice1[T]) Back() (T, bool) {
	var zero T
	if f.len == 0 {
		return zero, false
	}

	tailIndex := (f.tail - 1 + len(f.values)) % len(f.values)
	return f.values[tailIndex], true
}

// Set sets the element at the given index in the flex slice to the specified value.
func (f *FlexSlice1[T]) Set(index int, value T) bool {
	if f.withoutRange(index) {
		return false
	}

	actualIndex := (f.head + index) % len(f.values)
	f.values[actualIndex] = value
	return true
}

// Remove removes and returns the element at the specified index in the flex slice.
func (f *FlexSlice1[T]) Remove(index int) (T, bool) {
	var zero T
	if f.withoutRange(index) {
		return zero, false
	}

	//if index == 0 {
	//	return f.Shift()
	//}
	//
	//if index == f.len-1 {
	//	return f.Pop()
	//}

	actualIndex := (f.head + index) % len(f.values)
	value := f.values[actualIndex]

	// select move less expensive way to move elements
	if index < f.len/2 {
		// move front part
		for i := actualIndex; i != f.head; i = (i - 1 + len(f.values)) % len(f.values) {
			prevIndex := (i - 1 + len(f.values)) % len(f.values)
			f.values[i] = f.values[prevIndex]
		}
		f.values[f.head] = zero
		f.head = (f.head + 1) % len(f.values)
	} else {
		// move back part
		for i := actualIndex; i != f.tail; i = (i + 1) % len(f.values) {
			nextIndex := (i + 1) % len(f.values)
			f.values[i] = f.values[nextIndex]
		}
		f.tail = (f.tail - 1 + len(f.values)) % len(f.values)
		f.values[f.tail] = zero
	}
	f.len--

	return value, true
}

// InsertAt inserts a value at the specified index in the flex slice, shifting elements as necessary.
func (f *FlexSlice1[T]) InsertAt(index int, value T) {
	//if index <= 0 {
	//	f.Prepend(value)
	//	return
	//}
	//
	//if index >= f.len {
	//	f.Append(value)
	//	return
	//}

	actualIndex := (f.head + index) % len(f.values)
	f.len++
	if f.len == len(f.values) {
		newCap := f.calCapacity(1)
		newValues := make([]T, newCap)
		if f.head < f.tail {
			copy(newValues, f.values[f.head:actualIndex])
			newValues[index] = value
			copy(newValues[index+1:], f.values[actualIndex:f.tail])
		} else {
			if actualIndex >= f.head {
				copy(newValues, f.values[f.head:actualIndex])
				newValues[index] = value
				n := copy(newValues[index+1:], f.values[actualIndex:])
				copy(newValues[n+index+1:], f.values[:f.tail])
			} else {
				n := copy(newValues, f.values[f.head:])
				copy(newValues[n:], f.values[:actualIndex])
				newValues[index] = value
				copy(newValues[index+1:], f.values[:f.tail])
			}
		}

		f.values = newValues
		f.head = 0
		f.tail = f.len
		return
	}

	if index < f.len/2 {
		actualIndex = (f.head + index - 1) % len(f.values)
		// move front part
		f.head = (f.head - 1 + len(f.values)) % len(f.values)
		for i := f.head; i != actualIndex; i = (i + 1) % len(f.values) {
			nextIndex := (i + 1) % len(f.values)
			f.values[i] = f.values[nextIndex]
		}
	} else {
		// move back part
		for i := f.tail; i != actualIndex; i = (i - 1 + len(f.values)) % len(f.values) {
			prevIndex := (i - 1 + len(f.values)) % len(f.values)
			f.values[i] = f.values[prevIndex]
		}
		f.tail = (f.tail + 1) % len(f.values)
	}

	f.values[actualIndex] = value

	return
}

// ToSlice returns a slice containing all elements in the flex slice in order.
func (f *FlexSlice1[T]) ToSlice() []T {
	result := make([]T, f.len)
	for i := 0; i < f.len; i++ {
		result[i] = f.values[(f.head+i)%len(f.values)]
	}

	return result
}

// Range iterates over the elements in the flex slice, calling the provided function for each element.
func (f *FlexSlice1[T]) Range(fn func(index int, value T) bool) {
	for i := 0; i < f.len; i++ {
		if !fn(i, f.values[(f.head+i)%len(f.values)]) {
			break
		}
	}
}

// Shrink reduces the capacity of the flex slice if it is more than 4 times larger than the number of elements.
func (f *FlexSlice1[T]) Shrink() bool {
	if len(f.values) <= 8 || f.len > len(f.values)/4 {
		return false
	}

	newCap := len(f.values) / 2
	if newCap < 8 {
		newCap = 8
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
	return true
}

// Clear clears the flex slice, optionally shrinking its capacity.
func (f *FlexSlice1[T]) Clear(shrink bool) {
	oldLen := f.len
	f.len = 0

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

func (f *FlexSlice1[T]) calCapacity(n int) int {
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

func (f *FlexSlice1[T]) grow(n int) {
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

func (f *FlexSlice1[T]) withoutRange(index int) bool {
	return index < 0 || index >= f.len
}
