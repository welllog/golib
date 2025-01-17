package slicez

import (
	"reflect"
	"testing"
)

func TestDiff(t *testing.T) {
	type args struct {
		s1 []int
		s2 []int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "Test case 1",
			args: args{
				s1: []int{1, 2, 3, 4, 5},
				s2: []int{3, 4, 5, 6, 7},
			},
			want: []int{1, 2},
		},
		{
			name: "Test case 2",
			args: args{
				s1: []int{1, 2, 3},
				s2: []int{4, 5, 6},
			},
			want: []int{1, 2, 3},
		},
		{
			name: "Test case 3",
			args: args{
				s1: []int{1, 2, 3},
				s2: []int{1, 2, 3},
			},
			want: []int{},
		},
		{
			name: "Test case 4",
			args: args{
				s1: []int{1, 2, 3},
				s2: []int{},
			},
			want: []int{1, 2, 3},
		},
		{
			name: "Test case 5",
			args: args{
				s1: []int{},
				s2: []int{2, 3},
			},
			want: []int{},
		},
		{
			name: "Test case 6",
			args: args{
				s1: []int{},
				s2: []int{},
			},
			want: []int{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var dst []int
			if got := Diff(dst, tt.args.s1, tt.args.s2); !Equal(got, tt.want) {
				t.Errorf("Diff() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDiffInPlaceFirst(t *testing.T) {
	type args struct {
		s1 []int
		s2 []int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "Test case 1",
			args: args{
				s1: []int{1, 2, 3, 4, 5},
				s2: []int{3, 4, 5, 6, 7},
			},
			want: []int{1, 2},
		},
		{
			name: "Test case 2",
			args: args{
				s1: []int{1, 2, 3},
				s2: []int{4, 5, 6},
			},
			want: []int{1, 2, 3},
		},
		{
			name: "Test case 3",
			args: args{
				s1: []int{1, 2, 3},
				s2: []int{1, 2, 3},
			},
			want: []int{},
		},
		{
			name: "Test case 4",
			args: args{
				s1: []int{1, 2, 3},
				s2: []int{},
			},
			want: []int{1, 2, 3},
		},
		{
			name: "Test case 5",
			args: args{
				s1: []int{},
				s2: []int{2, 3},
			},
			want: []int{},
		},
		{
			name: "Test case 6",
			args: args{
				s1: []int{},
				s2: []int{},
			},
			want: []int{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DiffInPlaceFirst(tt.args.s1, tt.args.s2); !Equal(got, tt.want) {
				t.Errorf("DiffReuse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIntersect(t *testing.T) {
	type args struct {
		s1 []int
		s2 []int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "Test case 1",
			args: args{
				s1: []int{1, 2, 3, 4, 5},
				s2: []int{3, 4, 5, 6, 7},
			},
			want: []int{3, 4, 5},
		},
		{
			name: "Test case 2",
			args: args{
				s1: []int{1, 2, 3},
				s2: []int{4, 5, 6},
			},
			want: []int{},
		},
		{
			name: "Test case 3",
			args: args{
				s1: []int{1, 2, 3},
				s2: []int{1, 2, 3},
			},
			want: []int{1, 2, 3},
		},
		{
			name: "Test case 4",
			args: args{
				s1: []int{1, 2, 3},
				s2: []int{},
			},
			want: []int{},
		},
		{
			name: "test case 5",
			args: args{
				s1: []int{},
				s2: []int{1, 2, 3},
			},
			want: []int{},
		},
		{
			name: "test case 6",
			args: args{
				s1: []int{},
				s2: []int{},
			},
			want: []int{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var dst []int
			if got := Intersect(dst, tt.args.s1, tt.args.s2); !Equal(got, tt.want) {
				t.Errorf("Intersect() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIntersectInPlaceFirst(t *testing.T) {
	type args struct {
		s1 []int
		s2 []int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "Test case 1",
			args: args{
				s1: []int{1, 2, 3, 4, 5},
				s2: []int{3, 4, 5, 6, 7},
			},
			want: []int{3, 4, 5},
		},
		{
			name: "Test case 2",
			args: args{
				s1: []int{1, 2, 3},
				s2: []int{4, 5, 6},
			},
			want: []int{},
		},
		{
			name: "Test case 3",
			args: args{
				s1: []int{1, 2, 3},
				s2: []int{1, 2, 3},
			},
			want: []int{1, 2, 3},
		},
		{
			name: "Test case 4",
			args: args{
				s1: []int{1, 2, 3},
				s2: []int{},
			},
			want: []int{},
		},
		{
			name: "test case 5",
			args: args{
				s1: []int{},
				s2: []int{1, 2, 3},
			},
			want: []int{},
		},
		{
			name: "test case 6",
			args: args{
				s1: []int{},
				s2: []int{},
			},
			want: []int{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IntersectInPlaceFirst(tt.args.s1, tt.args.s2); !Equal(got, tt.want) {
				t.Errorf("Intersect() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnique(t *testing.T) {
	type args struct {
		s []int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "Test case 1",
			args: args{
				s: []int{1, 2, 2, 3, 3, 3, 4, 4, 4, 4},
			},
			want: []int{1, 2, 3, 4},
		},
		{
			name: "Test case 2",
			args: args{
				s: []int{1, 1, 1, 1},
			},
			want: []int{1},
		},
		{
			name: "Test case 3",
			args: args{
				s: []int{},
			},
			want: []int{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var dst []int
			if got := Unique(dst, tt.args.s); !Equal(got, tt.want) {
				t.Errorf("Unique() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUniqueInPlace(t *testing.T) {
	type args struct {
		s []int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "Test case 1",
			args: args{
				s: []int{1, 2, 2, 3, 3, 3, 4, 4, 4, 4},
			},
			want: []int{1, 2, 3, 4},
		},
		{
			name: "Test case 2",
			args: args{
				s: []int{1, 1, 1, 1},
			},
			want: []int{1},
		},
		{
			name: "Test case 3",
			args: args{
				s: []int{},
			},
			want: []int{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := UniqueInPlace(tt.args.s); !Equal(got, tt.want) {
				t.Errorf("UniqueInPlace() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUniqueFunc(t *testing.T) {
	type args struct {
		s []int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "Test case 1",
			args: args{
				s: []int{1, 2, 2, 3, 3, 3, 4, 4, 4, 4},
			},
			want: []int{1, 2, 3, 4},
		},
		{
			name: "Test case 2",
			args: args{
				s: []int{1, 1, 1, 1},
			},
			want: []int{1},
		},
		{
			name: "Test case 3",
			args: args{
				s: []int{},
			},
			want: []int{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var dst []int
			if got := UniqueByKey(dst, tt.args.s, func(n int) int {
				return 2 * n
			}); !Equal(got, tt.want) {
				t.Errorf("Unique() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUniqueFuncInPlace(t *testing.T) {
	type args struct {
		s []int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "Test case 1",
			args: args{
				s: []int{1, 2, 2, 3, 3, 3, 4, 4, 4, 4},
			},
			want: []int{1, 2, 3, 4},
		},
		{
			name: "Test case 2",
			args: args{
				s: []int{1, 1, 1, 1},
			},
			want: []int{1},
		},
		{
			name: "Test case 3",
			args: args{
				s: []int{},
			},
			want: []int{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := UniqueByKeyInPlace(tt.args.s, func(n int) int {
				return 2 * n
			}); !Equal(got, tt.want) {
				t.Errorf("UniqueInPlace() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFilter(t *testing.T) {
	type args struct {
		s         []int
		predicate func(int) bool
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "Test case 1",
			args: args{
				s: []int{1, 2, 3, 4, 5},
				predicate: func(i int) bool {
					return i%2 == 0
				},
			},
			want: []int{2, 4},
		},
		{
			name: "Test case 2",
			args: args{
				s: []int{1, 2, 3},
				predicate: func(i int) bool {
					return i > 3
				},
			},
			want: []int{},
		},
		{
			name: "Test case 3",
			args: args{
				s: []int{1, 2, 3},
				predicate: func(i int) bool {
					return i > 0
				},
			},
			want: []int{1, 2, 3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var dst []int
			if got := Filter(dst, tt.args.s, tt.args.predicate); !Equal(got, tt.want) {
				t.Errorf("Filter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFilterInPlace(t *testing.T) {
	type args struct {
		s         []int
		predicate func(int) bool
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "Test case 1",
			args: args{
				s: []int{1, 2, 3, 4, 5},
				predicate: func(i int) bool {
					return i%2 == 0
				},
			},
			want: []int{2, 4},
		},
		{
			name: "Test case 2",
			args: args{
				s: []int{1, 2, 3},
				predicate: func(i int) bool {
					return i > 3
				},
			},
			want: []int{},
		},
		{
			name: "Test case 3",
			args: args{
				s: []int{1, 2, 3},
				predicate: func(i int) bool {
					return i > 0
				},
			},
			want: []int{1, 2, 3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FilterInPlace(tt.args.s, tt.args.predicate); !Equal(got, tt.want) {
				t.Errorf("FilterInPlace() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEqual(t *testing.T) {
	type args struct {
		s1 []int
		s2 []int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Test case 1",
			args: args{
				s1: []int{1, 2, 3},
				s2: []int{1, 2, 3},
			},
			want: true,
		},
		{
			name: "Test case 2",
			args: args{
				s1: []int{1, 2, 3},
				s2: []int{3, 2, 1},
			},
			want: false,
		},
		{
			name: "Test case 3",
			args: args{
				s1: []int{},
				s2: []int{},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Equal(tt.args.s1, tt.args.s2); got != tt.want {
				t.Errorf("Equal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIndex(t *testing.T) {
	type args struct {
		s []int
		v int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "Test case 1",
			args: args{
				s: []int{1, 2, 3},
				v: 2,
			},
			want: 1,
		},
		{
			name: "Test case 2",
			args: args{
				s: []int{1, 2, 3},
				v: 4,
			},
			want: -1,
		},
		{
			name: "Test case 3",
			args: args{
				s: []int{},
				v: 1,
			},
			want: -1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Index(tt.args.s, tt.args.v); got != tt.want {
				t.Errorf("Index() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIndexFunc(t *testing.T) {
	type args struct {
		s []int
		v int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "Test case 1",
			args: args{
				s: []int{1, 2, 3},
				v: 2,
			},
			want: 1,
		},
		{
			name: "Test case 2",
			args: args{
				s: []int{1, 2, 3},
				v: 4,
			},
			want: -1,
		},
		{
			name: "Test case 3",
			args: args{
				s: []int{},
				v: 1,
			},
			want: -1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IndexFunc(tt.args.s, func(n int) bool {
				return n == tt.args.v
			}); got != tt.want {
				t.Errorf("Index() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContains(t *testing.T) {
	type args struct {
		s []int
		v int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Test case 1",
			args: args{
				s: []int{1, 2, 3},
				v: 2,
			},
			want: true,
		},
		{
			name: "Test case 2",
			args: args{
				s: []int{1, 2, 3},
				v: 4,
			},
			want: false,
		},
		{
			name: "Test case 3",
			args: args{
				s: []int{},
				v: 1,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Contains(tt.args.s, tt.args.v); got != tt.want {
				t.Errorf("Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContainsFunc(t *testing.T) {
	type args struct {
		s []int
		v int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Test case 1",
			args: args{
				s: []int{1, 2, 3},
				v: 2,
			},
			want: true,
		},
		{
			name: "Test case 2",
			args: args{
				s: []int{1, 2, 3},
				v: 4,
			},
			want: false,
		},
		{
			name: "Test case 3",
			args: args{
				s: []int{},
				v: 1,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ContainsFunc(tt.args.s, func(n int) bool {
				return n == tt.args.v
			}); got != tt.want {
				t.Errorf("Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChunk(t *testing.T) {
	tests := []struct {
		name      string
		s         interface{}
		chunkSize int
		want      interface{}
	}{
		{
			name:      "empty slice",
			s:         []int{},
			chunkSize: 2,
			want:      [][]int{},
		},
		{
			name:      "slice smaller than chunk size",
			s:         []int{1, 2},
			chunkSize: 3,
			want:      [][]int{{1, 2}},
		},
		{
			name:      "slice equal to chunk size",
			s:         []int{1, 2, 3},
			chunkSize: 3,
			want:      [][]int{{1, 2, 3}},
		},
		{
			name:      "slice larger than chunk size",
			s:         []int{1, 2, 3, 4, 5},
			chunkSize: 2,
			want:      [][]int{{1, 2}, {3, 4}, {5}},
		},
		{
			name:      "slice length not a multiple of chunk size",
			s:         []int{1, 2, 3, 4, 5},
			chunkSize: 3,
			want:      [][]int{{1, 2, 3}, {4, 5}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Chunk(tt.s.([]int), tt.chunkSize)
			if len(tt.s.([]int)) == 0 {
				if len(got) > 0 {
					t.Errorf("Chunk() = %v, want %v", got, tt.want)
				}
				return
			}
			if !reflect.DeepEqual(got, tt.want.([][]int)) {
				t.Errorf("Chunk() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChunkProcess(t *testing.T) {
	tests := []struct {
		args   []int
		chunk  int
		expect [][]int
	}{
		{
			[]int{1, 2, 3, 4, 5, 6},
			2,
			[][]int{{1, 2}, {3, 4}, {5, 6}},
		},
		{
			[]int{1, 2, 3, 4, 5, 6},
			3,
			[][]int{{1, 2, 3}, {4, 5, 6}},
		},
		{
			[]int{1, 2, 3, 4, 5, 6},
			4,
			[][]int{{1, 2, 3, 4}, {5, 6}},
		},
		{
			[]int{1, 2, 3, 4, 5, 6},
			6,
			[][]int{{1, 2, 3, 4, 5, 6}},
		},
		{
			[]int{1, 2, 3, 4, 5, 6},
			7,
			[][]int{{1, 2, 3, 4, 5, 6}},
		},
		{
			[]int{},
			1,
			[][]int{{}},
		},
	}

	for _, tt := range tests {
		var i int
		_ = ChunkProcess(tt.args, tt.chunk, func(arr []int) error {
			if !reflect.DeepEqual(arr, tt.expect[i]) {
				t.Errorf("Chunk = %v, want %v", arr, tt.expect[i])
			}
			i++
			return nil
		})
	}
}

func TestCopy(t *testing.T) {
	type args[T any] struct {
		arr    []T
		start  int
		length int
	}
	type testCase[T any] struct {
		name string
		args args[T]
		want []T
	}
	tests := []testCase[int]{
		{
			name: "case 1",
			args: args[int]{
				arr:    []int{1, 2, 3},
				start:  0,
				length: -2,
			},
			want: []int{1, 2, 3},
		},
		{
			name: "case 2",
			args: args[int]{
				arr:    []int{1, 2, 3},
				start:  3,
				length: 2,
			},
			want: []int{},
		},
		{
			name: "case 3",
			args: args[int]{
				arr:    []int{1, 2, 3},
				start:  1,
				length: 3,
			},
			want: []int{2, 3},
		},
		{
			name: "case 4",
			args: args[int]{
				arr:    []int{1, 2, 3},
				start:  1,
				length: 2,
			},
			want: []int{2, 3},
		},
		{
			name: "case 5",
			args: args[int]{
				arr:    []int{},
				start:  1,
				length: 2,
			},
			want: []int{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Copy(tt.args.arr, tt.args.start, tt.args.length)
			if !Equal(got, tt.want) {
				t.Errorf("Clone() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValues(t *testing.T) {
	type args struct {
		fn func(int) int
		ss [][]int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "Values with multiple slices",
			args: args{
				fn: func(i int) int { return i * 2 },
				ss: [][]int{{1, 2, 3}, {4, 5, 6}},
			},
			want: []int{2, 4, 6, 8, 10, 12},
		},
		{
			name: "Values with single slice",
			args: args{
				fn: func(i int) int { return i * 2 },
				ss: [][]int{{1, 2, 3}},
			},
			want: []int{2, 4, 6},
		},
		{
			name: "Values with empty slice",
			args: args{
				fn: func(i int) int { return i * 2 },
				ss: [][]int{{}},
			},
			want: []int{},
		},
		{
			name: "Values with no slices",
			args: args{
				fn: func(i int) int { return i * 2 },
				ss: [][]int{},
			},
			want: []int{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Values(tt.args.fn, tt.args.ss...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Values() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSubSlice(t *testing.T) {
	type args struct {
		s     []int
		start int
		end   int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "Normal case",
			args: args{
				s:     []int{1, 2, 3, 4, 5},
				start: 1,
				end:   3,
			},
			want: []int{2, 3},
		},
		{
			name: "Start index greater than length",
			args: args{
				s:     []int{1, 2, 3, 4, 5},
				start: 6,
				end:   8,
			},
			want: nil,
		},
		{
			name: "Start index less than 0",
			args: args{
				s:     []int{1, 2, 3, 4, 5},
				start: -1,
				end:   3,
			},
			want: []int{1, 2, 3},
		},
		{
			name: "End index less than 0",
			args: args{
				s:     []int{1, 2, 3, 4, 5},
				start: 1,
				end:   -1,
			},
			want: []int{2, 3, 4, 5},
		},
		{
			name: "End index greater than length",
			args: args{
				s:     []int{1, 2, 3, 4, 5},
				start: 1,
				end:   10,
			},
			want: []int{2, 3, 4, 5},
		},
		{
			name: "Start index equal to end index",
			args: args{
				s:     []int{1, 2, 3, 4, 5},
				start: 2,
				end:   2,
			},
			want: nil,
		},
		{
			name: "Start index greater than end index",
			args: args{
				s:     []int{1, 2, 3, 4, 5},
				start: 3,
				end:   2,
			},
			want: nil,
		},
		{
			name: "Empty slice",
			args: args{
				s:     []int{},
				start: 0,
				end:   1,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SubSlice(tt.args.s, tt.args.start, tt.args.end); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SubSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestRemove(t *testing.T) {
	type args struct {
		s     []int
		index int
	}
	tests := []struct {
		name      string
		args      args
		wantSlice []int
		wantElem  int
		wantOk    bool
	}{
		{
			name: "Remove element at index 0",
			args: args{
				s:     []int{1, 2, 3, 4, 5},
				index: 0,
			},
			wantSlice: []int{2, 3, 4, 5},
			wantElem:  1,
			wantOk:    true,
		},
		{
			name: "Remove element at index 2",
			args: args{
				s:     []int{1, 2, 3, 4, 5},
				index: 2,
			},
			wantSlice: []int{1, 2, 4, 5},
			wantElem:  3,
			wantOk:    true,
		},
		{
			name: "Remove element at last index",
			args: args{
				s:     []int{1, 2, 3, 4, 5},
				index: 4,
			},
			wantSlice: []int{1, 2, 3, 4},
			wantElem:  5,
			wantOk:    true,
		},
		{
			name: "Remove element from empty slice",
			args: args{
				s:     []int{},
				index: 0,
			},
			wantSlice: []int{},
			wantElem:  0,
			wantOk:    false,
		},
		{
			name: "Remove element with negative index",
			args: args{
				s:     []int{1, 2, 3, 4, 5},
				index: -1,
			},
			wantSlice: []int{1, 2, 3, 4, 5},
			wantElem:  0,
			wantOk:    false,
		},
		{
			name: "Remove element with out of range index",
			args: args{
				s:     []int{1, 2, 3, 4, 5},
				index: 5,
			},
			wantSlice: []int{1, 2, 3, 4, 5},
			wantElem:  0,
			wantOk:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSlice, gotElem, gotOk := Remove(tt.args.s, tt.args.index)
			if !reflect.DeepEqual(gotSlice, tt.wantSlice) {
				t.Errorf("Remove() gotSlice = %v, want %v", gotSlice, tt.wantSlice)
			}
			if gotElem != tt.wantElem {
				t.Errorf("Remove() gotElem = %v, want %v", gotElem, tt.wantElem)
			}
			if gotOk != tt.wantOk {
				t.Errorf("Remove() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}
