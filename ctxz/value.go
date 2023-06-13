package ctxz

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/welllog/golib/strz"
)

func String(ctx context.Context, key any) string {
	value := ctx.Value(key)
	return strz.ToString(value)
}

func Int(ctx context.Context, key any) int {
	value := ctx.Value(key)
	if value == nil {
		return 0
	}
	switch v := value.(type) {
	case int:
		return v
	case int8:
		return int(v)
	case int16:
		return int(v)
	case int32:
		return int(v)
	case int64:
		return int(v)
	case uint:
		return int(v)
	case uint8:
		return int(v)
	case uint16:
		return int(v)
	case uint32:
		return int(v)
	case uint64:
		return int(v)
	case float32:
		return int(v)
	case float64:
		return int(v)
	case json.Number:
		i, _ := v.Int64()
		return int(i)
	case string:
		i, _ := strconv.ParseInt(v, 10, 64)
		return int(i)
	default:
		return 0
	}
}

func Int64(ctx context.Context, key any) int64 {
	value := ctx.Value(key)
	if value == nil {
		return 0
	}
	switch v := value.(type) {
	case int64:
		return v
	case int:
		return int64(v)
	case int8:
		return int64(v)
	case int16:
		return int64(v)
	case int32:
		return int64(v)
	case uint:
		return int64(v)
	case uint8:
		return int64(v)
	case uint16:
		return int64(v)
	case uint32:
		return int64(v)
	case uint64:
		return int64(v)
	case float32:
		return int64(v)
	case float64:
		return int64(v)
	case json.Number:
		i, _ := v.Int64()
		return i
	case string:
		i, _ := strconv.ParseInt(v, 10, 64)
		return i
	default:
		return 0
	}
}

func Bool(ctx context.Context, key any) bool {
	value := ctx.Value(key)
	if value == nil {
		return false
	}
	switch v := value.(type) {
	case bool:
		return v
	case json.Number:
		i, _ := v.Int64()
		return i != 0
	case string:
		if v == "true" {
			return true
		} else if v == "false" {
			return false
		}

		i, _ := strconv.ParseInt(v, 10, 64)
		return i != 0
	default:
		return false
	}
}

func Float64(ctx context.Context, key any) float64 {
	value := ctx.Value(key)
	if value == nil {
		return 0
	}
	switch v := value.(type) {
	case float64:
		return v
	case float32:
		return float64(v)
	case int:
		return float64(v)
	case int8:
		return float64(v)
	case int16:
		return float64(v)
	case int32:
		return float64(v)
	case int64:
		return float64(v)
	case uint:
		return float64(v)
	case uint8:
		return float64(v)
	case uint16:
		return float64(v)
	case uint32:
		return float64(v)
	case uint64:
		return float64(v)
	case json.Number:
		f, _ := v.Float64()
		return f
	case string:
		f, _ := strconv.ParseFloat(v, 64)
		return f
	default:
		return 0
	}
}
