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

// All returns an iterator that yields all values in the roaring bitmap.
func (r *RoaringBitmap) All() iter.Seq[uint32] {
	return func(yield func(uint32) bool) {
		node := r.containers.Head()
		for node != nil {
			high := node.Key()
			c := node.Value()
			if c.Type() == 1 {
				ac := c.(*arrayContainer)
				for _, low := range ac.values {
					if !yield(uint32(high)<<16 | uint32(low)) {
						return
					}
				}
			} else {
				bc := c.(*bitmapContainer)
				for i := 0; i < len(bc.Bitmap.set); i++ {
					for j := 0; j < 64; j++ {
						if bc.Bitmap.set[i]&(1<<j) != 0 {
							if !yield(uint32(high)<<16 | uint32(i<<6+j)) {
								return
							}
						}
					}
				}
			}

			node = node.Next()
		}
	}
}
