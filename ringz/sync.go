package ringz

import (
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/welllog/golib/mathz"
)

// cache line usually is 64 bytes, but Apple Silicon, a.k.a. M1, has 128-byte cache line size.
const cacheLinePadSize = unsafe.Sizeof(uint64(0)) * 16

type item[T any] struct {
	value T
	pos   uint32
}

type SyncRing[T any] struct {
	values []item[T]
	cap    uint32
	mask   uint32
	// _      [cacheLinePadSize - unsafe.Sizeof(uint64(0))*3 - 8]byte
	head uint32
	// _      [cacheLinePadSize - 4]byte
	tail uint32
	// _      [cacheLinePadSize - 4]byte
}

func NewSync[T any](cap int) SyncRing[T] {
	var r SyncRing[T]
	r.Init(cap)
	return r
}

func (r *SyncRing[T]) Init(cap int) {
	c := uint32(cap)
	if c&(c-1) > 0 {
		c = roundupPowOfTwo(c)
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
	return mathz.Max(int(atomic.LoadUint32(&r.tail)-atomic.LoadUint32(&r.head)), 0)
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
	atomic.AddUint32(&holder.pos, 1)
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
	atomic.AddUint32(&holder.pos, r.mask)
	return value, true
}

// Peek returns the value from queue head without removing it.
// Notice if return false, means the queue is empty or concurrent Pop operation.
func (r *SyncRing[T]) Peek() (T, bool) {
	pos := atomic.LoadUint32(&r.head)
	holder := &r.values[pos&r.mask]
	seq := atomic.LoadUint32(&holder.pos)

	var zero T
	if pos+1 != seq {
		return zero, false
	}

	return holder.value, true
}

// PushWait pushes the value to queue tail with max wait duration.
func (r *SyncRing[T]) PushWait(value T, maxWait time.Duration) bool {
	if r.Push(value) {
		return true
	}

	if maxWait == 0 {
		maxWait = 500 * time.Millisecond
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
func (r *SyncRing[T]) PopWait(maxWait time.Duration) (T, bool) {
	if v, ok := r.Pop(); ok {
		return v, true
	}

	if maxWait == 0 {
		maxWait = 500 * time.Millisecond
	}

	ticker := time.NewTicker(10 * time.Millisecond)
	begin := time.Now()

	for {
		if v, ok := r.Pop(); ok {
			ticker.Stop()
			return v, true
		}

		now := <-ticker.C

		if now.Sub(begin) >= maxWait {
			ticker.Stop()
			var zero T
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
