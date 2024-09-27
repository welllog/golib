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
