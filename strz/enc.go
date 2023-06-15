package strz

import (
	"encoding/base64"
	"encoding/hex"
	"net"
	"strconv"
	"strings"

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
	j := 0
	t := make([]byte, 0, 3)

	for i := 0; i < len(s); i++ {
		b[j] = '\\'

		t = strconv.AppendInt(t, int64(s[i]), 8)

		switch len(t) {
		case 1:
			b[j+1] = 48
			b[j+2] = 48
			b[j+3] = t[0]
		case 2:
			b[j+1] = 48
			b[j+2] = t[0]
			b[j+3] = t[1]
		case 3:
			b[j+1] = t[0]
			b[j+2] = t[1]
			b[j+3] = t[2]
		default:
		}
		j += 4
		t = t[:0]
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

// OctalEncodeToString returns the octal encoding of s
func OctalEncodeToString[T typez.StrOrBytes](s T) string {
	return UnsafeString(OctalEncode(s))
}

// OctalDecodeToString returns the string represented by the octal s
func OctalDecodeToString[T typez.StrOrBytes](s T) string {
	return UnsafeString(OctalDecode(s))
}
