package algz

import (
	"testing"
)

func TestKnapsack(t *testing.T) {
	type item struct {
		w     int
		value int
	}

	tests := []struct {
		maxW     int
		items    []item
		expected []item
	}{
		{
			maxW: 10,
			items: []item{
				{w: 5, value: 10},
				{w: 4, value: 9},
				{w: 3, value: 8},
			},
			expected: []item{
				{w: 5, value: 10},
				{w: 4, value: 9},
			},
		},
		{
			maxW: 3,
			items: []item{
				{w: 5, value: 10},
				{w: 4, value: 9},
			},
			expected: []item{},
		},
		{
			maxW: 12,
			items: []item{
				{w: 5, value: 10},
				{w: 4, value: 9},
				{w: 3, value: 8},
			},
			expected: []item{
				{w: 5, value: 10},
				{w: 4, value: 9},
				{w: 3, value: 8},
			},
		},
		{
			maxW: 12,
			items: []item{
				{w: 5, value: 10},
				{w: 4, value: 9},
				{w: 3, value: 8},
				{w: 2, value: 7},
				{w: 1, value: 6},
			},
			expected: []item{
				{w: 5, value: 10},
				{w: 4, value: 9},
				{w: 2, value: 7},
				{w: 1, value: 6},
			},
		},
		{
			maxW: 12,
			items: []item{
				{w: 5, value: 10},
				{w: 4, value: 9},
				{w: 3, value: 8},
				{w: 2, value: 7},
				{w: 1, value: 6},
				{w: 1, value: 5},
			},
			expected: []item{
				{w: 5, value: 10},
				{w: 3, value: 8},
				{w: 2, value: 7},
				{w: 1, value: 6},
				{w: 1, value: 5},
			},
		},
		{
			maxW: 4,
			items: []item{
				{w: 1, value: 7},
				{w: 1, value: 6},
				{w: 2, value: 9},
				{w: 4, value: 9},
				{w: 1, value: 8},
			},
			expected: []item{
				{w: 1, value: 7},
				{w: 2, value: 9},
				{w: 1, value: 8},
			},
		},
	}

	for _, tt := range tests {
		actual := Knapsack(tt.maxW, tt.items, func(i item) int { return i.w }, func(i item) int { return i.value })
		if len(actual) != len(tt.expected) {
			t.Fatalf("expected %v, got %v", tt.expected, actual)
		}
		for i := range actual {
			if actual[i] != tt.expected[i] {
				t.Fatalf("expected %v, got %v", tt.expected, actual)
			}
		}
	}
}

func TestKnapsack2(t *testing.T) {
	type item struct {
		w     int
		value int
	}

	tests := []struct {
		maxW      int
		items     []item
		expected  []item
		expected2 []item
	}{
		{6,
			[]item{{3, 4}, {3, 2}, {6, 6}},
			[]item{{3, 4}, {3, 2}},
			[]item{{6, 6}},
		},
		{
			10,
			[]item{{5, 10}, {4, 9}, {3, 8}, {6, 10}},
			[]item{{5, 10}, {4, 9}},
			[]item{{5, 10}, {4, 9}},
		},
		{
			10,
			[]item{{5, 10}, {4, 9}, {3, 8}, {6, 10}, {9, 19}},
			[]item{{5, 10}, {4, 9}},
			[]item{{9, 19}},
		},
	}

	for _, tt := range tests {
		actual := Knapsack(tt.maxW, tt.items, func(i item) int { return i.w }, func(i item) int { return i.value })
		if len(actual) != len(tt.expected) {
			t.Fatalf("expected %v, got %v", tt.expected, actual)
		}
		for i := range actual {
			if actual[i] != tt.expected[i] {
				t.Fatalf("expected %v, got %v", tt.expected, actual)
			}
		}

		actual2 := Knapsack(tt.maxW, tt.items, func(i item) int { return i.w }, func(i item) int { return i.value }, func(old, new []item) bool {
			return len(old) > len(new)
		})
		if len(actual2) != len(tt.expected2) {
			t.Fatalf("expected %v, got %v", tt.expected2, actual2)
		}
		for i := range actual2 {
			if actual2[i] != tt.expected2[i] {
				t.Fatalf("expected %v, got %v", tt.expected2, actual2)
			}
		}
	}
}

func TestFindDpSolvers(t *testing.T) {
	type item struct {
		name  string
		value int
		sort  int
	}

	tests := []struct {
		maxW      int
		items     []item
		expected  []item
		expected2 []item
	}{
		{
			24,
			[]item{
				{"12a", 12, 89},
				{"12b", 12, 100},
				{"12c", 12, 96},
				{"5a", 5, 100},
				{"5b", 5, 89},
			},
			[]item{{"12a", 12, 89}, {"12c", 12, 96}},
			[]item{{"12a", 12, 89}, {"12c", 12, 96}},
		},
		{
			20,
			[]item{
				{"12a", 12, 89},
				{"12b", 12, 100},
				{"12c", 12, 96},
				{"5a", 5, 100},
				{"5b", 5, 89},
				{"5c", 5, 89},
			},
			[]item{{"12a", 12, 89}, {"5b", 5, 89}},
			[]item{{"12a", 12, 89}, {"5b", 5, 89}, {"5c", 5, 89}},
		},
		{
			18,
			[]item{
				{"12a", 12, 89},
				{"12b", 12, 100},
				{"12c", 12, 96},
				{"5a", 5, 100},
				{"5b", 5, 89},
				{"5c", 5, 89},
			},
			[]item{{"12a", 12, 89}, {"5b", 5, 89}},
			[]item{{"12a", 12, 89}, {"5b", 5, 89}, {"5c", 5, 89}},
		},
		{
			17,
			[]item{
				{"12a", 12, 89},
				{"12b", 12, 100},
				{"12c", 12, 99},
				{"5a", 5, 100},
				{"5b", 5, 89},
				{"5c", 5, 89},
			},
			[]item{{"12a", 12, 89}, {"5b", 5, 89}},
			[]item{{"12a", 12, 89}, {"5b", 5, 89}},
		},
		{
			10,
			[]item{
				{"12a", 12, 89},
				{"12b", 12, 100},
				{"12c", 12, 99},
				{"5a", 5, 100},
				{"5b", 5, 89},
				{"5c", 5, 89},
			},
			[]item{{"5b", 5, 89}, {"5c", 5, 89}},
			[]item{{"5b", 5, 89}, {"5c", 5, 89}},
		},
		{
			10,
			[]item{
				{"5a", 5, 100},
				{"5b", 5, 94},
				{"2a", 2, 94},
				{"2b", 2, 94},
				{"2c", 2, 94},
				{"2d", 2, 94},
				{"2e", 2, 94},
				{"2f", 2, 22},
			},
			[]item{{"2a", 2, 94}, {"2b", 2, 94}, {"2c", 2, 94}, {"2d", 2, 94}, {"2f", 2, 22}},
			[]item{{"2a", 2, 94}, {"2b", 2, 94}, {"2c", 2, 94}, {"2d", 2, 94}, {"2f", 2, 22}},
		},
		{
			14,
			[]item{
				{"5b", 5, 94},
				{"2a", 2, 94},
				{"2b", 2, 94},
				{"2c", 2, 94},
				{"2d", 2, 94},
				{"2e", 2, 94},
				{"2f", 2, 22},
			},
			[]item{{"5b", 5, 94}, {"2a", 2, 94}, {"2b", 2, 94}, {"2c", 2, 94}, {"2f", 2, 22}},
			[]item{{"5b", 5, 94}, {"2a", 2, 94}, {"2b", 2, 94}, {"2c", 2, 94}, {"2d", 2, 94}, {"2f", 2, 22}},
		},
	}

	tieBreaker := func(old, new []item) (replace bool) {
		var oldSum, newSum int
		for _, v := range old {
			oldSum += v.sort
		}

		for _, v := range new {
			newSum += v.sort
		}

		if newSum/len(new) < oldSum/len(old) {
			return true
		}
		return false
	}
	for _, tt := range tests {
		solvers := FindDpSolvers(tt.maxW, tt.items, func(i item) int { return i.value }, true, tieBreaker)
		actual := solvers.Best(tt.maxW)
		if len(actual) != len(tt.expected) {
			t.Fatalf("expected %v, got %v", tt.expected, actual)
		}
		for i := range actual {
			if actual[i] != tt.expected[i] {
				t.Fatalf("expected %v, got %v", tt.expected, actual)
			}
		}

		actual2 := solvers.BestAllowMinOverflow(tt.maxW)
		if len(actual2) != len(tt.expected2) {
			t.Fatalf("expected %v, got %v", tt.expected2, actual2)
		}
		for i := range actual2 {
			if actual2[i] != tt.expected2[i] {
				t.Fatalf("expected %v, got %v", tt.expected2, actual2)
			}
		}
	}
}
