package randz

import (
	"fmt"
	"testing"
)

func TestCountGenerator(t *testing.T) {
	counter := &CountGenerator{}
	// 30 minutes, 1~2 times every 3 seconds, 1~100 times after half an hour
	counter.AddRule(1800, 100, 3, 2)
	// 1 day, 1~3 times every 15 seconds, 1~300 times after one day
	counter.AddRule(86400, 300, 15, 3)
	// 5 days, 1~4 times every 180 seconds
	counter.AddRule(5*86400, 0, 180, 4)
	// 60 days, 1~5 times every 600 seconds
	counter.AddRule(60*86400, 0, 600, 5)

	id := "test"
	t.Log(counter.Max(50 * 86400))
	t.Log(counter.Generate(id, 50*86400))
	t.Log(counter.Min(50 * 86400))

	id = "world"
	t.Log(counter.Generate(id, 50*86400))
	t.Log(counter.Generate(id, 63*86400))
	t.Log(counter.Generate(id, 63*86400))
}

func TestCountGenerator_MinMax(t *testing.T) {
	counter := &CountGenerator{}
	counter.AddRule(1800, 100, 3, 2)
	counter.AddRule(86400, 300, 15, 3)
	counter.AddRule(5*86400, 0, 180, 4)
	counter.AddRule(60*86400, 0, 600, 5)

	for _, diff := range []int{1, 3, 100, 1799, 1800, 1801, 86399, 86400, 431999, 432000, 50 * 86400, 5184000, 6000000} {
		min, max := counter.Min(diff), counter.Max(diff)
		for i := 0; i < 1000; i++ {
			id := fmt.Sprintf("user%d", i)
			got := counter.Generate(id, diff)
			if got < min {
				t.Fatalf("Generate(%q, %d) = %d, less than Min %d", id, diff, got, min)
			}
			if got > max {
				t.Fatalf("Generate(%q, %d) = %d, greater than Max %d", id, diff, got, max)
			}
		}
	}
}

func TestCountGenerator_EdgeCases(t *testing.T) {
	counter := &CountGenerator{}
	counter.AddRule(100, 0, 10, 0)
	counter.AddRule(200, 0, 10, 0)

	// diff <= 0
	if got := counter.Min(0); got != 0 {
		t.Errorf("Min(0) = %d, want 0", got)
	}
	if got := counter.Min(-10); got != 0 {
		t.Errorf("Min(-10) = %d, want 0", got)
	}
	if got := counter.Max(0); got != 0 {
		t.Errorf("Max(0) = %d, want 0", got)
	}
	if got := counter.Generate("test", 0); got != 0 {
		t.Errorf("Generate(0) = %d, want 0", got)
	}

	// zero increment rules
	if got := counter.Min(150); got != 0 {
		t.Errorf("Min(150) = %d, want 0", got)
	}
	if got := counter.Max(150); got != 0 {
		t.Errorf("Max(150) = %d, want 0", got)
	}
	if got := counter.Generate("test", 150); got != 0 {
		t.Errorf("Generate(150) = %d, want 0", got)
	}
}

func TestCountGenerator_IntervalZeroPanic(t *testing.T) {
	var c CountGenerator
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("expected panic when interval <= 0")
		}
	}()
	c.AddRule(100, 10, 0, 1)
}
