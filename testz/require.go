package testz

import (
	"reflect"
	"strconv"
	"strings"
	"testing"
)

func Equal(t *testing.T, expected, actual any, msgAndArgs ...any) {
	t.Helper()

	if expected == nil && actual == nil {
		return
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
		requireLog(t, nil, actual, msgAndArgs)
	}
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
