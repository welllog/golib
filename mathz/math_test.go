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
		n    int8
		want int8
	}{
		{n: 0, want: 0},
		{n: 1, want: 1},
		{n: -1, want: 1},
		{n: 2, want: 2},
		{n: -2, want: 2},
		{n: 125, want: 125},
		{n: -125, want: 125},
		{n: 126, want: 126},
		{n: -126, want: 126},
		{n: -127, want: 127},
		{n: 127, want: 127},
		{n: -128, want: -128}, // overflow
	}

	for _, tt := range tests {
		testz.Equal(t, tt.want, Abs(tt.n))
	}

	tests2 := []struct {
		n    int16
		want int16
	}{
		{n: 0, want: 0},
		{n: 1, want: 1},
		{n: -1, want: 1},
		{n: 2, want: 2},
		{n: -2, want: 2},
		{n: -32765, want: 32765},
		{n: 32765, want: 32765},
		{n: -32766, want: 32766},
		{n: 32766, want: 32766},
		{n: -32767, want: 32767},
		{n: 32767, want: 32767},
	}

	for _, tt := range tests2 {
		testz.Equal(t, tt.want, Abs(tt.n))
	}

	tests3 := []struct {
		n    int32
		want int32
	}{
		{n: 0, want: 0},
		{n: 1, want: 1},
		{n: -1, want: 1},
		{n: 2, want: 2},
		{n: -2, want: 2},
		{n: -2147483647, want: 2147483647},
		{n: 2147483647, want: 2147483647},
		{n: -2147483646, want: 2147483646},
		{n: 2147483646, want: 2147483646},
		{n: -2147483645, want: 2147483645},
		{n: 2147483645, want: 2147483645},
	}

	for _, tt := range tests3 {
		testz.Equal(t, tt.want, Abs(tt.n))
	}

	tests4 := []struct {
		n    int64
		want int64
	}{
		{n: 0, want: 0},
		{n: 1, want: 1},
		{n: -1, want: 1},
		{n: 2, want: 2},
		{n: -2, want: 2},
		{n: -9223372036854775807, want: 9223372036854775807},
		{n: 9223372036854775807, want: 9223372036854775807},
		{n: -9223372036854775806, want: 9223372036854775806},
		{n: 9223372036854775806, want: 9223372036854775806},
		{n: -9223372036854775805, want: 9223372036854775805},
		{n: 9223372036854775805, want: 9223372036854775805},
	}

	for _, tt := range tests4 {
		testz.Equal(t, tt.want, Abs(tt.n))
	}
}

func TestBitOnesCount(t *testing.T) {
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
		testz.Equal(t, tt.want, BitOnesCount(tt.n))
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
		testz.Equal(t, tt.want, BitOnesCount(tt.n))
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
		testz.Equal(t, tt.want, BitOnesCount(tt.n))
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
		testz.Equal(t, tt.want, BitOnesCount(tt.n))
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

func TestBinaryInt64(t *testing.T) {
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
		testz.Equal(t, tt.want, BinaryInt64(int64(tt.n)), tt.n)
	}
}

func TestBinaryFloat64(t *testing.T) {
	tests := []struct {
		n    float64
		want string
	}{
		{n: 0.0, want: "0"},
		{n: 1.0, want: "11111111110000000000000000000000000000000000000000000000000000"},
		{n: -1.0, want: "1011111111110000000000000000000000000000000000000000000000000000"},
		{n: 2.0, want: "100000000000000000000000000000000000000000000000000000000000000"},
		{n: 3.0, want: "100000000001000000000000000000000000000000000000000000000000000"},
		{n: 0.1, want: "11111110111001100110011001100110011001100110011001100110011010"},
		{n: 0.2, want: "11111111001001100110011001100110011001100110011001100110011010"},
		{n: 0.3, want: "11111111010011001100110011001100110011001100110011001100110011"},
		{n: -0.1, want: "1011111110111001100110011001100110011001100110011001100110011010"},
		{n: -0.2, want: "1011111111001001100110011001100110011001100110011001100110011010"},
	}

	for _, tt := range tests {
		testz.Equal(t, tt.want, BinaryFloat64(tt.n), tt.n)
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

func TestEnumToBitMask(t *testing.T) {
	tests := []struct {
		nums []int
		want int
	}{
		{nums: []int{}, want: 0},
		{nums: []int{1}, want: 1},
		{nums: []int{2}, want: 2},
		{nums: []int{3}, want: 4},
		{nums: []int{4}, want: 8},
		{nums: []int{1, 2}, want: 3},
		{nums: []int{1, 3}, want: 5},
		{nums: []int{2, 3}, want: 6},
		{nums: []int{2, 4}, want: 10},
		{nums: []int{1, 2, 3}, want: 7},
		{nums: []int{1, 2, 3, 4}, want: 15},
	}

	for _, tt := range tests {
		testz.Equal(t, tt.want, EnumToBitMask(tt.nums), tt.nums)
	}
}

func TestBitMaskToEnum(t *testing.T) {
	tests := []struct {
		mask int
		want []int
	}{
		{mask: 0, want: nil},
		{mask: 1, want: []int{1}},
		{mask: 2, want: []int{2}},
		{mask: 3, want: []int{1, 2}},
		{mask: 4, want: []int{3}},
		{mask: 5, want: []int{1, 3}},
		{mask: 6, want: []int{2, 3}},
		{mask: 7, want: []int{1, 2, 3}},
	}

	for _, tt := range tests {
		testz.Equal(t, tt.want, BitMaskToEnum(tt.mask), tt.mask)
	}
}

func TestBitMaskToPower2Enum(t *testing.T) {
	tests := []struct {
		mask int
		want []int
	}{
		{mask: 0, want: nil},
		{mask: 1, want: []int{1}},
		{mask: 2, want: []int{2}},
		{mask: 3, want: []int{1, 2}},
		{mask: 4, want: []int{4}},
		{mask: 5, want: []int{1, 4}},
		{mask: 6, want: []int{2, 4}},
		{mask: 7, want: []int{1, 2, 4}},
	}

	for _, tt := range tests {
		testz.Equal(t, tt.want, BitMaskToPower2Enum(tt.mask), tt.mask)
	}
}
