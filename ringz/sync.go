package ringz

import (
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

func NewSync[T any](cap int) SyncRing[T] {
	var r SyncRing[T]
	r.Init(cap)
	return r
}

func (r *SyncRing[T]) Init(cap int) {
	var c uint32
	switch {
	case cap <= 0:
		panic("ringz.SyncRing Init: invalid capacity: " + strconv.Itoa(cap))
	case 1 == cap:
		c = 2
	default:
		c = uint32(cap)
		if c&(c-1) > 0 {
			c = roundupPowOfTwo(c)
		}
	}

	r.cap = c
	r.mask = c - 1
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

// Push pushes the value to queue tail.
// Notice if return false, means the queue is full or concurrent Push operation.
func (r *SyncRing[T]) Push(value T) bool {
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
	// atomic.AddUint32(&holder.pos, 1)
	atomic.StoreUint32(&holder.pos, seq+1)
	return true
}

// Pop removes and returns the value from queue head.
// Notice if return false, means the queue is empty or concurrent Pop operation.
func (r *SyncRing[T]) Pop() (T, bool) {
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
	// atomic.AddUint32(&holder.pos, r.mask)
	atomic.StoreUint32(&holder.pos, seq+r.mask)
	return value, true
}

// PushWait pushes the value to queue tail with max wait duration.
// If maxWait is negative, it will block until the value is pushed.
func (r *SyncRing[T]) PushWait(value T, maxWait time.Duration) bool {
	if maxWait < 0 {
		for {
			if r.Push(value) {
				return true
			}

			runtime.Gosched()
		}
	}

	if r.Push(value) {
		return true
	}

	if maxWait == 0 {
		return false
	}

	ticker := time.NewTicker(10 * time.Millisecond)
	begin := time.Now()

	for {
		now := <-ticker.C

		if r.Push(value) {
			ticker.Stop()
			return true
		}

		if now.Sub(begin) >= maxWait {
			ticker.Stop()
			return false
		}
	}
}

// PopWait removes and returns the value from queue head with max wait duration.
// If maxWait is negative, it will block until the value is popped.
func (r *SyncRing[T]) PopWait(maxWait time.Duration) (T, bool) {
	if maxWait < 0 {
		for {
			if v, ok := r.Pop(); ok {
				return v, true
			}

			runtime.Gosched()
		}
	}

	if v, ok := r.Pop(); ok {
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

		if v, ok := r.Pop(); ok {
			ticker.Stop()
			return v, true
		}

		if now.Sub(begin) >= maxWait {
			ticker.Stop()
			return zero, false
		}
	}
}

func roundupPowOfTwo(x uint32) uint32 {
	var pos int
	for i := x; i != 0; pos++ {
		i >>= 1
	}
	return 1 << pos
}
