package randz

import "testing"

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
