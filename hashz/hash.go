package hashz

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"hash"
	"io"

	"github.com/welllog/golib/strz"
)

func Md5(s string) string {
	h := md5.Sum(strz.Bytes(s))
	return strz.HexEncodeToString(h[:])
}

func Sha1(s string) string {
	h := sha1.Sum(strz.Bytes(s))
	return strz.HexEncodeToString(h[:])
}

func Sha256(s string) string {
	h := sha256.Sum256(strz.Bytes(s))
	return strz.HexEncodeToString(h[:])
}

func Hmac(key, data string, h func() hash.Hash) string {
	hh := hmac.New(h, strz.Bytes(key))
	hh.Write(strz.Bytes(data))
	return strz.HexEncodeToString(hh.Sum(nil))
}

func Md5Stream(s io.Reader) (string, error) {
	h := md5.New()
	_, err := io.Copy(h, s)
	if err != nil {
		return "", err
	}

	return strz.HexEncodeToString(h.Sum(nil)), nil
}

func Sha1Stream(s io.Reader) (string, error) {
	h := sha1.New()
	_, err := io.Copy(h, s)
	if err != nil {
		return "", err
	}

	return strz.HexEncodeToString(h.Sum(nil)), nil
}

func Sha256Stream(s io.Reader) (string, error) {
	h := sha256.New()
	_, err := io.Copy(h, s)
	if err != nil {
		return "", err
	}

	return strz.HexEncodeToString(h.Sum(nil)), nil
}
