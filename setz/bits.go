package setz

import (
	"fmt"
	"strings"
)

// Bits is a set of bits.
type Bits struct {
	length int
	Bitmap
}

// Add adds a number to the set.
func (b *Bits) Add(num uint) {
	if b.Bitmap.Add(num) {
		b.length++
	}
}

// Remove removes a number from the set.
func (b *Bits) Remove(num uint) {
	if b.Bitmap.Remove(num) {
		b.length--
	}
}

// Len returns the length of the set.
func (b *Bits) Len() int {
	return b.length
}

// BitsIter is an iterator for Bits.
type BitsIter struct {
	bits *Bitmap
	i    int
	j    int
	read bool
}

// Next returns true if there is a next value.
func (bi *BitsIter) Next() bool {
	if bi.read {
		bi.read = false
		bi.j++
	}

	for bi.i < len(bi.bits.set) {
		for bi.j < 64 {
			if bi.bits.set[bi.i]&(1<<bi.j) != 0 {
				bi.read = true
				return true
			}
			bi.j++
		}

		bi.i++
		bi.j = 0
	}

	return false
}

// Value returns the current value.
func (bi *BitsIter) Value() uint {
	return uint(bi.i<<6 + bi.j)
}

type Bitmap struct {
	set []uint64
}

// Grow grows the set to the given size.
func (b *Bitmap) Grow(n uint) {
	index := int(n >> 6)
	if index >= len(b.set) {
		grow := index + 1 - len(b.set)
		b.set = append(b.set, make([]uint64, grow)...)
	}
}

// Add adds a number to the set.
func (b *Bitmap) Add(num uint) bool {
	// num/64, num%64
	index, bit := int(num>>6), num&63
	if index >= len(b.set) {
		grow := index + 1 - len(b.set)
		b.set = append(b.set, make([]uint64, grow)...)
		b.set[index] |= 1 << bit
		return true
	}

	if b.set[index]&(1<<bit) == 0 {
		b.set[index] |= 1 << bit
		return true
	}
	return false
}

// Remove removes a number from the set.
func (b *Bitmap) Remove(num uint) bool {
	// num / 64, num % 64
	index, bit := int(num>>6), num&63
	if index < len(b.set) && (b.set[index]&(1<<bit)) != 0 {
		b.set[index] &= ^(1 << bit)
		return true
	}
	return false
}

// Contains returns true if the set contains the number.
func (b *Bitmap) Contains(num uint) bool {
	// num / 64, num % 64
	index, bit := int(num>>6), num&63
	return index < len(b.set) && (b.set[index]&(1<<bit)) != 0
}

// Cap returns the capacity of the set.
func (b *Bitmap) Cap() int {
	return len(b.set) << 6
}

// Iter returns a new BitsIter.
func (b *Bitmap) Iter() BitsIter {
	return BitsIter{bits: b}
}

// String returns a string representation of the set.
func (b *Bitmap) String() string {
	var buf strings.Builder
	buf.WriteByte('{')
	for i, v := range b.set {
		if v == 0 {
			continue
		}
		for j := uint(0); j < 64; j++ {
			if v&(1<<j) != 0 {
				if buf.Len() > len("{") {
					buf.WriteByte(' ')
				}
				_, _ = fmt.Fprintf(&buf, "%d", 64*uint(i)+j)
			}
		}
	}
	buf.WriteByte('}')
	return buf.String()
}
