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

func TestBits_Grow(t *testing.T) {
	m := NewBits()
	tests := []struct {
		grow int
		cap  int
	}{
		{-1, 0},
		{0, 64},
		{64, 128},
		{128, 192},
		{0, 192},
		{-1, 192},
		{100, 192},
		{192, 256},
	}

	for _, tt := range tests {
		m.Grow(tt.grow)
		testz.Equal(t, tt.cap, m.Cap())
	}
}

func TestBits_Iter(t *testing.T) {
	m := NewBits()
	iter := m.Iter()

	var count int
	for iter.Next() {
		count++
	}
	testz.Equal(t, 0, count)

	m.Add(1)
	m.Add(2)
	m.Add(3)

	iter = m.Iter()
	count = 0
	for iter.Next() {
		count++
		t.Log(iter.Value())
	}
	testz.Equal(t, 3, count)

	m.Add(128)
	iter = m.Iter()
	count = 0
	for iter.Next() {
		count++
		t.Log(iter.Value())
	}
	testz.Equal(t, 4, count)

	m.Add(256)
	m.Add(999)
	for iter.Next() {
		count++
		t.Log(iter.Value())
	}
	testz.Equal(t, 6, count)
}

func BenchmarkBits_Add(b *testing.B) {
	b.Run("add", func(b *testing.B) {
		b.ReportAllocs()
		m := NewBits()
		for i := 0; i < b.N; i++ {
			m.Add(uint(i))
		}
	})

	b.Run("add_with_grow", func(b *testing.B) {
		b.ReportAllocs()
		m := NewBits()
		m.Grow(500000000)
		for i := 0; i < b.N; i++ {
			m.Add(uint(i))
		}
	})
}
