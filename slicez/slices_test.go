package slicez

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/welllog/golib/testz"
)

func TestDiff(t *testing.T) {
	var dst []int
	testCases := []struct {
		s1   []int
		s2   []int
		want []int
	}{
		{
			s1:   []int{1, 2, 3},
			s2:   []int{1, 2, 4},
			want: []int{3},
		},
		{
			s1:   []int{1, 2, 3},
			s2:   []int{4, 5, 6},
			want: []int{1, 2, 3},
		},
		{
			s1:   []int{1, 2, 3},
			s2:   []int{1, 2, 3},
			want: []int{},
		},
		{
			s1:   []int{1, 2, 3},
			s2:   []int{},
			want: []int{1, 2, 3},
		},
		{
			s1:   []int{},
			s2:   []int{1, 2, 3},
			want: []int{},
		},
		{
			s1:   []int{},
			s2:   []int{},
			want: []int{},
		},
		{
			s1:   []int{1, 1, 2, 3},
			s2:   []int{1, 4},
			want: []int{2, 3},
		},
		{
			s1:   []int{1, 2, 3, 4, 5, 6},
			s2:   []int{1, 2, 3, 4, 7},
			want: []int{5, 6},
		},
		{
			s1:   []int{1, 2, 2, 3, 4, 4, 5, 5},
			s2:   []int{1, 3, 5},
			want: []int{2, 2, 4, 4},
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("s1=%v,s2=%v", tc.s1, tc.s2), func(t *testing.T) {
			got := Diff(dst, tc.s1, tc.s2)
			testz.Equal(t, tc.want, got)
		})
	}
}

func TestDiffInPlaceFirst(t *testing.T) {
	testCases := []struct {
		s1   []int
		s2   []int
		want []int
	}{
		{
			s1:   []int{1, 2, 3},
			s2:   []int{1, 2, 4},
			want: []int{3},
		},
		{
			s1:   []int{1, 2, 3},
			s2:   []int{4, 5, 6},
			want: []int{1, 2, 3},
		},
		{
			s1:   []int{1, 2, 3},
			s2:   []int{1, 2, 3},
			want: []int{},
		},
		{
			s1:   []int{1, 2, 3},
			s2:   []int{},
			want: []int{1, 2, 3},
		},
		{
			s1:   []int{},
			s2:   []int{1, 2, 3},
			want: []int{},
		},
		{
			s1:   []int{},
			s2:   []int{},
			want: []int{},
		},
		{
			s1:   []int{1, 1, 2, 3},
			s2:   []int{1, 4},
			want: []int{2, 3},
		},
		{
			s1:   []int{1, 2, 3, 4, 5, 6},
			s2:   []int{1, 2, 3, 4, 7},
			want: []int{5, 6},
		},
		{
			s1:   []int{1, 2, 2, 3, 4, 4, 5, 5},
			s2:   []int{1, 3, 5},
			want: []int{2, 2, 4, 4},
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("s1=%v,s2=%v", tc.s1, tc.s2), func(t *testing.T) {
			s1 := Copy(tc.s1, 0, -1)
			got := DiffInPlaceFirst(s1, tc.s2)
			sort.Ints(got)
			testz.Equal(t, tc.want, got)
		})
	}
}

func TestDiffSorted(t *testing.T) {
	var dst []int
	testCases := []struct {
		s1   []int
		s2   []int
		want []int
	}{
		{
			s1:   []int{1, 2, 3},
			s2:   []int{1, 2, 4},
			want: []int{3},
		},
		{
			s1:   []int{1, 2, 3},
			s2:   []int{4, 5, 6},
			want: []int{1, 2, 3},
		},
		{
			s1:   []int{1, 2, 3},
			s2:   []int{1, 2, 3},
			want: []int{},
		},
		{
			s1:   []int{1, 2, 3},
			s2:   []int{},
			want: []int{1, 2, 3},
		},
		{
			s1:   []int{},
			s2:   []int{1, 2, 3},
			want: []int{},
		},
		{
			s1:   []int{},
			s2:   []int{},
			want: []int{},
		},
		{
			s1:   []int{1, 1, 2, 3},
			s2:   []int{1, 4},
			want: []int{2, 3},
		},
		{
			s1:   []int{1, 2, 3, 5},
			s2:   []int{2, 4, 5, 6},
			want: []int{1, 3},
		},
		{
			s1:   []int{1, 2, 3, 4, 5, 6},
			s2:   []int{1, 2, 3, 4, 7},
			want: []int{5, 6},
		},
		{
			s1:   []int{1, 2, 2, 3, 4, 4, 5, 5},
			s2:   []int{1, 3, 5},
			want: []int{2, 2, 4, 4},
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("s1=%v,s2=%v", tc.s1, tc.s2), func(t *testing.T) {
			got := DiffSorted(dst, tc.s1, tc.s2)
			testz.Equal(t, tc.want, got)
		})
	}
}

func TestDiffSortedInPlaceFirst(t *testing.T) {
	testCases := []struct {
		s1   []int
		s2   []int
		want []int
	}{
		{
			s1:   []int{1, 2, 3},
			s2:   []int{1, 2, 4},
			want: []int{3},
		},
		{
			s1:   []int{1, 2, 3},
			s2:   []int{4, 5, 6},
			want: []int{1, 2, 3},
		},
		{
			s1:   []int{1, 2, 3},
			s2:   []int{1, 2, 3},
			want: []int{},
		},
		{
			s1:   []int{1, 2, 3},
			s2:   []int{},
			want: []int{1, 2, 3},
		},
		{
			s1:   []int{},
			s2:   []int{1, 2, 3},
			want: []int{},
		},
		{
			s1:   []int{},
			s2:   []int{},
			want: []int{},
		},
		{
			s1:   []int{1, 1, 2, 3},
			s2:   []int{1, 4},
			want: []int{2, 3},
		},
		{
			s1:   []int{1, 2, 3, 5},
			s2:   []int{2, 4, 5, 6},
			want: []int{1, 3},
		},
		{
			s1:   []int{1, 2, 3, 4, 5, 6},
			s2:   []int{1, 2, 3, 4, 7},
			want: []int{5, 6},
		},
		{
			s1:   []int{1, 2, 2, 3, 4, 4, 5, 5},
			s2:   []int{1, 3, 5},
			want: []int{2, 2, 4, 4},
		},
		{
			s1:   []int{1, 3, 5, 7, 9, 10},
			s2:   []int{3, 7},
			want: []int{1, 5, 9, 10},
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("s1=%v,s2=%v", tc.s1, tc.s2), func(t *testing.T) {
			s1 := Copy(tc.s1, 0, -1)
			got := DiffSortedInPlaceFirst(s1, tc.s2)
			testz.Equal(t, tc.want, got)
		})
	}
}

func TestIntersect(t *testing.T) {
	var dst []int
	testCases := []struct {
		s1   []int
		s2   []int
		want []int
	}{
		{
			s1:   []int{1, 2, 3},
			s2:   []int{1, 2, 4},
			want: []int{1, 2},
		},
		{
			s1:   []int{1, 2, 3},
			s2:   []int{4, 5, 6},
			want: []int{},
		},
		{
			s1:   []int{1, 2, 3},
			s2:   []int{1, 2, 3},
			want: []int{1, 2, 3},
		},
		{
			s1:   []int{1, 2, 3},
			s2:   []int{},
			want: []int{},
		},
		{
			s1:   []int{},
			s2:   []int{1, 2, 3},
			want: []int{},
		},
		{
			s1:   []int{},
			s2:   []int{},
			want: []int{},
		},
		{
			s1:   []int{1, 1, 2, 3},
			s2:   []int{1, 4},
			want: []int{1, 1},
		},
		{
			s1:   []int{1, 2, 3, 4, 5, 6},
			s2:   []int{1, 2, 3, 4, 7},
			want: []int{1, 2, 3, 4},
		},
		{
			s1:   []int{1, 1, 2, 3},
			s2:   []int{1, 2, 4},
			want: []int{1, 1, 2},
		},
		{
			s1:   []int{1, 1, 2, 3},
			s2:   []int{1, 1, 2, 4},
			want: []int{1, 1, 2},
		},
		{
			s1:   []int{1, 1, 2, 3},
			s2:   []int{1, 1, 1, 2, 4},
			want: []int{1, 1, 2},
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("s1=%v,s2=%v", tc.s1, tc.s2), func(t *testing.T) {
			got := Intersect(dst, tc.s1, tc.s2)
			testz.Equal(t, tc.want, got)
		})
	}
}

func TestIntersectSortedMultiset(t *testing.T) {
	var dst []int
	testCases := []struct {
		s1   []int
		s2   []int
		want []int
	}{
		{
			s1:   []int{1, 2, 3},
			s2:   []int{1, 2, 4},
			want: []int{1, 2},
		},
		{
			s1:   []int{1, 2, 3},
			s2:   []int{4, 5, 6},
			want: []int{},
		},
		{
			s1:   []int{1, 2, 3},
			s2:   []int{1, 2, 3},
			want: []int{1, 2, 3},
		},
		{
			s1:   []int{1, 2, 3},
			s2:   []int{},
			want: []int{},
		},
		{
			s1:   []int{},
			s2:   []int{1, 2, 3},
			want: []int{},
		},
		{
			s1:   []int{},
			s2:   []int{},
			want: []int{},
		},
		{
			s1:   []int{1, 1, 2, 3},
			s2:   []int{1, 4},
			want: []int{1},
		},
		{
			s1:   []int{1, 1, 2, 3},
			s2:   []int{1, 1, 4},
			want: []int{1, 1},
		},
		{
			s1:   []int{1, 1, 2, 3},
			s2:   []int{1, 1, 1, 4},
			want: []int{1, 1},
		},
		{
			s1:   []int{1, 2, 3, 5},
			s2:   []int{2, 4, 5, 6},
			want: []int{2, 5},
		},
		{
			s1:   []int{1, 2, 3, 4, 5, 6},
			s2:   []int{1, 2, 3, 4, 7},
			want: []int{1, 2, 3, 4},
		},
		{
			s1:   []int{1, 2, 2, 3},
			s2:   []int{2, 2, 2, 4},
			want: []int{2, 2},
		},
		{
			s1:   []int{1, 1, 2, 2},
			s2:   []int{1, 2, 2, 2},
			want: []int{1, 2, 2},
		},
		{
			s1:   []int{1, 1, 2, 3},
			s2:   []int{1, 2, 4},
			want: []int{1, 2},
		},
		{
			s1:   []int{1, 1, 2, 3},
			s2:   []int{1, 1, 2, 4},
			want: []int{1, 1, 2},
		},
		{
			s1:   []int{1, 1, 2, 3},
			s2:   []int{1, 1, 1, 2, 4},
			want: []int{1, 1, 2},
		},
		{
			s1:   []int{1, 2, 2, 3},
			s2:   []int{2, 2, 2, 4},
			want: []int{2, 2},
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("s1=%v,s2=%v", tc.s1, tc.s2), func(t *testing.T) {
			got := IntersectSortedMultiset(dst, tc.s1, tc.s2)
			testz.Equal(t, tc.want, got)
		})
	}
}

func TestIntersectSortedSet(t *testing.T) {
	var dst []int
	testCases := []struct {
		s1   []int
		s2   []int
		want []int
	}{
		{
			s1:   []int{1, 2, 3},
			s2:   []int{1, 2, 4},
			want: []int{1, 2},
		},
		{
			s1:   []int{1, 2, 3},
			s2:   []int{4, 5, 6},
			want: []int{},
		},
		{
			s1:   []int{1, 2, 3},
			s2:   []int{1, 2, 3},
			want: []int{1, 2, 3},
		},
		{
			s1:   []int{1, 2, 3},
			s2:   []int{},
			want: []int{},
		},
		{
			s1:   []int{},
			s2:   []int{1, 2, 3},
			want: []int{},
		},
		{
			s1:   []int{},
			s2:   []int{},
			want: []int{},
		},
		{
			s1:   []int{1, 1, 2, 3},
			s2:   []int{1, 4},
			want: []int{1},
		},
		{
			s1:   []int{1, 1, 2, 3},
			s2:   []int{1, 1, 4},
			want: []int{1},
		},
		{
			s1:   []int{1, 1, 2, 3},
			s2:   []int{1, 1, 1, 4},
			want: []int{1},
		},
		{
			s1:   []int{1, 1, 2, 3},
			s2:   []int{1, 1, 1, 2, 4},
			want: []int{1, 2},
		},
		{
			s1:   []int{1, 2, 3, 4, 5, 6},
			s2:   []int{1, 2, 3, 4, 7},
			want: []int{1, 2, 3, 4},
		},
		{
			s1:   []int{1, 1, 2, 3},
			s2:   []int{1, 2, 4},
			want: []int{1, 2},
		},
		{
			s1:   []int{1, 1, 2, 3},
			s2:   []int{1, 1, 2, 4},
			want: []int{1, 2},
		},
		{
			s1:   []int{1, 1, 2, 3},
			s2:   []int{1, 1, 1, 2, 4},
			want: []int{1, 2},
		},
		{
			s1:   []int{1, 2, 2, 3},
			s2:   []int{2, 2, 2, 4},
			want: []int{2},
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("s1=%v,s2=%v", tc.s1, tc.s2), func(t *testing.T) {
			got := IntersectSortedSet(dst, tc.s1, tc.s2)
			testz.Equal(t, tc.want, got)
		})
	}
}

func TestIntersectInPlaceFirst(t *testing.T) {
	testCases := []struct {
		s1   []int
		s2   []int
		want []int
	}{
		{
			s1:   []int{1, 2, 3},
			s2:   []int{1, 2, 4},
			want: []int{1, 2},
		},
		{
			s1:   []int{1, 2, 3},
			s2:   []int{4, 5, 6},
			want: []int{},
		},
		{
			s1:   []int{1, 2, 3},
			s2:   []int{1, 2, 3},
			want: []int{1, 2, 3},
		},
		{
			s1:   []int{1, 2, 3},
			s2:   []int{},
			want: []int{},
		},
		{
			s1:   []int{},
			s2:   []int{1, 2, 3},
			want: []int{},
		},
		{
			s1:   []int{},
			s2:   []int{},
			want: []int{},
		},
		{
			s1:   []int{1, 1, 2, 3},
			s2:   []int{1, 4},
			want: []int{1, 1},
		},
		{
			s1:   []int{1, 2, 3, 4, 5, 6},
			s2:   []int{1, 2, 3, 4, 7},
			want: []int{1, 2, 3, 4},
		},
		{
			s1:   []int{1, 1, 2, 3},
			s2:   []int{1, 2, 4},
			want: []int{1, 1, 2},
		},
		{
			s1:   []int{1, 1, 2, 3},
			s2:   []int{1, 1, 2, 4},
			want: []int{1, 1, 2},
		},
		{
			s1:   []int{1, 1, 2, 3},
			s2:   []int{1, 1, 1, 2, 4},
			want: []int{1, 1, 2},
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("s1=%v,s2=%v", tc.s1, tc.s2), func(t *testing.T) {
			s1 := Copy(tc.s1, 0, -1)
			got := IntersectInPlaceFirst(s1, tc.s2)
			sort.Ints(got)
			testz.Equal(t, tc.want, got)
		})
	}
}

func TestUnique(t *testing.T) {
	var dst []int
	testCases := []struct {
		s    []int
		want []int
	}{
		{
			s:    []int{1, 2, 3, 2, 1},
			want: []int{1, 2, 3},
		},
		{
			s:    []int{1, 2, 3, 4, 5, 1, 2},
			want: []int{1, 2, 3, 4, 5},
		},
		{
			s:    []int{1, 2, 3},
			want: []int{1, 2, 3},
		},
		{
			s:    []int{1, 1, 1},
			want: []int{1},
		},
		{
			s:    []int{},
			want: []int{},
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("s=%v", tc.s), func(t *testing.T) {
			got := Unique(dst, tc.s)
			testz.Equal(t, tc.want, got)
		})
	}
}

func TestUniqueSorted(t *testing.T) {
	var dst []int
	testCases := []struct {
		s    []int
		want []int
	}{
		{
			s:    []int{1, 1, 2, 2, 3},
			want: []int{1, 2, 3},
		},
		{
			s:    []int{1, 2, 3},
			want: []int{1, 2, 3},
		},
		{
			s:    []int{1, 1, 1},
			want: []int{1},
		},
		{
			s:    []int{},
			want: []int{},
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("s=%v", tc.s), func(t *testing.T) {
			got := UniqueSorted(dst, tc.s)
			testz.Equal(t, tc.want, got)
		})
	}
}

func TestUniqueInPlace(t *testing.T) {
	testCases := []struct {
		s    []int
		want []int
	}{
		{
			s:    []int{1, 2, 3, 2, 1},
			want: []int{1, 2, 3},
		},
		{
			s:    []int{1, 2, 3, 4, 5, 1, 2},
			want: []int{1, 2, 3, 4, 5},
		},
		{
			s:    []int{1, 2, 3},
			want: []int{1, 2, 3},
		},
		{
			s:    []int{1, 1, 1},
			want: []int{1},
		},
		{
			s:    []int{},
			want: []int{},
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("s=%v", tc.s), func(t *testing.T) {
			s := Copy(tc.s, 0, -1)
			got := UniqueInPlace(s)
			sort.Ints(got)
			testz.Equal(t, tc.want, got)
		})
	}
}

func TestUniqueInPlaceSorted(t *testing.T) {
	testCases := []struct {
		s    []int
		want []int
	}{
		{
			s:    []int{1, 1, 2, 2, 3},
			want: []int{1, 2, 3},
		},
		{
			s:    []int{1, 2, 3},
			want: []int{1, 2, 3},
		},
		{
			s:    []int{1, 1, 1},
			want: []int{1},
		},
		{
			s:    []int{},
			want: []int{},
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("s=%v", tc.s), func(t *testing.T) {
			s := Copy(tc.s, 0, -1)
			got := UniqueInPlaceSorted(s)
			testz.Equal(t, tc.want, got)
		})
	}
}

type uniqueObject struct {
	id   int
	name string
}

func TestUniqueByKey(t *testing.T) {
	var dst []uniqueObject
	testCases := []struct {
		s    []uniqueObject
		key  func(uniqueObject) int
		want []uniqueObject
	}{
		{
			s: []uniqueObject{
				{1, "a"},
				{2, "b"},
				{1, "c"},
			},
			key: func(o uniqueObject) int {
				return o.id
			},
			want: []uniqueObject{
				{1, "a"},
				{2, "b"},
			},
		},
		{
			s: []uniqueObject{
				{1, "a"}, {2, "b"}, {3, "c"}, {4, "d"}, {5, "e"}, {1, "f"},
			},
			key: func(o uniqueObject) int {
				return o.id
			},
			want: []uniqueObject{
				{1, "a"}, {2, "b"}, {3, "c"}, {4, "d"}, {5, "e"},
			},
		},
		{
			s: []uniqueObject{
				{1, "a"},
				{2, "b"},
				{1, "c"},
			},
			key: func(o uniqueObject) int {
				return o.id
			},
			want: []uniqueObject{
				{1, "a"},
				{2, "b"},
			},
		},
		{
			s: []uniqueObject{
				{1, "a"}, {2, "b"}, {3, "c"}, {4, "d"}, {5, "e"}, {1, "f"},
			},
			key: func(o uniqueObject) int {
				return o.id
			},
			want: []uniqueObject{
				{1, "a"}, {2, "b"}, {3, "c"}, {4, "d"}, {5, "e"},
			},
		},
		{
			s: []uniqueObject{
				{1, "a"},
				{2, "b"},
				{1, "c"},
			},
			key: func(o uniqueObject) int {
				return o.id
			},
			want: []uniqueObject{
				{1, "a"},
				{2, "b"},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("s=%v", tc.s), func(t *testing.T) {
			got := UniqueByKey(dst, tc.s, tc.key)
			testz.Equal(t, tc.want, got)
		})
	}
}

func TestUniqueByKeyInPlace(t *testing.T) {
	testCases := []struct {
		s    []uniqueObject
		key  func(uniqueObject) int
		want []uniqueObject
	}{
		{
			s: []uniqueObject{
				{1, "a"},
				{2, "b"},
				{1, "c"},
			},
			key: func(o uniqueObject) int {
				return o.id
			},
			want: []uniqueObject{
				{1, "a"},
				{2, "b"},
			},
		},
		{
			s: []uniqueObject{
				{1, "a"}, {2, "b"}, {3, "c"}, {4, "d"}, {5, "e"}, {1, "f"},
			},
			key: func(o uniqueObject) int {
				return o.id
			},
			want: []uniqueObject{
				{1, "a"}, {2, "b"}, {3, "c"}, {4, "d"}, {5, "e"},
			},
		},
		{
			s: []uniqueObject{
				{1, "a"},
				{2, "b"},
				{1, "c"},
			},
			key: func(o uniqueObject) int {
				return o.id
			},
			want: []uniqueObject{
				{1, "a"},
				{2, "b"},
			},
		},
		{
			s: []uniqueObject{
				{1, "a"}, {2, "b"}, {3, "c"}, {4, "d"}, {5, "e"}, {1, "f"},
			},
			key: func(o uniqueObject) int {
				return o.id
			},
			want: []uniqueObject{
				{1, "a"}, {2, "b"}, {3, "c"}, {4, "d"}, {5, "e"},
			},
		},
		{
			s: []uniqueObject{
				{1, "a"},
				{2, "b"},
				{1, "c"},
			},
			key: func(o uniqueObject) int {
				return o.id
			},
			want: []uniqueObject{
				{1, "a"},
				{2, "b"},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("s=%v", tc.s), func(t *testing.T) {
			s := Copy(tc.s, 0, -1)
			got := UniqueByKeyInPlace(s, tc.key)
			sort.Slice(got, func(i, j int) bool {
				return got[i].id < got[j].id
			})
			testz.Equal(t, tc.want, got)
		})
	}
}

func TestFilter(t *testing.T) {
	var dst []int
	testCases := []struct {
		s    []int
		f    func(int) bool
		want []int
	}{
		{
			s:    []int{1, 2, 3, 4, 5},
			f:    func(i int) bool { return i%2 == 0 },
			want: []int{2, 4},
		},
		{
			s:    []int{1, 2, 3, 4, 5},
			f:    func(i int) bool { return i > 5 },
			want: []int{},
		},
		{
			s:    []int{},
			f:    func(i int) bool { return i > 0 },
			want: []int{},
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("s=%v", tc.s), func(t *testing.T) {
			got := Filter(dst, tc.s, tc.f)
			testz.Equal(t, tc.want, got)
		})
	}
}

func TestFilterInPlace(t *testing.T) {
	testCases := []struct {
		s    []int
		f    func(int) bool
		want []int
	}{
		{
			s:    []int{1, 2, 3, 4, 5},
			f:    func(i int) bool { return i%2 == 0 },
			want: []int{2, 4},
		},
		{
			s:    []int{1, 2, 3, 4, 5},
			f:    func(i int) bool { return i > 5 },
			want: []int{},
		},
		{
			s:    []int{},
			f:    func(i int) bool { return i > 0 },
			want: []int{},
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("s=%v", tc.s), func(t *testing.T) {
			s := Copy(tc.s, 0, -1)
			got := FilterInPlace(s, tc.f)
			sort.Ints(got)
			testz.Equal(t, tc.want, got)
		})
	}
}

func TestEqual(t *testing.T) {
	testCases := []struct {
		s1   []int
		s2   []int
		want bool
	}{
		{
			s1:   []int{1, 2, 3},
			s2:   []int{1, 2, 3},
			want: true,
		},
		{
			s1:   []int{1, 2, 3},
			s2:   []int{1, 2, 4},
			want: false,
		},
		{
			s1:   []int{1, 2, 3},
			s2:   []int{1, 2},
			want: false,
		},
		{
			s1:   []int{},
			s2:   []int{},
			want: true,
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("s1=%v,s2=%v", tc.s1, tc.s2), func(t *testing.T) {
			got := Equal(tc.s1, tc.s2)
			testz.Equal(t, tc.want, got)
		})
	}
}

func TestIndex(t *testing.T) {
	testCases := []struct {
		s    []int
		v    int
		want int
	}{
		{
			s:    []int{1, 2, 3},
			v:    2,
			want: 1,
		},
		{
			s:    []int{1, 2, 3},
			v:    4,
			want: -1,
		},
		{
			s:    []int{},
			v:    1,
			want: -1,
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("s=%v,v=%d", tc.s, tc.v), func(t *testing.T) {
			got := Index(tc.s, tc.v)
			testz.Equal(t, tc.want, got)
		})
	}
}

func TestIndexFunc(t *testing.T) {
	testCases := []struct {
		s    []int
		f    func(int) bool
		want int
	}{
		{
			s:    []int{1, 2, 3},
			f:    func(i int) bool { return i == 2 },
			want: 1,
		},
		{
			s:    []int{1, 2, 3},
			f:    func(i int) bool { return i == 4 },
			want: -1,
		},
		{
			s:    []int{},
			f:    func(i int) bool { return i == 1 },
			want: -1,
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("s=%v", tc.s), func(t *testing.T) {
			got := IndexFunc(tc.s, tc.f)
			testz.Equal(t, tc.want, got)
		})
	}
}

func TestSubSlice(t *testing.T) {
	testCases := []struct {
		s     []int
		start int
		end   int
		want  []int
	}{
		{
			s:     []int{1, 2, 3, 4, 5},
			start: 1,
			end:   3,
			want:  []int{2, 3},
		},
		{
			s:     []int{1, 2, 3, 4, 5},
			start: 3,
			end:   -1,
			want:  []int{4, 5},
		},
		{
			s:     []int{1, 2, 3, 4, 5},
			start: 5,
			end:   -1,
			want:  []int{},
		},
		{
			s:     []int{1, 2, 3, 4, 5},
			start: -1,
			end:   3,
			want:  []int{1, 2, 3},
		},
		{
			s:     []int{1, 2, 3, 4, 5},
			start: 3,
			end:   1,
			want:  []int{},
		},
		{
			s:     []int{1, 2, 3, 4, 5},
			start: -2,
			end:   -1,
			want:  []int{1, 2, 3, 4, 5},
		},
		{
			s:     []int{},
			start: 0,
			end:   -1,
			want:  []int{},
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("s=%v,start=%d,end=%d", tc.s, tc.start, tc.end), func(t *testing.T) {
			got := SubSlice(tc.s, tc.start, tc.end)
			testz.Equal(t, tc.want, got)
		})
	}
}

func TestContains(t *testing.T) {
	testCases := []struct {
		s    []int
		v    int
		want bool
	}{
		{
			s:    []int{1, 2, 3},
			v:    2,
			want: true,
		},
		{
			s:    []int{1, 2, 3},
			v:    4,
			want: false,
		},
		{
			s:    []int{},
			v:    1,
			want: false,
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("s=%v,v=%d", tc.s, tc.v), func(t *testing.T) {
			got := Contains(tc.s, tc.v)
			testz.Equal(t, tc.want, got)
		})
	}
}

func TestContainsFunc(t *testing.T) {
	testCases := []struct {
		s    []int
		f    func(int) bool
		want bool
	}{
		{
			s:    []int{1, 2, 3},
			f:    func(i int) bool { return i == 2 },
			want: true,
		},
		{
			s:    []int{1, 2, 3},
			f:    func(i int) bool { return i == 4 },
			want: false,
		},
		{
			s:    []int{},
			f:    func(i int) bool { return i == 1 },
			want: false,
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("s=%v", tc.s), func(t *testing.T) {
			got := ContainsFunc(tc.s, tc.f)
			testz.Equal(t, tc.want, got)
		})
	}
}

func TestChunk(t *testing.T) {
	testCases := []struct {
		s         []int
		chunkSize int
		want      [][]int
	}{
		{
			s:         []int{1, 2, 3, 4, 5},
			chunkSize: 2,
			want:      [][]int{{1, 2}, {3, 4}, {5}},
		},
		{
			s:         []int{1, 2, 3, 4, 5},
			chunkSize: 5,
			want:      [][]int{{1, 2, 3, 4, 5}},
		},
		{
			s:         []int{1, 2, 3, 4, 5},
			chunkSize: 6,
			want:      [][]int{{1, 2, 3, 4, 5}},
		},
		{
			s:         []int{},
			chunkSize: 2,
			want:      nil,
		},
		{
			s:         []int{1, 2, 3, 4, 5},
			chunkSize: 0,
			want:      [][]int{{1, 2, 3, 4, 5}},
		},
		{
			s:         []int{1, 2, 3, 4, 5},
			chunkSize: -1,
			want:      [][]int{{1, 2, 3, 4, 5}},
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("s=%v,chunkSize=%d", tc.s, tc.chunkSize), func(t *testing.T) {
			got := Chunk(tc.s, tc.chunkSize)
			testz.Equal(t, tc.want, got)
		})
	}
}

func TestChunkProcess(t *testing.T) {
	testCases := []struct {
		s         []int
		chunkSize int
		want      string
	}{
		{
			s:         []int{1, 2, 3, 4, 5},
			chunkSize: 2,
			want:      "[1 2][3 4][5]",
		},
		{
			s:         []int{1, 2, 3, 4, 5},
			chunkSize: 5,
			want:      "[1 2 3 4 5]",
		},
		{
			s:         []int{1, 2, 3, 4, 5},
			chunkSize: 6,
			want:      "[1 2 3 4 5]",
		},
		{
			s:         []int{},
			chunkSize: 2,
			want:      "",
		},
		{
			s:         []int{1, 2, 3, 4, 5},
			chunkSize: 0,
			want:      "[1 2 3 4 5]",
		},
		{
			s:         []int{1, 2, 3, 4, 5},
			chunkSize: -1,
			want:      "[1 2 3 4 5]",
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("s=%v,chunkSize=%d", tc.s, tc.chunkSize), func(t *testing.T) {
			var builder strings.Builder
			err := ChunkProcess(tc.s, tc.chunkSize, func(chunk []int) error {
				builder.WriteString(fmt.Sprintf("%v", chunk))
				return nil
			})
			testz.Nil(t, err)
			testz.Equal(t, tc.want, builder.String())
		})
	}
}

func TestCopy(t *testing.T) {
	testCases := []struct {
		s      []int
		start  int
		length int
		want   []int
	}{
		{
			s:      []int{1, 2, 3, 4, 5},
			start:  1,
			length: 3,
			want:   []int{2, 3, 4},
		},
		{
			s:      []int{1, 2, 3, 4, 5},
			start:  3,
			length: -1,
			want:   []int{4, 5},
		},
		{
			s:      []int{1, 2, 3, 4, 5},
			start:  5,
			length: -1,
			want:   nil,
		},
		{
			s:      []int{1, 2, 3, 4, 5},
			start:  -1,
			length: 3,
			want:   []int{1, 2, 3},
		},
		{
			s:      []int{1, 2, 3, 4, 5},
			start:  3,
			length: 0,
			want:   nil,
		},
		{
			s:      []int{1, 2, 3, 4, 5},
			start:  1,
			length: 10,
			want:   []int{2, 3, 4, 5},
		},
		{
			s:      []int{},
			start:  0,
			length: -1,
			want:   nil,
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("s=%v,start=%d,length=%d", tc.s, tc.start, tc.length), func(t *testing.T) {
			got := Copy(tc.s, tc.start, tc.length)
			testz.Equal(t, tc.want, got)
		})
	}
}

func TestValues(t *testing.T) {
	testCases := []struct {
		ss   [][]int
		fn   func(int) string
		want []string
	}{
		{
			ss:   [][]int{{1, 2}, {3, 4, 5}},
			fn:   strconv.Itoa,
			want: []string{"1", "2", "3", "4", "5"},
		},
		{
			ss:   [][]int{{}},
			fn:   strconv.Itoa,
			want: []string{},
		},
		{
			ss:   [][]int{},
			fn:   strconv.Itoa,
			want: []string{},
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("ss=%v", tc.ss), func(t *testing.T) {
			got := Values(tc.fn, tc.ss...)
			testz.Equal(t, tc.want, got)
		})
	}
}

func TestRemove(t *testing.T) {
	testCases := []struct {
		s     []int
		index int
		wantS []int
		wantV int
		wantB bool
	}{
		{
			s:     []int{1, 2, 3},
			index: 1,
			wantS: []int{1, 3},
			wantV: 2,
			wantB: true,
		},
		{
			s:     []int{1, 2, 3},
			index: 0,
			wantS: []int{2, 3},
			wantV: 1,
			wantB: true,
		},
		{
			s:     []int{1, 2, 3},
			index: 2,
			wantS: []int{1, 2},
			wantV: 3,
			wantB: true,
		},
		{
			s:     []int{1, 2, 3},
			index: 3,
			wantS: []int{1, 2, 3},
			wantV: 0,
			wantB: false,
		},
		{
			s:     []int{1, 2, 3},
			index: -1,
			wantS: []int{1, 2, 3},
			wantV: 0,
			wantB: false,
		},
		{
			s:     []int{},
			index: 0,
			wantS: []int{},
			wantV: 0,
			wantB: false,
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("s=%v,index=%d", tc.s, tc.index), func(t *testing.T) {
			s := Copy(tc.s, 0, -1)
			gotS, gotV, gotB := Remove(s, tc.index)
			testz.Equal(t, tc.wantS, gotS)
			testz.Equal(t, tc.wantV, gotV)
			testz.Equal(t, tc.wantB, gotB)
		})
	}
}

func TestBinarySearch(t *testing.T) {
	testCases := []struct {
		s     []int
		v     int
		wantI int
		wantF bool
	}{
		{
			s:     []int{1, 2, 3, 4, 5},
			v:     3,
			wantI: 2,
			wantF: true,
		},
		{
			s:     []int{1, 2, 4, 5},
			v:     3,
			wantI: 2,
			wantF: false,
		},
		{
			s:     []int{1, 2, 3, 4, 5},
			v:     0,
			wantI: 0,
			wantF: false,
		},
		{
			s:     []int{1, 2, 3, 4, 5},
			v:     6,
			wantI: 5,
			wantF: false,
		},
		{
			s:     []int{},
			v:     1,
			wantI: 0,
			wantF: false,
		},
		{
			s:     []int{1, 3, 5},
			v:     2,
			wantI: 1,
			wantF: false,
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("s=%v,v=%d", tc.s, tc.v), func(t *testing.T) {
			gotI, gotF := BinarySearch(tc.s, tc.v)
			testz.Equal(t, tc.wantI, gotI)
			testz.Equal(t, tc.wantF, gotF)
		})
	}
}

func TestBinarySearchFunc(t *testing.T) {
	testCases := []struct {
		s     []int
		v     int
		fn    func(int, int) int
		wantI int
		wantF bool
	}{
		{
			s:     []int{1, 2, 3, 4, 5},
			v:     3,
			fn:    func(a, b int) int { return a - b },
			wantI: 2,
			wantF: true,
		},
		{
			s:     []int{1, 2, 4, 5},
			v:     3,
			fn:    func(a, b int) int { return a - b },
			wantI: 2,
			wantF: false,
		},
		{
			s:     []int{1, 2, 3, 4, 5},
			v:     0,
			fn:    func(a, b int) int { return a - b },
			wantI: 0,
			wantF: false,
		},
		{
			s:     []int{1, 2, 3, 4, 5},
			v:     6,
			fn:    func(a, b int) int { return a - b },
			wantI: 5,
			wantF: false,
		},
		{
			s:     []int{},
			v:     1,
			fn:    func(a, b int) int { return a - b },
			wantI: 0,
			wantF: false,
		},
		{
			s:     []int{1, 3, 5},
			v:     2,
			fn:    func(a, b int) int { return a - b },
			wantI: 1,
			wantF: false,
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("s=%v,v=%d", tc.s, tc.v), func(t *testing.T) {
			gotI, gotF := BinarySearchFunc(tc.s, tc.v, tc.fn)
			testz.Equal(t, tc.wantI, gotI)
			testz.Equal(t, tc.wantF, gotF)
		})
	}
}
