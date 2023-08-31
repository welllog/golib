package strz

import (
	"encoding/base64"
	"encoding/hex"
	"net"
	"strconv"
	"strings"
	"unicode/utf16"
	"unicode/utf8"

	"github.com/welllog/golib/typez"
)

// HexEncode returns the hex encoding of s
func HexEncode[T typez.StrOrBytes](s T) []byte {
	dst := make([]byte, hex.EncodedLen(len(s)))
	hexEncode(dst, s)
	return dst
}

// HexDecode returns the bytes represented by the hexadecimal s
func HexDecode[T typez.StrOrBytes](s T) ([]byte, error) {
	dst := make([]byte, hex.DecodedLen(len(s)))
	n, err := hexDecode(dst, s)
	return dst[:n], err
}

// HexDecodeInPlace decodes the hexadecimal s in place
func HexDecodeInPlace(b []byte) (int, error) {
	return hex.Decode(b, b)
}

// HexEncodeToString returns the hex encoding of s
func HexEncodeToString[T typez.StrOrBytes](s T) string {
	return UnsafeString(HexEncode(s))
}

// HexDecodeToString returns the string represented by the hexadecimal s
func HexDecodeToString[T typez.StrOrBytes](s T) (string, error) {
	b, err := HexDecode(s)
	return UnsafeString(b), err
}

// Base64Encode returns the base64 encoding of s
func Base64Encode[T typez.StrOrBytes](s T, enc *base64.Encoding) []byte {
	dst := make([]byte, enc.EncodedLen(len(s)))
	enc.Encode(dst, UnsafeStrOrBytesToBytes(s))
	return dst
}

// Base64Decode returns the bytes represented by the base64 s
func Base64Decode[T typez.StrOrBytes](s T, enc *base64.Encoding) ([]byte, error) {
	dst := make([]byte, enc.DecodedLen(len(s)))
	n, err := enc.Decode(dst, UnsafeStrOrBytesToBytes(s))
	return dst[:n], err
}

// Base64EncodeToString returns the base64 encoding of s
func Base64EncodeToString[T typez.StrOrBytes](s T, enc *base64.Encoding) string {
	return UnsafeString(Base64Encode(s, enc))
}

// Base64DecodeToString returns the string represented by the base64 s
func Base64DecodeToString[T typez.StrOrBytes](s T, enc *base64.Encoding) (string, error) {
	b, err := Base64Decode(s, enc)
	return UnsafeString(b), err
}

// IPv4ToLong converts an IPv4 address to a uint32
func IPv4ToLong(ip string) uint32 {
	var long uint32
	for _, v := range strings.Split(ip, ".") {
		n, _ := strconv.ParseInt(v, 10, 32)
		long = long<<8 + uint32(n)
	}
	return long
}

// LongToIPv4 converts a uint32 to an IPv4 address
func LongToIPv4(long uint32) string {
	return net.IPv4(byte(long>>24), byte(long>>16), byte(long>>8), byte(long)).String()
}

// OctalEncode returns the octal encoding of s
func OctalEncode[T typez.StrOrBytes](s T) []byte {
	b := make([]byte, len(s)*4)
	var j, f int

	for i := 0; i < len(s); i++ {
		b[j] = '\\'
		f = j + 1
		j += 4
		appendInt64(int64(s[i]), 8, b[f:j])
	}
	return b
}

// OctalDecode returns the bytes represented by the octal s
func OctalDecode[T typez.StrOrBytes](s T) []byte {
	b := make([]byte, len(s))
	j := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\\' {
			if i+1 < len(s) && s[i+1] == '\\' {
				b[j] = '\\'
				j++
				i++
				continue
			}
			if i+3 < len(s) && s[i+1] >= '0' && s[i+1] <= '7' && s[i+2] >= '0' && s[i+2] <= '7' && s[i+3] >= '0' && s[i+3] <= '7' {
				n, _ := ParseUint(s[i+1:i+4], 8, 8)
				b[j] = byte(n)
				j++
				i += 3
				continue
			}
		}
		b[j] = s[i]
		j++
	}
	return b[:j]
}

// OctalDecodeInPlace decodes the octal s in place
func OctalDecodeInPlace(b []byte) int {
	j := 0
	for i := 0; i < len(b); i++ {
		if b[i] == '\\' {
			if i+1 < len(b) && b[i+1] == '\\' {
				b[j] = '\\'
				j++
				i++
				continue
			}
			if i+3 < len(b) && b[i+1] >= '0' && b[i+1] <= '7' && b[i+2] >= '0' && b[i+2] <= '7' && b[i+3] >= '0' && b[i+3] <= '7' {
				n, _ := ParseUint(b[i+1:i+4], 8, 8)
				b[j] = byte(n)
				j++
				i += 3
				continue
			}
		}
		b[j] = b[i]
		j++
	}
	return j
}

// HexEncodeWithPrefix returns the hex encode format \xXX
func HexEncodeWithPrefix[T typez.StrOrBytes](s T) []byte {
	b := make([]byte, len(s)*4)
	var j, f int

	for i := 0; i < len(s); i++ {
		b[j] = '\\'
		b[j+1] = 'x'

		f = j + 2
		j += 4
		appendInt64(int64(s[i]), 16, b[f:j])
		toUpper(b[f:j])
	}
	return b
}

func HexDecodeWithPrefix(src, dst []byte) int {
	var e, f int
	for i := 0; i < len(src); {
		if len(src)-i < 4 {
			break
		}

		if src[i] != '\\' || src[i+1] != 'x' {
			i++
			continue
		}

		n, j, ok := parseUint(src[i+2:i+4], 16, 8)
		if !ok {
			i += 2 + j
			continue
		}

		if f < i {
			incr := copy(dst[e:], src[f:i])
			e += incr
		}
		dst[e] = byte(n)

		e++
		i += 4
		f = i
	}

	if f < len(src) {
		incr := copy(dst[e:], src[f:])
		e += incr
	}

	return e
}

func UnicodeEncode[T typez.StrOrBytes](s T) []byte {
	src := UnsafeStrOrBytesToString(s)

	b := make([]byte, utf8.RuneCountInString(src)*10)
	var j, f int

	for i := 0; i < len(src); {
		b[j] = '\\'
		b[j+1] = 'U'

		f = j + 2
		j += 10

		if bt := src[i]; bt < utf8.RuneSelf {
			appendInt64(int64(bt), 16, b[f:j])
			toUpper(b[f:j])

			i++
			continue
		}

		c, size := utf8.DecodeRuneInString(src[i:])
		if c != utf8.RuneError {
			appendInt64(int64(c), 16, b[f:j])
			toUpper(b[f:j])
		}

		i += size
	}

	return b
}

func UnicodeDecode(src, dst []byte) int {
	var e, f int
	for i := 0; i < len(src); {
		if len(src)-i < 10 {
			break
		}

		if src[i] != '\\' || src[i+1] != 'U' {
			i++
			continue
		}

		n, j, ok := parseUint(src[i+2:i+10], 16, 32)
		if !ok {
			i += 2 + j
			continue
		}

		if f < i {
			e += copy(dst[e:], src[f:i])
		}

		if n < utf8.RuneSelf {
			dst[e] = byte(n)
			e++
		} else {
			e += utf8.EncodeRune(dst[e:], rune(n))
		}
		i += 10
		f = i
	}

	if f < len(src) {
		e += copy(dst[e:], src[f:])
	}

	return e
}

func Utf16Encode[T typez.StrOrBytes](s T) []byte {
	src := UnsafeStrOrBytesToString(s)

	b := make([]byte, 0, utf8.RuneCountInString(src)*6)
	var j, f int

	for i := 0; i < len(src); {
		b = append(b, '\\', 'u', '0', '0', '0', '0')

		f = j + 2
		j += 6
		if bt := src[i]; bt < utf8.RuneSelf {
			appendInt64(int64(bt), 16, b[f:j])
			toUpper(b[f:j])

			i++
			continue
		}

		c, size := utf8.DecodeRuneInString(src[i:])
		if c == utf8.RuneError {
			copy(b[f:j], "FFFD")
			i += size
			continue
		}

		if (c >= 0 && c < 0xd800) || (c >= 0xe000 && c < 0x10000) {
			appendInt64(int64(c), 16, b[f:j])
			toUpper(b[f:j])
		} else if c >= 0x10000 && c <= '\U0010FFFF' {
			r1, r2 := utf16.EncodeRune(c)
			appendInt64(int64(r1), 16, b[f:j])
			toUpper(b[f:j])

			b = append(b, '\\', 'u', '0', '0', '0', '0')
			f = j + 2
			j += 6

			appendInt64(int64(r2), 16, b[f:j])
			toUpper(b[f:j])
		} else {
			copy(b[f:j], "FFFD")
		}

		i += size
	}

	return b
}

func Utf16Decode(src, dst []byte) int {
	var e, f int
	for i := 0; i < len(src); {
		if len(src)-i < 6 {
			break
		}

		if src[i] != '\\' || src[i+1] != 'u' {
			i++
			continue
		}

		n1, j, ok := parseUint(src[i+2:i+6], 16, 16)
		if !ok {
			i += 2 + j
			continue
		}

		if f < i {
			e += copy(dst[e:], src[f:i])
			f = i
		}

		if n1 < 0xd800 || n1 >= 0xe000 {
			e += utf8.EncodeRune(dst[e:], rune(n1))
			i += 6
			f = i
			continue
		}

		if n1 >= 0xd800 && n1 < 0xdc00 {
			i += 6
			if len(src)-i < 6 {
				break
			}

			if src[i] != '\\' || src[i+1] != 'u' {
				i++
				continue
			}

			n2, j, ok := parseUint(src[i+2:i+6], 16, 16)
			if !ok {
				i += 2 + j
				continue
			}

			if n2 >= 0xdc00 && n2 < 0xe000 {
				r := utf16.DecodeRune(rune(n1), rune(n2))
				e += utf8.EncodeRune(dst[e:], r)

				i += 6
				f = i
				continue
			}

			i += 6
		}

	}

	if f < len(src) {
		e += copy(dst[e:], src[f:])
	}

	return e
}

// OctalEncodeToString returns the octal encoding of s
func OctalEncodeToString[T typez.StrOrBytes](s T) string {
	return UnsafeString(OctalEncode(s))
}

// OctalDecodeToString returns the string represented by the octal s
func OctalDecodeToString[T typez.StrOrBytes](s T) string {
	return UnsafeString(OctalDecode(s))
}
