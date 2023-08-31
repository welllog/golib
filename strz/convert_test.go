package strz

import (
	"fmt"
	"testing"
)

func TestByte2Octal(t *testing.T) {
	tests := []struct {
		b    byte
		want string
	}{
		{
			b:    0,
			want: "\\000",
		},
		{
			b:    255,
			want: "\\377",
		},
		{
			b:    1,
			want: "\\001",
		},
		{
			b:    2,
			want: "\\002",
		},
		{
			b:    42,
			want: "\\052",
		},
	}

	buf := make([]byte, 4)
	for _, tt := range tests {
		byteToOctal(tt.b, buf)
		if string(buf) != tt.want {
			t.Errorf("byte2Octal(%v) = %v; want %v", tt.b, string(buf), tt.want)
		}
	}
}

func TestByteToHex(t *testing.T) {
	tests := []struct {
		b    byte
		want string
	}{
		{
			b:    0,
			want: "\\x00",
		},
		{
			b:    255,
			want: "\\xff",
		},
		{
			b:    1,
			want: "\\x01",
		},
		{
			b:    2,
			want: "\\x02",
		},
		{
			b:    42,
			want: "\\x2a",
		},
	}

	buf := make([]byte, 4)
	for _, tt := range tests {
		byteToHex(tt.b, buf)
		if string(buf) != tt.want {
			t.Errorf("byteToHex(%v) = %v; want %v", tt.b, string(buf), tt.want)
		}
	}
}

func TestDemo(t *testing.T) {
	fmt.Println(ParseUint("9A", 16, 16))
	fmt.Println(ParseUint("9a", 16, 16))
	fmt.Println(ParseUint("9F", 16, 16))
	fmt.Println(ParseUint("9g", 16, 16))
	fmt.Println(ParseUint(",2", 16, 16))
	fmt.Println(ParseUint("z2", 16, 16))

	fmt.Println(parseUint("9A", 16, 16))
	fmt.Println(parseUint("9a", 16, 16))
	fmt.Println(parseUint("9F", 16, 16))
	fmt.Println(parseUint("9g", 16, 16))
	fmt.Println(parseUint(",2", 16, 16))
	fmt.Println(parseUint("z2", 16, 16))
}

func TestDemo1(t *testing.T) {
	b := HexEncodeWithPrefix("哈哈😄")
	fmt.Println(string(b))
	n := HexDecodeWithPrefix(b, b)
	fmt.Println(string(b[:n]))

	s := []byte("\\xE5\\x93\\x88hello\\xE5\\x93\\x88你好\\xF0\\x9F\\x98\\x84")
	n = HexDecodeWithPrefix(s, s)
	fmt.Println(string(s[:n]))

	s = []byte("\\xE5\\x93\\x88hello\\x\\x1\\xE5\\x93\\x88你好\\xF0\\x9F\\x98\\x84")
	n = HexDecodeWithPrefix(s, s)
	fmt.Println(string(s[:n]))
}

func TestUpper(t *testing.T) {
	for i := uint8(33); i <= 126; i++ {
		fmt.Printf("%s: (%s)\n", string(i), string(upper(i)))
	}
}
