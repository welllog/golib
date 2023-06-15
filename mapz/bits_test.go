package mapz

import (
	"testing"

	"github.com/welllog/golib/testz"
)

func TestBits_Add(t *testing.T) {
	m := NewBits()
	testz.Equal(t, 0, m.Len(), "init bit map len must be zero")

	m.Add(1)
	testz.Equal(t, 1, m.Len())

	m.Add(2)
	testz.Equal(t, 2, m.Len())

	m.Add(2)
	testz.Equal(t, 2, m.Len())

	t.Log(m.String())
}

func TestBits_Contains(t *testing.T) {
	m := NewBits()

	tests := []struct {
		add uint
		has uint
		res bool
	}{
		{0, 0, true},
		{0, 1, false},
		{1, 1, true},
		{1, 2, false},
		{100, 100, true},
		{922, 923, false},
		{923, 922, true},
		{10000000000, 10000000000, true},
	}

	for _, tt := range tests {
		m.Add(tt.add)
		testz.Equal(t, tt.res, m.Contains(tt.has), "", tt.add, tt.has)
	}
}

func TestBits_Remove(t *testing.T) {
	m := NewBits()

	tests := []struct {
		add uint
		has uint
		res bool
	}{
		{0, 0, true},
		{0, 1, false},
		{1, 1, true},
		{1, 2, false},
		{100, 100, true},
		{922, 923, false},
		{923, 922, true},
		{10000000000, 10000000000, true},
	}

	for _, tt := range tests {
		m.Add(tt.add)
		testz.Equal(t, tt.res, m.Contains(tt.has), "", tt.add, tt.has)
	}

	for _, tt := range tests {
		m.Remove(tt.add)
	}

	for _, tt := range tests {
		testz.Equal(t, false, m.Contains(tt.has), "", tt.add, tt.has)
	}
}
