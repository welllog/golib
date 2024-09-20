package slicez

type FlexSlice[T any] struct {
	Values []T
}

// Append adds elements to the end of the slice.
func (f *FlexSlice[T]) Append(v ...T) {
	f.Values = append(f.Values, v...)
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
func (f *FlexSlice[T]) Remove(index int) (T, bool) {
	var v, zero T

	if !f.withinRange(index) {
		return zero, false
	}

	v = f.Values[index]
	f.Values[index] = zero
	copy(f.Values[index:], f.Values[index+1:])

	f.shrink()

	return v, true
}

// SubSlice returns a slice of the FlexSlice from start to end.
func (f *FlexSlice[T]) SubSlice(start, end int) FlexSlice[T] {
	nf := FlexSlice[T]{Values: SubSlice(f.Values, start, end)}
	nf.shrink()
	return nf
}

// Pop removes and returns the last element from the slice.
func (f *FlexSlice[T]) Pop() (T, bool) {
	return f.Remove(len(f.Values) - 1)
}

// Shift removes and returns the first element from the slice.
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
