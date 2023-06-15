package hashz

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"testing"

	"github.com/welllog/golib/testz"
)

func TestHash(t *testing.T) {
	s := "abcdefg"
	b := []byte(s)

	h1 := md5.New()
	h1.Write(b)
	testz.Equal(t, hex.EncodeToString(h1.Sum(nil)), Md5ToString(s))

	h2 := sha1.New()
	h2.Write(b)
	testz.Equal(t, hex.EncodeToString(h2.Sum(nil)), Sha1ToString(s))

	h3 := sha256.New()
	h3.Write(b)
	testz.Equal(t, hex.EncodeToString(h3.Sum(nil)), Sha256ToString(s))
}
