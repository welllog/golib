package mathz

import (
	"testing"

	"github.com/welllog/golib/testz"
)

func TestSum(t *testing.T) {
	type args struct {
		s []int
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
			},
			want: 6,
		},
		{
			name: "Test case 2",
			args: args{
				s: []int{1, -1, 2, -2},
			},
			want: 0,
		},
		{
			name: "Test case 3",
			args: args{
				s: []int{},
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Sum(tt.args.s...); got != tt.want {
				t.Errorf("Sum() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMax(t *testing.T) {
	type args struct {
		s []int
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
			},
			want: 3,
		},
		{
			name: "Test case 2",
			args: args{
				s: []int{1, -1, 2, -2},
			},
			want: 2,
		},
		{
			name: "Test case 3",
			args: args{
				s: []int{},
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Max(tt.args.s...); got != tt.want {
				t.Errorf("Max() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMin(t *testing.T) {
	type args struct {
		s []int
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
			},
			want: 1,
		},
		{
			name: "Test case 2",
			args: args{
				s: []int{1, -1, 2, -2},
			},
			want: -2,
		},
		{
			name: "Test case 3",
			args: args{
				s: []int{},
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Min(tt.args.s...); got != tt.want {
				t.Errorf("Max() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPow(t *testing.T) {
	tests := []struct {
		n    int
		p    uint
		want int
	}{
		{n: 2, p: 0, want: 1},
		{n: 2, p: 1, want: 2},
		{n: 2, p: 2, want: 4},
		{n: 2, p: 3, want: 8},
		{n: 2, p: 4, want: 16},
		{n: -1, p: 0, want: 1},
		{n: -1, p: 1, want: -1},
		{n: -1, p: 2, want: 1},
		{n: -1, p: 3, want: -1},
		{n: -2, p: 0, want: 1},
		{n: -2, p: 1, want: -2},
		{n: -2, p: 2, want: 4},
		{n: -2, p: 3, want: -8},
	}

	for _, tt := range tests {
		testz.Equal(t, tt.want, Pow(tt.n, tt.p))
	}
}

func TestAbs(t *testing.T) {
	tests := []struct {
		n    int
		want int
	}{
		{n: 0, want: 0},
		{n: 1, want: 1},
		{n: -1, want: 1},
		{n: 2, want: 2},
		{n: -2, want: 2},
	}

	for _, tt := range tests {
		testz.Equal(t, tt.want, Abs(tt.n))
	}
}

func TestBitCount(t *testing.T) {
	testsInt64 := []struct {
		n    int64
		want int
	}{
		{n: 0, want: 0},
		{n: 1, want: 1},
		{n: 2, want: 1},
		{n: 3, want: 2},
		{n: -1, want: 64},
		{n: -2, want: 63},
		{n: -3, want: 63},
		{n: -4, want: 62},
	}

	for _, tt := range testsInt64 {
		testz.Equal(t, tt.want, BitCount(tt.n))
	}

	testsInt8 := []struct {
		n    int8
		want int
	}{
		{n: 0, want: 0},
		{n: 1, want: 1},
		{n: 2, want: 1},
		{n: 3, want: 2},
		{n: -1, want: 8},
		{n: -2, want: 7},
		{n: -3, want: 7},
		{n: -4, want: 6},
	}

	for _, tt := range testsInt8 {
		testz.Equal(t, tt.want, BitCount(tt.n))
	}

	testsInt16 := []struct {
		n    int16
		want int
	}{
		{n: 0, want: 0},
		{n: 1, want: 1},
		{n: 2, want: 1},
		{n: 3, want: 2},
		{n: -1, want: 16},
		{n: -2, want: 15},
		{n: -3, want: 15},
		{n: -4, want: 14},
	}

	for _, tt := range testsInt16 {
		testz.Equal(t, tt.want, BitCount(tt.n))
	}

	testsInt32 := []struct {
		n    int32
		want int
	}{
		{n: 0, want: 0},
		{n: 1, want: 1},
		{n: 2, want: 1},
		{n: 3, want: 2},
		{n: -1, want: 32},
		{n: -2, want: 31},
		{n: -3, want: 31},
		{n: -4, want: 30},
	}

	for _, tt := range testsInt32 {
		testz.Equal(t, tt.want, BitCount(tt.n))
	}
}

func TestIsPower2(t *testing.T) {
	tests := []struct {
		n    uint
		want bool
	}{
		{n: 0, want: false},
		{n: 1, want: true},
		{n: 2, want: true},
		{n: 3, want: false},
		{n: 4, want: true},
		{n: 5, want: false},
		{n: 6, want: false},
		{n: 7, want: false},
		{n: 8, want: true},
	}

	for _, tt := range tests {
		testz.Equal(t, tt.want, IsPower2(tt.n), tt.n)
	}
}

func TestIsEven(t *testing.T) {
	tests := []struct {
		n    int
		want bool
	}{
		{n: 0, want: true},
		{n: 1, want: false},
		{n: 2, want: true},
		{n: 3, want: false},
		{n: -1, want: false},
		{n: -2, want: true},
		{n: -3, want: false},
		{n: -8, want: true},
	}

	for _, tt := range tests {
		testz.Equal(t, tt.want, IsEven(tt.n))
	}
}

func TestSwap(t *testing.T) {
	tests := []struct {
		a, b int
	}{
		{a: 0, b: 0},
		{a: 1, b: 0},
		{a: 0, b: 1},
		{a: 1, b: 1},
		{a: 2, b: 1},
		{a: 1, b: 2},
		{a: -1, b: 0},
		{a: 0, b: -1},
		{a: -1, b: -1},
		{a: -2, b: -1},
		{a: -1, b: -2},
	}

	for _, tt := range tests {
		a, b := tt.a, tt.b
		Swap(&a, &b)
		testz.Equal(t, a, tt.b)
		testz.Equal(t, b, tt.a)
	}
}

func TestBinary(t *testing.T) {
	tests := []struct {
		n    int
		want string
	}{
		{n: 0, want: "0"},
		{n: 1, want: "1"},
		{n: 2, want: "10"},
		{n: 3, want: "11"},
		{n: 4, want: "100"},
		{n: 5, want: "101"},
		{n: 6, want: "110"},
		{n: -1, want: "1111111111111111111111111111111111111111111111111111111111111111"},
		{n: -2, want: "1111111111111111111111111111111111111111111111111111111111111110"},
		{n: -3, want: "1111111111111111111111111111111111111111111111111111111111111101"},
		{n: -4, want: "1111111111111111111111111111111111111111111111111111111111111100"},
	}

	for _, tt := range tests {
		testz.Equal(t, tt.want, Binary(tt.n), tt.n)
	}
}

func TestMaxBitApprox(t *testing.T) {
	tests := []struct {
		n    int
		want int
	}{
		{n: 0, want: 0},
		{n: 1, want: 1},
		{n: 2, want: 2},
		{n: 3, want: 2},
		{n: 4, want: 4},
		{n: 5, want: 4},
		{n: 6, want: 4},
		{n: 7, want: 4},
		{n: 8, want: 8},
		{n: 9, want: 8},
		{n: 10, want: 8},
		{n: 11, want: 8},
		{n: -1, want: -9223372036854775808},
		{n: -100, want: -9223372036854775808},
	}

	for _, tt := range tests {
		testz.Equal(t, tt.want, MaxBitApprox(tt.n), tt.n)
	}
}

func TestMinBitApprox(t *testing.T) {
	tests := []struct {
		n    int
		want int
	}{
		{n: 0, want: 0},
		{n: 1, want: 1},
		{n: 2, want: 2},
		{n: 3, want: 1},
		{n: 4, want: 4},
		{n: 5, want: 1},
		{n: 6, want: 2},
		{n: 7, want: 1},
		{n: 8, want: 8},
		{n: -1, want: 1},
		{n: -2, want: 2},
		{n: -3, want: 1},
		{n: -4, want: 4},
		{n: -100, want: 4},
	}

	for _, tt := range tests {
		testz.Equal(t, tt.want, MinBitApprox(tt.n), tt.n)
	}
}
