package slicez

type FlexSlice[T any] struct {
	Values []T
}

// Append adds elements to the end of the slice.
func (f *FlexSlice[T]) Append(v ...T) {
	f.Values = append(f.Values, v...)
}

// Prepend adds elements to the beginning of the slice.
func (f *FlexSlice[T]) Prepend(v ...T) {
	n1 := len(v)
	n2 := len(f.Values)
	c := cap(f.Values)
	nc := n1 + n2
	if c >= nc {
		f.Values = f.Values[:nc]
		copy(f.Values[n1:], f.Values[:n2])
		copy(f.Values, v)
		return
	}

	if 2*c >= nc {
		c = 2 * c
	} else {
		c = nc
	}

	newValues := make([]T, nc, c)
	copy(newValues, v)
	copy(newValues[n1:], f.Values)
	f.Values = newValues
}

// Get returns the element at the given index.
func (f *FlexSlice[T]) Get(index int) (T, bool) {
	if f.withinRange(index) {
		return f.Values[index], true
	}

	var zero T
	return zero, false
}

// Remove removes the element at the given index.
// This operation is O(n) because it shifts all elements to the left.
func (f *FlexSlice[T]) Remove(index int) (T, bool) {
	var v, zero T

	if !f.withinRange(index) {
		return zero, false
	}

	l := len(f.Values) - 1
	v = f.Values[index]
	if index < l {
		copy(f.Values[index:], f.Values[index+1:])
	}
	f.Values[l] = zero
	f.Values = f.Values[:l]

	f.shrink()

	return v, true
}

// SubSlice returns a slice of the FlexSlice from start to end.
// If end is negative, it will return all elements from start to the end of the slice.
func (f *FlexSlice[T]) SubSlice(start, end int) FlexSlice[T] {
	nf := FlexSlice[T]{Values: SubSlice(f.Values, start, end)}
	nf.shrink()
	return nf
}

// Pop removes and returns the last element from the slice.
// This operation is O(1).
func (f *FlexSlice[T]) Pop() (T, bool) {
	return f.Remove(len(f.Values) - 1)
}

// Shift removes and returns the first element from the slice.
// This operation is O(n) because it shifts all elements to the left.
func (f *FlexSlice[T]) Shift() (T, bool) {
	return f.Remove(0)
}

// Len returns the number of elements in the slice.
func (f *FlexSlice[T]) Len() int {
	return len(f.Values)
}

func (f *FlexSlice[T]) shrink() {
	if cap(f.Values) <= 8 {
		return
	}

	if len(f.Values) <= cap(f.Values)/4 {
		newCap := len(f.Values) * 2
		if newCap < 8 {
			newCap = 8
		}

		newValues := make([]T, len(f.Values), newCap)
		copy(newValues, f.Values)
		f.Values = newValues
	}
}

func (f *FlexSlice[T]) withinRange(index int) bool {
	return index >= 0 && index < len(f.Values)
}
