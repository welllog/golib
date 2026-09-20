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
	if i != num+1 {
		t.Fatalf("iterated %d values, expected %d", i-1, num)
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

func TestRoaringBitmap_Iter(t *testing.T) {
	var m RoaringBitmap
	m.Add(1)
	m.Add(65537)
	m.Add(131073)

	if m.Len() != 3 {
		t.Fatalf("invalid len: %d", m.Len())
	}

	// values span multiple containers, all of them must be iterated
	want := []uint32{1, 65537, 131073}
	var got []uint32
	iter := m.Iter()
	for iter.Next() {
		got = append(got, iter.Value())
	}
	if len(got) != len(want) {
		t.Fatalf("iterated %d values %v, expected: %d values %v", len(got), got, len(want), want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("invalid value: %d, expected: %d", got[i], want[i])
		}
	}
}
