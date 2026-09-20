package testz

import (
	"testing"
)

type mockT struct {
	failed bool
}

func (m *mockT) Helper() {}

func (m *mockT) Fatalf(format string, args ...any) {
	m.failed = true
}

func TestEqual_DifferentSliceTypes(t *testing.T) {
	m := &mockT{}
	Equal(m, []int(nil), []string(nil))
	if !m.failed {
		t.Fatalf("Equal([]int(nil), []string(nil)) should fail, but it passed")
	}

	m2 := &mockT{}
	Equal(m2, []int{}, []string{})
	if !m2.failed {
		t.Fatalf("Equal([]int{}, []string{}) should fail, but it passed")
	}

	m3 := &mockT{}
	Equal(m3, []int(nil), []int{})
	if m3.failed {
		t.Fatalf("Equal([]int(nil), []int{}) should pass, but it failed")
	}
}
