package listz

import (
	"context"
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

// Enqueue adds a value to the end of the list.
func (l *SyncList[T]) Enqueue(value T) {
	node := unsafe.Pointer(&syncNode[T]{value: value})

	for {
		tail := atomic.LoadPointer(&l.tail)
		tailNode := (*syncNode[T])(tail)
		next := atomic.LoadPointer(&tailNode.next)

		if next == nil && atomic.CompareAndSwapPointer(&tailNode.next, next, node) {
			atomic.StorePointer(&l.tail, node)
			atomic.AddInt64(&l.len, 1)
			return
		}

		runtime.Gosched()
	}
}

// Push adds a value to the end of the list.
//
// Deprecated: Use Enqueue instead.
func (l *SyncList[T]) Push(value T) {
	l.Enqueue(value)
}

// EnqueueContext attempts to add a value to the end of the list.
// If the context is canceled or exceeds its deadline, it returns false.
func (l *SyncList[T]) EnqueueContext(ctx context.Context, value T) bool {
	node := unsafe.Pointer(&syncNode[T]{value: value})

	for {
		select {
		case <-ctx.Done():
			return false
		default:
		}

		tail := atomic.LoadPointer(&l.tail)
		tailNode := (*syncNode[T])(tail)
		next := atomic.LoadPointer(&tailNode.next)

		if next == nil && atomic.CompareAndSwapPointer(&tailNode.next, next, node) {
			atomic.StorePointer(&l.tail, node)
			atomic.AddInt64(&l.len, 1)
			return true
		}

		runtime.Gosched()
	}
}

// PushContext attempts to add a value to the end of the list with context.
func (l *SyncList[T]) PushContext(ctx context.Context, value T) bool {
	return l.EnqueueContext(ctx, value)
}

// Dequeue removes and returns the value at the front of the list.
// If the list is empty or concurrent Dequeue is in progress, it returns false.
func (l *SyncList[T]) Dequeue() (T, bool) {
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

// Pop removes and returns the value at the front of the list.
// If the list is empty or concurrent Pop is in progress, it returns false.
//
// Deprecated: Use Dequeue instead.
func (l *SyncList[T]) Pop() (T, bool) {
	return l.Dequeue()
}

// DequeueWait removes and returns the value at the front of the list.
// If maxWait is negative, it will block until the value is dequeued.
func (l *SyncList[T]) DequeueWait(maxWait time.Duration) (T, bool) {
	if maxWait < 0 {
		for {
			if value, ok := l.Dequeue(); ok {
				return value, ok
			}

			runtime.Gosched()
		}
	}

	if v, ok := l.Dequeue(); ok {
		return v, true
	}

	var zero T
	if maxWait == 0 {
		return zero, false
	}

	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()
	begin := time.Now()

	for {
		now := <-ticker.C

		if v, ok := l.Dequeue(); ok {
			return v, true
		}

		if now.Sub(begin) >= maxWait {
			return zero, false
		}
	}
}

// PopWait removes and returns the value at the front of the list.
// If maxWait is negative, it will block until the value is popped.
//
// Deprecated: Use DequeueWait instead.
func (l *SyncList[T]) PopWait(maxWait time.Duration) (T, bool) {
	return l.DequeueWait(maxWait)
}

// DequeueContext removes and returns the value at the front of the list.
// It blocks until a value is available or until ctx is done.
func (l *SyncList[T]) DequeueContext(ctx context.Context) (T, bool) {
	var zero T
	select {
	case <-ctx.Done():
		return zero, false
	default:
	}

	if v, ok := l.Dequeue(); ok {
		return v, true
	}

	// Spin a few times before sleeping on ticker to handle immediate producers
	for i := 0; i < 10; i++ {
		select {
		case <-ctx.Done():
			return zero, false
		default:
		}
		runtime.Gosched()
		if v, ok := l.Dequeue(); ok {
			return v, true
		}
	}

	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return zero, false
		case <-ticker.C:
			if v, ok := l.Dequeue(); ok {
				return v, true
			}
		}
	}
}

// PopContext removes and returns the value at the front of the list with context support.
func (l *SyncList[T]) PopContext(ctx context.Context) (T, bool) {
	return l.DequeueContext(ctx)
}
