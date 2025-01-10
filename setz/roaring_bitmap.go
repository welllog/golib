package setz

import (
	"unsafe"

	"github.com/welllog/golib/listz"
)

type container interface {
	Add(x uint16, buf []uint16) (container, bool)
	Remove(x uint16) bool
	Contains(x uint16) bool
	Iter() uint16Iter
	Type() int
	Len() int
}

type uint16Iter interface {
	Next() bool
	Value() uint16
}

type RoaringBitmap struct {
	containers listz.SkipList[uint16, container]
	buf        [4096]uint16
	len        int
}

type RoaringBitmapIter struct {
	node *listz.SkipNode[uint16, container]
	iter uint16Iter
}

// Add adds a number to the roaring bitmap.
func (r *RoaringBitmap) Add(num uint32) bool {
	high := uint16(num >> 16)
	low := uint16(num)

	node := r.containers.GetNode(high)
	if node == nil {
		ac := &arrayContainer{}
		ac.Add(low, r.buf[:])
		r.containers.Set(high, ac)
		r.len++
		return true
	}

	nc, ok := node.Value().Add(low, r.buf[:])
	if ok {
		r.len++
	}
	node.SetValue(nc)
	return ok
}

// Remove removes a number from the roaring bitmap.
func (r *RoaringBitmap) Remove(num uint32) bool {
	high := uint16(num >> 16)
	low := uint16(num)

	c, ok := r.containers.Get(high)
	if !ok {
		return false
	}

	ok = c.Remove(low)
	if ok {
		r.len--
		if c.Len() == 0 {
			r.containers.Remove(high)
		}
	}
	return ok
}

// Contains returns true if the roaring bitmap contains the number.
func (r *RoaringBitmap) Contains(num uint32) bool {
	high := uint16(num >> 16)
	low := uint16(num)

	c, ok := r.containers.Get(high)
	if !ok {
		return false
	}

	return c.Contains(low)
}

// Len returns the number of numbers in the roaring bitmap.
func (r *RoaringBitmap) Len() int {
	return r.len
}

// Iter returns an iterator for the roaring bitmap.
func (r *RoaringBitmap) Iter() RoaringBitmapIter {
	return RoaringBitmapIter{node: r.containers.Head(), iter: nil}
}

// Range calls fn for each number in the roaring bitmap.
func (r *RoaringBitmap) Range(fn func(num uint32) bool) {
	node := r.containers.Head()
	for node != nil {
		high := node.Key()
		c := node.Value()
		if c.Type() == 1 {
			ac := c.(*arrayContainer)
			for _, low := range ac.values {
				if !fn(uint32(high)<<16 | uint32(low)) {
					return
				}
			}
		} else {
			bc := c.(*bitmapContainer)
			for i := 0; i < len(bc.Bitmap.set); i++ {
				for j := 0; j < 64; j++ {
					if bc.Bitmap.set[i]&(1<<j) != 0 {
						if !fn(uint32(high)<<16 | uint32(i<<6+j)) {
							return
						}
					}
				}
			}
		}

		node = node.Next()
	}
}

func (i *RoaringBitmapIter) Next() bool {
	for i.node != nil {
		if i.iter == nil {
			i.iter = i.node.Value().Iter()
		}

		if i.iter.Next() {
			return true
		}

		i.node = i.node.Next()
	}

	return false
}

func (i *RoaringBitmapIter) Value() uint32 {
	return uint32(i.node.Key())<<16 | uint32(i.iter.Value())
}

type arrayContainer struct {
	values []uint16
}

type arrayContainerIter struct {
	c *arrayContainer
	i int
}

func (ac *arrayContainer) Remove(x uint16) bool {
	pos := search(ac.values, x)
	if pos < len(ac.values) && ac.values[pos] == x {
		ac.values = append(ac.values[:pos], ac.values[pos+1:]...)
		return true
	}
	return false
}

func (ac *arrayContainer) Contains(x uint16) bool {
	pos := search(ac.values, x)
	return pos < len(ac.values) && ac.values[pos] == x
}

func (ac *arrayContainer) Add(x uint16, buf []uint16) (container, bool) {
	pos := search(ac.values, x)

	if pos < len(ac.values) && ac.values[pos] == x {
		return ac, false
	}

	if len(ac.values) < 4096 {
		ac.values = append(ac.values, 0)
		copy(ac.values[pos+1:], ac.values[pos:])
		ac.values[pos] = x
		return ac, true
	}

	copy(buf, ac.values)
	// special handling for avoid alloc memory
	t := *(*[1024]uint64)(unsafe.Pointer(&ac.values[0]))
	newContainer := bitmapContainer{Bitmap: Bitmap{set: t[:]}}
	newContainer.setZero()
	for _, v := range buf {
		newContainer.add(uint(v))
	}
	newContainer.add(uint(x))
	// add not update len, so we need update it
	newContainer.length = 4097

	return &newContainer, true
}

func (ac *arrayContainer) Type() int {
	return 1
}

func (ac *arrayContainer) Len() int {
	return len(ac.values)
}

func (ac *arrayContainer) Iter() uint16Iter {
	return &arrayContainerIter{c: ac, i: -1}
}

func (i *arrayContainerIter) Next() bool {
	if i.i < len(i.c.values)-1 {
		i.i++
		return true
	}
	return false
}

func (i *arrayContainerIter) Value() uint16 {
	return i.c.values[i.i]
}

type bitmapContainer Bits

type bitmapContainerIter BitsIter

func (b *bitmapContainer) Add(x uint16, buf []uint16) (container, bool) {
	return b, (*Bits)(b).Add(uint(x))
}

func (b *bitmapContainer) Remove(x uint16) bool {
	return (*Bits)(b).Remove(uint(x))
}

func (b *bitmapContainer) Contains(x uint16) bool {
	return b.Bitmap.Contains(uint(x))
}

func (b *bitmapContainer) Type() int {
	return 2
}

func (b *bitmapContainer) Len() int {
	return (*Bits)(b).Len()
}

func (b *bitmapContainer) Iter() uint16Iter {
	iter := (*Bits)(b).Iter()
	return (*bitmapContainerIter)(&iter)
}

func (i *bitmapContainerIter) Next() bool {
	return (*BitsIter)(i).Next()
}

func (i *bitmapContainerIter) Value() uint16 {
	return uint16((*BitsIter)(i).Value())
}

func (b *bitmapContainer) setZero() {
	for i := 0; i < 1024; i += 32 {
		b.set[i] = 0
		b.set[i+1] = 0
		b.set[i+2] = 0
		b.set[i+3] = 0
		b.set[i+4] = 0
		b.set[i+5] = 0
		b.set[i+6] = 0
		b.set[i+7] = 0
		b.set[i+8] = 0
		b.set[i+9] = 0
		b.set[i+10] = 0
		b.set[i+11] = 0
		b.set[i+12] = 0
		b.set[i+13] = 0
		b.set[i+14] = 0
		b.set[i+15] = 0
		b.set[i+16] = 0
		b.set[i+17] = 0
		b.set[i+18] = 0
		b.set[i+19] = 0
		b.set[i+20] = 0
		b.set[i+21] = 0
		b.set[i+22] = 0
		b.set[i+23] = 0
		b.set[i+24] = 0
		b.set[i+25] = 0
		b.set[i+26] = 0
		b.set[i+27] = 0
		b.set[i+28] = 0
		b.set[i+29] = 0
		b.set[i+30] = 0
		b.set[i+31] = 0
	}
}

func search(values []uint16, x uint16) int {
	low, high := 0, len(values)
	for low < high {
		mid := int(uint(low+high) >> 1)
		if values[mid] < x {
			low = mid + 1
		} else {
			high = mid
		}
	}
	return low
}
