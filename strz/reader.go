package strz

import (
	"errors"
	"io"
	"unsafe"

	"github.com/welllog/golib/typez"
)

// Reader is a wrapper around a string or byte slice that implements io.Reader,
type Reader[T typez.StrOrBytes] struct {
	s T
	i int64
}

// NewReader returns a new Reader reading from s.
func NewReader[T typez.StrOrBytes](s T) *Reader[T] {
	return &Reader[T]{
		s: s,
		i: 0,
	}
}

// Len returns the number of bytes of the unread portion of the
func (r *Reader[T]) Len() int {
	if r.i >= int64(len(r.s)) {
		return 0
	}
	return int(int64(len(r.s)) - r.i)
}

// Size returns the original length of the underlying string or byte slice.
func (r *Reader[T]) Size() int64 { return int64(len(r.s)) }

// Reset resets the Reader to be reading from s.
func (r *Reader[T]) Read(p []byte) (n int, err error) {
	if r.i >= int64(len(r.s)) {
		return 0, io.EOF
	}
	n = copy(p, r.s[r.i:])
	r.i += int64(n)
	return
}

// ReadAt reads len(p) bytes from the Reader starting at byte offset off.
func (r *Reader[T]) ReadAt(b []byte, off int64) (n int, err error) {
	// cannot modify state - see io.ReaderAt
	if off < 0 {
		return 0, errors.New("strz.Reader.ReadAt: negative offset")
	}
	if off >= int64(len(r.s)) {
		return 0, io.EOF
	}
	n = copy(b, r.s[off:])
	if n < len(b) {
		err = io.EOF
	}
	return
}

// ReadByte reads and returns the next byte from the Reader.
func (r *Reader[T]) ReadByte() (byte, error) {
	if r.i >= int64(len(r.s)) {
		return 0, io.EOF
	}
	b := r.s[r.i]
	r.i++
	return b, nil
}

// UnreadByte unread the last byte. Only the most recently read byte can be unread.
func (r *Reader[T]) UnreadByte() error {
	if r.i <= 0 {
		return errors.New("strz.Reader.UnreadByte: at beginning of slice")
	}
	r.i--
	return nil
}

// WriteTo implements io.WriterTo.
func (r *Reader[T]) WriteTo(w io.Writer) (n int64, err error) {
	if r.i >= int64(len(r.s)) {
		return 0, nil
	}
	b := r.s[r.i:]
	m, err := w.Write(UnsafeStrOrBytesToBytes(b))
	if m > len(b) {
		panic("strz.Reader.WriteTo: invalid Write count")
	}
	r.i += int64(m)
	n = int64(m)
	if m != len(b) && err == nil {
		err = io.ErrShortWrite
	}
	return
}

// Close closes the Reader, preventing further reading.
func (r *Reader[T]) Close() error {
	return nil
}

// Seek implements io.Seeker.
func (r *Reader[T]) Seek(offset int64, whence int) (int64, error) {
	var abs int64
	switch whence {
	case io.SeekStart:
		abs = offset
	case io.SeekCurrent:
		abs = r.i + offset
	case io.SeekEnd:
		abs = int64(len(r.s)) + offset
	default:
		return 0, errors.New("strz.Reader.Seek: invalid whence")
	}
	if abs < 0 {
		return 0, errors.New("strz.Reader.Seek: negative position")
	}
	r.i = abs
	return abs, nil
}

// Reset resets the Reader to read from s.
func (r *Reader[T]) Reset(s T) {
	*r = Reader[T]{s, 0}
}

// Bytes returns a slice of the underlying string or byte slice.
func (r *Reader[T]) Bytes() []byte {
	if r.Len() == 0 {
		return nil
	}

	if unsafe.Sizeof(r.s)/typez.WordBytes == 2 {
		b := make([]byte, r.Len())
		copy(b, r.s[r.i:])
		return b
	}

	s := r.s[r.i:]
	return []byte(s)
}
