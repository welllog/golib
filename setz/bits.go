package setz

import (
	"fmt"
	"math/bits"
	"strings"
)

// Bits is a set of bits.
type Bits struct {
	length int
	Bitmap
}

// Add adds a number to the set.
func (b *Bits) Add(num uint) bool {
	if b.Bitmap.Add(num) {
		b.length++
		return true
	}
	return false
}

// Remove removes a number from the set.
func (b *Bits) Remove(num uint) bool {
	if b.Bitmap.Remove(num) {
		b.length--
		return true
	}
	return false
}

// Len returns the length of the set.
func (b *Bits) Len() int {
	return b.length
}

// Diff removes numbers in other from the set.
func (b *Bits) Diff(other Bits) {
	b.Bitmap.Diff(other.Bitmap)
	b.length = b.Bitmap.Len()
}

// Intersect removes numbers not in other from the set.
func (b *Bits) Intersect(other Bits) {
	b.Bitmap.Intersect(other.Bitmap)
	b.length = b.Bitmap.Len()
}

// Merge adds numbers in other to the set.
func (b *Bits) Merge(other Bits) {
	b.Bitmap.Merge(other.Bitmap)
	b.length = b.Bitmap.Len()
}

// BitmapIter is an iterator for Bitmap.
type BitmapIter struct {
	bm   *Bitmap
	i    int
	j    int
	read bool
}

// Next returns true if there is a next value.
func (bi *BitmapIter) Next() bool {
	if bi.read {
		bi.read = false
		bi.j++
	}

	for bi.i < len(bi.bm.set) {
		for bi.j < 64 {
			if bi.bm.set[bi.i]&(1<<bi.j) != 0 {
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
func (bi *BitmapIter) Value() uint {
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

// Len returns the length of the set.
func (b *Bitmap) Len() int {
	var count int
	for _, v := range b.set {
		count += bits.OnesCount64(v)
	}
	return count
}

// Cap returns the capacity of the set.
func (b *Bitmap) Cap() int {
	return len(b.set) << 6
}

// Iter returns a new BitmapIter.
func (b *Bitmap) Iter() BitmapIter {
	return BitmapIter{bm: b}
}

// Range calls fn for each number in the set.
func (b *Bitmap) Range(fn func(uint) bool) {
	for i := 0; i < len(b.set); i++ {
		for j := 0; j < 64; j++ {
			if b.set[i]&(1<<j) != 0 {
				if !fn(uint(i<<6 + j)) {
					return
				}
			}
		}
	}
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

// Diff removes numbers in other from the set.
func (b *Bitmap) Diff(other Bitmap) {
	for i := 0; i < len(b.set); i++ {
		if i >= len(other.set) {
			break
		}
		b.set[i] &= ^other.set[i]
	}
}

// Intersect removes numbers not in other from the set.
func (b *Bitmap) Intersect(other Bitmap) {
	for i := 0; i < len(b.set); i++ {
		if i >= len(other.set) {
			b.set[i] = 0
			continue
		}
		b.set[i] &= other.set[i]
	}
}

// Merge adds numbers in other to the set.
func (b *Bitmap) Merge(other Bitmap) {
	for i := 0; i < len(other.set); i++ {
		if i >= len(b.set) {
			b.set = append(b.set, other.set[i])
			continue
		}
		b.set[i] |= other.set[i]
	}
}

// Clone returns a copy of the set.
func (b *Bitmap) Clone() Bitmap {
	set := make([]uint64, len(b.set))
	copy(set, b.set)
	return Bitmap{set: set}
}

func (b *Bitmap) add(num uint) {
	// num/64, num%64
	index, bit := int(num>>6), num&63
	b.set[index] |= 1 << bit
}
