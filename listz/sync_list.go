package listz

import (
	"runtime"
	"sync"
	"sync/atomic"
	"unsafe"
)

type syncNode[T any] struct {
	value T
	next  unsafe.Pointer
}

type SyncList[T any] struct {
	len  int64
	head unsafe.Pointer
	tail unsafe.Pointer
	pool sync.Pool
}

func NewSync[T any]() *SyncList[T] {
	// Initialize with a dummy node
	dummy := unsafe.Pointer(&syncNode[T]{})
	return &SyncList[T]{
		head: dummy,
		tail: dummy,
		pool: sync.Pool{
			New: func() any {
				return &syncNode[T]{}
			},
		},
	}
}

func (l *SyncList[T]) Len() int {
	return int(atomic.LoadInt64(&l.len))
}

func (l *SyncList[T]) Push(value T) {
	// node := unsafe.Pointer(&syncNode[T]{value: value})
	node := unsafe.Pointer(l.getNode(value))

	for {
		tail := atomic.LoadPointer(&l.tail)
		tailNode := (*syncNode[T])(tail)
		next := atomic.LoadPointer(&tailNode.next)

		if tail == atomic.LoadPointer(&l.tail) {
			if next == nil { // tail is really pointing to the last node
				// Try to link the node at the end of the list
				if atomic.CompareAndSwapPointer(&tailNode.next, next, node) {
					// CAS succeeded, try to swing the tail to the inserted node
					atomic.CompareAndSwapPointer(&l.tail, tail, node)
					atomic.AddInt64(&l.len, 1)
					return
				}
			} else {
				// The tail was not pointing to the last node, try to swing the tail to the next node
				atomic.CompareAndSwapPointer(&l.tail, tail, next)
			}
		}

		runtime.Gosched()
	}
}

func (l *SyncList[T]) Pop() (T, bool) {
	var zero T

	for {
		head := atomic.LoadPointer(&l.head)
		tail := atomic.LoadPointer(&l.tail)
		next := atomic.LoadPointer(&((*syncNode[T])(head)).next)

		if head == atomic.LoadPointer(&l.head) {
			if head == tail { // The list is empty or tail is lagging behind
				if next == nil {
					return zero, false
				}
				// Try to swing the tail to the next node
				atomic.CompareAndSwapPointer(&l.tail, tail, next)
			} else {
				value := (*syncNode[T])(next).value
				// Try to swing the head to the next node
				if atomic.CompareAndSwapPointer(&l.head, head, next) {
					headNode := (*syncNode[T])(head)
					l.releaseNode(headNode)

					atomic.AddInt64(&l.len, -1)
					return value, true
				}
			}
		}

		runtime.Gosched()
	}
}

func (l *SyncList[T]) push(value T) {
	node := unsafe.Pointer(l.getNode(value))
	// node := unsafe.Pointer(&syncNode[T]{value: value})

	for {
		tail := atomic.LoadPointer(&l.tail)
		tailNode := (*syncNode[T])(tail)
		next := atomic.LoadPointer(&tailNode.next)

		if next == nil {
			if atomic.CompareAndSwapPointer(&tailNode.next, next, node) {
				// atomic.CompareAndSwapPointer(&l.tail, tail, node)
				atomic.StorePointer(&l.tail, node)
				atomic.AddInt64(&l.len, 1)
				return
			}
		}

		runtime.Gosched()
	}
}

func (l *SyncList[T]) pop() (T, bool) {
	head := atomic.LoadPointer(&l.head)
	tail := atomic.LoadPointer(&l.tail)

	var zero T
	if head == tail {
		return zero, false
	}

	headNode := (*syncNode[T])(head)
	next := atomic.LoadPointer(&headNode.next)
	if atomic.CompareAndSwapPointer(&l.head, head, next) {
		l.releaseNode(headNode)

		node := (*syncNode[T])(next)
		value := node.value
		node.value = zero
		atomic.AddInt64(&l.len, -1)
		return value, true
	}

	return zero, false
}

func (l *SyncList[T]) getNode(v T) *syncNode[T] {
	node := l.pool.Get().(*syncNode[T])
	node.value = v
	return node
}

func (l *SyncList[T]) releaseNode(node *syncNode[T]) {
	node.next = unsafe.Pointer(nil)
	l.pool.Put(node)
}
