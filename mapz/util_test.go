package mapz

import "testing"

func TestKeys(t *testing.T) {
	m := map[string]int{
		"a": 1,
		"b": 2,
	}
	keys := Keys(m)
	if len(keys) != 2 {
		t.Fatal("keys len not 2")
	}
	if keys[0] != "a" || keys[1] != "b" {
		t.Fatal("keys not match")
	}
}

func TestValues(t *testing.T) {
	m := map[string]int{
		"a": 1,
		"b": 2,
	}
	values := Values(m)
	if len(values) != 2 {
		t.Fatal("values len not 2")
	}
	if values[0] != 1 || values[1] != 2 {
		t.Fatal("values not match")
	}
}
