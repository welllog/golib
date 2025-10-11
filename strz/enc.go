package strz

import (
	"encoding/base64"
	"encoding/hex"
	"net"
	"regexp"
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

// IPv4ToLong converts an IPv4 address to an uint32
func IPv4ToLong(ip string) uint32 {
	var long uint32
	for _, v := range strings.Split(ip, ".") {
		n, _ := strconv.ParseInt(v, 10, 32)
		long = long<<8 + uint32(n)
	}
	return long
}

// LongToIPv4 converts an uint32 to an IPv4 address
func LongToIPv4(long uint32) string {
	return net.IPv4(byte(long>>24), byte(long>>16), byte(long>>8), byte(long)).String()
}

// OctalFormat returns the octal encode format \ooo
func OctalFormat[T typez.StrOrBytes](s T) []byte {
	b := make([]byte, len(s)*4)
	var j, f int

	for i := 0; i < len(s); i++ {
		b[j] = '\\'
		f = j + 1
		j += 4
		appendUint(b[f:j], uint64(s[i]), 8)
	}
	return b
}

// OctalParse fill dst with the bytes represented by the octal s like \ooo
func OctalParse(dst, src []byte) int {
	var e, f int
	for i := 0; i < len(src); {
		if len(src)-i < 4 {
			break
		}

		if src[i] != '\\' {
			i++
			continue
		}

		n, j, ok := parseUint(src[i+1:i+4], 8, 8)
		if !ok {
			i += 1 + j
			continue
		}

		if f < i {
			e += copy(dst[e:], src[f:i])
		}
		dst[e] = byte(n)

		e++
		i += 4
		f = i
	}

	if f < len(src) {
		e += copy(dst[e:], src[f:])
	}

	return e
}

// HexFormat returns the hex encode format \xXX
func HexFormat[T typez.StrOrBytes](s T) []byte {
	b := make([]byte, len(s)*4)
	var j, f int

	for i := 0; i < len(s); i++ {
		b[j] = '\\'
		b[j+1] = 'x'

		f = j + 2
		j += 4
		appendUint(b[f:j], uint64(s[i]), 16)
		toUpper(b[f:j])
	}
	return b
}

// HexParse fill dst with the bytes represented by the hex s like \xXX
func HexParse(dst, src []byte) int {
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
			e += copy(dst[e:], src[f:i])
		}
		dst[e] = byte(n)

		e++
		i += 4
		f = i
	}

	if f < len(src) {
		e += copy(dst[e:], src[f:])
	}

	return e
}

// UnicodeFormat returns the unicode encode format \UXXXXXXXX
func UnicodeFormat[T typez.StrOrBytes](s T) []byte {
	src := UnsafeStrOrBytesToString(s)
	b := make([]byte, utf8.RuneCountInString(src)*10)
	var j, f int

	for i := 0; i < len(src); {
		b[j] = '\\'
		b[j+1] = 'U'

		f = j + 2
		j += 10

		if bt := src[i]; bt < utf8.RuneSelf {
			appendUint(b[f:j], uint64(bt), 16)
			toUpper(b[f:j])

			i++
			continue
		}

		c, size := utf8.DecodeRuneInString(src[i:])
		if c == utf8.RuneError {
			copy(b[f:j], "0000FFFD")
			i += size
			continue
		}

		appendUint(b[f:j], uint64(c), 16)
		toUpper(b[f:j])
		i += size
	}

	return b
}

// UnicodeParse fill dst with the bytes represented by the Unicode s like \UXXXXXXXX
func UnicodeParse(dst, src []byte) int {
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

		if n > utf8.MaxRune {
			i += 10
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

// Utf16Format returns the utf16 encode format \uXXXX
func Utf16Format[T typez.StrOrBytes](s T) []byte {
	src := UnsafeStrOrBytesToString(s)
	b := make([]byte, 0, utf8.RuneCountInString(src)*6)
	var j, f int

	for i := 0; i < len(src); {
		b = append(b, '\\', 'u', '0', '0', '0', '0')

		f = j + 2
		j += 6
		if bt := src[i]; bt < utf8.RuneSelf {
			appendUint(b[f:j], uint64(bt), 16)
			toUpper(b[f:j])

			i++
			continue
		}

		c, size := utf8.DecodeRuneInString(src[i:])
		switch {
		case c == utf8.RuneError:
			copy(b[f:j], "FFFD")
		case c >= 0 && c < 0xd800, c >= 0xe000 && c < 0x10000:
			appendUint(b[f:j], uint64(c), 16)
			toUpper(b[f:j])
		case c >= 0x10000 && c <= utf8.MaxRune:
			r1, r2 := utf16.EncodeRune(c)
			appendUint(b[f:j], uint64(r1), 16)
			toUpper(b[f:j])

			b = append(b, '\\', 'u', '0', '0', '0', '0')
			f = j + 2
			j += 6

			appendUint(b[f:j], uint64(r2), 16)
			toUpper(b[f:j])
		default:
			copy(b[f:j], "FFFD")
		}

		i += size
	}

	return b
}

// Utf16Parse fill dst with the bytes represented by the utf16 s like \uXXXX
func Utf16Parse(dst, src []byte) int {
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
		}

		// English: n1 not satisfied skip current parse, n2 not satisfied skip again on the basis of if condition
		i += 6
	}

	if f < len(src) {
		e += copy(dst[e:], src[f:])
	}

	return e
}

// OctalFormatToString returns the octal encode format \ooo
func OctalFormatToString[T typez.StrOrBytes](s T) string {
	return UnsafeString(OctalFormat(s))
}

// OctalParseToString returns the string represented by the octal s like \ooo
func OctalParseToString[T typez.StrOrBytes](s T) string {
	b := make([]byte, len(s))
	n := OctalParse(b, UnsafeStrOrBytesToBytes(s))
	return UnsafeString(b[:n])
}

// HexFormatToString returns the hex encode format \xXX
func HexFormatToString[T typez.StrOrBytes](s T) string {
	return UnsafeString(HexFormat(s))
}

// HexParseToString returns the string represented by the hex s like \xXX
func HexParseToString[T typez.StrOrBytes](s T) string {
	b := make([]byte, len(s))
	n := HexParse(b, UnsafeStrOrBytesToBytes(s))
	return UnsafeString(b[:n])
}

// UnicodeFormatToString returns the unicode encode format \UXXXXXXXX
func UnicodeFormatToString[T typez.StrOrBytes](s T) string {
	return UnsafeString(UnicodeFormat(s))
}

// UnicodeParseToString returns the string represented by the Unicode s like \UXXXXXXXX
func UnicodeParseToString[T typez.StrOrBytes](s T) string {
	b := make([]byte, len(s))
	n := UnicodeParse(b, UnsafeStrOrBytesToBytes(s))
	return UnsafeString(b[:n])
}

// Utf16FormatToString returns the utf16 encode format \uXXXX
func Utf16FormatToString[T typez.StrOrBytes](s T) string {
	return UnsafeString(Utf16Format(s))
}

// Utf16ParseToString returns the string represented by the utf16 s like \uXXXX
func Utf16ParseToString[T typez.StrOrBytes](s T) string {
	b := make([]byte, len(s))
	n := Utf16Parse(b, UnsafeStrOrBytesToBytes(s))
	return UnsafeString(b[:n])
}

func Base64ParseToString[T typez.StrOrBytes](s T) string {
	// match base64 strings, including those with padding
	// Base64 characters: A-Z, a-z, 0-9, +, / (or - and _ for URL-safe)
	// Padding: = (0, 1, or 2 at the end)
	re := regexp.MustCompile(`[A-Za-z0-9+/_-]+(?:={0,2})?`)

	input := UnsafeStrOrBytesToString(s)
	return re.ReplaceAllStringFunc(input, func(match string) string {
		var decoded []byte
		var err error

		if len(match)%4 == 0 {
			// 1. try Standard Encoding (with padding)
			decoded, err = base64.StdEncoding.DecodeString(match)
			if err == nil {
				return UnsafeString(decoded)
			}

			// 2. try URL Encoding (with padding)
			decoded, err = base64.URLEncoding.DecodeString(match)
			if err == nil {
				return UnsafeString(decoded)
			}
		}

		// 3. try RawStdEncoding (no padding)
		decoded, err = base64.RawStdEncoding.DecodeString(match)
		if err == nil {
			return UnsafeString(decoded)
		}

		// 4. try RawURLEncoding (no padding)
		decoded, err = base64.RawURLEncoding.DecodeString(match)
		if err == nil {
			return UnsafeString(decoded)
		}

		return match
	})
}

var zeroPadding = []byte{'0', '0', '0', '0', '0', '0', '0', '0'}

// appendUint dst size not enough will panic
func appendUint(dst []byte, i uint64, base int) {
	b := strconv.AppendUint(dst[:0], i, base)
	x := len(dst) - len(b)
	copy(dst[x:], b)
	copy(dst[:x], zeroPadding)
}

func toUpper(dst []byte) {
	for i, b := range dst {
		dst[i] = upper(b)
	}
}
