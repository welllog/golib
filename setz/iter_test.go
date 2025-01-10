//go:build go1.23

package setz

import (
	"testing"

	"github.com/welllog/golib/testz"
)

func TestBits_All(t *testing.T) {
	b := Bits{}

	for i := 1; i <= 10000000; i += 8 {
		b.Add(uint(i))
	}

	i := uint(1)
	for v := range b.All() {
		testz.Equal(t, i, v)
		i += 8
	}
}

func TestRoaringBitmap_All(t *testing.T) {
	b := RoaringBitmap{}
	for i := 1; i <= 10000000; i++ {
		b.Add(uint32(i))
	}

	i := 1
	for v := range b.All() {
		testz.Equal(t, uint32(i), v)
		i++
	}
}
