package strz

import (
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"
)

// Mask returns a string with the first `start` and last `end` characters
func Mask(str, mask string, start, end int) string {
	l := utf8.RuneCountInString(str)
	ml := l - start - end
	if ml <= 0 {
		return str
	}

	if utf8.RuneCountInString(mask) == 1 {
		mask = strings.Repeat(mask, ml)
	}
	if ml == l {
		return mask
	}

	end = l - end
	var startIndex, endIndex, count int
	for i := 0; i < len(str); {
		if count == start {
			startIndex = i
		} else if count == end {
			endIndex = i
		}

		if b := str[i]; b < utf8.RuneSelf {
			i++
			count++
			continue
		}

		_, size := utf8.DecodeRuneInString(str[i:])
		i += size
		count++
	}

	if endIndex == 0 {
		endIndex = len(str)
	}

	return str[:startIndex] + mask + str[endIndex:]
}

// UcFirst returns a string with the first character converted to uppercase
func UcFirst(s string) string {
	if len(s) == 0 {
		return s
	}

	b := s[0]
	if b >= 'a' && b <= 'z' {
		b -= 32
		return string(b) + s[1:]
	}
	return s
}

// LcFirst returns a string with the first character converted to lowercase
func LcFirst(s string) string {
	if len(s) == 0 {
		return s
	}

	b := s[0]
	if b >= 'A' && b <= 'Z' {
		b += 32
		return string(b) + s[1:]
	}
	return s
}

// Rev returns a string with the characters in reverse order
func Rev(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// Len returns the number of runes in a string
func Len(s string) int {
	return utf8.RuneCountInString(s)
}

// Sub return a substring of a string by start position and length.
// If length is -1, the substring starting from start position to the end of the string will be returned.
func Sub(s string, start, length int) string {
	if start < 0 || length < -1 || s == "" {
		return s
	}

	if length == 0 {
		return ""
	}

	begin, count := -1, 0
	for i := 0; i < len(s); {
		if count == start {
			if length == -1 {
				return s[i:]
			}
			begin = i
		} else if begin >= 0 && start+length == count {
			return s[begin:i]
		}

		if bt := s[i]; bt < utf8.RuneSelf {
			i++
		} else {
			_, size := utf8.DecodeRuneInString(s[i:])
			i += size
		}
		count++
	}

	if begin < 0 {
		return ""
	}
	return s[begin:]
}

// SubByDisplay returns a substring of a string by display width
func SubByDisplay(s string, length int) string {
	if len(s) <= length {
		return s
	}

	var dpl, end int
	for _, v := range s {
		if v < utf8.RuneSelf {
			dpl += 1
		} else {
			dpl += 2
		}

		if dpl > length {
			break
		}

		end += utf8.RuneLen(v)
	}
	return s[:end]
}

// RemoveRunes removes the specified characters from the string
func RemoveRunes(s string, shouldRemove func(rune) bool) string {
	var buf strings.Builder
	for i, v := range s {
		if buf.Cap() > 0 {
			if !shouldRemove(v) {
				buf.WriteRune(v)
			}
			continue
		}

		if shouldRemove(v) {
			buf.Grow(len(s))
			buf.WriteString(s[:i])
		}
	}

	if buf.Cap() == 0 {
		return s
	}

	return buf.String()
}

// SnakeToCamelCase converts a snake case string to a camel case string
func SnakeToCamelCase(str string, firstUp bool) string {
	var buf strings.Builder

	start := 0
	for i := 0; i < len(str); {
		if b := str[i]; b < utf8.RuneSelf {
			if firstUp {
				firstUp = false

				if b >= 'a' && b <= 'z' {
					if buf.Cap() == 0 {
						buf.Grow(len(str))
					}

					if start < i {
						buf.WriteString(str[start:i])
					}
					buf.WriteByte(b - 32)
					i++
					start = i
					continue
				}

				i++
				continue
			}

			if i > 0 && b == '_' {
				if buf.Cap() == 0 {
					buf.Grow(len(str))
				}

				if start < i {
					buf.WriteString(str[start:i])
				}

				firstUp = true
				i++
				start = i
				continue
			}

			i++
			continue
		}

		_, size := utf8.DecodeRuneInString(str[i:])
		i += size
		firstUp = false
	}

	if buf.Len() == 0 {
		return str
	}

	if start < len(str) {
		buf.WriteString(str[start:])
	}
	return buf.String()
}

// CamelCaseToSnake converts a camel case string to a snake case string
func CamelCaseToSnake(str string) string {
	var buf strings.Builder

	start := 0
	for i := 0; i < len(str); {
		if b := str[i]; b < utf8.RuneSelf {
			if b >= 'A' && b <= 'Z' {
				if buf.Len() == 0 {
					buf.Grow(len(str))
				}

				if start < i {
					buf.WriteString(str[start:i])
				}
				if i > 0 {
					buf.WriteByte('_')
				}
				buf.WriteByte(b + 32)
				i++
				start = i
				continue
			}

			i++
			continue
		}

		_, size := utf8.DecodeRuneInString(str[i:])
		i += size
	}

	if buf.Len() == 0 {
		return str
	}

	if start < len(str) {
		buf.WriteString(str[start:])
	}
	return buf.String()
}

// ToString converts a value to a string
func ToString(value any) string {
	switch v := value.(type) {
	case []byte:
		return UnsafeString(v)
	case string:
		return v
	case nil:
		return ""
	case int:
		return strconv.Itoa(v)
	case int8:
		return strconv.FormatInt(int64(v), 10)
	case int16:
		return strconv.FormatInt(int64(v), 10)
	case int32:
		return strconv.FormatInt(int64(v), 10)
	case int64:
		return strconv.FormatInt(v, 10)
	case uint:
		return strconv.FormatUint(uint64(v), 10)
	case uint8:
		return strconv.FormatUint(uint64(v), 10)
	case uint16:
		return strconv.FormatUint(uint64(v), 10)
	case uint32:
		return strconv.FormatUint(uint64(v), 10)
	case uint64:
		return strconv.FormatUint(v, 10)
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(v)
	default:
		return fmt.Sprintf("%+v", value)
	}
}
