package listz

import (
	"context"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestSyncList_Push(t *testing.T) {
	l := NewSync[int]()
	for i := 1; i <= 10; i++ {
		l.Push(i)
		if l.Len() != i {
			t.Errorf("expected length %d, got %d", i, l.Len())
		}
	}

	var num int
	for {
		n, ok := l.Pop()
		if !ok {
			break
		}
		num++
		if n != num {
			t.Errorf("expected %d, got %d", num, n)
		}
	}
}

func TestSyncList_EnqueueDequeue(t *testing.T) {
	l := NewSync[int]()
	for i := 1; i <= 10; i++ {
		l.Enqueue(i)
		if l.Len() != i {
			t.Errorf("expected length %d, got %d", i, l.Len())
		}
	}

	var num int
	for {
		n, ok := l.Dequeue()
		if !ok {
			break
		}
		num++
		if n != num {
			t.Errorf("expected %d, got %d", num, n)
		}
	}
	if num != 10 {
		t.Errorf("expected 10 items dequeued, got %d", num)
	}
}

func TestSyncList_Push2(t *testing.T) {
	l := NewSync[int]()

	maxNum := 10000
	s := make([]uint32, maxNum)
	var w sync.WaitGroup

	w.Add(maxNum * 2)
	for i := 0; i < maxNum; i++ {
		go func(n int) {
			l.Enqueue(n)
			w.Done()
		}(i)
	}

	for i := 0; i < maxNum; i++ {
		go func() {
			for {
				n, ok := l.Dequeue()
				if ok {
					atomic.AddUint32(&s[n], 1)
					w.Done()
					break
				}

				runtime.Gosched()
			}
		}()
	}

	w.Wait()
	for _, v := range s {
		if v != 1 {
			t.Fatalf("Expected value 1, got %d", v)
		}
	}

	if l.Len() != 0 {
		t.Errorf("expected length 0, got %d", l.Len())
	}
}

func TestSyncList_Pop(t *testing.T) {
	l := NewSync[int]()
	_, ok := l.Dequeue()
	if ok {
		t.Errorf("expected false, got true")
	}

	c := runtime.GOMAXPROCS(0)
	seg := 100000
	s := make([]uint32, c*seg)
	var w sync.WaitGroup
	w.Add(2 * c)
	for i := 0; i < c; i++ {
		go func(n int) {
			for i := n * seg; i < (n+1)*seg; i++ {
				l.Enqueue(i)
			}
			w.Done()
		}(i)
	}

	for i := 0; i < c; i++ {
		go func() {
			var count int
			for {
				n, ok := l.Dequeue()
				if ok {
					atomic.AddUint32(&s[n], 1)
					count++
					if count == seg/2 {
						break
					}
				}

				runtime.Gosched()
			}
			w.Done()
		}()
	}

	w.Wait()

	if l.Len() != seg/2*c {
		t.Errorf("expected length %d, got %d", seg/2*c, l.Len())
	}

	w.Add(c)
	for i := 0; i < c; i++ {
		go func() {
			for {
				if l.Len() == 0 {
					break
				}

				n, ok := l.Dequeue()
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

	for i := 0; i < c*seg; i++ {
		if s[i] != 1 {
			t.Errorf("index %d expected 1, got %d", i, s[i])
		}
	}
}

func TestSyncList_DequeueWait(t *testing.T) {
	l := NewSync[int]()

	// PopWait on empty with 0 timeout
	val, ok := l.PopWait(0)
	if ok || val != 0 {
		t.Errorf("expected false and 0 on immediate timeout, got %v, %v", val, ok)
	}

	// DequeueWait on empty with 20ms timeout
	start := time.Now()
	val, ok = l.DequeueWait(20 * time.Millisecond)
	elapsed := time.Since(start)
	if ok || val != 0 {
		t.Errorf("expected false on timeout, got %v, %v", val, ok)
	}
	if elapsed < 20*time.Millisecond {
		t.Errorf("expected to wait at least 20ms, elapsed: %v", elapsed)
	}

	// Producer-consumer via DequeueWait
	go func() {
		time.Sleep(20 * time.Millisecond)
		l.Enqueue(42)
	}()

	val, ok = l.DequeueWait(200 * time.Millisecond)
	if !ok || val != 42 {
		t.Errorf("expected 42 and true, got %v, %v", val, ok)
	}
}

func TestSyncList_Context(t *testing.T) {
	t.Run("EnqueueContext and PushContext canceled", func(t *testing.T) {
		l := NewSync[int]()
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // cancel immediately

		ok := l.EnqueueContext(ctx, 1)
		if ok {
			t.Error("expected EnqueueContext to return false on canceled ctx")
		}

		ok = l.PushContext(ctx, 2)
		if ok {
			t.Error("expected PushContext to return false on canceled ctx")
		}
	})

	t.Run("EnqueueContext success", func(t *testing.T) {
		l := NewSync[int]()
		ctx := context.Background()

		if !l.EnqueueContext(ctx, 100) {
			t.Error("expected EnqueueContext to succeed")
		}
		if !l.PushContext(ctx, 200) {
			t.Error("expected PushContext to succeed")
		}
		if l.Len() != 2 {
			t.Errorf("expected len 2, got %d", l.Len())
		}
	})

	t.Run("DequeueContext and PopContext canceled immediately", func(t *testing.T) {
		l := NewSync[int]()
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		val, ok := l.DequeueContext(ctx)
		if ok || val != 0 {
			t.Errorf("expected false and zero, got %v, %v", val, ok)
		}

		val, ok = l.PopContext(ctx)
		if ok || val != 0 {
			t.Errorf("expected false and zero, got %v, %v", val, ok)
		}
	})

	t.Run("DequeueContext timeout", func(t *testing.T) {
		l := NewSync[int]()
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
		defer cancel()

		start := time.Now()
		val, ok := l.DequeueContext(ctx)
		elapsed := time.Since(start)

		if ok || val != 0 {
			t.Errorf("expected false on timeout, got %v, %v", val, ok)
		}
		if elapsed < 15*time.Millisecond {
			t.Errorf("expected elapsed >= 15ms, got %v", elapsed)
		}
	})

	t.Run("DequeueContext and PopContext receive value", func(t *testing.T) {
		l := NewSync[int]()
		ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
		defer cancel()

		go func() {
			time.Sleep(30 * time.Millisecond)
			l.Enqueue(999)
			l.Push(888)
		}()

		val, ok := l.DequeueContext(ctx)
		if !ok || val != 999 {
			t.Errorf("expected 999, true, got %v, %v", val, ok)
		}

		val, ok = l.PopContext(ctx)
		if !ok || val != 888 {
			t.Errorf("expected 888, true, got %v, %v", val, ok)
		}
	})
}

func BenchmarkSyncList(b *testing.B) {
	l := NewSync[int]()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			l.Enqueue(1)
			l.DequeueWait(-1)
		}
	})
}
