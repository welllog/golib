package ctxz

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestWithoutCancel(t *testing.T) {
	ctx := WithoutCancel(context.Background())
	select {
	case <-ctx.Done():
		t.Fatal("ctx should not done")
	case <-time.After(10 * time.Millisecond):
	}
}

func TestWithoutCancel2(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	ctx = WithoutCancel(ctx)
	select {
	case <-ctx.Done():
		t.Fatal("ctx should not done")
	case <-time.After(10 * time.Millisecond):
	}
}

func TestWithoutCancel3(t *testing.T) {
	ctx1, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	ctx2 := WithoutCancel(ctx1)
	time.Sleep(15 * time.Millisecond)

	select {
	case <-ctx2.Done():
		t.Fatal("ctx should not done")
	case <-time.After(20 * time.Millisecond):
	}

	if errors.Is(ctx1.Err(), ctx2.Err()) {
		t.Fatal("ctx1.Err() should not equal ctx2.Err()")
	}

	if ctx2.Err() != nil {
		t.Fatal("ctx2.Err() should be nil")
	}
}
