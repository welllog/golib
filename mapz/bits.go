package mapz

import (
	"bytes"
	"fmt"
)

// Bits is a set of bits.
type Bits struct {
	length int
	set    []uint64
}

// NewBits returns a new Bits.
func NewBits() *Bits {
	return &Bits{}
}

// Add adds a number to the set.
func (b *Bits) Add(num uint) {
	index, bit := num/64, num%64
	grow := int(index) - len(b.set) + 1
	if grow > 0 {
		b.set = append(b.set, make([]uint64, grow)...)
	}
	if b.set[index]&(1<<bit) == 0 {
		b.set[index] |= 1 << bit
		b.length++
	}
}

// Remove removes a number from the set.
func (b *Bits) Remove(num uint) {
	index, bit := num/64, num%64
	if int(index) < len(b.set) && (b.set[index]&(1<<bit)) != 0 {
		b.set[index] &= ^(1 << bit)
		b.length--
	}
}

// Contains returns true if the set contains the number.
func (b *Bits) Contains(num uint) bool {
	index, bit := num/64, num%64
	return int(index) < len(b.set) && (b.set[index]&(1<<bit)) != 0
}

// Len returns the length of the set.
func (b *Bits) Len() int {
	return b.length
}

// String returns a string representation of the set.
func (b *Bits) String() string {
	var buf bytes.Buffer
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
	_, _ = fmt.Fprintf(&buf, "\nLength: %d", b.length)
	return buf.String()
}
