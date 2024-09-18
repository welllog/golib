package heapz

import (
	"math/rand"
	"testing"
)

func intCmp(a, b int) bool {
	return a < b
}

func TestFromSlice(t *testing.T) {
	s := FromSlice([]int{3, 1, 2}, intCmp)
	if len(s.Values) != 3 {
		t.Errorf("expected length 3, got %d", len(s.Values))
	}
	if n, _ := s.Peek(); n != 1 {
		t.Errorf("expected root 1, got %d", s.Values[0])
	}
}

func TestPush(t *testing.T) {
	s := FromSlice([]int{3, 1, 2}, intCmp)
	s.Push(0)
	if len(s.Values) != 4 {
		t.Errorf("expected length 4, got %d", len(s.Values))
	}
	if n, _ := s.Peek(); n != 0 {
		t.Errorf("expected root 0, got %d", s.Values[0])
	}
}

func TestPop(t *testing.T) {
	s := FromSlice([]int{3, 1, 2}, intCmp)
	val, ok := s.Pop()
	if !ok || val != 1 {
		t.Errorf("expected 1, got %d", val)
	}
	if len(s.Values) != 2 {
		t.Errorf("expected length 2, got %d", len(s.Values))
	}
}

func TestPeek(t *testing.T) {
	s := FromSlice([]int{3, 1, 2}, intCmp)
	val, ok := s.Peek()
	if !ok || val != 1 {
		t.Errorf("expected 1, got %d", val)
	}
}

func TestRemove(t *testing.T) {
	s := FromSlice([]int{3, 1, 2}, intCmp)
	rv := s.Values[1]
	val, ok := s.Remove(1)
	if !ok || val != rv {
		t.Errorf("expected 1, got %d", val)
	}
	if len(s.Values) != 2 {
		t.Errorf("expected length 2, got %d", len(s.Values))
	}
}

func TestSliceFix(t *testing.T) {
	s := FromSlice([]int{3, 1, 2}, intCmp)
	s.Values[0] = 4
	s.Fix(0)
	if s.Values[0] != 2 {
		t.Errorf("expected root 2, got %d", s.Values[0])
	}
}

func verifyIntSlice(t *testing.T, s []int, i int) {
	t.Helper()
	n := len(s)
	j1 := 2*i + 1
	j2 := 2*i + 2
	if j1 < n {
		if s[j1] < s[i] {
			t.Errorf("heap invariant invalidated [%d] = %d > [%d] = %d", i, s[i], j1, s[j1])
			return
		}
		verifyIntSlice(t, s, j1)
	}
	if j2 < n {
		if s[j2] < s[i] {
			t.Errorf("heap invariant invalidated [%d] = %d > [%d] = %d", i, s[i], j1, s[j2])
			return
		}
		verifyIntSlice(t, s, j2)
	}
}

func TestFromSlice2(t *testing.T) {
	var s []int
	for i := 20; i > 0; i-- {
		s = append(s, 0)
	}

	h := FromSlice(s, intCmp)
	verifyIntSlice(t, h.Values, 0)

	for i := 1; h.Len() > 0; i++ {
		x, _ := h.Pop()
		verifyIntSlice(t, h.Values, 0)
		if x != 0 {
			t.Errorf("%d.th pop got %d; want %d", i, x, 0)
		}
	}
}

func TestFromSlice3(t *testing.T) {
	var s []int
	for i := 20; i > 0; i-- {
		s = append(s, i)
	}
	h := FromSlice(s, intCmp)
	verifyIntSlice(t, h.Values, 0)

	for i := 1; h.Len() > 0; i++ {
		x, _ := h.Pop()
		verifyIntSlice(t, h.Values, 0)
		if x != i {
			t.Errorf("%d.th pop got %d; want %d", i, x, i)
		}
	}
}

func TestSlice(t *testing.T) {
	var s []int
	verifyIntSlice(t, s, 0)

	for i := 20; i > 10; i-- {
		s = append(s, i)
	}
	h := FromSlice(s, intCmp)
	verifyIntSlice(t, h.Values, 0)

	for i := 10; i > 0; i-- {
		h.Push(i)
		verifyIntSlice(t, h.Values, 0)
	}

	for i := 1; h.Len() > 0; i++ {
		x, _ := h.Pop()
		if i < 20 {
			h.Push(20 + i)
		}
		verifyIntSlice(t, h.Values, 0)
		if x != i {
			t.Errorf("%d.th pop got %d; want %d", i, x, i)
		}
	}
}

func TestSliceRemove0(t *testing.T) {
	h := NewSlice(10, intCmp)
	for i := 0; i < 10; i++ {
		h.Push(i)
	}
	verifyIntSlice(t, h.Values, 0)

	for h.Len() > 0 {
		i := h.Len() - 1
		x, _ := h.Remove(i)
		if x != i {
			t.Errorf("Remove(%d) got %d; want %d", i, x, i)
		}
		verifyIntSlice(t, h.Values, 0)
	}
}

func TestSliceRemove1(t *testing.T) {
	h := NewSlice(10, intCmp)
	for i := 0; i < 10; i++ {
		h.Push(i)
	}
	verifyIntSlice(t, h.Values, 0)

	for i := 0; h.Len() > 0; i++ {
		x, _ := h.Remove(0)
		if x != i {
			t.Errorf("Remove(0) got %d; want %d", x, i)
		}
		verifyIntSlice(t, h.Values, 0)
	}
}

func TestSliceRemove2(t *testing.T) {
	N := 10

	h := NewSlice(N, intCmp)
	for i := 0; i < N; i++ {
		h.Push(i)
	}
	verifyIntSlice(t, h.Values, 0)

	m := make(map[int]bool)
	for h.Len() > 0 {
		n, _ := h.Remove((h.Len() - 1) / 2)
		m[n] = true
		verifyIntSlice(t, h.Values, 0)
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

func TestSliceFix1(t *testing.T) {
	h := NewSlice(10, intCmp)
	h.Remove((h.Len() - 1) / 2)

	for i := 200; i > 0; i -= 10 {
		h.Push(i)
	}
	h.Remove((h.Len() - 1) / 2)

	if h.Values[0] != 10 {
		t.Fatalf("Expected head to be 10, was %d", h.Values[0])
	}
	h.Values[0] = 210
	h.Fix(0)
	verifyIntSlice(t, h.Values, 0)

	for i := 100; i > 0; i-- {
		elem := rand.Intn(h.Len())
		if i&1 == 0 {
			h.Values[elem] *= 2
		} else {
			h.Values[elem] /= 2
		}
		h.Fix(elem)
		verifyIntSlice(t, h.Values, 0)
	}
}
