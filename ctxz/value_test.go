package ctxz

import (
	"context"
	"testing"
)

func TestBool(t *testing.T) {
	zeroValues := []any{
		float64(0), float32(0),
		int(0), int8(0), int16(0), int32(0), int64(0),
		uint(0), uint8(0), uint16(0), uint32(0), uint64(0),
	}
	for _, v := range zeroValues {
		ctx := context.WithValue(context.Background(), "k", v)
		if Bool(ctx, "k") {
			t.Errorf("Bool(key=%T(%v)) = true, want false", v, v)
		}
	}

	nonZeroValues := []any{
		float64(1), float32(1),
		int(1), int8(1), int16(1), int32(1), int64(1),
		uint(1), uint8(1), uint16(1), uint32(1), uint64(1),
	}
	for _, v := range nonZeroValues {
		ctx := context.WithValue(context.Background(), "k", v)
		if !Bool(ctx, "k") {
			t.Errorf("Bool(key=%T(%v)) = false, want true", v, v)
		}
	}

	tests := []struct {
		value any
		want  bool
	}{
		{nil, false},
		{true, true},
		{false, false},
		{"true", true},
		{"false", false},
		{"1", true},
		{"0", false},
		{"x", false},
	}
	for _, tt := range tests {
		ctx := context.WithValue(context.Background(), "k", tt.value)
		if got := Bool(ctx, "k"); got != tt.want {
			t.Errorf("Bool(key=%#v) = %v, want %v", tt.value, got, tt.want)
		}
	}
}
