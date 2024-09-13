package strz

import (
	"testing"
	"unsafe"

	"github.com/welllog/golib/testz"
)

func TestUnsafeString(t *testing.T) {
	tests := []string{
		"hello",
		"world ",
		"&%^@*#",
		"",
	}
	for _, tt := range tests {
		testz.Equal(t, tt, UnsafeString([]byte(tt)), "", tt)
	}

	testz.Equal(t, "", UnsafeString(nil))
	s1 := UnsafeString(nil)
	testz.Equal(t, uintptr(0), *(*uintptr)(unsafe.Pointer(&s1)), "string underlying array should be nil")

	b1 := []byte("hello world")
	b2 := UnsafeString(b1)
	testz.Equal(t, *(*uintptr)(unsafe.Pointer(&b1)), *(*uintptr)(unsafe.Pointer(&b2)), "should not copy data")
}

func TestUnsafeBytes(t *testing.T) {
	tests := []string{
		"hello",
		"world ",
		"&%^@*#",
		"\"sds测试",
	}
	for _, tt := range tests {
		testz.Equal(t, []byte(tt), UnsafeBytes(tt), "", tt)
	}

	testz.Equal(t, []byte(nil), UnsafeBytes(""))
	b1 := UnsafeBytes("")
	testz.Equal(t, uintptr(0), *(*uintptr)(unsafe.Pointer(&b1)), "bytes underlying array should be nil")

	s1 := "hello world"
	b2 := UnsafeBytes(s1)
	testz.Equal(t, *(*uintptr)(unsafe.Pointer(&s1)), *(*uintptr)(unsafe.Pointer(&b2)), "should not copy data")
}

func TestUnsafeStrOrBytesToBytes(t *testing.T) {
	tests := []string{
		"hello",
		"world",
		"👋,??, what happen",
	}

	for _, tt := range tests {
		testz.Equal(t, []byte(tt), UnsafeStrOrBytesToBytes(tt), tt)
	}

	tests2 := [][]byte{
		[]byte("hello"),
		[]byte("world"),
		[]byte("👋,??, what happen"),
		{},
		nil,
	}
	for _, tt := range tests2 {
		testz.Equal(t, tt, UnsafeStrOrBytesToBytes(tt), string(tt))
	}

	for i := range tests2 {
		tests2[i] = append(tests2[i], 't', 'e', 's', 't')
	}
	for _, tt := range tests2 {
		testz.Equal(t, tt, UnsafeStrOrBytesToBytes(tt), string(tt))
	}

	testz.Equal(t, []byte(nil), UnsafeStrOrBytesToBytes(""))
	testz.Equal(t, []byte(nil), UnsafeStrOrBytesToBytes([]byte(nil)))

	s1 := "hello world"
	b1 := []byte("hello world")
	b2 := UnsafeStrOrBytesToBytes(s1)
	b3 := UnsafeStrOrBytesToBytes(b1)
	testz.Equal(t, *(*uintptr)(unsafe.Pointer(&s1)), *(*uintptr)(unsafe.Pointer(&b2)), "should not copy data")
	testz.Equal(t, *(*uintptr)(unsafe.Pointer(&b1)), *(*uintptr)(unsafe.Pointer(&b3)), "should not copy data")

	b4 := make([]byte, 0, 10)
	b5 := UnsafeStrOrBytesToBytes(b4)
	testz.Equal(t, b4, b5)
	testz.Equal(t, cap(b4), cap(b5))
}

func TestUnsafeStrOrBytesToString(t *testing.T) {
	tests := []string{
		"hello",
		"world",
		"👋,??, what happen",
		"",
	}

	for _, tt := range tests {
		testz.Equal(t, tt, UnsafeStrOrBytesToString(tt), tt)
	}

	tests2 := [][]byte{
		[]byte("hello"),
		[]byte("world"),
		[]byte("👋,??, what happen"),
		{},
		nil,
	}
	for _, tt := range tests2 {
		testz.Equal(t, string(tt), UnsafeStrOrBytesToString(tt), string(tt))
	}

	s1 := "hello world"
	b1 := []byte("hello world")
	s2 := UnsafeStrOrBytesToString(s1)
	s3 := UnsafeStrOrBytesToString(b1)
	testz.Equal(t, *(*uintptr)(unsafe.Pointer(&s1)), *(*uintptr)(unsafe.Pointer(&s2)), "should not copy data")
	testz.Equal(t, *(*uintptr)(unsafe.Pointer(&b1)), *(*uintptr)(unsafe.Pointer(&s3)), "should not copy data")
}
