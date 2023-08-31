package goz

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Limiter struct {
	c            chan struct{}
	w            sync.WaitGroup
	panicHandler func(any)
}

func NewLimiter(limit int) *Limiter {
	if limit < 1 {
		limit = 3
	}
	return &Limiter{
		c: make(chan struct{}, limit),
	}
}

func (l *Limiter) SetPanicHandler(fn func(any)) *Limiter {
	l.panicHandler = fn
	return l
}

func (l *Limiter) Go(fn func()) *Limiter {
	l.add()

	go Recover(fn, l.panicHandler, l.done)

	return l
}

func (l *Limiter) Wait(waitTime ...time.Duration) {
	if len(waitTime) > 0 {
		quit := make(chan struct{}, 1)
		go func(ch chan<- struct{}) {
			l.w.Wait()
			ch <- struct{}{}
		}(quit)

		select {
		case <-quit:
		case <-time.After(waitTime[0]):
		}
		return
	}

	l.w.Wait()
}

func (l *Limiter) add() {
	l.c <- struct{}{}
	l.w.Add(1)
}

func (l *Limiter) done() {
	l.w.Done()
	<-l.c
}

func Recover(fn func(), panicFn func(any), cleanups ...func()) {
	defer func() {
		if p := recover(); p != nil {
			if panicFn != nil {
				panicFn(p)
			} else {
				var buf strings.Builder
				buf.Grow(200)

				buf.WriteString(fmt.Sprintf("panic: %v  Traceback:", p))
				stack(&buf, 4, 6)

				fmt.Println(buf.String())
			}
		}

		for _, cleanup := range cleanups {
			cleanup()
		}
	}()

	fn()
}

func stack(buf *strings.Builder, skip, deep int) {
	callers := make([]uintptr, deep)
	n := runtime.Callers(skip, callers)
	frames := runtime.CallersFrames(callers[:n])
	for {
		frame, more := frames.Next()
		_, _ = buf.WriteString("\n\t")
		_, _ = buf.WriteString(frame.File)
		_, _ = buf.WriteString(":")
		_, _ = buf.WriteString(strconv.Itoa(frame.Line))
		if !more {
			break
		}
	}
}

type Logger interface {
	Error(args ...any)
}

func LogPanic(l Logger, deep int) func(any) {
	return func(a any) {
		var buf strings.Builder
		buf.Grow(200)

		buf.WriteString(fmt.Sprintf("panic: %v  Traceback:", a))
		stack(&buf, 5, deep)

		l.Error(buf.String())
	}
}
