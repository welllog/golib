//go:build go1.21

package strz

import (
	"unsafe"

	"github.com/welllog/golib/typez"
)

// UnsafeString converts byte slice to string.
func UnsafeString(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
}

// UnsafeBytes converts string to byte slice. maybe safe risk
func UnsafeBytes(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

// UnsafeStrOrBytesToBytes converts string or byte slice to byte slice. maybe safe risk
func UnsafeStrOrBytesToBytes[T typez.StrOrBytes](s T) []byte {
	if unsafe.Sizeof(s)/typez.WordBytes == 2 {
		return unsafe.Slice(unsafe.StringData(string(s)), len(s))
	}
	return []byte(s)
}

// UnsafeStrOrBytesToString converts string or byte slice to string.
func UnsafeStrOrBytesToString[T typez.StrOrBytes](s T) string {
	if unsafe.Sizeof(s)/typez.WordBytes == 2 {
		return string(s)
	}
	return unsafe.String(unsafe.SliceData([]byte(s)), len(s))
}
