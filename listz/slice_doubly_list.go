package listz

import (
	"math"
)

const nullIdx uint16 = math.MaxUint16

type SliceDList[T any] struct {
	nodes []sdNode[T]
	head  uint16
	tail  uint16
	free  uint16
	len   int
}

type sdNode[T any] struct {
	value T
	prev  uint16
	next  uint16
}

// Init initializes or clears the list.
func (l *SliceDList[T]) Init(cap int) {
	if cap < 0 {
		cap = 0
	}
	if cap > math.MaxUint16 {
		cap = math.MaxUint16
	}
	l.nodes = make([]sdNode[T], cap)
	l.head = nullIdx
	l.tail = nullIdx
	l.len = 0
	l.free = nullIdx
	if cap > 0 {
		l.free = 0
		for i := 0; i < cap-1; i++ {
			l.nodes[i].next = uint16(i + 1)
		}
		l.nodes[cap-1].next = nullIdx
	}
}

// Len returns the number of elements in the list.
func (l *SliceDList[T]) Len() int {
	return l.len
}

// Cap returns the capacity of the list.
func (l *SliceDList[T]) Cap() int {
	return len(l.nodes)
}

// HasFree checks if there are free nodes available in the list.
func (l *SliceDList[T]) HasFree() bool {
	return l.free != nullIdx
}

// PushFront insert a new element at the front of the list.
func (l *SliceDList[T]) PushFront(value T) (uint16, bool) {
	idx, ok := l.allocNode()
	if !ok {
		return nullIdx, false
	}

	l.nodes[idx] = sdNode[T]{
		value: value,
		prev:  nullIdx,
		next:  l.head,
	}

	if l.head != nullIdx {
		l.nodes[l.head].prev = idx
	}
	l.head = idx

	if l.tail == nullIdx {
		l.tail = idx
	}

	l.len++
	return idx, true
}

// PushBack inserts a new element at the back of the list.
func (l *SliceDList[T]) PushBack(value T) (uint16, bool) {
	idx, ok := l.allocNode()
	if !ok {
		return nullIdx, false
	}

	l.nodes[idx] = sdNode[T]{
		value: value,
		prev:  l.tail,
		next:  nullIdx,
	}

	if l.tail != nullIdx {
		l.nodes[l.tail].next = idx
	}
	l.tail = idx

	if l.head == nullIdx {
		l.head = idx
	}

	l.len++
	return idx, true
}

// Remove removes the element at the specified index from the list.
func (l *SliceDList[T]) Remove(idx uint16) bool {
	if !l.isLiveNode(idx) {
		return false
	}

	node := l.nodes[idx]
	prevIdx := node.prev
	nextIdx := node.next

	if prevIdx != nullIdx {
		l.nodes[prevIdx].next = nextIdx
	} else {
		l.head = nextIdx
	}

	if nextIdx != nullIdx {
		l.nodes[nextIdx].prev = prevIdx
	} else {
		l.tail = prevIdx
	}

	l.nodes[idx] = sdNode[T]{prev: nullIdx, next: l.free}
	l.free = idx
	l.len--
	return true
}

// Front return the index and value of the first element in the list.
func (l *SliceDList[T]) Front() (idx uint16, value T, ok bool) {
	if l.head == nullIdx {
		return nullIdx, value, false
	}
	return l.head, l.nodes[l.head].value, true
}

// Back return the index and value of the last element in the list.
func (l *SliceDList[T]) Back() (idx uint16, value T, ok bool) {
	if l.tail == nullIdx {
		return nullIdx, value, false
	}
	return l.tail, l.nodes[l.tail].value, true
}

// Get returns the value at the specified index in the list.
func (l *SliceDList[T]) Get(idx uint16) (value T, ok bool) {
	if !l.isLiveNode(idx) {
		return value, false
	}
	return l.nodes[idx].value, true
}

// InsertAfter at the specified index `mark` inserts a new element after the node at `mark`.
func (l *SliceDList[T]) InsertAfter(value T, mark uint16) (uint16, bool) {
	if !l.isLiveNode(mark) {
		return nullIdx, false
	}

	if mark == l.tail {
		return l.PushBack(value)
	}

	idx, ok := l.allocNode()
	if !ok {
		return nullIdx, false
	}
	oldNextIdx := l.nodes[mark].next

	l.nodes[idx] = sdNode[T]{
		value: value,
		prev:  mark,
		next:  oldNextIdx,
	}

	l.nodes[mark].next = idx
	l.nodes[oldNextIdx].prev = idx
	l.len++

	return idx, true
}

// InsertBefore at the specified index `mark` inserts a new element before the node at `mark`.
func (l *SliceDList[T]) InsertBefore(value T, mark uint16) (uint16, bool) {
	if !l.isLiveNode(mark) {
		return nullIdx, false
	}

	if mark == l.head {
		return l.PushFront(value)
	}

	idx, ok := l.allocNode()
	if !ok {
		return nullIdx, false
	}
	oldPrevIdx := l.nodes[mark].prev

	l.nodes[idx] = sdNode[T]{
		value: value,
		prev:  oldPrevIdx,
		next:  mark,
	}

	l.nodes[mark].prev = idx
	l.nodes[oldPrevIdx].next = idx
	l.len++

	return idx, true
}

// MoveToFront moves the node at the specified index to the front of the list.
func (l *SliceDList[T]) MoveToFront(idx uint16) bool {
	if !l.isLiveNode(idx) {
		return false
	}

	if l.head == idx {
		return true // Already at the front
	}

	// unlink the node from its current position
	l.unlink(idx)

	// link it to the front
	l.nodes[idx].prev = nullIdx
	l.nodes[idx].next = l.head
	l.nodes[l.head].prev = idx
	l.head = idx

	return true
}

// MoveToBack moves the node at the specified index to the back of the list.
func (l *SliceDList[T]) MoveToBack(idx uint16) bool {
	if !l.isLiveNode(idx) {
		return false
	}

	if l.tail == idx {
		return true
	}

	// unlink the node from its current position
	l.unlink(idx)

	// link it to the back
	l.nodes[idx].prev = l.tail
	l.nodes[idx].next = nullIdx
	l.nodes[l.tail].next = idx
	l.tail = idx

	return true
}

// MoveBefore moves the node at the specified index to before the node at the specified mark index.
func (l *SliceDList[T]) MoveBefore(idx, mark uint16) bool {
	if !l.isLiveNode(idx) || !l.isLiveNode(mark) {
		return false
	}

	if idx == mark {
		return true
	}

	if mark == l.head {
		return l.MoveToFront(idx)
	}

	// unlink the node from its current position
	l.unlink(idx)

	// insert the node before mark
	markPrevIdx := l.nodes[mark].prev
	l.nodes[idx].prev = markPrevIdx
	l.nodes[idx].next = mark
	l.nodes[markPrevIdx].next = idx
	l.nodes[mark].prev = idx

	return true
}

// MoveAfter moves the node at the specified index to after the node at the specified mark index.
func (l *SliceDList[T]) MoveAfter(idx, mark uint16) bool {
	if !l.isLiveNode(idx) || !l.isLiveNode(mark) {
		return false
	}

	if idx == mark {
		return true
	}

	if mark == l.tail {
		return l.MoveToBack(idx)
	}

	// unlink the node from its current position
	l.unlink(idx)

	// insert the node after mark
	markNextIdx := l.nodes[mark].next
	l.nodes[idx].prev = mark
	l.nodes[idx].next = markNextIdx
	l.nodes[markNextIdx].prev = idx
	l.nodes[mark].next = idx

	return true
}

// TODO Range

// unlink unlinks the node at the specified index from the list.
func (l *SliceDList[T]) unlink(idx uint16) {
	node := l.nodes[idx]
	prevIdx := node.prev
	nextIdx := node.next

	if prevIdx != nullIdx {
		l.nodes[prevIdx].next = nextIdx
	} else {
		l.head = nextIdx
	}

	if nextIdx != nullIdx {
		l.nodes[nextIdx].prev = prevIdx
	} else {
		l.tail = prevIdx
	}
}

// allocNode allocates a new node for the list.
func (l *SliceDList[T]) allocNode() (uint16, bool) {
	if l.free != nullIdx {
		idx := l.free
		l.free = l.nodes[idx].next
		return idx, true
	}
	if len(l.nodes) >= math.MaxUint16 {
		return nullIdx, false
	}

	idx := uint16(len(l.nodes))
	l.nodes = append(l.nodes, sdNode[T]{})
	return idx, true

	//idx := uint16(len(l.nodes))
	//l.nodes = append(l.nodes, sdNode[T]{})
	//if capacity := cap(l.nodes); capacity > len(l.nodes) {
	//	if capacity > math.MaxUint16 {
	//		capacity = math.MaxUint16
	//	}
	//
	//	l.nodes = l.nodes[:capacity]
	//
	//	if len(l.nodes) - 1 > int(idx) {
	//		l.free = idx + 1
	//		for i := int(idx + 1); i < len(l.nodes) - 1; i++ {
	//			l.nodes[i].next = uint16(i + 1)
	//		}
	//		l.nodes[len(l.nodes)-1].next = nullIdx
	//	}
	//
	//}
	//return idx, true
}

// isLiveNode checks if the node at the given index is a live node in the list.
func (l *SliceDList[T]) isLiveNode(idx uint16) bool {
	if idx >= uint16(len(l.nodes)) {
		return false
	}
	// if node is head or has a previous node, it is in the list.
	return l.head == idx || l.nodes[idx].prev != nullIdx
}
