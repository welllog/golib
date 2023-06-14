package hashz

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"hash"
	"io"

	"github.com/welllog/golib/strz"
	"github.com/welllog/golib/typez"
)

// Md5 returns the MD5 checksum of the data.
func Md5[T typez.StrOrBytes](s T) []byte {
	h := md5.Sum(strz.UnsafeStrOrBytesToBytes(s))
	return strz.HexEncode(h[:])
}

// Md5ToString returns the MD5 checksum of the data as a string.
func Md5ToString[T typez.StrOrBytes](s T) string {
	return strz.UnsafeString(Md5(s))
}

// Sha1 returns the SHA-1 checksum of the data.
func Sha1[T typez.StrOrBytes](s T) []byte {
	h := sha1.Sum(strz.UnsafeStrOrBytesToBytes(s))
	return strz.HexEncode(h[:])
}

// Sha1ToString returns the SHA-1 checksum of the data as a string.
func Sha1ToString[T typez.StrOrBytes](s T) string {
	return strz.UnsafeString(Sha1(s))
}

// Sha256 returns the SHA-256 checksum of the data.
func Sha256[T typez.StrOrBytes](s T) []byte {
	h := sha256.Sum256(strz.UnsafeStrOrBytesToBytes(s))
	return strz.HexEncode(h[:])
}

// Sha256ToString returns the SHA-256 checksum of the data as a string.
func Sha256ToString[T typez.StrOrBytes](s T) string {
	return strz.UnsafeString(Sha256(s))
}

// Hmac returns the HMAC-HASH of the data.
func Hmac[T, E typez.StrOrBytes](key T, data E, h func() hash.Hash) []byte {
	hh := hmac.New(h, strz.UnsafeStrOrBytesToBytes(key))
	hh.Write(strz.UnsafeStrOrBytesToBytes(data))
	return strz.HexEncode(hh.Sum(nil))
}

// HmacToString returns the HMAC-HASH of the data as a string.
func HmacToString[T, E typez.StrOrBytes](key T, data E, h func() hash.Hash) string {
	return strz.UnsafeString(Hmac(key, data, h))
}

// Md5Stream returns the MD5 checksum of the data.
func Md5Stream(s io.Reader) ([]byte, error) {
	h := md5.New()
	_, err := io.Copy(h, s)
	if err != nil {
		return nil, err
	}

	return strz.HexEncode(h.Sum(nil)), nil
}

// Sha1Stream returns the SHA-1 checksum of the data.
func Sha1Stream(s io.Reader) ([]byte, error) {
	h := sha1.New()
	_, err := io.Copy(h, s)
	if err != nil {
		return nil, err
	}

	return strz.HexEncode(h.Sum(nil)), nil
}

// Sha256Stream returns the SHA-256 checksum of the data.
func Sha256Stream(s io.Reader) ([]byte, error) {
	h := sha256.New()
	_, err := io.Copy(h, s)
	if err != nil {
		return nil, err
	}

	return strz.HexEncode(h.Sum(nil)), nil
}
