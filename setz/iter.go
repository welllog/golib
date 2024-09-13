//go:build go1.23

package setz

import "iter"

// IterValues returns an iterator that yields all values in the bit set.
func (b *Bits) IterValues() iter.Seq[uint] {
	bi := b.Iter()

	return func(yield func(uint) bool) {
		for bi.Next() {
			if !yield(bi.Value()) {
				break
			}
		}
	}
}
