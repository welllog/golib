package listz

import (
	"fmt"
	"testing"
)

func TestSkipList_Get(t *testing.T) {
	var l SkipList[string, int]
	_, ok := l.Get("a")
	if ok {
		t.Fatalf("get key from empty list should return false")
	}
	_, ok = l.Remove("a")
	if ok {
		t.Fatalf("remove key from empty list should return false")
	}

	start := 'a'
	for i := 0; i < 26; i++ {
		l.Set(string(start+int32(i)), i)
	}
	for i := 0; i < 26; i++ {
		v, ok := l.Get(string(start + int32(i)))
		if !ok {
			t.Fatalf("get %s but not found", string(start+int32(i)))
		}
		if v != i {
			t.Fatalf("get %s expected %d, got %d", string(start+int32(i)), i, v)
		}
	}

	for i := 0; i < 26; i++ {
		ok := l.SetNx(string(start+int32(i)), i+1)
		if ok {
			t.Fatalf("setnx %s should return false", string(start+int32(i)))
		}

		v, ok := l.Get(string(start + int32(i)))
		if !ok {
			t.Fatalf("get %s but not found", string(start+int32(i)))
		}
		if v != i {
			t.Fatalf("get %s expected %d, got %d", string(start+int32(i)), i, v)
		}
	}

	for i := 0; i < 26; i++ {
		ok := l.SetX(string(start+int32(i)), i+1)
		if !ok {
			t.Fatalf("setx %s should return true", string(start+int32(i)))
		}

		v, ok := l.Get(string(start + int32(i)))
		if !ok {
			t.Fatalf("get %s but not found", string(start+int32(i)))
		}
		if v != i+1 {
			t.Fatalf("get %s expected %d, got %d", string(start+int32(i)), i+1, v)
		}
	}

	if l.Len() != 26 {
		t.Fatalf("list length expected 26, got %d", l.Len())
	}

	var line bool
	for i := l.level - 1; i >= 0; i-- {
		cur := &l.head
		for cur.next[i] != nil {
			fmt.Printf("%s [%d] -> ", cur.next[i].key, cur.next[i].val)
			cur = cur.next[i]
			line = true
		}
		if line {
			fmt.Printf("##level%d##\n", i+1)
			line = false
		}
	}

	for i := 0; i < 26; i++ {
		v, ok := l.Remove(string(start + int32(i)))
		if !ok {
			t.Fatalf("remove %s but not found", string(start+int32(i)))
		}
		if v != i+1 {
			t.Fatalf("remove %s expected %d, got %d", string(start+int32(i)), i+1, v)
		}
	}

	if l.Len() != 0 {
		t.Fatalf("list length expected 0, got %d", l.Len())
	}

	if l.level != 1 {
		t.Fatalf("list level expected 1, got %d", l.level)
	}

	line = false
	for i := l.level - 1; i >= 0; i-- {
		cur := &l.head
		for cur.next[i] != nil {
			fmt.Printf("%s [%d] -> ", cur.next[i].key, cur.next[i].val)
			cur = cur.next[i]
			line = true
		}
		if line {
			fmt.Printf("##level%d##\n", i+1)
			line = false
		}
	}
}

func TestSkipList_Range(t *testing.T) {
	var l SkipList[string, int]
	l.Range(func(key string, val int) bool {
		t.Fatalf("range should not be called")
		return false
	})

	start := 'a'
	for i := 0; i < 26; i++ {
		l.SetNx(string(start+int32(i)), i)
	}

	var idx int
	l.Range(func(key string, val int) bool {
		if key != string(start+int32(idx)) && val != idx {
			t.Fatalf("range expected %s %d, got %s %d", string(start+int32(idx)), idx, key, val)
		}
		idx++
		return true
	})

	if idx != 26 {
		t.Fatalf("range expected 26, got %d", idx)
	}

	idx = 0
	l.Range(func(key string, val int) bool {
		if key != string(start+int32(idx)) && val != idx {
			t.Fatalf("range expected %s %d, got %s %d", string(start+int32(idx)), idx, key, val)
		}
		idx++
		if idx == 13 {
			return false
		}
		return true
	})

	if idx != 13 {
		t.Fatalf("range expected 13, got %d", idx)
	}
}

func TestSkipList_RangeWithStart(t *testing.T) {
	var l SkipList[string, int]
	arr := []string{"a", "b", "d", "e", "f"}
	for i, v := range arr {
		l.Set(v, i)
	}

	idx := 1
	l.RangeWithStart("b", func(key string, val int) bool {
		if key != arr[idx] && val != idx {
			t.Fatalf("range expected %s %d, got %s %d", arr[idx], idx, key, val)
		}
		idx++
		return true
	})
	if idx != 5 {
		t.Fatalf("range expected 5, got %d", idx)
	}

	idx = 2
	l.RangeWithStart("c", func(key string, val int) bool {
		if key != arr[idx] && val != idx {
			t.Fatalf("range expected %s %d, got %s %d", arr[idx], idx, key, val)
		}
		idx++
		return true
	})
	if idx != 5 {
		t.Fatalf("range expected 5, got %d", idx)
	}

	l.RangeWithStart("f", func(key string, val int) bool {
		if key != "f" && val != 4 {
			t.Fatalf("range expected f 4, got %s %d", key, val)
		}
		return true
	})

	l.RangeWithStart("g", func(key string, val int) bool {
		t.Fatalf("range should not be called")
		return true
	})

	l.RangeWithStart("e", func(key string, val int) bool {
		if key != "e" && val != 3 {
			t.Fatalf("range expected e 3, got %s %d", key, val)
		}
		return false
	})
}

func TestSkipList_RangeWithRange(t *testing.T) {
	var l SkipList[string, int]
	arr := []string{"a", "b", "d", "e", "f"}
	for i, v := range arr {
		l.Set(v, i)
	}

	idx := 1
	l.RangeWithRange("b", "e", func(key string, val int) bool {
		if key != arr[idx] && val != idx {
			t.Fatalf("range expected %s %d, got %s %d", arr[idx], idx, key, val)
		}
		idx++
		return true
	})
	if idx != 3 {
		t.Fatalf("range expected 3, got %d", idx)
	}

	idx = 2
	l.RangeWithRange("c", "f", func(key string, val int) bool {
		if key != arr[idx] && val != idx {
			t.Fatalf("range expected %s %d, got %s %d", arr[idx], idx, key, val)
		}
		idx++
		return true
	})
	if idx != 4 {
		t.Fatalf("range expected 4, got %d", idx)
	}

	l.RangeWithRange("f", "g", func(key string, val int) bool {
		if key != "f" && val != 4 {
			t.Fatalf("range expected f 4, got %s %d", key, val)
		}
		return true
	})

	l.RangeWithRange("g", "h", func(key string, val int) bool {
		t.Fatalf("range should not be called")
		return true
	})

	l.RangeWithRange("e", "f", func(key string, val int) bool {
		if key != "e" && val != 3 {
			t.Fatalf("range expected e 3, got %s %d", key, val)
		}
		return false
	})
}

func BenchmarkSkipList_Set(b *testing.B) {
	b.Run("less_data", func(b *testing.B) {
		b.Run("map", func(b *testing.B) {
			m := make(map[int]int)
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				for n := 0; n < 100; n += 3 {
					m[n+2] = n
					m[n+1] = n
					m[n] = n
				}
			}
		})

		b.Run("skip_list", func(b *testing.B) {
			var l SkipList[int, int]
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				for n := 0; n < 100; n += 3 {
					l.Set(n+2, n)
					l.Set(n+1, n)
					l.Set(n, n)
				}
			}
		})
	})

	b.Run("medium_data", func(b *testing.B) {
		b.Run("map", func(b *testing.B) {
			m := make(map[int]int)
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				for n := 0; n < 10000; n += 3 {
					m[n+2] = n
					m[n+1] = n
					m[n] = n
				}
			}
		})

		b.Run("skip_list", func(b *testing.B) {
			var l SkipList[int, int]
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				for n := 0; n < 10000; n += 3 {
					l.Set(n+2, n)
					l.Set(n+1, n)
					l.Set(n, n)
				}
			}
		})
	})

	b.Run("big_data", func(b *testing.B) {
		b.Run("map", func(b *testing.B) {
			m := make(map[int]int)
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				for n := 0; n < 1000000; n += 3 {
					m[n+2] = n
					m[n+1] = n
					m[n] = n
				}
			}
		})

		b.Run("skip_list", func(b *testing.B) {
			var l SkipList[int, int]
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				for n := 0; n < 1000000; n += 3 {
					l.Set(n+2, n)
					l.Set(n+1, n)
					l.Set(n, n)
				}
			}
		})
	})
}

func BenchmarkSkipList_Get(b *testing.B) {
	b.Run("less_data", func(b *testing.B) {
		b.Run("map", func(b *testing.B) {
			m := make(map[int]int)
			num := 100
			for n := 0; n < num; n++ {
				m[n] = n
			}

			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				for n := 0; n < num; n++ {
					_, _ = m[n]
				}
			}
		})

		b.Run("skip_list", func(b *testing.B) {
			var l SkipList[int, int]
			num := 100
			for n := 0; n < num; n++ {
				l.Set(n, n)
			}

			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				for n := 0; n < num; n++ {
					_, _ = l.Get(n)
				}
			}
		})
	})

	b.Run("medium_data", func(b *testing.B) {
		b.Run("map", func(b *testing.B) {
			m := make(map[int]int)
			num := 10000
			for n := 0; n < num; n++ {
				m[n] = n
			}

			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				for n := 0; n < num; n++ {
					_, _ = m[n]
				}
			}
		})

		b.Run("skip_list", func(b *testing.B) {
			var l SkipList[int, int]
			num := 10000
			for n := 0; n < num; n++ {
				l.Set(n, n)
			}

			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				for n := 0; n < num; n++ {
					_, _ = l.Get(n)
				}
			}
		})
	})

	b.Run("big_data", func(b *testing.B) {
		b.Run("map", func(b *testing.B) {
			m := make(map[int]int)
			num := 1000000
			for n := 0; n < num; n++ {
				m[n] = n
			}

			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				for n := 0; n < num; n++ {
					_, _ = m[n]
				}
			}
		})

		b.Run("skip_list", func(b *testing.B) {
			var l SkipList[int, int]
			num := 1000000
			for n := 0; n < num; n++ {
				l.Set(n, n)
			}

			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				for n := 0; n < num; n++ {
					_, _ = l.Get(n)
				}
			}
		})
	})
}

func BenchmarkSkipList_Remove(b *testing.B) {
	b.Run("less_data", func(b *testing.B) {
		b.Run("map", func(b *testing.B) {
			m := make(map[int]int)
			num := 100
			for n := 0; n < num; n++ {
				m[n] = n
			}

			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				for n := 0; n < num; n++ {
					delete(m, n)
				}
			}
		})

		b.Run("skip_list", func(b *testing.B) {
			var l SkipList[int, int]
			num := 100
			for n := 0; n < num; n++ {
				l.Set(n, n)
			}

			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				for n := 0; n < num; n++ {
					l.Remove(n)
				}
			}
		})
	})

	b.Run("medium_data", func(b *testing.B) {
		b.Run("map", func(b *testing.B) {
			m := make(map[int]int)
			num := 10000
			for n := 0; n < num; n++ {
				m[n] = n
			}

			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				for n := 0; n < num; n++ {
					delete(m, n)
				}
			}
		})

		b.Run("skip_list", func(b *testing.B) {
			var l SkipList[int, int]
			num := 10000
			for n := 0; n < num; n++ {
				l.Set(n, n)
			}

			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				for n := 0; n < num; n++ {
					l.Remove(n)
				}
			}
		})
	})

	b.Run("big_data", func(b *testing.B) {
		b.Run("map", func(b *testing.B) {
			m := make(map[int]int)
			num := 1000000
			for n := 0; n < num; n++ {
				m[n] = n
			}

			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				for n := 0; n < num; n++ {
					delete(m, n)
				}
			}
		})

		b.Run("skip_list", func(b *testing.B) {
			var l SkipList[int, int]
			num := 1000000
			for n := 0; n < num; n++ {
				l.Set(n, n)
			}

			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				for n := 0; n < num; n++ {
					l.Remove(n)
				}
			}
		})
	})
}

func BenchmarkSkipList_Set2(b *testing.B) {
	b.Run("ordered", func(b *testing.B) {
		b.Run("map", func(b *testing.B) {
			m := make(map[int]int)
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				m[i] = i
			}
		})

		b.Run("skip_list", func(b *testing.B) {
			var l SkipList[int, int]
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				l.Set(i, i)
			}
		})

		b.Run("skip_list_with_cmp", func(b *testing.B) {
			l := NewSkipListWithCmp[int, int](cmpInt)
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				l.Set(i, i)
			}
		})
	})

	b.Run("not_ordered", func(b *testing.B) {
		b.Run("map", func(b *testing.B) {
			m := make(map[int]int)
			b.ReportAllocs()
			for i := 0; i < b.N; i += 4 {
				m[i+3] = i
				m[i+2] = i
				m[i+1] = i
				m[i] = i
			}
		})

		b.Run("skip_list", func(b *testing.B) {
			var l SkipList[int, int]
			b.ReportAllocs()
			for i := 0; i < b.N; i += 4 {
				l.Set(i+3, i)
				l.Set(i+2, i)
				l.Set(i+1, i)
				l.Set(i, i)
			}
		})

		b.Run("skip_list_with_cmp", func(b *testing.B) {
			l := NewSkipListWithCmp[int, int](cmpInt)
			b.ReportAllocs()
			for i := 0; i < b.N; i += 4 {
				l.Set(i+3, i)
				l.Set(i+2, i)
				l.Set(i+1, i)
				l.Set(i, i)
			}
		})
	})

}

func cmpInt(a int, b int) int {
	if a < b {
		return -1
	}

	if a > b {
		return 1
	}

	return 0
}
