package hashz

import (
	"github.com/welllog/golib/typez"
)

// BKDRHash BKDR Hash Function
func BKDRHash[T typez.StrOrBytes](s T) uint32 {
	var seed uint32 = 131 // 31 131 1313 13131 131313 etc..
	var hash uint32 = 0
	for i := 0; i < len(s); i++ {
		hash = hash*seed + uint32(s[i])
	}
	return hash & 0x7FFFFFFF // 0x7FFFFFFF = 2^31 - 1
}

// BKDRHash64 BKDR Hash Function
func BKDRHash64[T typez.StrOrBytes](s T) uint64 {
	var seed uint64 = 131 // 31 131 1313 13131 131313 etc..
	var hash uint64 = 0
	for i := 0; i < len(s); i++ {
		hash = hash*seed + uint64(s[i])
	}
	return hash & 0x7FFFFFFFFFFFFFFF // 0x7FFFFFFFFFFFFFFF = 2^63 - 1
}

// APHash AP Hash Function
func APHash[T typez.StrOrBytes](s T) uint32 {
	var hash uint32 = 0
	for i := 0; i < len(s); i++ {
		if (i & 1) == 0 {
			hash ^= (hash << 7) ^ uint32(s[i]) ^ (hash >> 3)
		} else {
			hash ^= ^((hash << 11) ^ uint32(s[i]) ^ (hash >> 5)) + 1
		}
	}
	return hash & 0x7FFFFFFF
}

// APHash64 AP Hash Function
func APHash64[T typez.StrOrBytes](s T) uint64 {
	var hash uint64 = 0
	for i := 0; i < len(s); i++ {
		if (i & 1) == 0 {
			hash ^= (hash << 7) ^ uint64(s[i]) ^ (hash >> 3)
		} else {
			hash ^= ^((hash << 11) ^ uint64(s[i]) ^ (hash >> 5)) + 1
		}
	}
	return hash & 0x7FFFFFFFFFFFFFFF
}

// DJBHash DJB Hash Function
func DJBHash[T typez.StrOrBytes](s T) uint32 {
	var hash uint32 = 5381
	for i := 0; i < len(s); i++ {
		hash += (hash << 5) + uint32(s[i])
	}
	return hash & 0x7FFFFFFF
}

// DJBHash64 DJB Hash Function
func DJBHash64[T typez.StrOrBytes](s T) uint64 {
	var hash uint64 = 5381
	for i := 0; i < len(s); i++ {
		hash += (hash << 5) + uint64(s[i])
	}
	return hash & 0x7FFFFFFFFFFFFFFF
}

// JSHash JS Hash Function
func JSHash[T typez.StrOrBytes](s T) uint32 {
	var hash uint32 = 1315423911
	for i := 0; i < len(s); i++ {
		hash ^= (hash << 5) + uint32(s[i]) + (hash >> 2)
	}
	return hash & 0x7FFFFFFF
}

// JSHash64 JS Hash Function
func JSHash64[T typez.StrOrBytes](s T) uint64 {
	var hash uint64 = 1315423911
	for i := 0; i < len(s); i++ {
		hash ^= (hash << 5) + uint64(s[i]) + (hash >> 2)
	}
	return hash & 0x7FFFFFFFFFFFFFFF
}

// RSHash RS Hash Function
func RSHash[T typez.StrOrBytes](s T) uint32 {
	var b uint32 = 378551
	var a uint32 = 63689
	var hash uint32 = 0
	for i := 0; i < len(s); i++ {
		hash = hash*a + uint32(s[i])
		a *= b
	}
	return hash & 0x7FFFFFFF
}

// RSHash64 RS Hash Function
func RSHash64[T typez.StrOrBytes](s T) uint64 {
	var b uint64 = 378551
	var a uint64 = 63689
	var hash uint64 = 0
	for i := 0; i < len(s); i++ {
		hash = hash*a + uint64(s[i])
		a *= b
	}
	return hash & 0x7FFFFFFFFFFFFFFF
}

// SDBMHash SDBM Hash Function
func SDBMHash[T typez.StrOrBytes](s T) uint32 {
	var hash uint32 = 0
	for i := 0; i < len(s); i++ {
		// equivalent to: hash = 65599*hash + uint32(str[i]);
		hash = uint32(s[i]) + (hash << 6) + (hash << 16) - hash
	}
	return hash & 0x7FFFFFFF
}

// SDBMHash64 SDBM Hash Function
func SDBMHash64[T typez.StrOrBytes](s T) uint64 {
	var hash uint64 = 0
	for i := 0; i < len(s); i++ {
		// equivalent to: hash = 65599*hash + uint32(str[i]);
		hash = uint64(s[i]) + (hash << 6) + (hash << 16) - hash
	}
	return hash & 0x7FFFFFFFFFFFFFFF
}

// PJWHash PJW Hash Function
func PJWHash[T typez.StrOrBytes](s T) uint32 {
	var BitsInUnignedInt uint32 = 4 * 8
	var ThreeQuarters = (BitsInUnignedInt * 3) / 4
	var OneEighth = BitsInUnignedInt / 8
	var HighBits uint32 = (0xFFFFFFFF) << (BitsInUnignedInt - OneEighth)
	var hash uint32 = 0
	var test uint32 = 0
	for i := 0; i < len(s); i++ {
		hash = (hash << OneEighth) + uint32(s[i])
		if test = hash & HighBits; test != 0 {
			hash = (hash ^ (test >> ThreeQuarters)) & (^HighBits + 1)
		}
	}
	return hash & 0x7FFFFFFF
}

// PJWHash64 PJW Hash Function
func PJWHash64[T typez.StrOrBytes](s T) uint64 {
	var BitsInUnignedInt uint64 = 4 * 8
	var ThreeQuarters = (BitsInUnignedInt * 3) / 4
	var OneEighth = BitsInUnignedInt / 8
	var HighBits uint64 = (0xFFFFFFFF) << (BitsInUnignedInt - OneEighth)
	var hash uint64 = 0
	var test uint64 = 0
	for i := 0; i < len(s); i++ {
		hash = (hash << OneEighth) + uint64(s[i])
		if test = hash & HighBits; test != 0 {
			hash = (hash ^ (test >> ThreeQuarters)) & (^HighBits + 1)
		}
	}
	return hash & 0x7FFFFFFFFFFFFFFF
}

// ELFHash ELF Hash Function
func ELFHash[T typez.StrOrBytes](s T) uint32 {
	var hash uint32 = 0
	var x uint32 = 0
	for i := 0; i < len(s); i++ {
		hash = (hash << 4) + uint32(s[i])
		if x = hash & 0xF0000000; x != 0 {
			hash ^= x >> 24
			hash &= ^x + 1
		}
	}
	return hash & 0x7FFFFFFF
}

// ELFHash64 ELF Hash Function
func ELFHash64[T typez.StrOrBytes](s T) uint64 {
	var hash uint64 = 0
	var x uint64 = 0
	for i := 0; i < len(s); i++ {
		hash = (hash << 4) + uint64(s[i])
		if x = hash & 0xF0000000; x != 0 {
			hash ^= x >> 24
			hash &= ^x + 1
		}
	}
	return hash & 0x7FFFFFFFFFFFFFFF
}
