package strz

import (
	"encoding/base64"
	"encoding/hex"
	"net"
	"strconv"
	"strings"
)

// HexEncodeToString returns the hexadecimal encoding of b
func HexEncodeToString(b []byte) string {
	dst := make([]byte, hex.EncodedLen(len(b)))
	hex.Encode(dst, b)
	return String(dst)
}

// HexDecodeString returns the bytes represented by the hexadecimal string s
func HexDecodeString(s string) ([]byte, error) {
	return hex.DecodeString(s)
}

// Base64EncodeToString returns the base64 encoding of b
func Base64EncodeToString(b []byte) string {
	dst := make([]byte, base64.StdEncoding.EncodedLen(len(b)))
	base64.StdEncoding.Encode(dst, b)
	return String(dst)
}

// Base64DecodeString returns the bytes represented by the base64 string s
func Base64DecodeString(s string) ([]byte, error) {
	dst := make([]byte, base64.StdEncoding.DecodedLen(len(s)))
	n, err := base64.StdEncoding.Decode(dst, Bytes(s))
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
