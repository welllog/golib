//go:build go1.23

package setz

import (
	"testing"

	"github.com/welllog/golib/testz"
)

func TestBits_IterValues(t *testing.T) {
	s := []uint{1, 2, 3, 4, 5}
	b := Bits{}

	for _, s := range s {
		b.Add(s)
	}

	var i int
	for v := range b.IterValues() {
		testz.Equal(t, s[i], v)
		i++
	}
}
