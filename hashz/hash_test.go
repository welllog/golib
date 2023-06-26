package hashz

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
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

	h3 := sha256.New224()
	h3.Write(b)
	testz.Equal(t, hex.EncodeToString(h3.Sum(nil)), Sha224ToString(s))

	h4 := sha256.New()
	h4.Write(b)
	testz.Equal(t, hex.EncodeToString(h4.Sum(nil)), Sha256ToString(s))

	h5 := sha512.New384()
	h5.Write(b)
	testz.Equal(t, hex.EncodeToString(h5.Sum(nil)), Sha384ToString(s))

	h6 := sha512.New()
	h6.Write(b)
	testz.Equal(t, hex.EncodeToString(h6.Sum(nil)), Sha512ToString(s))

	h7 := sha512.New512_224()
	h7.Write(b)
	testz.Equal(t, hex.EncodeToString(h7.Sum(nil)), Sha512_224ToString(s))

	h8 := sha512.New512_256()
	h8.Write(b)
	testz.Equal(t, hex.EncodeToString(h8.Sum(nil)), Sha512_256ToString(s))
}
