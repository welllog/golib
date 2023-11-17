package randz

import (
	"math/rand"
	"testing"
	"time"
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
