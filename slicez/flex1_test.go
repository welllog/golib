package slicez

import (
	"github.com/welllog/golib/testz"
	"testing"
)

func TestFlexSlice1_Basic(t *testing.T) {
	var f FlexSlice1[int]
	f.Append(0, 1, 2, 3, 4)
	f.Append(5, 6, 7, 8, 9, 10)
	f.Prepend(-4, -3, -2, -1)

	start := -4
	end := 10
	for i := 0; start+i <= end; i++ {
		if v, ok := f.Get(i); !ok || v != start+i {
			t.Errorf("Append/Prepend failed at index %d: expected %d, got %d", i, start+i, v)
		}
	}

	size := end - start + 1
	if f.Len() != end-start+1 {
		t.Errorf("Append/Prepend failed: expected length %d, got %d", end-start+1, f.Len())
	}
	capacity := f.Cap()

	for i := 0; i < 4; i++ {
		n, ok := f.Pop()
		if !ok || n != end-i {
			t.Errorf("Pop failed: expected %d, got %d", end-i, n)
		}

		n, ok = f.Shift()
		if !ok || n != start+i {
			t.Errorf("Shift failed: expected %d, got %d", start+i, n)
		}
	}

	if f.Len() != size-8 {
		t.Errorf("Pop/Shift failed: expected length %d, got %d", size-8, f.Len())
	}

	if f.Cap() != capacity {
		t.Errorf("Pop/Shift failed: expected capacity %d, got %d", capacity, f.Cap())
	}

	start = 0
	end = 6
	arr := []int{}
	for i := 0; start+i <= end; i++ {
		if v, ok := f.Get(i); !ok || v != start+i {
			t.Errorf("Get failed at index %d: expected %d, got %d", i, start+i, v)
		}
		n, _ := f.Get(i)
		arr = append(arr, n)
	}

	testz.Equal(t, arr, f.ToSlice())

	n, ok := f.Remove(f.Len() - 1)
	if !ok || n != end {
		t.Errorf("Remove failed: expected %d, got %d", end, n)
	}

	n, ok = f.Remove(0)
	if !ok || n != start {
		t.Errorf("Remove failed: expected %d, got %d", start, n)
	}

	n, ok = f.Remove(f.Len())
	if ok {
		t.Error("Remove should have failed for out-of-bounds index")
	}

	// 1 2 3 4 5
	n, ok = f.Remove(1)
	if !ok || n != 2 {
		t.Errorf("Remove failed: expected %d, got %d", 2, n)
	}

	if f.Len() != 4 {
		t.Errorf("Remove failed: expected length %d, got %d", 4, f.Len())
	}

	arr = []int{1, 3, 4, 5}
	for i, v := range arr {
		if val, ok := f.Get(i); !ok || val != v {
			t.Errorf("Get failed at index %d: expected %d, got %d", i, v, val)
		}
	}

	f.InsertAt(2, 99)
	if f.Len() != 5 {
		t.Errorf("InsertAt failed: expected length %d, got %d", 5, f.Len())
	}

	arr = []int{1, 3, 99, 4, 5}
	for i, v := range arr {
		if val, ok := f.Get(i); !ok || val != v {
			t.Errorf("Get failed at index %d: expected %d, got %d", i, v, val)
		}
	}

	f.InsertAt(4, 10)
	if f.Len() != 6 {
		t.Errorf("InsertAt failed: expected length %d, got %d", 6, f.Len())
	}
	arr = []int{1, 3, 99, 4, 10, 5}
	for i, v := range arr {
		if val, ok := f.Get(i); !ok || val != v {
			t.Errorf("Get failed at index %d: expected %d, got %d", i, v, val)
		}
	}

	f.InsertAt(1, 7)
	arr = []int{1, 7, 3, 99, 4, 10, 5}
	for i, v := range arr {
		if val, ok := f.Get(i); !ok || val != v {
			t.Errorf("Get failed at index %d: expected %d, got %d", i, v, val)
		}
	}

	f.InsertAt(2, 12)
	arr = []int{1, 7, 12, 3, 99, 4, 10, 5}
	for i, v := range arr {
		if val, ok := f.Get(i); !ok || val != v {
			t.Errorf("Get failed at index %d: expected %d, got %d", i, v, val)
		}
	}
}
