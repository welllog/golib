package setz

import (
	"testing"

	"github.com/welllog/golib/testz"
)

func TestBits_Add(t *testing.T) {
	m := Bits{}
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
	m := Bits{}

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
	m := Bits{}

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
	m := Bits{}
	tests := []struct {
		grow uint
		cap  int
	}{
		{0, 64},
		{64, 128},
		{128, 192},
		{0, 192},
		{3, 192},
		{100, 192},
		{192, 256},
	}

	for _, tt := range tests {
		m.Grow(tt.grow)
		testz.Equal(t, tt.cap, m.Cap())
	}
}

func TestBits_Iter(t *testing.T) {
	m := Bits{}
	iter := m.Iter()

	var count int
	for iter.Next() {
		count++
	}
	testz.Equal(t, 0, count)

	for i := 1; i <= 100; i++ {
		m.Add(uint(i))
	}
	iter = m.Iter()
	count = 0
	for iter.Next() {
		count++
		testz.Equal(t, uint(count), iter.Value())
	}
	testz.Equal(t, 100, count)
}

func TestBits_Range(t *testing.T) {
	m := Bits{}
	m.Range(func(num uint) bool {
		t.Error("should not be called")
		return true
	})

	for i := 1; i <= 100; i++ {
		m.Add(uint(i))
	}
	count := 0
	m.Range(func(num uint) bool {
		count++
		testz.Equal(t, uint(count), num)
		return true
	})
	testz.Equal(t, 100, count)
}

func TestBits_Diff(t *testing.T) {
	m1 := Bits{}
	m2 := Bits{}

	for i := 1; i <= 100; i++ {
		m1.Add(uint(i))
	}

	for i := 50; i <= 150; i++ {
		m2.Add(uint(i))
	}

	m1.Diff(m2)
	testz.Equal(t, 49, m1.Len())

	i := 0
	m1.Range(func(num uint) bool {
		i++
		if i != int(num) {
			t.Fatalf("invalid value: %d, expected: %d", num, i)
		}
		return true
	})
	if i != 49 {
		t.Fatalf("invalid count: %d", i)
	}
}

func TestBits_Intersect(t *testing.T) {
	m1 := Bits{}
	m2 := Bits{}

	for i := 1; i <= 100; i++ {
		m1.Add(uint(i))
	}

	for i := 50; i <= 150; i++ {
		m2.Add(uint(i))
	}

	m1.Intersect(m2)
	testz.Equal(t, 51, m1.Len())

	i := 49
	m1.Range(func(num uint) bool {
		i++
		if i != int(num) {
			t.Fatalf("invalid value: %d, expected: %d", num, i)
		}
		return true
	})
	if i != 100 {
		t.Fatalf("invalid count: %d", i)
	}
}

func TestBits_Merge(t *testing.T) {
	m1 := Bits{}
	m2 := Bits{}

	for i := 1; i <= 100; i++ {
		m1.Add(uint(i))
	}

	for i := 50; i <= 150; i++ {
		m2.Add(uint(i))
	}

	m1.Merge(m2)
	testz.Equal(t, 150, m1.Len())

	i := 0
	m1.Range(func(num uint) bool {
		i++
		if i != int(num) {
			t.Fatalf("invalid value: %d, expected: %d", num, i)
		}
		return true
	})
	if i != 150 {
		t.Fatalf("invalid count: %d", i)
	}
}

func TestBits_MergeExt(t *testing.T) {
	// 1. 空合并非空，非空合并空
	empty := Bitmap{}
	nonEmpty := Bitmap{}
	nonEmpty.Add(10)
	nonEmpty.Add(200)

	m1 := empty.Clone()
	m1.Merge(nonEmpty)
	testz.Equal(t, 2, m1.Len())
	testz.Assert(t, m1.Contains(10))
	testz.Assert(t, m1.Contains(200))

	m2 := nonEmpty.Clone()
	m2.Merge(empty)
	testz.Equal(t, 2, m2.Len())
	testz.Assert(t, m2.Contains(10))
	testz.Assert(t, m2.Contains(200))

	// 2. 小集合合并大集合（触发批量 append）
	small := Bitmap{}
	small.Add(1)
	big := Bitmap{}
	for i := uint(1000); i <= 1010; i++ {
		big.Add(i)
	}
	small.Merge(big)
	testz.Equal(t, 12, small.Len())
	testz.Assert(t, small.Contains(1))
	for i := uint(1000); i <= 1010; i++ {
		testz.Assert(t, small.Contains(i))
	}

	// 3. 大集合合并小集合（不触发 append）
	bigClone := big.Clone()
	tiny := Bitmap{}
	tiny.Add(2)
	bigClone.Merge(tiny)
	testz.Equal(t, 12, bigClone.Len())
	testz.Assert(t, bigClone.Contains(2))

	// 4. 自合并
	selfMerge := small.Clone()
	selfMerge.Merge(selfMerge)
	testz.Equal(t, small.Len(), selfMerge.Len())

	// 5. Bits 包装类型 Merge 测试
	b1 := Bits{}
	b1.Add(1)
	b2 := Bits{}
	b2.Add(2)
	b1.Merge(b2)
	testz.Equal(t, 2, b1.Len())
	testz.Assert(t, b1.Contains(1))
	testz.Assert(t, b1.Contains(2))
}

func BenchmarkBits_Add(b *testing.B) {
	b.Run("add", func(b *testing.B) {
		b.ReportAllocs()
		m := Bits{}
		for i := 0; i < b.N; i++ {
			m.Add(uint(i))
		}
	})

	b.Run("add_with_grow", func(b *testing.B) {
		b.ReportAllocs()
		m := Bits{}
		m.Grow(500000000)
		for i := 0; i < b.N; i++ {
			m.Add(uint(i))
		}
	})
}
