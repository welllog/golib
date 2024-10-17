//go:build go1.23

package setz

import "iter"

// All returns an iterator that yields all values in the bit set.
func (b *Bits) All() iter.Seq[uint] {
	return func(yield func(uint) bool) {
		for i := 0; i < len(b.set); i++ {
			for j := 0; j < 64; j++ {
				if b.set[i]&(1<<j) != 0 {
					if !yield(uint(i<<6 + j)) {
						return
					}
				}
			}
		}
	}
}
