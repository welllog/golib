package listz

import (
	"reflect"
	"testing"
)

// checkSliceDList checks the content of a SliceDList against a slice of elements.
func checkSliceDList[T comparable](t *testing.T, l *SliceDList[T], es []T) {
	t.Helper()
	if l.Len() != len(es) {
		t.Errorf("l.Len() = %d, want %d", l.Len(), len(es))
		return
	}

	var result []T
	l.Range(func(idx int, value T) bool {
		result = append(result, value)
		return true
	})

	if !reflect.DeepEqual(result, es) {
		t.Errorf("list content = %v, want %v", result, es)
	}
}

func TestSliceDList_Init(t *testing.T) {
	var l SliceDList[int]
	l.Init(10)
	if l.Cap() != 10 {
		t.Errorf("Cap() = %d, want 10", l.Cap())
	}
	if l.Len() != 0 {
		t.Errorf("Len() = %d, want 0", l.Len())
	}
	if l.head != nullIdx || l.tail != nullIdx {
		t.Errorf("head/tail should be nullIdx on empty list")
	}
}

func TestSliceDList_Push(t *testing.T) {
	var l SliceDList[int]
	l.Init(3)

	// PushBack
	idx1, ok1 := l.PushBack(10)
	if !ok1 || idx1 != 0 {
		t.Errorf("PushBack(10) failed, got idx %d, ok %v", idx1, ok1)
	}
	checkSliceDList(t, &l, []int{10})

	// PushFront
	idx2, ok2 := l.PushFront(20)
	if !ok2 || idx2 != 1 {
		t.Errorf("PushFront(20) failed, got idx %d, ok %v", idx2, ok2)
	}
	checkSliceDList(t, &l, []int{20, 10})

	// PushBack again
	idx3, ok3 := l.PushBack(30)
	if !ok3 || idx3 != 2 {
		t.Errorf("PushBack(30) failed, got idx %d, ok %v", idx3, ok3)
	}
	checkSliceDList(t, &l, []int{20, 10, 30})

	// List is full
	if l.HasFree() {
		t.Error("HasFree() should be false when list is full")
	}
	idx4, ok4 := l.PushBack(40)
	if !ok4 || idx4 != 3 {
		t.Errorf("PushBack(40) should fail, got idx %d, ok %v", idx4, ok4)
	}
	checkSliceDList(t, &l, []int{20, 10, 30, 40})
}

func TestSliceDList_Remove(t *testing.T) {
	var l SliceDList[int]
	l.Init(5)
	idx1, _ := l.PushBack(10)
	idx2, _ := l.PushBack(20)
	idx3, _ := l.PushBack(30)

	// Remove middle
	if !l.Remove(idx2) {
		t.Error("Remove(idx2) failed")
	}
	checkSliceDList(t, &l, []int{10, 30})

	// Remove head
	if !l.Remove(idx1) {
		t.Error("Remove(idx1) failed")
	}
	checkSliceDList(t, &l, []int{30})

	// Remove tail (the only element)
	if !l.Remove(idx3) {
		t.Error("Remove(idx3) failed")
	}
	checkSliceDList(t, &l, nil)

	// Remove from empty list
	if l.Remove(0) {
		t.Error("Remove on empty list should fail")
	}
}

func TestSliceDList_Access(t *testing.T) {
	var l SliceDList[int]
	l.Init(5)
	idx1, _ := l.PushBack(10)
	_, _ = l.PushBack(20)
	idx3, _ := l.PushBack(30)

	// Front
	fIdx, fVal, fOk := l.Front()
	if !fOk || fIdx != idx1 || fVal != 10 {
		t.Errorf("Front() = %d, %d, %v; want %d, 10, true", fIdx, fVal, fOk, idx1)
	}

	// Back
	bIdx, bVal, bOk := l.Back()
	if !bOk || bIdx != idx3 || bVal != 30 {
		t.Errorf("Back() = %d, %d, %v; want %d, 30, true", bIdx, bVal, bOk, idx3)
	}

	// Get
	gVal, gOk := l.Get(idx1)
	if !gOk || gVal != 10 {
		t.Errorf("Get(%d) = %d, %v; want 10, true", idx1, gVal, gOk)
	}

	// Get invalid
	_, gOkInvalid := l.Get(99)
	if gOkInvalid {
		t.Error("Get(99) on invalid index should fail")
	}
}

func TestSliceDList_Insert(t *testing.T) {
	var l SliceDList[int]
	l.Init(5)
	idx1, _ := l.PushBack(10)
	idx2, _ := l.PushBack(30)

	// InsertAfter
	idx3, ok := l.InsertAfter(20, idx1)
	if !ok {
		t.Error("InsertAfter failed")
	}
	checkSliceDList(t, &l, []int{10, 20, 30})

	// InsertBefore
	_, ok = l.InsertBefore(15, idx3)
	if !ok {
		t.Error("InsertBefore failed")
	}
	checkSliceDList(t, &l, []int{10, 15, 20, 30})

	// InsertAfter tail
	_, ok = l.InsertAfter(40, idx2)
	if !ok {
		t.Error("InsertAfter tail failed")
	}
	checkSliceDList(t, &l, []int{10, 15, 20, 30, 40})
}

func TestSliceDList_Move(t *testing.T) {
	var l SliceDList[int]
	l.Init(5)
	idx1, _ := l.PushBack(10)
	idx2, _ := l.PushBack(20)
	idx3, _ := l.PushBack(30)
	idx4, _ := l.PushBack(40)

	// MoveToFront
	l.MoveToFront(idx3) // 30, 10, 20, 40
	checkSliceDList(t, &l, []int{30, 10, 20, 40})

	// MoveToBack
	l.MoveToBack(idx1) // 30, 20, 40, 10
	checkSliceDList(t, &l, []int{30, 20, 40, 10})

	// MoveAfter
	l.MoveAfter(idx3, idx4) // 20, 40, 30, 10
	checkSliceDList(t, &l, []int{20, 40, 30, 10})

	// MoveBefore
	l.MoveBefore(idx1, idx2) // 10, 20, 40, 30
	checkSliceDList(t, &l, []int{10, 20, 40, 30})
}

func TestSliceDList_Range(t *testing.T) {
	var l SliceDList[int]
	l.Init(4)
	_, _ = l.PushBack(10)
	idx2, _ := l.PushBack(20)
	_, _ = l.PushBack(30)

	t.Run("Range full", func(t *testing.T) {
		var result []int
		l.Range(func(idx int, value int) bool {
			result = append(result, value)
			return true
		})
		if !reflect.DeepEqual(result, []int{10, 20, 30}) {
			t.Errorf("Range result = %v, want [10 20 30]", result)
		}
	})

	t.Run("Range stopped", func(t *testing.T) {
		var result []int
		l.Range(func(idx int, value int) bool {
			result = append(result, value)
			return value != 20
		})
		if !reflect.DeepEqual(result, []int{10, 20}) {
			t.Errorf("Range stopped result = %v, want [10 20]", result)
		}
	})

	t.Run("RangeFrom", func(t *testing.T) {
		var result []int
		l.RangeFrom(idx2, func(idx int, value int) bool {
			result = append(result, value)
			return true
		})
		if !reflect.DeepEqual(result, []int{20, 30}) {
			t.Errorf("RangeFrom result = %v, want [20 30]", result)
		}
	})

	t.Run("RangeFrom invalid", func(t *testing.T) {
		called := false
		l.RangeFrom(99, func(idx int, value int) bool {
			called = true
			return true
		})
		if called {
			t.Error("RangeFrom with invalid index should not call the function")
		}
	})

	t.Run("RevRangeFrom", func(t *testing.T) {
		var result []int
		l.RevRangeFrom(idx2, func(idx int, value int) bool {
			result = append(result, value)
			return true
		})
		if !reflect.DeepEqual(result, []int{20, 10}) {
			t.Errorf("RevRangeFrom result = %v, want [20 10]", result)
		}
	})

	t.Run("RevRangeFrom invalid", func(t *testing.T) {
		called := false
		l.RevRangeFrom(99, func(idx int, value int) bool {
			called = true
			return true
		})
		if called {
			t.Error("RevRangeFrom with invalid index should not call the function")
		}
	})
}

func TestSliceDList_RevRange(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		var l SliceDList[int]
		l.Init(10)

		called := false
		l.RevRange(func(idx int, value int) bool {
			called = true
			return true
		})

		if called {
			t.Error("RevRange on empty list should not call the function")
		}
	})

	t.Run("full iteration", func(t *testing.T) {
		var l SliceDList[int]
		l.Init(10)
		l.PushBack(10)
		l.PushBack(20)
		l.PushBack(30)

		var result []int
		l.RevRange(func(idx int, value int) bool {
			result = append(result, value)
			return true
		})

		expected := []int{30, 20, 10}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("RevRange result = %v, want %v", result, expected)
		}
	})

	t.Run("stopped iteration", func(t *testing.T) {
		var l SliceDList[int]
		l.Init(10)
		l.PushBack(10)
		l.PushBack(20)
		l.PushBack(30)
		l.PushBack(40)

		var result []int
		l.RevRange(func(idx int, value int) bool {
			result = append(result, value)
			return value != 20
		})

		expected := []int{40, 30, 20}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("RevRange result = %v, want %v", result, expected)
		}
	})
}
