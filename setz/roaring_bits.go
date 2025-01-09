package setz

import (
	"github.com/welllog/golib/listz"
)

type container interface {
	Add(x uint16) container
	Remove(x uint16) container
	Contains(x uint16) bool
}

type RoaringBits struct {
	containers listz.SkipList[uint16, container]
}

func (r *RoaringBits) Add(num uint32) {
	// high := uint16
}

type arrayContainer struct {
	values []uint16
}

func (ac *arrayContainer) Remove(x uint16) container {
	// TODO implement me
	panic("implement me")
}

func (ac *arrayContainer) Contains(x uint16) bool {
	// TODO implement me
	panic("implement me")
}

func (ac *arrayContainer) Add(x uint16) container {
	pos := search(ac.values, x)

	if pos < len(ac.values) && ac.values[pos] == x {
		return ac
	}

	if len(ac.values) < 4096 {
		ac.values = append(ac.values, 0)
		copy(ac.values[pos+1:], ac.values[pos:])
		ac.values[pos] = x
		return ac
	}

	return nil
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
