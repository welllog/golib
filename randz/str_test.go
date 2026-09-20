package randz

import (
	"math/rand"
	"testing"
	"time"

	"github.com/welllog/golib/testz"
)

func TestRandStr_String(t *testing.T) {
	randStr := NewStrGenerator("这是一次测试随机字符串生成", rand.NewSource(time.Now().UnixNano()))
	t.Log("bits: ", randStr.charIdxBits)
	for i := 0; i < 3; i++ {
		t.Log(randStr.Generate(5))
	}
}

func TestRandStr_String2(t *testing.T) {
	randStr := NewStrGenerator(CHAR_SET, rand.NewSource(time.Now().UnixNano()))
	t.Log("bits: ", randStr.charIdxBits)
	for i := 0; i < 3; i++ {
		t.Log(randStr.Generate(5))
	}
}

func TestStrGenerator_BitsPowerOfTwo(t *testing.T) {
	tests := []struct {
		charSet      string
		expectedBits int
		expectedMask int64
	}{
		{"a", 1, 1},
		{"ab", 1, 1},
		{"abcd", 2, 3},
		{"01234567", 3, 7},
		{"0123456789abcdef", 4, 15}, // 16 chars: 4 bits
		{"abcdefghijklmnopqrstuvwxyz234567", 5, 31},                                 // 32 chars: 5 bits
		{"ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/", 6, 63}, // 64 chars: 6 bits
		{"ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789", 6, 63},   // 62 chars: 6 bits
		{CHAR_SET, 5, 31}, // 31 chars: 5 bits
	}

	for _, tt := range tests {
		gen := NewStrGenerator(tt.charSet, rand.NewSource(1))
		if gen.charIdxBits != tt.expectedBits {
			t.Errorf("charSet len %d: expected bits %d, got %d", len(tt.charSet), tt.expectedBits, gen.charIdxBits)
		}
		if gen.charIdxMask != tt.expectedMask {
			t.Errorf("charSet len %d: expected mask %d, got %d", len(tt.charSet), tt.expectedMask, gen.charIdxMask)
		}
		res := gen.Generate(20)
		if len([]rune(res)) != 20 {
			t.Errorf("expected length 20, got %d", len([]rune(res)))
		}
	}
}

func BenchmarkStrGenerator_Generate(b *testing.B) {
	b.Run("lock_free", func(b *testing.B) {
		randStr := NewStrGenerator(CHAR_SET, rand.NewSource(time.Now().UnixNano()))
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			randStr.Generate(10)
		}
	})

	b.Run("with_lock", func(b *testing.B) {
		randStr := NewStrGenerator(CHAR_SET, defRandSource)
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			randStr.Generate(10)
		}
	})
}

func BenchmarkString(b *testing.B) {
	b.Run("default", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				String(10)
			}
		})
	})

	b.Run("customer", func(b *testing.B) {
		r := NewStrGenerator(CHAR_SET, defRandSource)
		b.ReportAllocs()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				r.Generate(10)
			}
		})
	})
}

func TestStrGenerator_EmptyCharSetAndNegativeLength(t *testing.T) {
	genEmpty := NewStrGenerator("", rand.NewSource(1))
	testz.Equal(t, "", genEmpty.Generate(10))

	gen := NewStrGenerator("abc", rand.NewSource(1))
	testz.Equal(t, "", gen.Generate(-1))
	testz.Equal(t, "", gen.Generate(0))
}
