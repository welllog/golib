package listz

import (
	"runtime"
	"sync"
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
	s := make([]uint8, maxNum)
	var w sync.WaitGroup

	w.Add(maxNum * 2)
	for i := 0; i < maxNum; i++ {
		go func(n int) {
			l.Push(n)
			// l.push(n)
			w.Done()
		}(i)
	}

	for i := 0; i < maxNum; i++ {
		go func() {
			for {
				n, ok := l.Pop()
				// n, ok := l.pop()
				if ok {
					s[n] = 1
					w.Done()
					break
				}

				runtime.Gosched()
			}
		}()
	}

	w.Wait()
	for _, v := range s {
		if v == 0 {
			t.Errorf("expected 1, got 0")
		}
	}

	if l.Len() != 0 {
		t.Errorf("expected length 0, got %d", l.Len())
	}
}

func BenchmarkSyncList(b *testing.B) {
	l := NewSync[int]()
	// b.Run("std-push", func(b *testing.B) {
	// 	b.ResetTimer()
	// 	b.ReportAllocs()
	// 	b.RunParallel(func(pb *testing.PB) {
	// 		for pb.Next() {
	// 			l.Push(1)
	// 			for {
	// 				_, ok := l.Pop()
	// 				if ok {
	// 					break
	// 				}
	// 				runtime.Gosched()
	// 			}
	// 		}
	// 	})
	// })

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

func BenchmarkPool(b *testing.B) {
	l := NewSync[int]()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			n := l.getNode(1)
			l.releaseNode(n)
		}
	})
}
