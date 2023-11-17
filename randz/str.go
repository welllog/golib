package randz

import (
	"math/rand"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"
)

const (
	CHAR_SET       = "ABCDEFGHJKMNPQRSTUVWXYZ23456789"
	CHAR_LOWER_SET = "abcdefghjkmnpqrstuvwxyz23456789"
)

var defStrGen unsafe.Pointer
var defRandSource = NewLockRandSource(time.Now().UnixNano())

func init() {
	SetStrGeneratorCharSet(CHAR_SET)
}

func SetStrGeneratorCharSet(charSet string) {
	g := NewStrGenerator(charSet, defRandSource)
	atomic.StorePointer(&defStrGen, unsafe.Pointer(&g))
}

// String returns a random string with the specified length.
func String(n int) string {
	return (*StrGenerator)(atomic.LoadPointer(&defStrGen)).Generate(n)
}

type LockRandSource struct {
	s  rand.Source
	mu sync.Mutex
}

func NewLockRandSource(seed int64) *LockRandSource {
	return &LockRandSource{s: rand.NewSource(seed)}
}

func (r *LockRandSource) Int63() int64 {
	r.mu.Lock()
	n := r.s.Int63()
	r.mu.Unlock()
	return n
}

func (r *LockRandSource) Seed(seed int64) {
	r.mu.Lock()
	r.s.Seed(seed)
	r.mu.Unlock()
}

// StrGenerator is a random string generator.
type StrGenerator struct {
	charSet     []rune // character set
	charIdxBits int    // bit required to represent the number of character sets
	charIdxMask int64  // mask, get the last charIdxBits bits of an integer
	charIdxMax  int    // divide the random number into charIdxBits parts and use them respectively
	randSource  rand.Source
}

// NewStrGenerator returns a new StrGenerator.
func NewStrGenerator(charSet string, randSource rand.Source) StrGenerator {
	r := []rune(charSet)

	var bits int
	for l := len(r); l != 0; bits++ {
		l = l >> 1
	}

	return StrGenerator{
		charSet:     r,
		charIdxBits: bits,
		charIdxMask: 1<<bits - 1,
		charIdxMax:  63 / bits,
		randSource:  randSource,
	}
}

// Generate returns a random string with the specified length.
func (r *StrGenerator) Generate(n int) string {
	var buf strings.Builder
	buf.Grow(n)
	for i, cache, remain := n-1, r.randSource.Int63(), r.charIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = r.randSource.Int63(), r.charIdxMax
		}
		if idx := int(cache & r.charIdxMask); idx < len(r.charSet) {
			buf.WriteRune(r.charSet[idx])
			i--
		}
		cache >>= r.charIdxBits
		remain--
	}
	return buf.String()
}
