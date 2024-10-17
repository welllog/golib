// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package heapz

import (
	"math/rand"
	"testing"
)

type myIntHeap []int

func (h *myIntHeap) Less(i, j int) bool {
	return (*h)[i] < (*h)[j]
}

func (h *myIntHeap) Swap(i, j int) {
	(*h)[i], (*h)[j] = (*h)[j], (*h)[i]
}

func (h *myIntHeap) Len() int {
	return len(*h)
}

func (h *myIntHeap) Pop() (v int) {
	*h, v = (*h)[:h.Len()-1], (*h)[h.Len()-1]
	return
}

func (h *myIntHeap) Push(v int) {
	*h = append(*h, v)
}

func (h myIntHeap) verify(t *testing.T, i int) {
	t.Helper()
	n := h.Len()
	j1 := 2*i + 1
	j2 := 2*i + 2
	if j1 < n {
		if h.Less(j1, i) {
			t.Errorf("heap invariant invalidated [%d] = %d > [%d] = %d", i, h[i], j1, h[j1])
			return
		}
		h.verify(t, j1)
	}
	if j2 < n {
		if h.Less(j2, i) {
			t.Errorf("heap invariant invalidated [%d] = %d > [%d] = %d", i, h[i], j1, h[j2])
			return
		}
		h.verify(t, j2)
	}
}

func TestStdInit0(t *testing.T) {
	h := new(myIntHeap)
	for i := 20; i > 0; i-- {
		h.Push(0) // all elements are the same
	}
	Init[int](h)
	h.verify(t, 0)

	for i := 1; h.Len() > 0; i++ {
		x := Pop[int](h)
		h.verify(t, 0)
		if x != 0 {
			t.Errorf("%d.th pop got %d; want %d", i, x, 0)
		}
	}
}

func TestStdInit1(t *testing.T) {
	h := new(myIntHeap)
	for i := 20; i > 0; i-- {
		h.Push(i) // all elements are different
	}
	Init[int](h)
	h.verify(t, 0)

	for i := 1; h.Len() > 0; i++ {
		x := Pop[int](h)
		h.verify(t, 0)
		if x != i {
			t.Errorf("%d.th pop got %d; want %d", i, x, i)
		}
	}
}

func TestStd(t *testing.T) {
	h := new(myIntHeap)
	h.verify(t, 0)

	for i := 20; i > 10; i-- {
		h.Push(i)
	}
	Init[int](h)
	h.verify(t, 0)

	for i := 10; i > 0; i-- {
		Push[int](h, i)
		h.verify(t, 0)
	}

	for i := 1; h.Len() > 0; i++ {
		x := Pop[int](h)
		if i < 20 {
			Push[int](h, 20+i)
		}
		h.verify(t, 0)
		if x != i {
			t.Errorf("%d.th pop got %d; want %d", i, x, i)
		}
	}
}

func TestStdRemove0(t *testing.T) {
	h := new(myIntHeap)
	for i := 0; i < 10; i++ {
		h.Push(i)
	}
	h.verify(t, 0)

	for h.Len() > 0 {
		i := h.Len() - 1
		x := Remove[int](h, i)
		if x != i {
			t.Errorf("Remove(%d) got %d; want %d", i, x, i)
		}
		h.verify(t, 0)
	}
}

func TestStdRemove1(t *testing.T) {
	h := new(myIntHeap)
	for i := 0; i < 10; i++ {
		h.Push(i)
	}
	h.verify(t, 0)

	for i := 0; h.Len() > 0; i++ {
		x := Remove[int](h, 0).(int)
		if x != i {
			t.Errorf("Remove(0) got %d; want %d", x, i)
		}
		h.verify(t, 0)
	}
}

func TestStdRemove2(t *testing.T) {
	N := 10

	h := new(myIntHeap)
	for i := 0; i < N; i++ {
		h.Push(i)
	}
	h.verify(t, 0)

	m := make(map[int]bool)
	for h.Len() > 0 {
		m[Remove[int](h, (h.Len()-1)/2).(int)] = true
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

func TestStdFix(t *testing.T) {
	h := new(myIntHeap)
	h.verify(t, 0)

	for i := 200; i > 0; i -= 10 {
		Push[int](h, i)
	}
	h.verify(t, 0)

	if (*h)[0] != 10 {
		t.Fatalf("Expected head to be 10, was %d", (*h)[0])
	}
	(*h)[0] = 210
	Fix[int](h, 0)
	h.verify(t, 0)

	for i := 100; i > 0; i-- {
		elem := rand.Intn(h.Len())
		if i&1 == 0 {
			(*h)[elem] *= 2
		} else {
			(*h)[elem] /= 2
		}
		Fix[int](h, elem)
		h.verify(t, 0)
	}
}
