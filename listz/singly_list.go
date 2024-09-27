package listz

// SNode is a node of a singly linked list.
type SNode[T any] struct {
	Value T
	next  *SNode[T]
}

// Next returns the next list node or nil.
func (e *SNode[T]) Next() *SNode[T] {
	return e.next
}

// SList represents a singly linked list.
type SList[T any] struct {
	head *SNode[T]
	tail *SNode[T]
	len  int
}

// Len returns the number of nodes of list l.
func (l *SList[T]) Len() int {
	return l.len
}

// Front returns the first node of list l or nil if the list is empty.
func (l *SList[T]) Front() *SNode[T] {
	return l.head
}

// Back returns the last node of list l or nil if the list is empty.
func (l *SList[T]) Back() *SNode[T] {
	return l.tail
}

// Get returns the node at index i.
func (l *SList[T]) Get(i int) *SNode[T] {
	if !l.withinRange(i) {
		return nil
	}

	e := l.head
	for index := 0; index < i; index++ {
		e = e.next
	}
	return e
}

// Remove removes the node at index i from the list.
// This operation is O(n) where n is the len of the list.
func (l *SList[T]) Remove(i int) *SNode[T] {
	if !l.withinRange(i) {
		return nil
	}

	var before *SNode[T]
	e := l.head
	for index := 0; index < i; index++ {
		before = e
		e = e.next
	}

	if e == l.head {
		l.head = e.next
	}
	if e == l.tail {
		l.tail = before
	}
	if before != nil {
		before.next = e.next
	}
	e.next = nil
	l.len--
	return e
}

// RemoveFront removes the first node from the list and returns it.
func (l *SList[T]) RemoveFront() *SNode[T] {
	if l.len == 0 {
		return nil
	}

	e := l.head
	l.head = e.next
	e.next = nil
	if l.len == 1 {
		l.tail = nil
	}
	l.len--
	return e
}

// PushFront inserts a new node with value v at the front of list l.
func (l *SList[T]) PushFront(v T) {
	e := &SNode[T]{Value: v}
	l.PushFrontNode(e)
}

// PushBack inserts a new node with value v at the back of list l.
func (l *SList[T]) PushBack(v T) {
	e := &SNode[T]{Value: v}
	l.PushBackNode(e)
}

// InsertAt inserts a new node with value v at index i.
// The old node at index i will be pushed to the next index.
func (l *SList[T]) InsertAt(i int, v T) {
	e := &SNode[T]{Value: v}
	l.InsertNodeAt(i, e)
}

// PushFrontNode inserts a new node at the front of the list.
func (l *SList[T]) PushFrontNode(e *SNode[T]) {
	e.next = l.head
	l.head = e

	if l.len == 0 {
		l.tail = e
	}
	l.len++
}

// PushBackNode inserts a new node at the back of the list.
func (l *SList[T]) PushBackNode(e *SNode[T]) {
	if l.len == 0 {
		l.head = e
	} else {
		l.tail.next = e
	}
	l.tail = e
	l.len++
}

// InsertNodeAt inserts a new node at index i.
// The old node at index i will be pushed to the next index.
func (l *SList[T]) InsertNodeAt(i int, e *SNode[T]) {
	if i <= 0 {
		l.PushFrontNode(e)
		return
	}

	if i >= l.len {
		l.PushBackNode(e)
		return
	}

	before := l.head
	for index := 0; index < i-1; index++ {
		before = before.next
	}

	e.next = before.next
	before.next = e
	l.len++
}

// Swap swaps the values of two nodes at indexes i and j.
func (l *SList[T]) Swap(i, j int) {
	if l.withinRange(i) && l.withinRange(j) && i != j {
		var e1, e2 *SNode[T]
		for index, ce := 0, l.head; e1 == nil || e2 == nil; index, ce = index+1, ce.next {
			switch index {
			case i:
				e1 = ce
			case j:
				e2 = ce
			}
		}
		e1.Value, e2.Value = e2.Value, e1.Value
	}
}

func (l *SList[T]) withinRange(index int) bool {
	return index >= 0 && index < l.len
}
