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

func requireLog(t *testing.T, expected, actual any, msgAndArgs []any) {
	t.Helper()
	errLog(t, expected, actual, "expected: %v, actual: %v;", msgAndArgs)
}

func invalidOpLog(t *testing.T, expected, actual any, msgAndArgs []any) {
	t.Helper()
	errLog(t, expected, actual, "Invalid operation: %#v == %#v;", msgAndArgs)
}

func errLog(t *testing.T, expected, actual any, opLog string, msgAndArgs []any) {
	t.Helper()

	args := make([]interface{}, 0, len(msgAndArgs)+2)
	args = append(args, expected, actual)

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
