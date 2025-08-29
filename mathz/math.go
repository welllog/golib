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
func Pow[T typez.Integer](x T, n uint) T {
	var ret T = 1
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
// Note: min int will overflow and return min int.
func Abs[T typez.Signed](n T) T {
	i := n >> ((unsafe.Sizeof(n) << 3) - 1)
	//i := n >> 63 // go vet will complain
	return n ^ i - i
}

// BitOnesCount returns the number of bits that are set in n.
func BitOnesCount[T typez.Integer](n T) int {
	count := 0
	for n != 0 {
		count++
		n &= n - 1
	}
	return count
}

// IsPower2 returns true if n is a power of two.
func IsPower2[T typez.Unsigned](n T) bool {
	return n != 0 && (n&(n-1)) == 0
}

// IsEven returns true if n is even.
func IsEven[T typez.Integer](n T) bool {
	return 0 == (n & 1)
}

// Swap swaps the values of a and b.
func Swap[T typez.Integer](a, b *T) {
	*a ^= *b
	*b ^= *a
	*a ^= *b
}

// BinaryInt64 returns the binary representation of n.
func BinaryInt64[T typez.Int64](n T) string {
	// return strconv.FormatUint(uint64(*(*uint)(unsafe.Pointer(&n))), 2)
	return strconv.FormatUint(uint64(n), 2)
}

// BinaryFloat64 returns the binary representation of a float64 number.
func BinaryFloat64(n float64) string {
	return strconv.FormatUint(*(*uint64)(unsafe.Pointer(&n)), 2)
}

// MaxBitApprox return the highest 1 in n
func MaxBitApprox[T typez.Integer](n T) T {
	return 1 << uint(bits.Len64(uint64(n))-1)
}

// MinBitApprox returns the lowest 1 in n
func MinBitApprox[T typez.Signed](n T) T {
	return n & (-n)
}

// EnumToBitMask converts a slice of integers (starting from 1) to a bitmask integer.
// Caller must ensure that the integers are positive and within the range of the bitmask type.
func EnumToBitMask[T typez.Integer](nums []T) T {
	var mask T
	for _, num := range nums {
		if num > 0 {
			mask |= 1 << (num - 1)
		}
	}
	return mask
}

// BitMaskToEnum converts a bitmask integer to a slice of integers (starting from 1).
func BitMaskToEnum[T typez.Integer](mask T) []T {
	var nums []T
	for i := T(0); mask > 0; i++ {
		if mask&1 == 1 {
			nums = append(nums, i+1)
		}
		mask >>= 1
	}
	return nums
}

// BitMaskContains checks if a bitmask contains a specific number (starting from 1).
func BitMaskContains[T typez.Integer](mask T, num T) bool {
	if num <= 0 {
		return false
	}
	return (mask & (1 << (num - 1))) != 0
}

// BitMaskToPower2Enum converts a bitmask integer to a slice of integers that are powers of 2.
func BitMaskToPower2Enum[T typez.Integer](mask T) []T {
	var nums []T
	for i := T(0); mask > 0; i++ {
		if mask&1 == 1 {
			nums = append(nums, 1<<i)
		}
		mask >>= 1
	}
	return nums
}
