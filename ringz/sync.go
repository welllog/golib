package ringz

import (
	"context"
	"runtime"
	"strconv"
	"sync/atomic"
	"time"
)

type item[T any] struct {
	value T
	pos   uint32
}

type SyncRing[T any] struct {
	values []item[T]
	cap    uint32
	mask   uint32
	head   uint32
	tail   uint32
}

func NewSync[T any](capacity int) SyncRing[T] {
	var r SyncRing[T]
	r.Init(capacity)
	return r
}

func (r *SyncRing[T]) Init(capacity int) {
	var c uint32
	switch {
	case capacity <= 0:
		panic("ringz.SyncRing Init: invalid capacity: " + strconv.Itoa(capacity))
	case capacity > 1<<30:
		panic("ringz.SyncRing Init: capacity exceeds maximum limit: " + strconv.Itoa(capacity))
	case 1 == capacity:
		c = 2
	default:
		c = uint32(capacity)
		if c&(c-1) > 0 {
			c = roundUpPowOfTwo(c)
		}
	}

	r.cap = c
	r.mask = c - 1
	r.head = 0
	r.tail = 0
	r.values = make([]item[T], c)

	for i := range r.values {
		r.values[i].pos = uint32(i)
	}
}

// IsEmpty returns true if the ring is empty.
func (r *SyncRing[T]) IsEmpty() bool {
	return atomic.LoadUint32(&r.head) == atomic.LoadUint32(&r.tail)
}

// IsFull returns true if the ring is full.
func (r *SyncRing[T]) IsFull() bool {
	return atomic.LoadUint32(&r.tail)-atomic.LoadUint32(&r.head) == r.cap
}

// Len returns the length of the ring.
func (r *SyncRing[T]) Len() int {
	l := atomic.LoadUint32(&r.tail) - atomic.LoadUint32(&r.head)
	if l > r.cap {
		return int(r.cap)
	}

	return int(l)
}

// Cap returns the capacity of the ring.
func (r *SyncRing[T]) Cap() int {
	return int(r.cap)
}

// Enqueue pushes the value to queue tail.
// Notice if return false, means the queue is full or concurrent Enqueue operation.
func (r *SyncRing[T]) Enqueue(value T) bool {
	pos := atomic.LoadUint32(&r.tail)
	holder := &r.values[pos&r.mask]
	seq := atomic.LoadUint32(&holder.pos)

	if pos != seq {
		return false
	}

	if !atomic.CompareAndSwapUint32(&r.tail, pos, pos+1) {
		return false
	}

	holder.value = value
	atomic.StoreUint32(&holder.pos, seq+1)
	return true
}

// Push pushes the value to queue tail.
//
// Deprecated: Use Enqueue instead.
func (r *SyncRing[T]) Push(value T) bool {
	return r.Enqueue(value)
}

// Dequeue removes and returns the value from queue head.
// Notice if return false, means the queue is empty or concurrent Dequeue operation.
func (r *SyncRing[T]) Dequeue() (T, bool) {
	pos := atomic.LoadUint32(&r.head)
	holder := &r.values[pos&r.mask]
	seq := atomic.LoadUint32(&holder.pos)

	var zero T
	if pos+1 != seq {
		return zero, false
	}

	if !atomic.CompareAndSwapUint32(&r.head, pos, pos+1) {
		return zero, false
	}

	value := holder.value
	holder.value = zero
	atomic.StoreUint32(&holder.pos, seq+r.mask)
	return value, true
}

// Pop removes and returns the value from queue head.
//
// Deprecated: Use Dequeue instead.
func (r *SyncRing[T]) Pop() (T, bool) {
	return r.Dequeue()
}

// EnqueueWait pushes the value to queue tail with max wait duration.
// If maxWait is negative, it will block until the value is enqueued.
func (r *SyncRing[T]) EnqueueWait(value T, maxWait time.Duration) bool {
	if maxWait < 0 {
		for {
			if r.Enqueue(value) {
				return true
			}

			runtime.Gosched()
		}
	}

	if r.Enqueue(value) {
		return true
	}

	if maxWait == 0 {
		return false
	}

	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()
	begin := time.Now()

	for {
		now := <-ticker.C

		if r.Enqueue(value) {
			return true
		}

		if now.Sub(begin) >= maxWait {
			return false
		}
	}
}

// PushWait pushes the value to queue tail with max wait duration.
//
// Deprecated: Use EnqueueWait instead.
func (r *SyncRing[T]) PushWait(value T, maxWait time.Duration) bool {
	return r.EnqueueWait(value, maxWait)
}

// DequeueWait removes and returns the value from queue head with max wait duration.
// If maxWait is negative, it will block until the value is dequeued.
func (r *SyncRing[T]) DequeueWait(maxWait time.Duration) (T, bool) {
	if maxWait < 0 {
		for {
			if v, ok := r.Dequeue(); ok {
				return v, true
			}

			runtime.Gosched()
		}
	}

	if v, ok := r.Dequeue(); ok {
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

		if v, ok := r.Dequeue(); ok {
			return v, true
		}

		if now.Sub(begin) >= maxWait {
			return zero, false
		}
	}
}

// PopWait removes and returns the value from queue head with max wait duration.
//
// Deprecated: Use DequeueWait instead.
func (r *SyncRing[T]) PopWait(maxWait time.Duration) (T, bool) {
	return r.DequeueWait(maxWait)
}

// EnqueueContext attempts to push the value to queue tail until success or ctx is done.
func (r *SyncRing[T]) EnqueueContext(ctx context.Context, value T) bool {
	select {
	case <-ctx.Done():
		return false
	default:
	}

	if r.Enqueue(value) {
		return true
	}

	// Spin a few times before sleeping on ticker
	for i := 0; i < 10; i++ {
		select {
		case <-ctx.Done():
			return false
		default:
		}
		runtime.Gosched()
		if r.Enqueue(value) {
			return true
		}
	}

	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return false
		case <-ticker.C:
			if r.Enqueue(value) {
				return true
			}
		}
	}
}

// PushContext attempts to push the value to queue tail until success or ctx is done.
//
// Deprecated: Use EnqueueContext instead.
func (r *SyncRing[T]) PushContext(ctx context.Context, value T) bool {
	return r.EnqueueContext(ctx, value)
}

// DequeueContext removes and returns the value from queue head until available or ctx is done.
func (r *SyncRing[T]) DequeueContext(ctx context.Context) (T, bool) {
	var zero T
	select {
	case <-ctx.Done():
		return zero, false
	default:
	}

	if v, ok := r.Dequeue(); ok {
		return v, true
	}

	// Spin a few times before sleeping on ticker
	for i := 0; i < 10; i++ {
		select {
		case <-ctx.Done():
			return zero, false
		default:
		}
		runtime.Gosched()
		if v, ok := r.Dequeue(); ok {
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
			if v, ok := r.Dequeue(); ok {
				return v, true
			}
		}
	}
}

// PopContext removes and returns the value from queue head until available or ctx is done.
//
// Deprecated: Use DequeueContext instead.
func (r *SyncRing[T]) PopContext(ctx context.Context) (T, bool) {
	return r.DequeueContext(ctx)
}

func roundUpPowOfTwo(x uint32) uint32 {
	if x >= 1<<31 {
		return 1 << 31
	}
	var pos int
	for i := x; i != 0; pos++ {
		i >>= 1
	}
	return 1 << pos
}
