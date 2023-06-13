package strz

import (
	"testing"

	"github.com/welllog/golib/testz"
)

func TestIP2long(t *testing.T) {
	ip := "127.0.0.1"
	n := IPv4ToLong(ip)
	testz.Equal(t, ip, LongToIPv4(n))
}
