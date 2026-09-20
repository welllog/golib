package mapz

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

// a panic in a callback must not leak the lock
func TestSafeKV_CallbackPanic(t *testing.T) {
	run := func(name string, trigger func(*SafeKV[string, int])) {
		kv := NewSafeKV[string, int](4)
		kv.Set("a", 1)

		func() {
			defer func() { _ = recover() }()
			trigger(kv)
		}()

		done := make(chan struct{})
		go func() {
			kv.Set("b", 2)
			close(done)
		}()
		select {
		case <-done:
		case <-time.After(2 * time.Second):
			t.Errorf("%s: Set blocked forever after callback panic, lock leaked", name)
		}
	}

	run("Range", func(kv *SafeKV[string, int]) {
		kv.Range(func(string, int) bool { panic("boom") })
	})
	run("Map", func(kv *SafeKV[string, int]) {
		kv.Map(func(KV[string, int]) { panic("boom") })
	})
	run("ReadMap", func(kv *SafeKV[string, int]) {
		kv.ReadMap(func(KV[string, int]) { panic("boom") })
	})
	run("GetWithLock", func(kv *SafeKV[string, int]) {
		kv.GetWithLock("a", func(int) { panic("boom") })
	})
	run("All", func(kv *SafeKV[string, int]) {
		kv.All()(func(string, int) bool { panic("boom") })
	})
	run("SetIf", func(kv *SafeKV[string, int]) {
		kv.SetIf("a", 2, func(int) bool { panic("boom") })
	})
	run("RemoveIf", func(kv *SafeKV[string, int]) {
		kv.RemoveIf("a", func(int) bool { panic("boom") })
	})
}

func TestSafeKV_GetOrSetFunc_WaitErr(t *testing.T) {
	kv := NewSafeKV[string, int](4)

	ready := make(chan struct{})
	block := make(chan struct{})

	var g1Got, g2Got bool
	var g1Err, g2Err error
	var wg sync.WaitGroup

	wg.Add(2)
	go func() {
		defer wg.Done()
		_, g1Got, g1Err = kv.GetOrSetFunc("key", func() (int, error) {
			close(ready)
			<-block
			return 0, fmt.Errorf("load error")
		})
	}()

	<-ready
	go func() {
		defer wg.Done()
		_, g2Got, g2Err = kv.GetOrSetFunc("key", func() (int, error) {
			return 100, nil
		})
	}()

	// 稍微等待确保协程 2 已经进入 s.calls 等待
	time.Sleep(20 * time.Millisecond)
	close(block)

	wg.Wait()

	if g1Got != false {
		t.Errorf("g1 expected got=false, actual got=%v", g1Got)
	}
	if g1Err == nil {
		t.Errorf("g1 expected err, actual nil")
	}

	if g2Got != false {
		t.Errorf("g2 expected got=false on error, actual got=%v", g2Got)
	}
	if g2Err == nil {
		t.Errorf("g2 expected err, actual nil")
	}

	// 确认底层并未存储该 key
	if _, ok := kv.Get("key"); ok {
		t.Errorf("key should not be stored in SafeKV on error")
	}
}
