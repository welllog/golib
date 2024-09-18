//go:build go1.23

package mapz

import (
	"fmt"
	"runtime"
	"testing"
	"time"

	"github.com/welllog/golib/testz"
)

func TestSafeKV_All(t *testing.T) {
	s := NewSafeKV[string, int](4)
	m := KV[string, int]{
		"a": 1,
		"b": 2,
		"c": 3,
	}

	for k, v := range m {
		s.Set(k, v)
	}

	ch := make(chan error, 1)
	var n int
	for k, v := range s.All() {
		testz.Equal(t, m[k], v)
		if n == 0 {
			go func() {
				begin := time.Now()
				s.Set("d", 4)
				if time.Since(begin).Milliseconds() < 200 {
					ch <- fmt.Errorf("Set should block")
				} else {
					close(ch)
				}
			}()
			runtime.Gosched()
			time.Sleep(200 * time.Millisecond)
		}
		n++
	}

	err, ok := <-ch
	if ok {
		t.Fatal(err)
	}
}
