//go:build !go1.21

package strz

import (
	"reflect"
	"unsafe"

	"github.com/welllog/golib/typez"
)

// UnsafeString converts byte slice to string.
func UnsafeString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// UnsafeBytes converts string to byte slice. maybe safe risk
func UnsafeBytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(
		&struct {
			string
			Cap int
		}{s, len(s)},
	))
}

// UnsafeStrOrBytesToBytes converts string or byte slice to byte slice. maybe safe risk
func UnsafeStrOrBytesToBytes[T typez.StrOrBytes](s T) []byte {
	if unsafe.Sizeof(s)/typez.WordBytes == 2 {
		var b []byte
		hdr := (*reflect.SliceHeader)(unsafe.Pointer(&b))
		hdr.Data = *(*uintptr)(unsafe.Pointer(&s))
		hdr.Len = len(s)
		hdr.Cap = len(s)
		return b
	}

	return []byte(s)
}

// UnsafeStrOrBytesToString converts string or byte slice to string.
func UnsafeStrOrBytesToString[T typez.StrOrBytes](s T) string {
	return *(*string)(unsafe.Pointer(&s))
}
