package mathz

import (
	"math/bits"
	"strconv"
	"unsafe"

	"github.com/welllog/golib/typez"
)

const (
	MaxUint = ^uint(0)
	MaxInt  = int(^uint(0) >> 1)
)

// Max returns the maximum value in a slice of numbers.
func Max[T typez.Ordered](n ...T) T {
	var ret T
	if len(n) == 0 {
		return ret
	}

	ret = n[0]
	for _, v := range n[1:] {
		if v > ret {
			ret = v
		}
	}
	return ret
}

// Min returns the minimum value in a slice of numbers.
func Min[T typez.Ordered](n ...T) T {
	var ret T
	if len(n) == 0 {
		return ret
	}

	ret = n[0]
	for _, v := range n[1:] {
		if v < ret {
			ret = v
		}
	}
	return ret
}

// Sum returns the sum of a slice of numbers.
func Sum[T typez.Number](n ...T) T {
	var sum T
	for _, v := range n {
		sum += v
	}
	return sum
}

// Pow returns x**n, the base-x exponential of n.
func Pow(x int, n uint) int {
	var ret = 1
	for n != 0 {
		if (n & 1) != 0 { // n % 2 != 0
			ret = ret * x
		}
		n >>= 1
		x = x * x
	}
	return ret
}

// Abs returns the absolute value of n.
func Abs(n int) int {
	i := n >> (typez.WordBits - 1)
	return n ^ i - i
}

// BitCount returns the number of bits that are set in n.
func BitCount[T typez.Integer](n T) int {
	count := 0
	for n != 0 {
		count++
		n &= n - 1
	}
	return count
}

// IsPower2 returns true if n is a power of two.
func IsPower2(n uint) bool {
	return n != 0 && (n&(n-1)) == 0
}

// IsEven returns true if n is even.
func IsEven(n int) bool {
	return 0 == (n & 1)
}

// Swap swaps the values of a and b.
func Swap(a, b *int) {
	*a ^= *b
	*b ^= *a
	*a ^= *b
}

// Binary returns the binary representation of n.
func Binary(n int) string {
	return strconv.FormatUint(uint64(*(*uint)(unsafe.Pointer(&n))), 2)
}

// MaxBitApprox return the highest bit of n that is 1
func MaxBitApprox(n int) int {
	return 1 << uint(bits.Len(uint(n))-1)
}

// MinBitApprox return the lowest bit of n that is 1
func MinBitApprox(n int) int {
	return n & (-n)
}
