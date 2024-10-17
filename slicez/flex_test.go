package slicez

import (
	"testing"
	"unsafe"

	"github.com/welllog/golib/testz"
)

func TestShrink(t *testing.T) {
	tests := []struct {
		name     string
		initial  []int
		expected []int
	}{
		{
			name:     "Capacity less than or equal to 8",
			initial:  make([]int, 4, 8),
			expected: make([]int, 4, 8),
		},
		{
			name:     "Length greater than one-fourth of capacity",
			initial:  make([]int, 11, 40),
			expected: make([]int, 11, 40),
		},
		{
			name:     "Length less than or equal to one-fourth of capacity",
			initial:  make([]int, 5, 40),
			expected: make([]int, 5, 10),
		},
		{
			name:     "Length exactly one-fourth of capacity",
			initial:  make([]int, 10, 40),
			expected: make([]int, 10, 20),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := FlexSlice[int]{Values: tt.initial}
			f.shrink()
			if cap(f.Values) != cap(tt.expected) {
				t.Errorf("expected capacity %d, got %d", cap(tt.expected), cap(f.Values))
			}
			if len(f.Values) != len(tt.expected) {
				t.Errorf("expected length %d, got %d", len(tt.expected), len(f.Values))
			}
		})
	}
}

func TestPop(t *testing.T) {
	tests := []struct {
		name     string
		initial  []int
		expected int
		ok       bool
	}{
		{
			name:     "Pop from non-empty slice",
			initial:  []int{1, 2, 3},
			expected: 3,
			ok:       true,
		},
		{
			name:     "Pop from empty slice",
			initial:  []int{},
			expected: 0,
			ok:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := FlexSlice[int]{Values: tt.initial}
			value, ok := f.Pop()
			if value != tt.expected || ok != tt.ok {
				t.Errorf("expected value %d and ok %v, got value %d and ok %v", tt.expected, tt.ok, value, ok)
			}
		})
	}
}

func TestPop1(t *testing.T) {
	s := []int{6, 5, 4, 3, 2, 1}
	f := FlexSlice[int]{Values: s}
	var v int
	for f.Len() > 0 {
		v++
		n, _ := f.Pop()
		testz.Equal(t, v, n)
	}
}

func TestShift(t *testing.T) {
	tests := []struct {
		name     string
		initial  []int
		expected int
		ok       bool
	}{
		{
			name:     "Shift from non-empty slice",
			initial:  []int{1, 2, 3},
			expected: 1,
			ok:       true,
		},
		{
			name:     "Shift from empty slice",
			initial:  []int{},
			expected: 0,
			ok:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := FlexSlice[int]{Values: tt.initial}
			value, ok := f.Shift()
			if value != tt.expected || ok != tt.ok {
				t.Errorf("expected value %d and ok %v, got value %d and ok %v", tt.expected, tt.ok, value, ok)
			}
		})
	}
}

func TestShift1(t *testing.T) {
	s := []int{1, 2, 3, 4, 5, 6}
	f := FlexSlice[int]{Values: s}
	var v int
	for f.Len() > 0 {
		v++
		n, _ := f.Shift()
		testz.Equal(t, v, n)
	}
}

func TestPrepend(t *testing.T) {
	s := []int{1, 2, 3, 4, 5, 6, 7}
	p := (*(*[1]uintptr)(unsafe.Pointer(&s)))[0]
	f := FlexSlice[int]{Values: s}

	n, _ := f.Pop()
	testz.Equal(t, 7, n)
	testz.Equal(t, 6, f.Len())

	n, _ = f.Pop()
	testz.Equal(t, 6, n)
	testz.Equal(t, 5, f.Len())

	n, _ = f.Shift()
	testz.Equal(t, 1, n)
	testz.Equal(t, 4, f.Len())

	p1 := (*(*[1]uintptr)(unsafe.Pointer(&f.Values)))[0]
	testz.Equal(t, p, p1)

	f.Prepend(0)
	f.Prepend(-1)
	testz.Equal(t, 6, f.Len())

	p2 := (*(*[1]uintptr)(unsafe.Pointer(&f.Values)))[0]
	testz.Equal(t, p, p2)

	f.Prepend(-3, -2)
	testz.Equal(t, 8, f.Len())

	p3 := (*(*[1]uintptr)(unsafe.Pointer(&f.Values)))[0]
	if p3 == p {
		t.Error("expected new slice")
	}

	testz.Equal(t, 14, cap(f.Values))

	n, _ = f.Shift()
	testz.Equal(t, -3, n)

	n, _ = f.Shift()
	testz.Equal(t, -2, n)

	testz.Equal(t, 6, f.Len())
}
