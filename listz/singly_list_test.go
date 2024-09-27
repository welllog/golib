package listz

import (
	"testing"
)

func TestSList(t *testing.T) {
	list := &SList[int]{}

	// Test PushFront
	list.PushFront(1)
	if list.Len() != 1 {
		t.Errorf("expected length 1, got %d", list.Len())
	}
	if list.Front().Value != 1 {
		t.Errorf("expected front value 1, got %d", list.Front().Value)
	}
	if list.Back().Value != 1 {
		t.Errorf("expected back value 1, got %d", list.Back().Value)
	}

	// Test PushBack
	list.PushBack(2)
	if list.Len() != 2 {
		t.Errorf("expected length 2, got %d", list.Len())
	}
	if list.Back().Value != 2 {
		t.Errorf("expected back value 2, got %d", list.Back().Value)
	}

	// Test InsertAt
	list.InsertAt(1, 3)
	if list.Len() != 3 {
		t.Errorf("expected length 3, got %d", list.Len())
	}
	if list.Front().Next().Value != 3 {
		t.Errorf("expected value 3 at index 1, got %d", list.Front().Next().Value)
	}

	// Test Remove
	removedNode := list.Remove(1)
	if removedNode.Value != 3 {
		t.Errorf("expected removed value 3, got %d", removedNode.Value)
	}
	if list.Len() != 2 {
		t.Errorf("expected length 2, got %d", list.Len())
	}

	// Test RemoveFront
	removedFront := list.RemoveFront()
	if removedFront.Value != 1 {
		t.Errorf("expected removed front value 1, got %d", removedFront.Value)
	}
	if list.Len() != 1 {
		t.Errorf("expected length 1, got %d", list.Len())
	}
	if list.Front().Value != 2 {
		t.Errorf("expected front value 2, got %d", list.Front().Value)
	}

	// Test RemoveFront on empty list
	list.RemoveFront()
	removedFront = list.RemoveFront()
	if removedFront != nil {
		t.Errorf("expected nil, got %v", removedFront)
	}
	if list.Len() != 0 {
		t.Errorf("expected length 0, got %d", list.Len())
	}
}
func TestSListGet(t *testing.T) {
	list := &SList[int]{}

	// Test Get on empty list
	if node := list.Get(0); node == nil {
		t.Log("Get on empty list returned nil as expected")
	} else {
		t.Errorf("expected nil, got %v", node)
	}

	// Add elements to the list
	list.PushBack(1)
	list.PushBack(2)
	list.PushBack(3)

	// Test Get on valid indices
	if node := list.Get(0); node == nil || node.Value != 1 {
		t.Errorf("expected value 1 at index 0, got %v", node)
	}
	if node := list.Get(1); node == nil || node.Value != 2 {
		t.Errorf("expected value 2 at index 1, got %v", node)
	}
	if node := list.Get(2); node == nil || node.Value != 3 {
		t.Errorf("expected value 3 at index 2, got %v", node)
	}

	// Test Get on invalid indices
	if node := list.Get(-1); node != nil {
		t.Errorf("expected nil for index -1, got %v", node)
	}
	if node := list.Get(3); node != nil {
		t.Errorf("expected nil for index 3, got %v", node)
	}
}
func TestSListSwap(t *testing.T) {
	list := &SList[int]{}

	// Add elements to the list
	list.PushBack(1)
	list.PushBack(2)
	list.PushBack(3)
	list.PushBack(4)

	// Test Swap with valid indices
	list.Swap(1, 3)
	if list.Get(1).Value != 4 {
		t.Errorf("expected value 4 at index 1, got %d", list.Get(1).Value)
	}
	if list.Get(3).Value != 2 {
		t.Errorf("expected value 2 at index 3, got %d", list.Get(3).Value)
	}

	// Test Swap with the same index
	list.Swap(2, 2)
	if list.Get(2).Value != 3 {
		t.Errorf("expected value 3 at index 2, got %d", list.Get(2).Value)
	}

	// Test Swap with out-of-range indices
	list.Swap(-1, 2)
	if list.Get(2).Value != 3 {
		t.Errorf("expected value 3 at index 2 after out-of-range swap, got %d", list.Get(2).Value)
	}
	list.Swap(2, 4)
	if list.Get(2).Value != 3 {
		t.Errorf("expected value 3 at index 2 after out-of-range swap, got %d", list.Get(2).Value)
	}

	// Test Swap with one valid and one invalid index
	list.Swap(1, 5)
	if list.Get(1).Value != 4 {
		t.Errorf("expected value 4 at index 1 after invalid swap, got %d", list.Get(1).Value)
	}
}
