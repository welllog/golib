package listz

import (
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
)

func TestSyncList_Push(t *testing.T) {
	l := NewSync[int]()
	for i := 1; i <= 10; i++ {
		// l.Push(i)
		l.push(i)
		if l.Len() != i {
			t.Errorf("expected length %d, got %d", i, l.Len())
		}
	}

	var num int
	for {
		// n, ok := l.Pop()
		n, ok := l.pop()
		if !ok {
			break
		}
		num++
		if n != num {
			t.Errorf("expected %d, got %d", num, n)
		}
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
			// l.Push(n)
			l.push(n)
			w.Done()
		}(i)
	}

	for i := 0; i < maxNum; i++ {
		go func() {
			for {
				// n, ok := l.Pop()
				n, ok := l.pop()
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
	_, ok := l.Pop()
	if ok {
		t.Errorf("expected false, got true")
	}

	c := runtime.GOMAXPROCS(0)
	seg := 1000000
	s := make([]uint32, c*seg)
	var w sync.WaitGroup
	w.Add(2 * c)
	for i := 0; i < c; i++ {
		go func(n int) {
			for i := n * seg; i < (n+1)*seg; i++ {
				l.push(i)
			}
			w.Done()
		}(i)
	}

	for i := 0; i < c; i++ {
		go func() {
			var count int
			for {
				n, ok := l.pop()
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

				n, ok := l.pop()
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

func TestSyncList_Pop2(t *testing.T) {
	l := NewSync[int]()
	_, ok := l.Pop()
	if ok {
		t.Errorf("expected false, got true")
	}

	c := runtime.GOMAXPROCS(0)
	seg := 10000
	s := make([]uint32, c*seg)
	var w sync.WaitGroup
	w.Add(2 * c)
	for i := 0; i < c; i++ {
		go func(n int) {
			for i := n * seg; i < (n+1)*seg; i++ {
				l.Push(i)
			}
			w.Done()
		}(i)
	}

	for i := 0; i < c; i++ {
		go func() {
			var count int
			for {
				n, ok := l.Pop()
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

				n, ok := l.Pop()
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

func BenchmarkSyncList(b *testing.B) {
	l := NewSync[int]()
	b.Run("std-push", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				l.Push(1)
				for {
					_, ok := l.Pop()
					if ok {
						break
					}
					runtime.Gosched()
				}
			}
		})
	})

	b.Run("push", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				l.push(1)
				for {
					_, ok := l.pop()
					if ok {
						break
					}
					runtime.Gosched()
				}
			}
		})
	})
}
