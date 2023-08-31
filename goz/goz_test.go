package goz

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestRecover(t *testing.T) {
	var w sync.WaitGroup
	w.Add(2)
	go Recover(func() {
		panic("test1")
	}, nil, w.Done)

	go Recover(func() {
		panic("test2")
	}, func(a any) {
		t.Log(a)
	}, w.Done)

	w.Wait()
}

func TestLimiter_Go(t *testing.T) {
	var n int32

	limiter := NewLimiter(2).SetPanicHandler(func(a any) {
		t.Log(a)
	})

	now := time.Now()
	limiter.Go(func() {
		atomic.AddInt32(&n, 1)
		time.Sleep(time.Second)
	}).Go(func() {
		atomic.AddInt32(&n, 1)
		time.Sleep(time.Second)
		panic("test panic")
	}).Go(func() {
		atomic.AddInt32(&n, 1)
		if time.Since(now).Milliseconds() < 1000 {
			t.Fatal("time should more than 1s")
		}
		time.Sleep(time.Second)
	}).Go(func() {
		atomic.AddInt32(&n, 1)
		time.Sleep(time.Second)
	}).Go(func() {
		atomic.AddInt32(&n, 1)
		if time.Since(now).Milliseconds() < 2000 {
			t.Fatal("time should more than 2s")
		}
	})

	limiter.Wait()
	if n != 5 {
		t.Error("n should be 5")
	}
}

func TestLimiter_Go2(t *testing.T) {
	limiter := NewLimiter(0)
	limiter.Go(func() {
		panic(1)
	}).Go(func() {
		panic(2)
	}).Go(func() {
		panic1(3)
	}).Wait()
}

func TestLimiter_Wait(t *testing.T) {
	now := time.Now()
	NewLimiter(1).Go(func() {
		time.Sleep(time.Second)
	}).Wait(500 * time.Millisecond)

	since := time.Since(now).Milliseconds()
	if since > 1000 {
		t.Fatal("wait must less than 1000")
	}
	t.Logf("wait %d ms", since)
}

func TestLogPanic(t *testing.T) {
	f := LogPanic(logger{}, 10)

	var w sync.WaitGroup
	w.Add(1)
	go Recover(func() {
		panic1(1)
	}, f, w.Done)

	NewLimiter(1).SetPanicHandler(f).Go(func() {
		panic1(2)
	}).Wait()

	w.Wait()
}

type logger struct{}

func (l logger) Error(args ...any) {
	fmt.Println(args...)
}

func panic1(n int) {
	panic2(n)
}

func panic2(n int) {
	panic3(n)
}

func panic3(n int) {
	panic(n)
}
