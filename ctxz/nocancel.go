package ctxz

import (
	"context"
	"time"
)

// noCancelContext is a context that never cancels.
type noCancelContext struct {
	ctx context.Context
}

func (n noCancelContext) Deadline() (deadline time.Time, ok bool) {
	return time.Time{}, false
}

func (n noCancelContext) Done() <-chan struct{} {
	return nil
}

func (n noCancelContext) Err() error {
	return nil
}

func (n noCancelContext) Value(key any) any {
	return n.ctx.Value(key)
}

// WithoutCancel returns a copy of the parent context that never cancels.
func WithoutCancel(ctx context.Context) (valueOnlyContext context.Context) {
	if ctx == nil {
		panic("cannot create context from nil parent")
	}
	if ctx == context.Background() {
		return ctx
	}
	return noCancelContext{ctx}
}
