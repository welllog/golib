package ipz

import (
	"testing"

	"github.com/welllog/golib/testz"
)

func TestIPv4ToLong(t *testing.T) {
	tests := []struct {
		s string
		e uint32
	}{
		{"127.0.0.1", 2130706433},
		{"192.108.1.1", 3228303617},
		{"10.10.10.2", 168430082},
	}

	for _, v := range tests {
		testz.Equal(t, v.e, IPv4ToLong(v.s), v.s)
	}
}

func TestLongToIPv4(t *testing.T) {
	tests := []struct {
		s uint32
		e string
	}{
		{2130706433, "127.0.0.1"},
		{3228303617, "192.108.1.1"},
		{168430082, "10.10.10.2"},
	}

	for _, v := range tests {
		testz.Equal(t, v.e, LongToIPv4(v.s), v.s)
	}
}
