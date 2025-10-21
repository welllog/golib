package testz

import (
	"reflect"
	"strconv"
	"strings"
	"testing"
)

func Assert(t *testing.T, condition bool, msgAndArgs ...any) {
	t.Helper()

	if !condition {
		errLog(t, "assertion failed", nil, msgAndArgs)
	}
}

func Equal(t *testing.T, expected, actual any, msgAndArgs ...any) {
	t.Helper()

	if expected == nil && actual == nil {
		return
	}

	// consider nil slice and empty slice to be equal
	if isSlice(expected) && isSlice(actual) {
		expectedValue := reflect.ValueOf(expected)
		actualValue := reflect.ValueOf(actual)
		// Check if both are effectively empty (nil or length 0)
		eIsNilOrEmpty := expectedValue.IsNil() || expectedValue.Len() == 0
		aIsNilOrEmpty := actualValue.IsNil() || actualValue.Len() == 0
		if eIsNilOrEmpty && aIsNilOrEmpty {
			return
		}
	}

	if expected == nil || actual == nil {
		requireLog(t, expected, actual, msgAndArgs)
	}

	e := reflect.TypeOf(expected)
	a := reflect.TypeOf(actual)
	if e.Kind() != a.Kind() {
		requireLog(t, expected, actual, msgAndArgs)
	}

	if e.Kind() == reflect.Func {
		invalidOpLog(t, expected, actual, msgAndArgs)
	}

	if !reflect.DeepEqual(expected, actual) {
		requireLog(t, expected, actual, msgAndArgs)
	}
}

func Nil(t *testing.T, actual any, msgAndArgs ...any) {
	t.Helper()

	if actual != nil {
		// Check if it's a typed nil (e.g., (*int)(nil))
		val := reflect.ValueOf(actual)
		isTypedNil := val.Kind() == reflect.Ptr ||
			val.Kind() == reflect.Map ||
			val.Kind() == reflect.Slice ||
			val.Kind() == reflect.Chan ||
			val.Kind() == reflect.Func ||
			val.Kind() == reflect.Interface
		if isTypedNil && val.IsNil() {
			return
		}
		requireLog(t, nil, actual, msgAndArgs)
	}
}

func isSlice(v any) bool {
	if v == nil {
		return false
	}
	return reflect.TypeOf(v).Kind() == reflect.Slice
}

func requireLog(t *testing.T, expected, actual any, msgAndArgs []any) {
	t.Helper()
	errLog(t, "expected: <%T> %v, actual: <%T> %v;", []any{
		expected, expected, actual, actual,
	}, msgAndArgs)
}

func invalidOpLog(t *testing.T, expected, actual any, msgAndArgs []any) {
	t.Helper()
	errLog(t, "Invalid operation: %#v == %#v;", []any{
		expected, actual,
	}, msgAndArgs)
}

func errLog(t *testing.T, opLog string, assertArgs, msgAndArgs []any) {
	t.Helper()

	args := make([]any, len(assertArgs), len(assertArgs)+len(msgAndArgs))
	copy(args, assertArgs)

	var format strings.Builder
	format.WriteString(opLog)

	for i, v := range msgAndArgs {
		if i == 0 {
			format.WriteString(" msg: %v;")
		} else {
			format.WriteString(" arg")
			format.WriteString(strconv.Itoa(i))
			format.WriteString(": %v;")
		}
		args = append(args, v)
	}

	t.Fatalf(format.String(), args...)
}
