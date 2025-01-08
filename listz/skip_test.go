package listz

import (
	"fmt"
	"testing"
)

func TestNewSkipList(t *testing.T) {
	l := NewSkipList[int, int]()
	for i := 0; i < 1000; i++ {
		level := l.randomLevel()
		if level > 12 {
			fmt.Println(level)
		}
	}
}
