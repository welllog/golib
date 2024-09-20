package slicez

import (
	"testing"
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
