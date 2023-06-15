package strz

import (
	"encoding/base64"
	"encoding/hex"
	"testing"

	"github.com/welllog/golib/testz"
)

func TestHexEncode(t *testing.T) {
	tests := []struct {
		s string
		e []byte
	}{
		{"hello", []byte(`68656c6c6f`)},
		{"hello world", []byte(`68656c6c6f20776f726c64`)},
		{"hello\nworld", []byte(`68656c6c6f0a776f726c64`)},
		{"👋，世界，welcome", []byte(`f09f918befbc8ce4b896e7958cefbc8c77656c636f6d65`)},
	}
	for _, v := range tests {
		testz.Equal(t, v.e, HexEncode(v.s), v.s, string(v.e), hex.EncodeToString([]byte(v.s)))
	}

	tests2 := []struct {
		s []byte
		e []byte
	}{
		{[]byte("hello"), []byte(`68656c6c6f`)},
		{[]byte("hello world"), []byte(`68656c6c6f20776f726c64`)},
		{[]byte("hello\nworld"), []byte(`68656c6c6f0a776f726c64`)},
		{[]byte("👋，世界，welcome"), []byte(`f09f918befbc8ce4b896e7958cefbc8c77656c636f6d65`)},
	}

	for _, v := range tests2 {
		testz.Equal(t, v.e, HexEncode(v.s), v.s, string(v.e), hex.EncodeToString(v.s))
	}
}

func TestHexDecode(t *testing.T) {
	tests := []struct {
		s string
		e []byte
	}{
		{"68656c6c6f", []byte(`hello`)},
		{"68656c6c6f20776f726c64", []byte(`hello world`)},
		{"68656c6c6f0a776f726c64", []byte("hello\nworld")},
		{"f09f918befbc8ce4b896e7958cefbc8c77656c636f6d65", []byte("👋，世界，welcome")},
	}

	for _, v := range tests {
		r, err := HexDecode(v.s)
		testz.Nil(t, err)
		testz.Equal(t, v.e, r, v.s, string(v.e))
	}

	tests2 := []struct {
		s []byte
		e []byte
	}{
		{[]byte(`68656c6c6f`), []byte(`hello`)},
		{[]byte(`68656c6c6f20776f726c64`), []byte(`hello world`)},
		{[]byte(`68656c6c6f0a776f726c64`), []byte("hello\nworld")},
		{[]byte(`f09f918befbc8ce4b896e7958cefbc8c77656c636f6d65`), []byte("👋，世界，welcome")},
	}

	for _, v := range tests2 {
		r, err := HexDecode(v.s)
		testz.Nil(t, err)
		testz.Equal(t, v.e, r, string(v.s), string(v.e))
	}
}

func TestHexDecodeInPlace(t *testing.T) {
	tests := []struct {
		s []byte
		e []byte
	}{
		{[]byte(`68656c6c6f`), []byte(`hello`)},
		{[]byte(`68656c6c6f20776f726c64`), []byte(`hello world`)},
		{[]byte(`68656c6c6f0a776f726c64`), []byte("hello\nworld")},
		{[]byte(`f09f918befbc8ce4b896e7958cefbc8c77656c636f6d65`), []byte("👋，世界，welcome")},
	}

	for _, v := range tests {
		r, err := HexDecodeInPlace(v.s)
		testz.Nil(t, err)
		testz.Equal(t, v.e, v.s[:r], string(v.e))
	}
}

func TestBase64StdEncode(t *testing.T) {
	tests1 := []struct {
		s string
		e []byte
	}{
		{"hello", []byte(`aGVsbG8=`)},
		{"hello world", []byte(`aGVsbG8gd29ybGQ=`)},
		{"hello\nworld", []byte(`aGVsbG8Kd29ybGQ=`)},
		{"👋，世界，welcome", []byte(`8J+Ri++8jOS4lueVjO+8jHdlbGNvbWU=`)},
	}

	for _, v := range tests1 {
		testz.Equal(t, v.e, Base64Encode(v.s, base64.StdEncoding), v.s, string(v.e), base64.StdEncoding.EncodeToString([]byte(v.s)))
	}

	tests2 := []struct {
		s []byte
		e []byte
	}{
		{[]byte("hello"), []byte(`aGVsbG8=`)},
		{[]byte("hello world"), []byte(`aGVsbG8gd29ybGQ=`)},
		{[]byte("hello\nworld"), []byte(`aGVsbG8Kd29ybGQ=`)},
		{[]byte("👋，世界，welcome"), []byte(`8J+Ri++8jOS4lueVjO+8jHdlbGNvbWU=`)},
	}

	for _, v := range tests2 {
		testz.Equal(t, v.e, Base64Encode(v.s, base64.StdEncoding), v.s, string(v.e), base64.StdEncoding.EncodeToString(v.s))
	}
}

func TestBse64StdDecode(t *testing.T) {
	tests1 := []struct {
		s string
		e []byte
	}{
		{"aGVsbG8=", []byte(`hello`)},
		{"aGVsbG8gd29ybGQ=", []byte(`hello world`)},
		{"aGVsbG8Kd29ybGQ=", []byte("hello\nworld")},
		{"8J+Ri++8jOS4lueVjO+8jHdlbGNvbWU=", []byte(`👋，世界，welcome`)},
	}

	for _, v := range tests1 {
		b, err := Base64Decode(v.s, base64.StdEncoding)
		if err != nil {
			t.Fatal(err)
		}
		testz.Equal(t, v.e, b, v.s, string(b))
	}

	tests2 := []struct {
		s []byte
		e []byte
	}{
		{[]byte(`aGVsbG8=`), []byte(`hello`)},
		{[]byte(`aGVsbG8gd29ybGQ=`), []byte(`hello world`)},
		{[]byte(`aGVsbG8Kd29ybGQ=`), []byte("hello\nworld")},
		{[]byte(`8J+Ri++8jOS4lueVjO+8jHdlbGNvbWU=`), []byte(`👋，世界，welcome`)},
	}

	for _, v := range tests2 {
		b, err := Base64Decode(v.s, base64.StdEncoding)
		if err != nil {
			t.Fatal(err)
		}
		testz.Equal(t, v.e, b, string(v.s), string(b))
	}
}

func TestIPv4ToLong(t *testing.T) {
	tests := []struct {
		s string
		e uint32
	}{
		{"127.0.0.1", 2130706433},
		{"192.108.1.1", 3228303617},
		{"10.10.10.2", 168430082},
	}

	for _, v := range tests {
		testz.Equal(t, v.e, IPv4ToLong(v.s), v.s)
	}
}

func TestLongToIPv4(t *testing.T) {
	tests := []struct {
		s uint32
		e string
	}{
		{2130706433, "127.0.0.1"},
		{3228303617, "192.108.1.1"},
		{168430082, "10.10.10.2"},
	}

	for _, v := range tests {
		testz.Equal(t, v.e, LongToIPv4(v.s), v.s)
	}
}

func TestOctalEncode(t *testing.T) {
	tests1 := []struct {
		s string
		e []byte
	}{
		{"hello", []byte(`\150\145\154\154\157`)},
		{"hello world", []byte(`\150\145\154\154\157\040\167\157\162\154\144`)},
		{"hello\nworld", []byte(`\150\145\154\154\157\012\167\157\162\154\144`)},
		{"👋，世界，welcome", []byte(`\360\237\221\213\357\274\214\344\270\226\347\225\214\357\274\214\167\145\154\143\157\155\145`)},
	}
	for _, v := range tests1 {
		testz.Equal(t, v.e, OctalEncode(v.s), v.s)
	}

	tests2 := []struct {
		s []byte
		e []byte
	}{
		{[]byte("hello"), []byte(`\150\145\154\154\157`)},
		{[]byte("hello world"), []byte(`\150\145\154\154\157\040\167\157\162\154\144`)},
		{[]byte("hello\nworld"), []byte(`\150\145\154\154\157\012\167\157\162\154\144`)},
		{[]byte("👋，世界，welcome"), []byte(`\360\237\221\213\357\274\214\344\270\226\347\225\214\357\274\214\167\145\154\143\157\155\145`)},
	}
	for _, v := range tests2 {
		testz.Equal(t, v.e, OctalEncode(v.s), string(v.s))
	}
}

func TestOctalDecode(t *testing.T) {
	tests1 := []struct {
		s string
		e []byte
	}{
		{`\150\145\154\154\157`, []byte("hello")},
		{`\150\145\154\154\157\040\167\157\162\154\144`, []byte("hello world")},
		{`\150\145\154\154\157\012\167\157\162\154\144`, []byte("hello\nworld")},
		{`\360\237\221\213\357\274\214\344\270\226\347\225\214\357\274\214\167\145\154\143\157\155\145`, []byte("👋，世界，welcome")},
	}
	for _, v := range tests1 {
		testz.Equal(t, v.e, OctalDecode(v.s), v.s)
	}

	tests2 := []struct {
		s []byte
		e []byte
	}{
		{[]byte(`\150\145\154\154\157`), []byte("hello")},
		{[]byte(`\150\145\154\154\157\040\167\157\162\154\144`), []byte("hello world")},
		{[]byte(`\150\145\154\154\157\012\167\157\162\154\144`), []byte("hello\nworld")},
		{[]byte(`\360\237\221\213\357\274\214\344\270\226\347\225\214\357\274\214\167\145\154\143\157\155\145`), []byte("👋，世界，welcome")},
	}
	for _, v := range tests2 {
		testz.Equal(t, v.e, OctalDecode(v.s), string(v.s))
	}
}

func TestOctalDecodeInPlace(t *testing.T) {
	tests := []struct {
		s []byte
		e []byte
	}{
		{[]byte(`\150\145\154\154\157`), []byte("hello")},
		{[]byte(`\150\145\154\154\157\040\167\157\162\154\144`), []byte("hello world")},
		{[]byte(`\150\145\154\154\157\012\167\157\162\154\144`), []byte("hello\nworld")},
		{[]byte(`\360\237\221\213\357\274\214\344\270\226\347\225\214\357\274\214\167\145\154\143\157\155\145`), []byte("👋，世界，welcome")},
	}
	for _, v := range tests {
		n := OctalDecodeInPlace(v.s)
		testz.Equal(t, v.e, v.s[:n], string(v.e))
	}
}

func BenchmarkHexEncode(b *testing.B) {
	s1 := "hello world, i love programming, i coding in golang"
	s2 := []byte("what will happen in the future, i don't know")

	b.Run("std.HexEncode", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			dst1 := make([]byte, hex.EncodedLen(len(s1)))
			hex.Encode(dst1, []byte(s1))

			dst2 := make([]byte, hex.EncodedLen(len(s2)))
			hex.Encode(dst2, s2)
		}
	})

	b.Run("strz.HexEncode", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			HexEncode(s1)
			HexEncode(s2)
		}
	})
}

func BenchmarkBase64Encode(b *testing.B) {
	s1 := "68656c6c6f20776f726c642c2069206c6f76652070726f6772616d6d696e672c206920636f64696e6720696e20676f6c616e67"
	s2 := []byte("68656c6c6f20776f726c642c2069206c6f76652070726f6772616d6d696e672c206920636f64696e6720696e20676f6c616e67")

	b.Run("std.Base64Encode", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			dst1 := make([]byte, base64.StdEncoding.EncodedLen(len(s1)))
			base64.StdEncoding.Encode(dst1, []byte(s1))

			dst2 := make([]byte, base64.StdEncoding.EncodedLen(len(s2)))
			base64.StdEncoding.Encode(dst2, s2)
		}
	})

	b.Run("strz.Base64Encode", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			Base64Encode(s1, base64.StdEncoding)
			Base64Encode(s2, base64.StdEncoding)
		}
	})
}

func BenchmarkHexEncodeToString(b *testing.B) {
	s := "68656c6c6f20776f726c642c2069206c6f76652070726f6772616d6d696e672c206920636f64696e6720696e20676f6c616e67"
	b.Run("std.HexEncodeToString", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			hex.EncodeToString([]byte(s))
		}
	})

	b.Run("strz.HexEncodeToString", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			HexEncodeToString(s)
		}
	})
}
