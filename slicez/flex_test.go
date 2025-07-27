package slicez

import (
	"fmt"
	"testing"

	"github.com/welllog/golib/testz"
)

func TestFlexSlice_Basic(t *testing.T) {
	var f FlexSlice[int]
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

	// 1, 3, 99, 4, 5
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

	// 1, 3, 99, 4, 10, 5
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

	// 1, 7, 3, 99, 4, 10, 5
	f.InsertAt(1, 7)
	arr = []int{1, 7, 3, 99, 4, 10, 5}
	testz.Equal(t, arr, f.ToSlice())

	f.InsertAt(2, 12)
	arr = []int{1, 7, 12, 3, 99, 4, 10, 5}
	testz.Equal(t, arr, f.ToSlice())
}

func TestFlexSlice_AppendPrependPopShift(t *testing.T) {
	t.Run("basic operations", func(t *testing.T) {
		var f FlexSlice[int]
		f.Append(1, 2, 3)
		testz.Equal(t, []int{1, 2, 3}, f.ToSlice(), "Append failed")
		testz.Equal(t, 3, f.Len(), "Length after Append is wrong")

		f.Prepend(0)
		testz.Equal(t, []int{0, 1, 2, 3}, f.ToSlice(), "Prepend failed")
		testz.Equal(t, 4, f.Len(), "Length after Prepend is wrong")

		v, ok := f.Pop()
		testz.Assert(t, ok, "Pop should succeed")
		testz.Equal(t, 3, v, "Pop returned wrong value")
		testz.Equal(t, []int{0, 1, 2}, f.ToSlice(), "Slice after Pop is wrong")
		testz.Equal(t, 3, f.Len(), "Length after Pop is wrong")

		v, ok = f.Shift()
		testz.Assert(t, ok, "Shift should succeed")
		testz.Equal(t, 0, v, "Shift returned wrong value")
		testz.Equal(t, []int{1, 2}, f.ToSlice(), "Slice after Shift is wrong")
		testz.Equal(t, 2, f.Len(), "Length after Shift is wrong")
	})

	t.Run("operations on empty slice", func(t *testing.T) {
		var f FlexSlice[int]
		v, ok := f.Pop()
		testz.Assert(t, !ok, "Pop on empty slice should fail")
		testz.Equal(t, 0, v, "Pop on empty slice should return zero value")

		v, ok = f.Shift()
		testz.Assert(t, !ok, "Shift on empty slice should fail")
		testz.Equal(t, 0, v, "Shift on empty slice should return zero value")

		f.Append() // Append nothing
		testz.Equal(t, 0, f.Len(), "Append nothing should not change length")

		f.Prepend() // Prepend nothing
		testz.Equal(t, 0, f.Len(), "Prepend nothing should not change length")
	})

	t.Run("operations causing growth", func(t *testing.T) {
		f := NewFlexSlice[int](2)
		f.Append(1, 2)
		testz.Equal(t, 2, f.Cap(), "Initial capacity should be 2")

		f.Append(3) // Should grow
		testz.Equal(t, []int{1, 2, 3}, f.ToSlice(), "Slice content is wrong after growing append")
		testz.Assert(t, f.Cap() > 2, "Capacity should grow after append")

		f.Prepend(0) // Should grow again
		testz.Equal(t, []int{0, 1, 2, 3}, f.ToSlice(), "Slice content is wrong after growing prepend")
		testz.Assert(t, f.Cap() > 3, "Capacity should grow after prepend")
	})

	t.Run("operations to empty the slice", func(t *testing.T) {
		var f FlexSlice[int]
		f.Append(1, 2)

		v, ok := f.Pop()
		testz.Assert(t, ok, "Pop should succeed")
		testz.Equal(t, 2, v, "Pop value is wrong")

		v, ok = f.Shift()
		testz.Assert(t, ok, "Shift should succeed")
		testz.Equal(t, 1, v, "Shift value is wrong")

		testz.Assert(t, f.IsEmpty(), "Slice should be empty")
		testz.Equal(t, 0, f.Len(), "Length should be 0")
	})

	t.Run("interleaved operations with wrap around", func(t *testing.T) {
		f := NewFlexSlice[int](5)
		f.Append(1, 2, 3) // [1, 2, 3, _, _], head=0, tail=3
		f.Shift()         // [_, 2, 3, _, _], head=1, tail=3
		f.Shift()         // [_, _, 3, _, _], head=2, tail=3
		f.Append(4, 5)    // [_, _, 3, 4, 5], head=2, tail=0
		f.Append(6)       // Should grow and re-order
		testz.Equal(t, []int{3, 4, 5, 6}, f.ToSlice(), "Slice content is wrong after wrap around and grow")

		v, ok := f.Pop()
		testz.Assert(t, ok)
		testz.Equal(t, 6, v)

		v, ok = f.Shift()
		testz.Assert(t, ok)
		testz.Equal(t, 3, v)

		testz.Equal(t, []int{4, 5}, f.ToSlice(), "Slice content is wrong after final pop/shift")
	})
}

func TestFlexSlice_GetSet(t *testing.T) {
	tests := []struct {
		name        string
		initialData []int
		initialCap  int
		head        int
	}{
		{
			name:        "standard",
			initialData: []int{10, 20, 30, 40, 50},
			initialCap:  10,
			head:        0,
		},
		{
			name:        "with head wrap",
			initialData: []int{10, 20, 30, 40, 50},
			initialCap:  7,
			head:        5,
		},
		{
			name:        "full slice",
			initialData: []int{10, 20, 30},
			initialCap:  3,
			head:        0,
		},
		{
			name:        "single element",
			initialData: []int{100},
			initialCap:  5,
			head:        3,
		},
		{
			name:        "empty slice",
			initialData: []int{},
			initialCap:  5,
			head:        0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test Get
			for i := -1; i <= len(tt.initialData); i++ {
				f := makeFlexSlice(tt.initialData, tt.head, tt.initialCap)
				val, ok := f.Get(i)

				if i >= 0 && i < len(tt.initialData) {
					if !ok {
						t.Errorf("Get(%d) expected ok=true, got false", i)
					}
					if val != tt.initialData[i] {
						t.Errorf("Get(%d) expected value %d, got %d", i, tt.initialData[i], val)
					}
				} else {
					if ok {
						t.Errorf("Get(%d) expected ok=false, got true", i)
					}
				}
			}

			// Test Set
			for i := -1; i <= len(tt.initialData); i++ {
				f := makeFlexSlice(tt.initialData, tt.head, tt.initialCap)
				setValue := 999
				ok := f.Set(i, setValue)

				if i >= 0 && i < len(tt.initialData) {
					if !ok {
						t.Errorf("Set(%d) expected ok=true, got false", i)
					}
					expectedData := make([]int, len(tt.initialData))
					copy(expectedData, tt.initialData)
					expectedData[i] = setValue
					testz.Equal(t, expectedData, f.ToSlice(), "data mismatch after Set")
				} else {
					if ok {
						t.Errorf("Set(%d) expected ok=false, got true", i)
					}
					// Ensure data was not modified
					testz.Equal(t, tt.initialData, f.ToSlice(), "data should not be modified on failed Set")
				}
			}
		})
	}
}

func TestFlexSlice_InsertAt(t *testing.T) {
	tests := []struct {
		data   []int
		values []int
	}{
		{
			data:   []int{1, 2, 3, 4, 5},
			values: []int{10},
		},
		{
			data:   []int{1, 2, 3, 4, 5},
			values: []int{10, 20},
		},
		{
			data:   []int{1, 2, 3, 4, 5},
			values: []int{10, 20, 30},
		},
		{
			data:   []int{1, 2, 3, 4, 5},
			values: []int{10, 20, 30, 40},
		},
		{
			data:   []int{},
			values: []int{1},
		},
		{
			data:   []int{},
			values: []int{1, 2},
		},
		{
			data:   []int{},
			values: []int{1, 2, 3},
		},
		{
			data:   []int{1},
			values: []int{2},
		},
		{
			data:   []int{1},
			values: []int{2, 3},
		},
		{
			data:   []int{1},
			values: []int{2, 3, 4},
		},
		{
			data:   []int{1, 2},
			values: []int{3},
		},
		{
			data:   []int{1, 2},
			values: []int{3, 4},
		},
		{
			data:   []int{1, 2},
			values: []int{3, 4, 5},
		},
		{
			data:   []int{1, 2, 3},
			values: []int{4},
		},
		{
			data:   []int{1, 2, 3},
			values: []int{4, 5},
		},
		{
			data:   []int{1, 2, 3},
			values: []int{4, 5, 6},
		},
		{
			data:   []int{1, 2, 3, 4},
			values: []int{5},
		},
		{
			data:   []int{1, 2, 3, 4},
			values: []int{5, 6},
		},
		{
			data:   []int{1, 2, 3, 4},
			values: []int{5, 6, 7},
		},
	}

	for _, tt := range tests {
		for capacity := len(tt.data); capacity < len(tt.data)+20; capacity++ {
			for head := 0; head < capacity; head++ {
				for at := 0; at <= len(tt.data); at++ {
					expected := make([]int, 0, len(tt.data)+len(tt.values))
					expected = append(expected, tt.data[:at]...)
					expected = append(expected, tt.values...)
					expected = append(expected, tt.data[at:]...)

					f := makeFlexSlice(tt.data, head, capacity)
					//fmt.Printf("%+v \n", f)
					f.InsertAt(at, tt.values...)
					//fmt.Printf("%+v \n", f)
					testz.Equal(t, expected, f.ToSlice(), fmt.Sprintf("insertAt: %d, values: %v, head: %d, cap: %d", at, tt.values, head, capacity))
				}
			}
		}
	}
}

func TestNewFlexSlice_Remove(t *testing.T) {
	tests := []struct {
		data []int
	}{
		{data: []int{1}},
		{data: []int{1, 2}},
		{data: []int{1, 2, 3}},
		{data: []int{1, 2, 3, 4}},
		{data: []int{1, 2, 3, 4, 5}},
		{data: []int{1, 2, 3, 4, 5, 6}},
	}

	for _, tt := range tests {
		for capacity := len(tt.data); capacity < len(tt.data)+20; capacity++ {
			for head := 0; head < capacity; head++ {
				for at := 0; at < len(tt.data); at++ {
					expected := make([]int, 0, len(tt.data)-1)
					expected = append(expected, tt.data[:at]...)
					expected = append(expected, tt.data[at+1:]...)

					f := makeFlexSlice(tt.data, head, capacity)
					//fmt.Printf("%+v \n", f)
					f.Remove(at)
					//fmt.Printf("%+v \n", f)
					testz.Equal(t, expected, f.ToSlice(), fmt.Sprintf("remove: %d, head: %d, cap: %d", at, head, capacity))
				}
			}
		}
	}
}

func TestFlexSlice_Join(t *testing.T) {
	tests := []struct {
		data1 []int
		data2 []int
	}{
		{
			data1: []int{1, 2, 3},
			data2: []int{4, 5, 6},
		},
		{
			data1: []int{},
			data2: []int{1, 2, 3},
		},
		{
			data1: []int{1, 2, 3},
			data2: []int{},
		},
		{
			data1: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			data2: []int{11, 12, 13, 14, 15},
		},
		{
			data1: []int{1},
			data2: []int{2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},
		},
	}

	for _, tt := range tests {
		for cap1 := len(tt.data1); cap1 < len(tt.data1)+len(tt.data2)+5; cap1++ {
			for head1 := 0; head1 < cap1; head1++ {
				for cap2 := len(tt.data2); cap2 < len(tt.data2)+5; cap2++ {
					for head2 := 0; head2 < cap2; head2++ {
						f1 := makeFlexSlice(tt.data1, head1, cap1)
						f2 := makeFlexSlice(tt.data2, head2, cap2)

						expected := append(append([]int{}, tt.data1...), tt.data2...)

						f1.Join(*f2)

						msg := fmt.Sprintf("data1: %v, data2: %v, cap1: %d, head1: %d, cap2: %d, head2: %d",
							tt.data1, tt.data2, cap1, head1, cap2, head2)
						testz.Equal(t, expected, f1.ToSlice(), msg)
					}
				}
			}
		}
	}
}

func TestFlexSlice_SubSlice(t *testing.T) {
	tests := []struct {
		data []int
	}{
		{data: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}},
		{data: []int{1}},
		{data: []int{}},
	}

	for _, tt := range tests {
		for capacity := len(tt.data); capacity < len(tt.data)+10; capacity++ {
			for head := 0; head < capacity; head++ {
				// Test various start and end combinations, including out-of-bounds
				for start := -1; start <= len(tt.data)+1; start++ {
					for end := start - 1; end <= len(tt.data)+2; end++ {
						f := makeFlexSlice(tt.data, head, capacity)

						expected := SubSlice(tt.data, start, end)

						subSlice := f.SubSlice(start, end)
						result := subSlice.ToSlice()

						msg := fmt.Sprintf("data: %v, cap: %d, head: %d, start: %d, end: %d",
							tt.data, capacity, head, start, end)
						testz.Equal(t, expected, result, msg)
					}
				}
			}
		}
	}
}

func TestFlexSlice_Range(t *testing.T) {
	tests := []struct {
		data []int
	}{
		{data: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}},
		{data: []int{1}},
		{data: []int{}},
	}

	for _, tt := range tests {
		for capacity := len(tt.data); capacity < len(tt.data)+10; capacity++ {
			for head := 0; head < capacity; head++ {
				f := makeFlexSlice(tt.data, head, capacity)
				msg := fmt.Sprintf("data: %v, cap: %d, head: %d", tt.data, capacity, head)

				// Test full iteration
				result := []int{}
				indices := []int{}
				f.Range(func(index int, value int) bool {
					indices = append(indices, index)
					result = append(result, value)
					return true
				})

				expectedIndices := make([]int, len(tt.data))
				for i := range expectedIndices {
					expectedIndices[i] = i
				}

				testz.Equal(t, tt.data, result, "Full range values: "+msg)
				testz.Equal(t, expectedIndices, indices, "Full range indices: "+msg)

				// Test early exit
				if len(tt.data) > 1 {
					stopAt := len(tt.data) / 2
					result = []int{}
					f.Range(func(index int, value int) bool {
						if index == stopAt {
							return false
						}
						result = append(result, value)
						return true
					})
					testz.Equal(t, tt.data[:stopAt], result, "Early exit range: "+msg)
				}
			}
		}
	}
}

func TestFlexSlice_RevRange(t *testing.T) {
	tests := []struct {
		data []int
	}{
		{data: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}},
		{data: []int{1}},
		{data: []int{}},
	}

	for _, tt := range tests {
		for capacity := len(tt.data); capacity < len(tt.data)+10; capacity++ {
			for head := 0; head < capacity; head++ {
				f := makeFlexSlice(tt.data, head, capacity)
				msg := fmt.Sprintf("data: %v, cap: %d, head: %d", tt.data, capacity, head)

				// Test full reverse iteration
				result := []int{}
				indices := []int{}
				f.RevRange(func(index int, value int) bool {
					indices = append(indices, index)
					result = append(result, value)
					return true
				})

				expected := make([]int, len(tt.data))
				expectedIndices := make([]int, len(tt.data))
				for i := 0; i < len(tt.data); i++ {
					expected[i] = tt.data[len(tt.data)-1-i]
					expectedIndices[i] = len(tt.data) - 1 - i
				}

				testz.Equal(t, expected, result, "Full rev-range values: "+msg)
				testz.Equal(t, expectedIndices, indices, "Full rev-range indices: "+msg)

				// Test early exit
				if len(tt.data) > 1 {
					stopAt := len(tt.data) / 2
					result = []int{}
					f.RevRange(func(index int, value int) bool {
						if index == stopAt {
							return false
						}
						result = append(result, value)
						return true
					})

					expected = tt.data[stopAt+1:]
					// reverse the expected slice
					for i, j := 0, len(expected)-1; i < j; i, j = i+1, j-1 {
						expected[i], expected[j] = expected[j], expected[i]
					}

					testz.Equal(t, expected, result, "Early exit rev-range: "+msg)
				}
			}
		}
	}
}

func TestFlexSlice_Shrink(t *testing.T) {
	tests := []struct {
		name         string
		initialCap   int
		initialLen   int
		head         int
		shouldShrink bool
		expectedCap  int
	}{
		{
			name:         "no shrink - cap too small",
			initialCap:   8,
			initialLen:   1,
			head:         0,
			shouldShrink: false,
			expectedCap:  8,
		},
		{
			name:         "no shrink - len too large",
			initialCap:   33,
			initialLen:   9, // 9 > 33/4
			head:         5,
			shouldShrink: false,
			expectedCap:  33,
		},
		{
			name:         "shrink - simple case",
			initialCap:   33,
			initialLen:   8, // 8 <= 33/4
			head:         0,
			shouldShrink: true,
			expectedCap:  16, // 33/2
		},
		{
			name:         "shrink - with head wrap",
			initialCap:   64,
			initialLen:   10,
			head:         60,
			shouldShrink: true,
			expectedCap:  32,
		},
		{
			name:         "shrink - to min cap",
			initialCap:   15,
			initialLen:   3, // 3 <= 15/4
			head:         2,
			shouldShrink: true,
			expectedCap:  8, // 15/2 = 7, but min is 8
		},
		{
			name:         "shrink - exact quarter",
			initialCap:   32,
			initialLen:   8, // 8 <= 32/4
			head:         10,
			shouldShrink: true,
			expectedCap:  16,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := make([]int, tt.initialLen)
			for i := 0; i < tt.initialLen; i++ {
				data[i] = i
			}

			f := makeFlexSlice(data, tt.head, tt.initialCap)
			originalSlice := f.ToSlice()

			shrunk := f.Shrink()

			if shrunk != tt.shouldShrink {
				t.Errorf("Shrink() returned %v, want %v", shrunk, tt.shouldShrink)
			}

			if f.Cap() != tt.expectedCap {
				t.Errorf("Cap() = %v, want %v", f.Cap(), tt.expectedCap)
			}

			if f.Len() != tt.initialLen {
				t.Errorf("Len() = %v, want %v", f.Len(), tt.initialLen)
			}

			testz.Equal(t, originalSlice, f.ToSlice(), "data should be preserved after shrink")
		})
	}
}

func TestFlexSlice_Grow(t *testing.T) {
	tests := []struct {
		name        string
		initialData []int
		initialCap  int
		head        int
		growN       uint
		expectedCap int
	}{
		{
			name:        "no grow needed",
			initialData: []int{1, 2, 3},
			initialCap:  10,
			head:        0,
			growN:       2, // 3 + 2 = 5 < 10, so no grow
			expectedCap: 10,
		},
		{
			name:        "grow empty slice",
			initialData: []int{},
			initialCap:  0,
			head:        0,
			growN:       5, // minCap is 5, default is 8
			expectedCap: 8,
		},
		{
			name:        "grow small slice",
			initialData: []int{1, 2, 3},
			initialCap:  4,
			head:        0,
			growN:       2, // 3 + 2 = 5. newCap is 4*2=8
			expectedCap: 8,
		},
		{
			name:        "grow with head wrap",
			initialData: []int{1, 2, 3, 4, 5},
			initialCap:  6,
			head:        3, // data is at indices 3,4,5,0,1
			growN:       2, // 5 + 2 = 7. newCap is 6*2=12
			expectedCap: 12,
		},
		{
			name:        "grow with large n",
			initialData: []int{1, 2, 3},
			initialCap:  4,
			head:        1,
			growN:       10, // 3 + 10 = 13. newCap is 13
			expectedCap: 13,
		},
		{
			name:        "grow with zero n",
			initialData: []int{1, 2, 3},
			initialCap:  3,
			head:        0,
			growN:       0, // 3 + 0 = 3, no grow
			expectedCap: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := makeFlexSlice(tt.initialData, tt.head, tt.initialCap)
			originalSlice := f.ToSlice()
			originalLen := f.Len()

			f.Grow(tt.growN)

			if f.Cap() != tt.expectedCap {
				t.Errorf("Cap() = %v, want %v", f.Cap(), tt.expectedCap)
			}
			if f.Len() != originalLen {
				t.Errorf("Len() = %v, want %v", f.Len(), originalLen)
			}
			testz.Equal(t, originalSlice, f.ToSlice(), "data should be preserved after grow")
		})
	}
}

func TestFlexSlice_Clear(t *testing.T) {
	tests := []struct {
		name        string
		initialData []int
		initialCap  int
		head        int
		shrink      bool
		expectedCap int
	}{
		{
			name:        "clear without shrink",
			initialData: []int{1, 2, 3, 4, 5},
			initialCap:  10,
			head:        2,
			shrink:      false,
			expectedCap: 10,
		},
		{
			name:        "clear with shrink",
			initialData: []int{1, 2},
			initialCap:  32,
			head:        10,
			shrink:      true,
			expectedCap: 16, // 32 / 2
		},
		{
			name:        "clear with shrink 2",
			initialData: []int{1},
			initialCap:  10,
			head:        0,
			shrink:      true,
			expectedCap: 8, // 10/2=5, but min is 8
		},
		{
			name:        "clear with shrink 3",
			initialData: []int{1, 2, 3},
			initialCap:  10,
			head:        5,
			shrink:      true,
			expectedCap: 8,
		},
		{
			name:        "clear empty slice",
			initialData: []int{},
			initialCap:  16,
			head:        0,
			shrink:      true,
			expectedCap: 8,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := makeFlexSlice(tt.initialData, tt.head, tt.initialCap)

			f.Clear(tt.shrink)

			if f.Len() != 0 {
				t.Errorf("Len() after Clear() = %v, want 0", f.Len())
			}
			if f.Cap() != tt.expectedCap {
				t.Errorf("Cap() after Clear() = %v, want %v", f.Cap(), tt.expectedCap)
			}
			if f.head != 0 {
				t.Errorf("head after Clear() = %v, want 0", f.head)
			}
			if f.tail != 0 {
				t.Errorf("tail after Clear() = %v, want 0", f.tail)
			}

			for i, v := range f.values {
				var zero int
				if v != zero {
					fmt.Println(f.values)
					t.Errorf("value at index %d should be zeroed, but got %v", i, v)
				}
			}
		})
	}
}

func makeFlexSlice(s []int, head, cap int) *FlexSlice[int] {
	if cap < len(s) {
		cap = len(s)
	}

	var f FlexSlice[int]
	if head < 0 || (head >= cap && cap > 0) || (head > cap && cap == 0) {
		panic("head must be in range [0, cap)")
	}

	if cap == len(s) {
		// full slice
		// tail must be equal to head
		f.tail = head

	} else {
		f.tail = (head + len(s)) % cap
	}

	f.values = make([]int, cap)
	for i, v := range s {
		f.values[(i+head)%cap] = v
	}
	f.head = head
	f.len = len(s)
	return &f

}
