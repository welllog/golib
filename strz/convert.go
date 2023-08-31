package strz

import (
	"strconv"
)

var zeroPadding = []byte{'0', '0', '0', '0', '0', '0', '0', '0'}

func appendInt64(n int64, base int, dst []byte) {
	s := strconv.AppendInt(dst[:0], n, base)
	i := len(dst) - len(s)
	copy(dst[i:], s)
	copy(dst[:i], zeroPadding)
}

func toUpper(dst []byte) {
	for i, b := range dst {
		dst[i] = upper(b)
	}
}

func byteToOctal(b byte, dst []byte) {
	appendInt64(int64(b), 8, dst)
}

func byteToHex(b byte, dst []byte) {
	appendInt64(int64(b), 16, dst)
}

func runeToHex(r rune, dst []byte) {
	appendInt64(int64(r), 16, dst)
}

func u16ToHex(r uint16, dst []byte) {
	appendInt64(int64(r), 16, dst)
}
