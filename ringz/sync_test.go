package ringz

import (
	"context"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestSyncRing_InvalidCapacity(t *testing.T) {
	for _, invalidCap := range []int{0, -1, -100, (1 << 30) + 1, 1 << 31, 2000000000 + 1000000000} {
		func() {
			defer func() {
				if r := recover(); r == nil {
					t.Errorf("expected panic for cap %d, but did not panic", invalidCap)
				}
			}()
			NewSync[int](invalidCap)
		}()
	}
}

func TestRoundUpPowOfTwo(t *testing.T) {
	tests := []struct {
		input    uint32
		expected uint32
	}{
		{1, 2},
		{2, 4},
		{3, 4},
		{5, 8},
		{8, 16},
		{100, 128},
		{1024, 2048},
		{(1 << 30) - 1, 1 << 30},
		{1 << 31, 1 << 31},
		{(1 << 31) + 1, 1 << 31},
	}
	for _, tt := range tests {
		actual := roundUpPowOfTwo(tt.input)
		if actual != tt.expected {
			t.Errorf("roundUpPowOfTwo(%d) = %d, expected %d", tt.input, actual, tt.expected)
		}
	}
}

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

	seg := 10000
	r := NewSync[int](c * seg)
	s := make([]uint32, c*seg)

	var w sync.WaitGroup
	w.Add(2 * c)
	for i := 0; i < c; i++ {
		go func(n int) {
			for i := n * seg; i < (n+1)*seg; i++ {
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
				if count == seg/2 {
					break
				}
			}
			w.Done()
		}()
	}

	w.Wait()

	if r.Len() != seg/2*c {
		t.Errorf("expected length %d, got %d", seg/2*c, r.Len())
	}

	w.Add(c)
	for i := 0; i < c; i++ {
		go func() {
			for {
				if r.Len() == 0 {
					break
				}

				n, ok := r.Pop()
				if ok {
					atomic.AddUint32(&s[n], 1)
				} else {
					runtime.Gosched()
				}
			}
			w.Done()
		}()
	}
	w.Wait()

	if r.Len() != 0 {
		t.Errorf("expected length 0, got %d", r.Len())
	}

	for i := 0; i < c*seg; i++ {
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

func TestSyncRing_Init_Reuse(t *testing.T) {
	r := NewSync[int](4)
	r.Push(1)
	r.Push(2)
	r.Init(4)
	if !r.IsEmpty() {
		t.Fatalf("expected ring to be empty after Init, got IsEmpty() = false, Len() = %d", r.Len())
	}
	if r.Len() != 0 {
		t.Fatalf("expected ring Len() == 0, got %d", r.Len())
	}
	if !r.Push(10) {
		t.Fatalf("failed to push after Init")
	}
	val, ok := r.Pop()
	if !ok || val != 10 {
		t.Fatalf("failed to pop after Init: got val=%v, ok=%v", val, ok)
	}
}

func TestSyncRing_EnqueueDequeue(t *testing.T) {
	r := NewSync[int](4)
	if !r.Enqueue(10) {
		t.Error("expected Enqueue(10) to succeed")
	}
	if !r.Push(20) {
		t.Error("expected Push(20) to succeed")
	}
	if r.Len() != 2 {
		t.Errorf("expected len 2, got %d", r.Len())
	}

	val, ok := r.Dequeue()
	if !ok || val != 10 {
		t.Errorf("expected (10, true), got (%v, %v)", val, ok)
	}
	val, ok = r.Pop()
	if !ok || val != 20 {
		t.Errorf("expected (20, true), got (%v, %v)", val, ok)
	}
}

func TestSyncRing_Wait(t *testing.T) {
	r := NewSync[int](2) // capacity becomes 2
	r.Enqueue(1)
	r.Enqueue(2)

	// Full: EnqueueWait with 20ms timeout should fail
	start := time.Now()
	ok := r.EnqueueWait(3, 20*time.Millisecond)
	if ok {
		t.Error("expected EnqueueWait to fail on full ring")
	}
	if time.Since(start) < 15*time.Millisecond {
		t.Error("expected to wait for timeout")
	}

	// DequeueWait on consumer
	go func() {
		time.Sleep(20 * time.Millisecond)
		r.Dequeue()
	}()

	ok = r.EnqueueWait(3, 200*time.Millisecond)
	if !ok {
		t.Error("expected EnqueueWait to succeed after Dequeue")
	}
}

func TestSyncRing_Context(t *testing.T) {
	t.Run("EnqueueContext and PushContext canceled", func(t *testing.T) {
		r := NewSync[int](2)
		r.Enqueue(1)
		r.Enqueue(2) // full

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		if r.EnqueueContext(ctx, 3) {
			t.Error("expected EnqueueContext to return false on canceled ctx")
		}
		if r.PushContext(ctx, 3) {
			t.Error("expected PushContext to return false on canceled ctx")
		}
	})

	t.Run("DequeueContext and PopContext canceled", func(t *testing.T) {
		r := NewSync[int](2) // empty

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		val, ok := r.DequeueContext(ctx)
		if ok || val != 0 {
			t.Errorf("expected false on canceled ctx, got %v, %v", val, ok)
		}
		val, ok = r.PopContext(ctx)
		if ok || val != 0 {
			t.Errorf("expected false on canceled ctx, got %v, %v", val, ok)
		}
	})

	t.Run("DequeueContext success from concurrent producer", func(t *testing.T) {
		r := NewSync[int](4)
		ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
		defer cancel()

		go func() {
			time.Sleep(25 * time.Millisecond)
			r.Enqueue(777)
			r.Push(888)
		}()

		val, ok := r.DequeueContext(ctx)
		if !ok || val != 777 {
			t.Errorf("expected (777, true), got (%v, %v)", val, ok)
		}
		val, ok = r.PopContext(ctx)
		if !ok || val != 888 {
			t.Errorf("expected (888, true), got (%v, %v)", val, ok)
		}
	})
}
