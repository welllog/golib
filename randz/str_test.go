package randz

import "testing"

func TestRandStr_String(t *testing.T) {
	randStr := NewStrGenerator("这是一次测试随机字符串生成")
	t.Log("bits: ", randStr.charIdxBits)
	for i := 0; i < 3; i++ {
		t.Log(randStr.Generate(5))
	}
}

func TestRandStr_String2(t *testing.T) {
	randStr := NewStrGenerator(CHAR_SET)
	t.Log("bits: ", randStr.charIdxBits)
	for i := 0; i < 3; i++ {
		t.Log(randStr.Generate(5))
	}
}

func BenchmarkRandStr_String(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		String(10)
	}
}

func BenchmarkRandStr_String2(b *testing.B) {
	randStr := NewStrGenerator(CHAR_LOWER_SET)
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			randStr.Generate(10)
		}
	})
}
