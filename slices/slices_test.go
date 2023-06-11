package slices

import "testing"

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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DiffInPlaceFirst(tt.args.s1, tt.args.s2); !Equal(got, tt.want) {
				t.Errorf("DiffReuse() = %v, want %v", got, tt.want)
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
