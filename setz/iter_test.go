//go:build go1.23

package setz

import (
	"slices"
	"testing"

	"github.com/welllog/golib/testz"
)

func TestBits_All(t *testing.T) {
	s := []uint{1, 2, 3, 4, 5, 100, 102, 500, 501, 400, 1000, 9, 10}
	b := Bits{}

	for _, s := range s {
		b.Add(s)
	}

	slices.Sort(s)

	var i int
	for v := range b.All() {
		testz.Equal(t, s[i], v)
		i++
	}
}
