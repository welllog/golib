package ringz

import (
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestSyncRing_Len(t *testing.T) {
	for c := 1; c <= 100; c++ {
		r := NewSync[int](c)
		testRing(&r, r.Cap(), t)
	}
}

func TestSyncRing_IsFull(t *testing.T) {
	// Test empty queue
	q := NewSync[int](8)
	if q.IsFull() {
		t.Errorf("Expected IsFull to return false, but it returned true")
	}

	// Test partially full queue
	q.Push(1)
	q.Push(2)
	if q.IsFull() {
		t.Errorf("Expected IsFull to return false, but it returned true")
	}

	// Test full queue
	for i := 3; i <= 8; i++ {
		q.Push(i)
	}
	if !q.IsFull() {
		t.Errorf("Expected IsFull to return true, but it returned false")
	}

	for i := uint32(0); i < 8; i++ {
		q.Pop()
	}
	if q.IsFull() {
		t.Errorf("Expected IsFull to return false, but it returned true")
	}

	// Test concurrent access
	var wg sync.WaitGroup
	var expected uint32 = 1000
	for i := uint32(0); i < expected; i++ {
		wg.Add(1)
		go func() {
			q.Push(0)
			wg.Done()
		}()
	}
	wg.Wait()
	if !q.IsFull() {
		t.Errorf("Expected IsFull to return true, but it returned false")
	}

	// Test concurrent access with removal and addition
	for i := uint32(0); i < expected; i++ {
		wg.Add(1)
		go func() {
			q.Push(0)
			q.Pop()
			wg.Done()
		}()
	}
	wg.Wait()
	if q.IsFull() {
		t.Errorf("Expected IsFull to return false, but it returned true")
	}
}

func TestSyncRing_IsEmpty(t *testing.T) {
	q := NewSync[int](8)

	// Test empty queue
	if !q.IsEmpty() {
		t.Errorf("Expected queue to be empty, but it is not")
	}

	// Test single element
	q.Push(1)
	if q.IsEmpty() {
		t.Errorf("Expected queue to not be empty, but it is")
	}

	// Test multiple elements
	q.Push(2)
	q.Push(3)
	if q.IsEmpty() {
		t.Errorf("Expected queue to not be empty, but it is")
	}

	// Test concurrent access
	var wg sync.WaitGroup
	var expected = 1000
	for i := 0; i < expected; i++ {
		wg.Add(1)
		go func(i int) {
			q.Push(i)
			wg.Done()
		}(i)
	}
	wg.Wait()
	if q.IsEmpty() {
		t.Errorf("Expected queue to not be empty, but it is")
	}

	// Test concurrent access with removal
	for i := uint32(0); i < uint32(expected); i++ {
		wg.Add(1)
		go func() {
			_, _ = q.Pop()
			wg.Done()
		}()
	}
	wg.Wait()
	if !q.IsEmpty() {
		t.Errorf("Expected queue to be empty, but it is not")
	}
}

func TestSyncRing_PushAndPop(t *testing.T) {
	maxNum := 4000
	q := NewSync[int](100)
	s := make([]uint32, maxNum)

	var wg sync.WaitGroup
	wg.Add(maxNum * 2)

	begin := time.Now()
	for i := 0; i < maxNum; i++ {
		go func(n int) {
			for {
				if q.Push(n) {
					break
				} else {
					runtime.Gosched()
				}
			}
			wg.Done()
		}(i)
	}

	for i := 0; i < maxNum; i++ {
		go func() {
			for {
				v, ok := q.Pop()
				if ok {
					atomic.AddUint32(&s[v], 1)
					break
				} else {
					runtime.Gosched()
				}
			}
			wg.Done()
		}()
	}

	wg.Wait()
	t.Logf("Time: %v", time.Since(begin))

	if q.Len() != 0 {
		t.Fatalf("Expected queue to be empty, but it is not")
	}

	if q.IsFull() {
		t.Fatalf("Expected queue to not be full, but it is")
	}

	if !q.IsEmpty() {
		t.Fatalf("Expected queue to be empty, but it is not")
	}

	for i := 0; i < maxNum; i++ {
		if s[i] != 1 {
			t.Fatalf("Expected value %d to be in the queue, but it is not", i)
		}
	}
}

func TestSyncRing_PushWaitAndPopWait(t *testing.T) {
	maxNum := 4000
	q := NewSync[int](100)
	s := make([]uint32, maxNum)

	var wg sync.WaitGroup
	wg.Add(maxNum * 2)

	begin := time.Now()
	for i := 0; i < maxNum; i++ {
		go func(n int) {
			q.PushWait(n, time.Hour)
			wg.Done()
		}(i)
	}

	for i := 0; i < maxNum; i++ {
		go func() {
			v, ok := q.PopWait(time.Hour)
			if ok {
				atomic.AddUint32(&s[v], 1)
			}
			wg.Done()
		}()
	}

	wg.Wait()
	t.Logf("Time: %v", time.Since(begin))

	if q.Len() != 0 {
		t.Fatalf("Expected queue to be empty, but it is not")
	}

	if q.IsFull() {
		t.Fatalf("Expected queue to not be full, but it is")
	}

	if !q.IsEmpty() {
		t.Fatalf("Expected queue to be empty, but it is not")
	}

	for i := 0; i < maxNum; i++ {
		if s[i] != 1 {
			t.Fatalf("Expected value %d to be in the queue, but it is not", i)
		}
	}
}

func TestSyncRing_PopWait(t *testing.T) {
	c := runtime.GOMAXPROCS(0)

	r := NewSync[int](c * 1000)
	s := make([]uint32, c*1000)

	var w sync.WaitGroup
	w.Add(2 * c)
	for i := 0; i < c; i++ {
		go func(n int) {
			for i := n * 1000; i < (n+1)*1000; i++ {
				r.PushWait(i, -1)
			}
			w.Done()
		}(i)
	}

	for i := 0; i < c; i++ {
		go func() {
			var count int
			for {
				n, _ := r.PopWait(-1)
				atomic.AddUint32(&s[n], 1)

				count++
				if count == 500 {
					break
				}
			}
			w.Done()
		}()
	}

	w.Wait()

	if r.Len() != 500*c {
		t.Errorf("expected length %d, got %d", 500*c, r.Len())
	}

	w.Add(c)
	for i := 0; i < c; i++ {
		go func() {
			for {
				n, ok := r.PopWait(3 * time.Second)
				if !ok {
					break
				}

				atomic.AddUint32(&s[n], 1)
			}
			w.Done()
		}()
	}
	w.Wait()

	for i := 0; i < c*1000; i++ {
		if s[i] != 1 {
			t.Errorf("expected 1, got %d", s[i])
		}
	}
}

type mutexRing[T any] struct {
	Ring[T]
	mu sync.Mutex
}

func newMutexRing[T any](cap int) *mutexRing[T] {
	var r mutexRing[T]
	r.Init(cap)
	return &r
}

func (r *mutexRing[T]) Push(value T) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.Ring.Push(value)
}

func (r *mutexRing[T]) Pop() (T, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()

	return r.Ring.Pop()
}

func BenchmarkSyncRing(b *testing.B) {
	b.Run("SyncRing", func(b *testing.B) {
		q := NewSync[int](8)
		b.ResetTimer()
		b.ReportAllocs()

		for i := 0; i < b.N; i++ {
			q.Push(1)
			q.Pop()
		}
	})

	b.Run("SyncRing-Parallel", func(b *testing.B) {
		q := NewSync[int](8)
		b.ResetTimer()
		b.ReportAllocs()

		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				q.PushWait(1, -1)
				q.PopWait(-1)
			}
		})
	})

	b.Run("MutexRing", func(b *testing.B) {
		q := newMutexRing[int](8)
		b.ResetTimer()
		b.ReportAllocs()

		for i := 0; i < b.N; i++ {
			q.Push(1)
			q.Pop()
		}
	})

	b.Run("MutexRing-Parallel", func(b *testing.B) {
		q := newMutexRing[int](8)
		b.ResetTimer()
		b.ReportAllocs()

		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				q.Push(1)
				q.Pop()
			}
		})
	})

	b.Run("Channel", func(b *testing.B) {
		q := make(chan int, 8)
		b.ResetTimer()
		b.ReportAllocs()

		for i := 0; i < b.N; i++ {
			q <- 1
			<-q
		}
	})

	b.Run("Channel-Parallel", func(b *testing.B) {
		q := make(chan int, 8)
		b.ResetTimer()
		b.ReportAllocs()

		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				q <- 1
				<-q
			}
		})
	})
}
