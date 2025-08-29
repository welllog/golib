package listz

// DNode is a node of a doubly linked list.
type DNode[T any] struct {
	// Next and previous pointers in the doubly-linked list of nodes.
	// To simplify the implementation, internally a list l is implemented
	// as a ring, such that &l.root is both the next element of the last
	// list element (l.Back()) and the previous element of the first list
	// element (l.Front()).
	next, prev *DNode[T]

	// The list to which this node belongs.
	list *DList[T]

	// The value stored with this node.
	Value T
}

// Next returns the next list node or nil.
func (e *DNode[T]) Next() *DNode[T] {
	if p := e.next; e.list != nil && p != &e.list.root {
		return p
	}
	return nil
}

// Prev returns the previous list node or nil.
func (e *DNode[T]) Prev() *DNode[T] {
	if p := e.prev; e.list != nil && p != &e.list.root {
		return p
	}
	return nil
}

// DList represents a doubly linked list.
// The zero value for DList is an empty list ready to use.
type DList[T any] struct {
	root DNode[T] // sentinel list node, only &root, root.prev, and root.next are used
	len  int      // current list length excluding (this) sentinel node
}

// Init initializes or clears list l.
func (l *DList[T]) Init() *DList[T] {
	l.root.next = &l.root
	l.root.prev = &l.root
	l.len = 0
	return l
}

// NewDoubly returns an initialized list.
func NewDoubly[T any]() *DList[T] { return new(DList[T]).Init() }

// Len returns the number of nodes of list l.
// The complexity is O(1).
func (l *DList[T]) Len() int { return l.len }

// Front returns the first node of list l or nil if the list is empty.
func (l *DList[T]) Front() *DNode[T] {
	if l.len == 0 {
		return nil
	}
	return l.root.next
}

// Back returns the last node of list l or nil if the list is empty.
func (l *DList[T]) Back() *DNode[T] {
	if l.len == 0 {
		return nil
	}
	return l.root.prev
}

// Remove removes e from l if e is a node of list l.
// It returns the node value e.Value.
// The node must not be nil.
func (l *DList[T]) Remove(e *DNode[T]) T {
	if e.list == l {
		// if e.list == l, l must have been initialized when e was inserted
		// in l or l == nil (e is a zero node) and l.remove will crash
		l.remove(e)
	}
	return e.Value
}

// PushFront inserts a new node e with value v at the front of list l and returns e.
func (l *DList[T]) PushFront(v T) *DNode[T] {
	l.lazyInit()
	return l.insertValue(v, &l.root)
}

// PushBack inserts a new node e with value v at the back of list l and returns e.
func (l *DList[T]) PushBack(v T) *DNode[T] {
	l.lazyInit()
	return l.insertValue(v, l.root.prev)
}

// InsertBefore inserts a new node e with value v immediately before mark and returns e.
// If mark is not a node of l, the list is not modified.
// The mark must not be nil.
func (l *DList[T]) InsertBefore(v T, mark *DNode[T]) *DNode[T] {
	if mark.list != l {
		return nil
	}
	// see comment in List.Remove about initialization of l
	return l.insertValue(v, mark.prev)
}

// InsertAfter inserts a new node e with value v immediately after mark and returns e.
// If mark is not a node of l, the list is not modified.
// The mark must not be nil.
func (l *DList[T]) InsertAfter(v T, mark *DNode[T]) *DNode[T] {
	if mark.list != l {
		return nil
	}
	// see comment in List.Remove about initialization of l
	return l.insertValue(v, mark)
}

// PushFrontNode inserts a new node e at the front of list l.
func (l *DList[T]) PushFrontNode(e *DNode[T]) {
	l.lazyInit()
	l.insert(e, &l.root)
}

// PushBackNode inserts a new node e at the back of list l.
func (l *DList[T]) PushBackNode(e *DNode[T]) {
	l.lazyInit()
	l.insert(e, l.root.prev)
}

// InsertNodeBefore inserts a new node e before mark and returns e.
func (l *DList[T]) InsertNodeBefore(e, mark *DNode[T]) {
	if mark.list != l {
		return
	}
	l.insert(e, mark.prev)
}

// InsertNodeAfter inserts a new node e after mark and returns e.
func (l *DList[T]) InsertNodeAfter(e, mark *DNode[T]) {
	if mark.list != l {
		return
	}
	l.insert(e, mark)
}

// MoveToFront moves node e to the front of list l.
// If e is not a node of l, the list is not modified.
// The node must not be nil.
func (l *DList[T]) MoveToFront(e *DNode[T]) {
	if e.list != l || l.root.next == e {
		return
	}
	// see comment in List.Remove about initialization of l
	l.move(e, &l.root)
}

// MoveToBack moves node e to the back of list l.
// If e is not a node of l, the list is not modified.
// The node must not be nil.
func (l *DList[T]) MoveToBack(e *DNode[T]) {
	if e.list != l || l.root.prev == e {
		return
	}
	// see comment in List.Remove about initialization of l
	l.move(e, l.root.prev)
}

// MoveBefore moves node e to its new position before mark.
// If e or mark is not a node of l, or e == mark, the list is not modified.
// The node and mark must not be nil.
func (l *DList[T]) MoveBefore(e, mark *DNode[T]) {
	if e.list != l || e == mark || mark.list != l {
		return
	}
	l.move(e, mark.prev)
}

// MoveAfter moves node e to its new position after mark.
// If e or mark is not a node of l, or e == mark, the list is not modified.
// The node and mark must not be nil.
func (l *DList[T]) MoveAfter(e, mark *DNode[T]) {
	if e.list != l || e == mark || mark.list != l {
		return
	}
	l.move(e, mark)
}

// PushBackDList inserts a copy of another list at the back of list l.
// The lists l and other may be the same. They must not be nil.
func (l *DList[T]) PushBackDList(other *DList[T]) {
	l.lazyInit()
	for i, e := other.Len(), other.Front(); i > 0; i, e = i-1, e.Next() {
		l.insertValue(e.Value, l.root.prev)
	}
}

// PushFrontDList inserts a copy of another list at the front of list l.
// The lists l and other may be the same. They must not be nil.
func (l *DList[T]) PushFrontDList(other *DList[T]) {
	l.lazyInit()
	for i, e := other.Len(), other.Back(); i > 0; i, e = i-1, e.Prev() {
		l.insertValue(e.Value, &l.root)
	}
}

// lazyInit lazily initializes a zero List value.
func (l *DList[T]) lazyInit() {
	if l.root.next == nil {
		l.Init()
	}
}

// insert inserts e after at, increments l.len, and returns e.
func (l *DList[T]) insert(e, at *DNode[T]) *DNode[T] {
	e.prev = at
	e.next = at.next
	e.prev.next = e
	e.next.prev = e
	e.list = l
	l.len++
	return e
}

// insertValue is a convenience wrapper for insert(&DNode{Value: v}, at).
func (l *DList[T]) insertValue(v T, at *DNode[T]) *DNode[T] {
	return l.insert(&DNode[T]{Value: v}, at)
}

// remove removes e from its list, decrements l.len
func (l *DList[T]) remove(e *DNode[T]) {
	e.prev.next = e.next
	e.next.prev = e.prev
	e.next = nil // avoid memory leaks
	e.prev = nil // avoid memory leaks
	e.list = nil
	l.len--
}

// move moves e to next to at.
func (l *DList[T]) move(e, at *DNode[T]) {
	if e == at {
		return
	}
	e.prev.next = e.next
	e.next.prev = e.prev

	e.prev = at
	e.next = at.next
	e.prev.next = e
	e.next.prev = e
}
