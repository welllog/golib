package strz

import (
	"encoding/base64"
	"encoding/hex"
)

// HexEncodeToString returns the hexadecimal encoding of b
func HexEncodeToString(b []byte) string {
	dst := make([]byte, hex.EncodedLen(len(b)))
	hex.Encode(dst, b)
	return String(dst)
}

// Base64EncodeToString returns the base64 encoding of b
func Base64EncodeToString(b []byte) string {
	dst := make([]byte, base64.StdEncoding.EncodedLen(len(b)))
	base64.StdEncoding.Encode(dst, b)
	return String(dst)
}
