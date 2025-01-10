package setz

import (
	"testing"
)

func TestRoaringBitmap_Add(t *testing.T) {
	num := 10000000
	var m RoaringBitmap
	for i := 1; i <= num; i++ {
		m.Add(uint32(i))
	}

	if m.Len() != num {
		t.Fatalf("invalid len: %d", m.Len())
	}

	for i := 1; i <= num; i++ {
		if !m.Contains(uint32(i)) {
			t.Fatalf("missing %d", i)
		}
	}

	if m.Contains(uint32(num + 1)) {
		t.Fatalf("unexpected %d", num+1)
	}

	iter := m.Iter()
	i := 1
	for iter.Next() {
		if iter.Value() != uint32(i) {
			t.Fatalf("invalid value: %d, expected: %d", iter.Value(), i)
		}
		i++
	}
	i = 1
	m.Range(func(num uint32) bool {
		if num != uint32(i) {
			t.Fatalf("invalid value: %d, expected: %d", num, i)
		}
		i++
		return true
	})

	for i := 1; i <= num; i++ {
		m.Remove(uint32(i))
	}
	if m.Len() != 0 {
		t.Fatalf("invalid len: %d", m.Len())
	}

	iter = m.Iter()
	for iter.Next() {
		t.Fatalf("unexpected value: %d", iter.Value())
	}
}
