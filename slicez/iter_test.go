//go:build go1.23

package slicez

import (
	"fmt"
	"slices"
	"testing"

	"github.com/welllog/golib/testz"
)

func TestFlexSlice_Iterators(t *testing.T) {
	tests := []struct {
		name        string
		initialData []int
		initialCap  int
		head        int
	}{
		{
			name:        "standard",
			initialData: []int{10, 20, 30, 40, 50},
			initialCap:  10,
			head:        0,
		},
		{
			name:        "with head wrap",
			initialData: []int{10, 20, 30, 40, 50},
			initialCap:  7,
			head:        5, // data at 5, 6, 0, 1, 2
		},
		{
			name:        "full slice",
			initialData: []int{10, 20, 30},
			initialCap:  3,
			head:        1, // data at 1, 2, 0
		},
		{
			name:        "single element",
			initialData: []int{100},
			initialCap:  5,
			head:        3,
		},
		{
			name:        "empty slice",
			initialData: []int{},
			initialCap:  5,
			head:        0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := makeFlexSlice(tt.initialData, tt.head, tt.initialCap)
			msg := fmt.Sprintf("data: %v, cap: %d, head: %d", tt.initialData, tt.initialCap, tt.head)

			// Test All()
			t.Run("All", func(t *testing.T) {
				result := []int{}
				for v := range f.All() {
					result = append(result, v)
				}
				testz.Equal(t, tt.initialData, result, msg)
			})

			// Test AllWithIndex()
			t.Run("AllWithIndex", func(t *testing.T) {
				indices := []int{}
				values := []int{}
				for i, v := range f.AllWithIndex() {
					indices = append(indices, i)
					values = append(values, v)
				}

				expectedIndices := make([]int, len(tt.initialData))
				for i := range expectedIndices {
					expectedIndices[i] = i
				}

				testz.Equal(t, tt.initialData, values, "values mismatch: "+msg)
				testz.Equal(t, expectedIndices, indices, "indices mismatch: "+msg)
			})

			// Test RevAll()
			t.Run("RevAll", func(t *testing.T) {
				result := []int{}
				for v := range f.RevAll() {
					result = append(result, v)
				}
				expected := slices.Clone(tt.initialData)
				slices.Reverse(expected)
				testz.Equal(t, expected, result, msg)
			})

			// Test RevAllWithIndex()
			t.Run("RevAllWithIndex", func(t *testing.T) {
				indices := []int{}
				values := []int{}
				for i, v := range f.RevAllWithIndex() {
					indices = append(indices, i)
					values = append(values, v)
				}

				expectedValues := slices.Clone(tt.initialData)
				slices.Reverse(expectedValues)

				expectedIndices := make([]int, len(tt.initialData))
				for i := range expectedIndices {
					expectedIndices[i] = len(tt.initialData) - 1 - i
				}

				testz.Equal(t, expectedValues, values, "values mismatch: "+msg)
				testz.Equal(t, expectedIndices, indices, "indices mismatch: "+msg)
			})
		})
	}
}
