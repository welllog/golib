package listz

import (
	"runtime"
	"sync/atomic"
	"time"
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
}

// NewSync creates a new SyncList.
func NewSync[T any]() *SyncList[T] {
	// Initialize with a dummy node
	dummy := unsafe.Pointer(&syncNode[T]{})
	return &SyncList[T]{
		head: dummy,
		tail: dummy,
	}
}

// Len returns the number of elements in the list.
func (l *SyncList[T]) Len() int {
	return int(atomic.LoadInt64(&l.len))
}

// Push adds a value to the end of the list.
func (l *SyncList[T]) Push(value T) {
	node := unsafe.Pointer(&syncNode[T]{value: value})

	for {
		tail := atomic.LoadPointer(&l.tail)
		tailNode := (*syncNode[T])(tail)
		next := atomic.LoadPointer(&tailNode.next)

		if next == nil && atomic.CompareAndSwapPointer(&tailNode.next, next, node) {
			// atomic.CompareAndSwapPointer(&l.tail, tail, node)
			atomic.StorePointer(&l.tail, node)
			atomic.AddInt64(&l.len, 1)
			return
		}

		runtime.Gosched()
	}
}

// Pop removes and returns the value at the front of the list.
// If the list is empty or concurrent Pop is in progress, it returns false.
func (l *SyncList[T]) Pop() (T, bool) {
	head := atomic.LoadPointer(&l.head)
	tail := atomic.LoadPointer(&l.tail)

	var zero T
	if head == tail {
		return zero, false
	}

	headNode := (*syncNode[T])(head)
	next := atomic.LoadPointer(&headNode.next)
	if atomic.CompareAndSwapPointer(&l.head, head, next) {
		node := (*syncNode[T])(next)
		value := node.value
		node.value = zero
		atomic.AddInt64(&l.len, -1)
		return value, true
	}

	return zero, false
}

// PopWait removes and returns the value at the front of the list.
// If maxWait is negative, it will block until the value is popped.
func (l *SyncList[T]) PopWait(maxWait time.Duration) (T, bool) {
	if maxWait < 0 {
		for {
			if value, ok := l.Pop(); ok {
				return value, ok
			}

			runtime.Gosched()
		}
	}

	if v, ok := l.Pop(); ok {
		return v, true
	}

	var zero T
	if maxWait == 0 {
		return zero, false
	}

	ticker := time.NewTicker(10 * time.Millisecond)
	begin := time.Now()

	for {
		now := <-ticker.C

		if v, ok := l.Pop(); ok {
			ticker.Stop()
			return v, true
		}

		if now.Sub(begin) >= maxWait {
			ticker.Stop()
			return zero, false
		}
	}
}
