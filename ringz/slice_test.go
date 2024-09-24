package ringz

import "testing"

func TestIsEmpty(t *testing.T) {
	// Test when the slice is newly created
	s := NewSlice[int](3)
	if !s.IsEmpty() {
		t.Errorf("expected IsEmpty to be true, got false")
	}

	// Test when the slice has one element
	s.Push(1)
	if s.IsEmpty() {
		t.Errorf("expected IsEmpty to be false, got true")
	}

	// Test when the slice has multiple elements
	s.Push(2)
	s.Push(3)
	if s.IsEmpty() {
		t.Errorf("expected IsEmpty to be false, got true")
	}

	// Test after popping all elements
	s.Pop()
	s.Pop()
	s.Pop()
	if !s.IsEmpty() {
		t.Errorf("expected IsEmpty to be true, got false")
	}

	// Test with wrap-around condition
	s.Push(1)
	s.Push(2)
	s.Push(3)
	s.Pop()
	s.Push(4)
	s.Pop()
	s.Pop()
	s.Pop()
	if !s.IsEmpty() {
		t.Errorf("expected IsEmpty to be true, got false")
	}
}

func TestIsFull(t *testing.T) {
	// Test when the slice is empty
	s := NewSlice[int](3)
	if s.IsFull() {
		t.Errorf("expected IsFull to be false, got true")
	}

	// Test when the slice is partially filled
	s.Push(1)
	if s.IsFull() {
		t.Errorf("expected IsFull to be false, got true")
	}

	// Test when the slice is full
	s.Push(2)
	s.Push(3)
	if !s.IsFull() {
		t.Errorf("expected IsFull to be true, got false")
	}

	// Test after popping an element from a full slice
	s.Pop()
	if s.IsFull() {
		t.Errorf("expected IsFull to be false, got true")
	}

	// Test after pushing an element to a nearly full slice
	s.Push(4)
	if !s.IsFull() {
		t.Errorf("expected IsFull to be true, got false")
	}

	// Test with wrap-around condition
	s.Pop()
	s.Push(5)
	if !s.IsFull() {
		t.Errorf("expected IsFull to be true, got false")
	}
}

func TestSlicePop(t *testing.T) {
	// Test Pop on an empty slice
	s := NewSlice[int](5)
	val, ok := s.Pop()
	if ok {
		t.Errorf("expected false, got true")
	}
	if val != 0 {
		t.Errorf("expected 0, got %d", val)
	}

	// Test Pop on a slice with one element
	s.Push(1)
	val, ok = s.Pop()
	if !ok {
		t.Errorf("expected true, got false")
	}
	if val != 1 {
		t.Errorf("expected 1, got %d", val)
	}
	if !s.IsEmpty() {
		t.Errorf("expected slice to be empty")
	}

	// Test Pop on a slice with multiple elements
	s.Push(1)
	s.Push(2)
	s.Push(3)
	val, ok = s.Pop()
	if !ok {
		t.Errorf("expected true, got false")
	}
	if val != 1 {
		t.Errorf("expected 1, got %d", val)
	}
	if s.Len() != 2 {
		t.Errorf("expected length 2, got %d", s.Len())
	}

	val, ok = s.Pop()
	if !ok {
		t.Errorf("expected true, got false")
	}
	if val != 2 {
		t.Errorf("expected 2, got %d", val)
	}
	if s.Len() != 1 {
		t.Errorf("expected length 1, got %d", s.Len())
	}

	val, ok = s.Pop()
	if !ok {
		t.Errorf("expected true, got false")
	}
	if val != 3 {
		t.Errorf("expected 3, got %d", val)
	}
	if !s.IsEmpty() {
		t.Errorf("expected slice to be empty")
	}

	// Test Pop on a wrapped-around slice
	s = NewSlice[int](3)
	s.Push(1)
	s.Push(2)
	s.Push(3)
	s.Pop()
	s.Push(4)
	val, ok = s.Pop()
	if !ok {
		t.Errorf("expected true, got false")
	}
	if val != 2 {
		t.Errorf("expected 2, got %d", val)
	}
	val, ok = s.Pop()
	if !ok {
		t.Errorf("expected true, got false")
	}
	if val != 3 {
		t.Errorf("expected 3, got %d", val)
	}
	val, ok = s.Pop()
	if !ok {
		t.Errorf("expected true, got false")
	}
	if val != 4 {
		t.Errorf("expected 4, got %d", val)
	}
	if !s.IsEmpty() {
		t.Errorf("expected slice to be empty")
	}
}

func TestSliceLen(t *testing.T) {
	// Test Len on an empty slice
	s := NewSlice[int](5)
	if s.Len() != 0 {
		t.Errorf("expected length 0, got %d", s.Len())
	}

	// Test Len on a slice with one element
	s.Push(1)
	if s.Len() != 1 {
		t.Errorf("expected length 1, got %d", s.Len())
	}

	// Test Len on a slice with multiple elements
	s.Push(2)
	s.Push(3)
	if s.Len() != 3 {
		t.Errorf("expected length 3, got %d", s.Len())
	}

	// Test Len after popping elements
	s.Pop()
	if s.Len() != 2 {
		t.Errorf("expected length 2, got %d", s.Len())
	}
	s.Pop()
	if s.Len() != 1 {
		t.Errorf("expected length 1, got %d", s.Len())
	}

	// Test Len with wrap-around condition
	s = NewSlice[int](3)
	s.Push(1)
	s.Push(2)
	s.Push(3)
	s.Pop()
	s.Push(4)
	if s.Len() != 3 {
		t.Errorf("expected length 3, got %d", s.Len())
	}
	s.Pop()
	s.Pop()
	if s.Len() != 1 {
		t.Errorf("expected length 1, got %d", s.Len())
	}
	s.Pop()
	if s.Len() != 0 {
		t.Errorf("expected length 0, got %d", s.Len())
	}
}

func TestSliceRecap(t *testing.T) {
	// Test Recap to a larger size
	s := NewSlice[int](3)
	s.Push(1)
	s.Push(2)
	s.Push(3)
	s.Recap(5)
	if s.Len() != 3 {
		t.Errorf("expected length 3, got %d", s.Len())
	}
	if s.cap != 5 {
		t.Errorf("expected size 5, got %d", s.cap)
	}
	if s.head != 0 || s.tail != 2 {
		t.Errorf("expected head 0 and tail 2, got head %d and tail %d", s.head, s.tail)
	}

	// Test Recap to a smaller size (no change expected)
	s.Recap(2)
	if s.Len() != 3 {
		t.Errorf("expected length 3, got %d", s.Len())
	}
	if s.cap != 5 {
		t.Errorf("expected size 5, got %d", s.cap)
	}

	// Test Recap on an empty slice
	s = NewSlice[int](3)
	s.Recap(5)
	if s.Len() != 0 {
		t.Errorf("expected length 0, got %d", s.Len())
	}
	if s.cap != 5 {
		t.Errorf("expected size 5, got %d", s.cap)
	}
	if s.head != -1 || s.tail != -1 {
		t.Errorf("expected head -1 and tail -1, got head %d and tail %d", s.head, s.tail)
	}

	// Test Recap with wrap-around condition
	s = NewSlice[int](3)
	s.Push(1)
	s.Push(2)
	s.Push(3)
	s.Pop()
	s.Push(4)
	s.Recap(5)
	if s.Len() != 3 {
		t.Errorf("expected length 3, got %d", s.Len())
	}
	if s.cap != 5 {
		t.Errorf("expected size 5, got %d", s.cap)
	}
	if s.head != 0 || s.tail != 2 {
		t.Errorf("expected head 0 and tail 2, got head %d and tail %d", s.head, s.tail)
	}
	expectedValues := []int{2, 3, 4, 0, 0}
	for i, v := range expectedValues {
		if s.values[i] != v {
			t.Errorf("expected value %d at index %d, got %d", v, i, s.values[i])
		}
	}
}
