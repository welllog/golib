package goz

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	defaultLimit     = 3
	defaultStackDeep = 32
)

type Limiter struct {
	c            chan struct{}
	w            sync.WaitGroup
	panicHandler func(any)
}

// NewLimiter creates a new Limiter with the specified limit of concurrent goroutines.
// if limit is 0, it defaults to 3.
// if limit is less than 0, it means no limit.
func NewLimiter(limit int) *Limiter {
	if limit < 0 {
		return &Limiter{}
	}

	if limit == 0 {
		limit = defaultLimit
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
	if l.c == nil {
		l.w.Add(1)
		go Recover(fn, l.panicHandler, l.w.Done)
		return l
	}

	l.add()
	go Recover(fn, l.panicHandler, l.done)
	return l
}

func (l *Limiter) Done() <-chan struct{} {
	quit := make(chan struct{})
	go func(ch chan<- struct{}) {
		l.w.Wait()
		close(ch)
	}(quit)
	return quit
}

func (l *Limiter) Wait(waitTime ...time.Duration) {
	if len(waitTime) > 0 {
		select {
		case <-l.Done():
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
				buf.Grow(1024)

				buf.WriteString(fmt.Sprintf("panic: %v  Traceback:", p))
				stack(&buf, 4, defaultStackDeep)

				fmt.Println(buf.String())
			}
		}

		if len(cleanups) == 0 {
			return
		}

		var index int
		defer func() {
			if p := recover(); p != nil {
				s := fmt.Sprintf("cleanup panic: %v, index: %d", p, index)
				if panicFn != nil {
					panicFn(s)
				} else {
					fmt.Println(s)
				}
			}
		}()

		for i, cleanup := range cleanups {
			index = i
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
		buf.Grow(512)

		buf.WriteString(fmt.Sprintf("panic: %v  Traceback:", a))
		stack(&buf, 5, deep)

		l.Error(buf.String())
	}
}
