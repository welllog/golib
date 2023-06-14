package strz

import (
	"encoding/hex"
	"fmt"

	"github.com/welllog/golib/typez"
)

const hextable = "0123456789abcdef"

func hexEncode[T typez.StrOrBytes](dst []byte, src T) int {
	j := 0
	for i := 0; i < len(src); i++ {
		dst[j] = hextable[src[i]>>4]
		dst[j+1] = hextable[src[i]&0x0f]
		j += 2
	}
	return len(src) * 2
}

func hexDecode[T typez.StrOrBytes](dst []byte, src T) (int, error) {
	i, j := 0, 1
	for ; j < len(src); j += 2 {
		a, ok := fromHexChar(src[j-1])
		if !ok {
			return i, fmt.Errorf("encoding/hex: invalid byte: %#U", rune(src[j-1]))
		}
		b, ok := fromHexChar(src[j])
		if !ok {
			return i, fmt.Errorf("encoding/hex: invalid byte: %#U", rune(src[j]))
		}
		dst[i] = (a << 4) | b
		i++
	}
	if len(src)%2 == 1 {
		// Check for invalid char before reporting bad length,
		// since the invalid char (if present) is an earlier problem.
		if _, ok := fromHexChar(src[j-1]); !ok {
			return i, fmt.Errorf("encoding/hex: invalid byte: %#U", rune(src[j-1]))
		}
		return i, hex.ErrLength
	}
	return i, nil
}

// fromHexChar converts a hex character into its value and a success flag.
func fromHexChar(c byte) (byte, bool) {
	switch {
	case '0' <= c && c <= '9':
		return c - '0', true
	case 'a' <= c && c <= 'f':
		return c - 'a' + 10, true
	case 'A' <= c && c <= 'F':
		return c - 'A' + 10, true
	}

	return 0, false
}
