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

	for _, k := range keys {
		if k != "a" && k != "b" {
			t.Fatal("keys not match")
		}
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

	for _, v := range values {
		if v != 1 && v != 2 {
			t.Fatal("values not match")
		}
	}
}
