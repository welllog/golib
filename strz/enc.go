package strz

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"io"
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

// Base64StdEncode returns the base64 encoding of s
func Base64StdEncode[T typez.StrOrBytes](s T) []byte {
	dst := make([]byte, 0, base64.StdEncoding.EncodedLen(len(s)))
	w := bytes.NewBuffer(dst)

	r := NewReader(s)
	e := base64.NewEncoder(base64.StdEncoding, w)

	_, _ = io.Copy(e, r)
	_ = e.Close()
	return w.Bytes()
}

// Base64StdDecode returns the bytes represented by the base64 s
func Base64StdDecode[T typez.StrOrBytes](s T) ([]byte, error) {
	dst := make([]byte, base64.StdEncoding.DecodedLen(len(s)))

	r := NewReader(s)
	e := base64.NewDecoder(base64.StdEncoding, r)

	n, err := e.Read(dst)
	return dst[:n], err
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
	b := make([]byte, 0, len(s)*2)
	for i := 0; i < len(s); i++ {
		b = append(b, '\\')

		l := len(b)
		b = strconv.AppendInt(b, int64(s[i]), 8)

		switch len(b) - l {
		case 1:
			d := b[len(b)-1]
			b = append(b[:l], 48, 48, d)
		case 2:
			buf := make([]byte, 2)
			copy(buf, b[len(b)-2:])
			b = append(b[:l], 48, buf[0], buf[1])
		default:
		}
	}
	return b
}

// OctalDecode returns the bytes represented by the octal s
func OctalDecode[T typez.StrOrBytes](s T) []byte {
	b := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		if s[i] == '\\' {
			if i+1 < len(s) && s[i+1] == '\\' {
				b = append(b, '\\')
				i++
				continue
			}
			if i+3 < len(s) && s[i+1] >= '0' && s[i+1] <= '7' && s[i+2] >= '0' && s[i+2] <= '7' && s[i+3] >= '0' && s[i+3] <= '7' {
				n, _ := ParseUint(s[i+1:i+4], 8, 8)
				b = append(b, byte(n))
				i += 3
				continue
			}
		}
		b = append(b, s[i])
	}
	return b
}
