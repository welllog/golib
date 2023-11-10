package dsz

import (
	"testing"
)

func TestSet_Add(t *testing.T) {
	s := make(Set[int])
	if !s.Add(1) {
		t.Fatal("add 1 failed")
	}
	if s.Add(1) {
		t.Fatal("add 1 again")
	}
	if !s.Add(2) {
		t.Fatal("add 2 failed")
	}
}

func TestSet_MultiAdd(t *testing.T) {
	s := make(Set[int])
	s.MultiAdd(1, 2, 3)
	if !s.Has(1) {
		t.Fatal("add 1 failed")
	}
	if !s.Has(2) {
		t.Fatal("add 2 failed")
	}
	if !s.Has(3) {
		t.Fatal("add 3 failed")
	}
	if s.Has(4) {
		t.Fatal("not has 4")
	}
}

func TestSet_Delete(t *testing.T) {
	s := make(Set[int])
	s.MultiAdd(1, 2, 3)
	if !s.Delete(1) {
		t.Fatal("delete 1 failed")
	}
	if s.Delete(1) {
		t.Fatal("delete 1 again")
	}
	if !s.Delete(2) {
		t.Fatal("delete 2 failed")
	}
	if !s.Delete(3) {
		t.Fatal("delete 3 failed")
	}
	if s.Delete(4) {
		t.Fatal("delete 4")
	}
}

func TestSet_Values(t *testing.T) {
	s := make(Set[int])
	s.MultiAdd(1, 2, 3)
	var vs []int
	vs = s.Values(vs)
	if len(vs) != 3 {
		t.Fatal("values failed")
	}
}

func TestSet_Diff(t *testing.T) {
	s := make(Set[int])
	s.MultiAdd(0, 1, 2, 3)
	s1 := make(Set[int])
	s1.MultiAdd(1, 2, 3, 4)

	s.Diff(s1)
	if len(s) != 1 {
		t.Fatal("diff failed")
	}
}

func TestSet_Intersect(t *testing.T) {
	s := make(Set[int])
	s.MultiAdd(0, 1, 2, 3)
	s1 := make(Set[int])
	s1.MultiAdd(1, 2, 3, 4)

	s.Intersect(s1)
	if len(s) != 3 {
		t.Fatal("intersect failed")
	}
}

func TestSet_DiffWithSlice(t *testing.T) {
	s := make(Set[int])
	s.MultiAdd(0, 1, 2, 3)
	s1 := []int{1, 2, 3, 4}

	s.DiffWithSlice(s1)
	if len(s) != 1 {
		t.Fatal("diff failed")
	}
}

func TestSet_IntersectWithSlice(t *testing.T) {
	s := make(Set[int])
	s.MultiAdd(0, 1, 2, 3)
	s1 := []int{1, 2, 3, 4}

	s.IntersectWithSlice(s1)
	if len(s) != 3 {
		t.Fatal("intersect failed")
	}
}

func TestSet_Filter(t *testing.T) {
	s := make(Set[int])
	s.MultiAdd(1, 2, 3)

	s.Filter(func(v int) bool {
		return v > 1
	})
	if len(s) != 2 {
		t.Fatal("filter failed")
	}
}
