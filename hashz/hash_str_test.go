package hashz

import (
	"fmt"
	"testing"

	"github.com/welllog/golib/testz"
)

// distinct inputs sharing a common suffix must not collapse to few hash values
func TestPJWHash_Distinct(t *testing.T) {
	const n = 1000
	set32 := make(map[uint32]struct{}, n)
	set64 := make(map[uint64]struct{}, n)
	for i := 0; i < n; i++ {
		s := fmt.Sprintf("%08dsuffix123", i)
		set32[PJWHash(s)] = struct{}{}
		set64[PJWHash64(s)] = struct{}{}
	}
	testz.Equal(t, n, len(set32))
	testz.Equal(t, n, len(set64))
}

// distinct inputs sharing a common suffix must not collapse to few hash values
func TestELFHash_Distinct(t *testing.T) {
	const n = 1000
	set32 := make(map[uint32]struct{}, n)
	set64 := make(map[uint64]struct{}, n)
	for i := 0; i < n; i++ {
		s := fmt.Sprintf("%08dsuffix123", i)
		set32[ELFHash(s)] = struct{}{}
		set64[ELFHash64(s)] = struct{}{}
	}
	testz.Equal(t, n, len(set32))
	testz.Equal(t, n, len(set64))
}

func TestAPHash_ClassicConsistency(t *testing.T) {
	testz.Equal(t, uint32(2147284991), APHash("ab"))
}

func TestPJWHash64_Spans64Bit(t *testing.T) {
	s := "test_longer_string_for_64bit_hash_distribution"
	h64 := PJWHash64(s)
	h32 := PJWHash(s)
	if h64 == uint64(h32) {
		t.Fatalf("PJWHash64 should not be truncated to 32-bit: got %d == %d", h64, h32)
	}
	if h64 <= 0xFFFFFFFF {
		t.Fatalf("PJWHash64 expected to exceed 32-bit space for long string, got %d", h64)
	}
}

func TestELFHash64_Spans64Bit(t *testing.T) {
	s := "test_longer_string_for_64bit_hash_distribution"
	h64 := ELFHash64(s)
	h32 := ELFHash(s)
	if h64 == uint64(h32) {
		t.Fatalf("ELFHash64 should not be truncated to 32-bit: got %d == %d", h64, h32)
	}
	if h64 <= 0xFFFFFFFF {
		t.Fatalf("ELFHash64 expected to exceed 32-bit space for long string, got %d", h64)
	}
}
