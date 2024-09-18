package heapz

import (
	"container/heap"
	"math/rand"
	"testing"
)

func (h *Heap[T]) verify(t *testing.T, i int) {
	t.Helper()
	n := h.Len()
	j1 := 2*i + 1
	j2 := 2*i + 2
	if j1 < n {
		if h.cmp(h.values[j1], h.values[i]) {
			t.Errorf("heap invariant invalidated [%d] = %v > [%d] = %v", i, h.values[i].Value, j1, h.values[j1])
			return
		}
		h.verify(t, j1)
	}
	if j2 < n {
		if h.cmp(h.values[j2], h.values[i]) {
			t.Errorf("heap invariant invalidated [%d] = %v > [%d] = %v", i, h.values[i].Value, j1, h.values[j2].Value)
			return
		}
		h.verify(t, j2)
	}
}

func TestInit0(t *testing.T) {
	s := []int{}
	for i := 20; i > 0; i-- {
		s = append(s, 0) // all elements are the same
	}

	h := Heap[int]{}
	h.Init(s, intCmp)
	h.verify(t, 0)

	for i := 1; h.Len() > 0; i++ {
		e := h.Pop()
		h.verify(t, 0)
		if e.Value != 0 {
			t.Errorf("%d.th pop got %d; want %d", i, e.Value, 0)
		}
	}
}

func TestInit1(t *testing.T) {
	s := []int{}
	for i := 20; i > 0; i-- {
		s = append(s, i) // all elements are different
	}

	var h Heap[int]
	h.Init(s, intCmp)
	h.verify(t, 0)

	for i := 1; h.Len() > 0; i++ {
		e := h.Pop()
		h.verify(t, 0)
		if e.Value != i {
			t.Errorf("%d.th pop got %d; want %d", i, e.Value, i)
		}
	}
}

func Test(t *testing.T) {
	s := []int{}

	for i := 20; i > 10; i-- {
		s = append(s, i)
	}
	var h Heap[int]
	h.Init(s, intCmp)
	h.verify(t, 0)

	for i := 10; i > 0; i-- {
		_ = h.Push(i)
		h.verify(t, 0)
	}

	for i := 1; h.Len() > 0; i++ {
		e := h.Pop()
		if i < 20 {
			h.Push(20 + i)
		}
		h.verify(t, 0)
		if e.Value != i {
			t.Errorf("%d.th pop got %d; want %d", i, e.Value, i)
		}
	}
}

func TestRemove0(t *testing.T) {
	h := New(10, intCmp)
	for i := 0; i < 10; i++ {
		h.Push(i)
	}
	h.verify(t, 0)

	for h.Len() > 0 {
		i := h.Len() - 1
		e := h.values[i]
		h.Remove(e)
		if e.Value != i {
			t.Errorf("Remove(%d) got %d; want %d", i, e.Value, i)
		}
		h.verify(t, 0)
	}
}

func TestRemove1(t *testing.T) {
	h := New(10, intCmp)
	for i := 0; i < 10; i++ {
		h.Push(i)
	}
	h.verify(t, 0)

	for i := 0; h.Len() > 0; i++ {
		e := h.values[0]
		h.Remove(e)
		if e.Value != i {
			t.Errorf("Remove(0) got %d; want %d", e.Value, i)
		}
		h.verify(t, 0)
	}
}

func TestRemove2(t *testing.T) {
	N := 10

	h := New(N, intCmp)
	for i := 0; i < N; i++ {
		h.Push(i)
	}
	h.verify(t, 0)

	m := make(map[int]bool)
	for h.Len() > 0 {
		e := h.values[(h.Len()-1)/2]
		h.Remove(e)
		m[e.Value] = true
		h.verify(t, 0)
	}

	if len(m) != N {
		t.Errorf("len(m) = %d; want %d", len(m), N)
	}
	for i := 0; i < len(m); i++ {
		if !m[i] {
			t.Errorf("m[%d] doesn't exist", i)
		}
	}
}

func BenchmarkDup(b *testing.B) {
	const n = 10000
	h := New(n, func(i int, j int) bool {
		return i < j
	})
	for i := 0; i < b.N; i++ {
		for j := 0; j < n; j++ {
			h.Push(0) // all elements are the same
		}
		for h.Len() > 0 {
			h.Pop()
		}
	}
}

func TestFix(t *testing.T) {
	h := New(20, intCmp)
	h.verify(t, 0)

	for i := 200; i > 0; i -= 10 {
		h.Push(i)
	}
	h.verify(t, 0)

	if h.values[0].Value != 10 {
		t.Fatalf("Expected head to be 10, was %d", h.values[0].Value)
	}
	h.values[0].Value = 210
	h.Fix(h.values[0])
	h.verify(t, 0)

	for i := 100; i > 0; i-- {
		elem := rand.Intn(h.Len())
		if i&1 == 0 {
			h.values[elem].Value *= 2
		} else {
			h.values[elem].Value /= 2
		}
		h.Fix(h.values[elem])
		h.verify(t, 0)
	}
}

type myHeap []int

func (h *myHeap) Less(i, j int) bool {
	return (*h)[i] < (*h)[j]
}

func (h *myHeap) Swap(i, j int) {
	(*h)[i], (*h)[j] = (*h)[j], (*h)[i]
}

func (h *myHeap) Len() int {
	return len(*h)
}

func (h *myHeap) Pop() (v any) {
	*h, v = (*h)[:h.Len()-1], (*h)[h.Len()-1]
	return
}

func (h *myHeap) Push(v any) {
	*h = append(*h, v.(int))
}

func BenchmarkHeap_PushPop(b *testing.B) {
	b.Run("std.Heap", func(b *testing.B) {
		n := 10000
		h := make(myHeap, 0, n)

		b.ResetTimer()
		b.ReportAllocs()

		for i := 0; i < b.N; i++ {
			for j := 0; j < n; j++ {
				heap.Push(&h, j&1000)
			}

			for h.Len() > 0 {
				heap.Pop(&h)
			}
		}
	})

	b.Run("heapz.Heap", func(b *testing.B) {
		n := 10000
		h := New[int](n, intCmp)

		b.ResetTimer()
		b.ReportAllocs()

		for i := 0; i < b.N; i++ {
			for j := 0; j < n; j++ {
				h.Push(j & 1000)
			}

			for h.Len() > 0 {
				h.Pop()
			}
		}
	})

	b.Run("heapz.Slice", func(b *testing.B) {
		n := 10000
		s := make([]int, 0, n)
		h := FromSlice(s, intCmp)

		b.ResetTimer()
		b.ReportAllocs()

		for i := 0; i < b.N; i++ {
			for j := 0; j < n; j++ {
				h.Push(j & 1000)
			}

			for len(h.Values) > 0 {
				h.Pop()
			}
		}
	})
}
