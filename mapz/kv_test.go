package mapz

import "testing"

func TestGetExistingKey(t *testing.T) {
	m := KV[string, int]{
		"a": 1,
		"b": 2,
	}
	value, ok := m.Get("a")
	if !ok || value != 1 {
		t.Fatal("Get existing key failed")
	}
}

func TestGetNonExistingKey(t *testing.T) {
	m := KV[string, int]{
		"a": 1,
		"b": 2,
	}
	_, ok := m.Get("c")
	if ok {
		t.Fatal("Get non-existing key failed")
	}
}

func TestSetNewKey(t *testing.T) {
	m := KV[string, int]{}
	m.Set("a", 1)
	if len(m) != 1 || m["a"] != 1 {
		t.Fatal("Set new key failed")
	}
}

func TestSetExistingKey(t *testing.T) {
	m := KV[string, int]{
		"a": 1,
	}
	m.Set("a", 2)
	if len(m) != 1 || m["a"] != 2 {
		t.Fatal("Set existing key failed")
	}
}

func TestSetNxExistingKey(t *testing.T) {
	m := KV[string, int]{
		"a": 1,
	}
	ok := m.SetNx("a", 2)
	if ok || m["a"] != 1 {
		t.Fatal("SetNx existing key failed")
	}
}

func TestSetNxNewKey(t *testing.T) {
	m := KV[string, int]{}
	ok := m.SetNx("a", 1)
	if !ok || m["a"] != 1 {
		t.Fatal("SetNx new key failed")
	}
}

func TestSetXExistingKey(t *testing.T) {
	m := KV[string, int]{
		"a": 1,
	}
	ok := m.SetX("a", 2)
	if !ok || m["a"] != 2 {
		t.Fatal("SetX existing key failed")
	}
}

func TestSetXNewKey(t *testing.T) {
	m := KV[string, int]{}
	ok := m.SetX("a", 1)
	if ok || len(m) != 0 {
		t.Fatal("SetX new key failed")
	}
}

func TestDeleteExistingKey(t *testing.T) {
	m := KV[string, int]{
		"a": 1,
	}
	m.Delete("a")
	if len(m) != 0 {
		t.Fatal("Delete existing key failed")
	}
}

func TestDeleteNonExistingKey(t *testing.T) {
	m := KV[string, int]{
		"a": 1,
	}
	m.Delete("b")
	if len(m) != 1 {
		t.Fatal("Delete non-existing key failed")
	}
}

func TestHasExistingKey(t *testing.T) {
	m := KV[string, int]{
		"a": 1,
	}
	if !m.Has("a") {
		t.Fatal("Has existing key failed")
	}
}

func TestHasNonExistingKey(t *testing.T) {
	m := KV[string, int]{
		"a": 1,
	}
	if m.Has("b") {
		t.Fatal("Has non-existing key failed")
	}
}

func TestLen(t *testing.T) {
	m := KV[string, int]{
		"a": 1,
		"b": 2,
	}
	if m.Len() != 2 {
		t.Fatal("Len failed")
	}
}

func TestKV_Keys(t *testing.T) {
	m := KV[string, int]{
		"a": 1,
		"b": 2,
	}
	keys := m.Keys()
	if len(keys) != 2 || (keys[0] != "a" && keys[0] != "b") || (keys[1] != "a" && keys[1] != "b") {
		t.Fatal("Keys failed")
	}
}

func TestKV_Values(t *testing.T) {
	m := KV[string, int]{
		"a": 1,
		"b": 2,
	}
	values := m.Values()
	if len(values) != 2 || (values[0] != 1 && values[0] != 2) || (values[1] != 1 && values[1] != 2) {
		t.Fatal("Values failed")
	}
}

func TestRange(t *testing.T) {
	m := KV[string, int]{
		"a": 1,
		"b": 2,
	}
	count := 0
	m.Range(func(key string, value int) bool {
		count++
		return true
	})
	if count != 2 {
		t.Fatal("Range failed")
	}
}
