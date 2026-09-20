package strz

import (
	"testing"
	"unicode/utf8"

	"github.com/welllog/golib/testz"
)

func TestMask(t *testing.T) {
	tests := []struct {
		str, mask, want string
		start, end      int
	}{
		{"1234567890", "*", "1********0", 1, 1},
		{"1234567890", "*", "12******90", 2, 2},
		{"1234567890", "*", "123****890", 3, 3},
		{"1234567890", "*", "1234**7890", 4, 4},
		{"1234567890", "*", "1234567890", 5, 5},
		{"1234567890", "*", "1234567890", 6, 6},
		{"1234567890", "*", "1*********", 1, 0},
		{"1234567890", "*", "**********", 0, 0},
		{"1234567890", "*", "*********0", 0, 1},
		{"1234567890", "*", "********90", 0, 2},
		{"1234567890", "*", "*******890", 0, 3},
		{"1234567890", "*", "******7890", 0, 4},
		{"1234567890", "*", "*****67890", 0, 5},
		{"1234567890", "*", "****567890", 0, 6},
		{"1234567890", "*", "***4567890", 0, 7},
		{"1234567890", "*", "**34567890", 0, 8},
		{"1234567890", "*", "*234567890", 0, 9},
		{"1234567890", "*", "1234567890", 0, 10},
		{"1234567890", "*", "1234567890", 0, 11},
		{"1234567890", "*", "1234567890", 0, 12},
		{"你好世界", "", "", 0, 0},
		{"你好世界", "", "你", 1, 0},
		{"你好世界", "", "你好", 2, 0},
		{"你好世界", "", "你好世", 3, 0},
		{"你好世界", "", "你好世界", 4, 0},
		{"你好世界", "", "你好世界", 4, 1},
		{"你好世界", "", "你好世界", 4, 2},
		{"你好世界", "", "你好世界", 4, 3},
		{"你好世界", "", "你好世界", 4, 4},
		{"你好世界", "", "你好世界", 3, 4},
		{"你好世界", "", "你好世界", 2, 4},
		{"你好世界", "", "你好世界", 1, 4},
		{"你好世界", "", "你好世界", 0, 4},
		{"你好世界", "", "好世界", 0, 3},
		{"你好世界", "", "世界", 0, 2},
		{"你好世界", "", "界", 0, 1},
		{"你好世界", "", "你界", 1, 1},
		{"你好世界", "", "你好界", 2, 1},
		{"你好世界", "", "你世界", 1, 2},
		{"你好世界", "😀", "你好😀😀", 2, 0},
		{"你好世界", "😀", "😀😀😀😀", 0, 0},
		{"你好世界", "😀😀", "你好世😀😀", 3, 0},
		{"你好世界", "😀😀", "你😀😀界", 1, 1},
		{"你好世界", "😀😀", "😀😀界", 0, 1},
		{"你好世界", "😀😀", "😀😀", 0, 0},
		{"你好世界", "😀😀", "😀😀", -1, -2},
		{"", "😀😀", "😀😀", -1, -2},
		{"", "😀😀", "", 0, 0},
		{"", "😀😀", "", 1, 3},
	}

	for _, tt := range tests {
		if got := Mask(tt.str, tt.mask, tt.start, tt.end); got != tt.want {
			t.Errorf("Mask(%q, %q, %d, %d) = %q, want %q", tt.str, tt.mask, tt.start, tt.end, got, tt.want)
		}
	}
}

func TestUcFirst(t *testing.T) {
	tests := []struct {
		o string
		n string
	}{
		{"test", "Test"},
		{"TEST", "TEST"},
		{"What", "What"},
		{"123yes", "123yes"},
	}
	for _, tt := range tests {
		testz.Equal(t, tt.n, UcFirst(tt.o), "", tt.o)
	}
}

func TestLcFirst(t *testing.T) {
	tests := []struct {
		o string
		n string
	}{
		{"test", "test"},
		{"TEST", "tEST"},
		{"What", "what"},
		{"123Yes", "123Yes"},
	}
	for _, tt := range tests {
		testz.Equal(t, tt.n, LcFirst(tt.o), "", tt.o)
	}
}

func TestRev(t *testing.T) {
	tests := []struct {
		o string
		n string
	}{
		{"test", "tset"},
		{"What", "tahW"},
		{"123&^!", "!^&321"},
	}
	for _, tt := range tests {
		testz.Equal(t, tt.n, Rev(tt.o), "", tt.o)
	}
}

func TestLen(t *testing.T) {
	tests := []struct {
		o string
		n int
	}{
		{"test", 4},
		{"What", 4},
		{"123&^!", 6},
		{"", 0},
		{"你好", 2},
		{"😀", 1},
	}
	for _, tt := range tests {
		testz.Equal(t, tt.n, Len(tt.o), "", tt.o)
		testz.Equal(t, tt.n, RuneCount(tt.o), "", tt.o)
	}
}

func TestSub(t *testing.T) {
	tests := []struct {
		str    string
		start  int
		length int
		result string
	}{
		{"test", 0, 2, "te"},
		{"test", 10, 5, ""},
		{"test", 3, 5, "t"},
		{"test", 2, 1, "s"},
		{"test", 1, -1, "est"},
		{"测试case", 1, 2, "试c"},
		{"测试case", 1, 10, "试case"},
		{"测试case", 5, 10, "e"},
		{"测试case", 5, -1, "e"},
		{"测试case", 5, 1, "e"},
		{"测试case", 6, 0, ""},
		{"测试case", 6, -1, ""},
		{"测试&案例1 33", 2, 5, "&案例1 "},
	}
	for _, tt := range tests {
		testz.Equal(t, tt.result, Sub(tt.str, tt.start, tt.length), "", tt.str, tt.start, tt.length)
	}
}

func TestSubByDisplay(t *testing.T) {
	tests := []struct {
		str    string
		length int
		result string
	}{
		{"test", 2, "te"},
		{"test", 2, "te"},
		{"test", 4, "test"},
		{"测试case", 2, "测"},
		{"测试case", 3, "测"},
		{"测试case", 5, "测试c"},
		// invalid utf-8 bytes count as 2 columns each, must not panic
		{"\xff\xff\xff\xffX", 4, "\xff\xff"},
		{"abc\xff\xff\xff\xffdef", 9, "abc\xff\xff\xff"},
	}

	for _, tt := range tests {
		testz.Equal(t, tt.result, SubByDisplay(tt.str, tt.length), "", tt.str, tt.length)
	}
}

func TestRemoveRunes(t *testing.T) {
	removeFunc := func(n int) func(rune) bool {
		return func(r rune) bool {
			return utf8.RuneLen(r) > n
		}
	}

	tests := []struct {
		str    string
		max    int
		result string
	}{
		{"test测试", 0, ""},
		{"test测试case", 1, "testcase"},
		{"test测试case", 3, "test测试case"},
		{"test测试😀😀,haha", 3, "test测试,haha"},
		{"test测试😀😀,haha", 4, "test测试😀😀,haha"},
	}

	for _, tt := range tests {
		maxn := tt.max
		testz.Equal(t, tt.result, RemoveRunes(tt.str, removeFunc(maxn)), "", tt.str, tt.max)
	}
}

func TestSnakeToCamelCase(t *testing.T) {
	tests := []struct {
		f bool
		a string
		b string
	}{
		{false, "test_snake", "testSnake"},
		{true, "test_snake", "TestSnake"},
		{false, "test_Snake", "testSnake"},
		{true, "test_Snake", "TestSnake"},
		{false, "a_b_c_d", "aBCD"},
		{true, "a_b_c_d", "ABCD"},
		{true, "a_b_c_🤣d", "ABC🤣d"},
		{false, "abcd", "abcd"},
		{true, "abcd", "Abcd"},
		{true, "ABCD", "ABCD"},
		{true, "_a_b_c_d", "_aBCD"},
		{false, "_a_b_c_d", "_aBCD"},
		{false, "_", "_"},
		{true, "a", "A"},
		{false, "a", "a"},
		{false, "a__b", "aB"},
		{true, "a__b", "AB"},
		{false, "a___b", "aB"},
		{false, "test___case", "testCase"},
		{false, "user_1_name", "user1Name"},
		{false, "a_1b", "a1b"},
		{false, "a_b_", "aB"},
		{false, "a_b__", "aB"},
	}
	for _, tt := range tests {
		testz.Equal(t, tt.b, SnakeToCamelCase(tt.a, tt.f), "", tt.a, tt.f)
	}
}

func TestCamelCaseToSnake(t *testing.T) {
	tests := []struct {
		a string
		b string
	}{
		{"test_snake", "testSnake"},
		{"test_snake", "TestSnake"},
		{"a_b_c_d", "aBCD"},
		{"a_b_c_d", "ABCD"},
		{"a_b_c_d", "a_b_c_d"},
		{"a_b_c_d", "A_b_c_d"},
		{"a_b_c_dedg", "a_b_c_dedg"},
		{"a", "A"},
		{"a_b", "AB"},
	}
	for _, tt := range tests {
		testz.Equal(t, tt.a, CamelCaseToSnake(tt.b))
	}
}

func TestToString(t *testing.T) {
	tests := []struct {
		a interface{}
		b string
	}{
		{1, "1"},
		{1.1, "1.1"},
		{true, "true"},
		{false, "false"},
		{nil, ""},
		{[]byte("test"), "test"},
		{[]byte(""), ""},
		{[]byte(nil), ""},
		{[]int{1, 2, 3}, "[1 2 3]"},
		{[]int{}, "[]"},
		{map[string]int{"a": 1, "b": 2}, "map[a:1 b:2]"},
		{map[string]int{}, "map[]"},
	}
	for _, tt := range tests {
		testz.Equal(t, tt.b, ToString(tt.a))
	}
}

func BenchmarkBytes(b *testing.B) {
	s := "hello world, i love you"
	b.Run("UnsafeBytes", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			UnsafeBytes(s)
		}
	})

	b.Run("UnsafeStrOrBytesToBytes", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			UnsafeStrOrBytesToBytes(s)
		}
	})

	b.Run("std.bytes", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = []byte(s)
		}
	})
}
