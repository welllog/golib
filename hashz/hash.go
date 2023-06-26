package hashz

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
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

// Sha224 returns the SHA-224 checksum of the data.
func Sha224[T typez.StrOrBytes](s T) []byte {
	h := sha256.Sum224(strz.UnsafeStrOrBytesToBytes(s))
	return strz.HexEncode(h[:])
}

// Sha224ToString returns the SHA-224 checksum of the data as a string.
func Sha224ToString[T typez.StrOrBytes](s T) string {
	return strz.UnsafeString(Sha224(s))
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

// Sha384 returns the SHA-384 checksum of the data.
func Sha384[T typez.StrOrBytes](s T) []byte {
	h := sha512.Sum384(strz.UnsafeStrOrBytesToBytes(s))
	return strz.HexEncode(h[:])
}

// Sha384ToString returns the SHA-384 checksum of the data as a string.
func Sha384ToString[T typez.StrOrBytes](s T) string {
	return strz.UnsafeString(Sha384(s))
}

// Sha512 returns the SHA-512 checksum of the data.
func Sha512[T typez.StrOrBytes](s T) []byte {
	h := sha512.Sum512(strz.UnsafeStrOrBytesToBytes(s))
	return strz.HexEncode(h[:])
}

// Sha512ToString returns the SHA-512 checksum of the data as a string.
func Sha512ToString[T typez.StrOrBytes](s T) string {
	return strz.UnsafeString(Sha512(s))
}

// Sha512_224 returns the SHA-512/224 checksum of the data.
func Sha512_224[T typez.StrOrBytes](s T) []byte {
	h := sha512.Sum512_224(strz.UnsafeStrOrBytesToBytes(s))
	return strz.HexEncode(h[:])
}

// Sha512_224ToString returns the SHA-512/224 checksum of the data as a string.
func Sha512_224ToString[T typez.StrOrBytes](s T) string {
	return strz.UnsafeString(Sha512_224(s))
}

// Sha512_256 returns the SHA-512/256 checksum of the data.
func Sha512_256[T typez.StrOrBytes](s T) []byte {
	h := sha512.Sum512_256(strz.UnsafeStrOrBytesToBytes(s))
	return strz.HexEncode(h[:])
}

// Sha512_256ToString returns the SHA-512/256 checksum of the data as a string.
func Sha512_256ToString[T typez.StrOrBytes](s T) string {
	return strz.UnsafeString(Sha512_256(s))
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

// Sha224Stream returns the SHA-224 checksum of the data.
func Sha224Stream(s io.Reader) ([]byte, error) {
	h := sha256.New224()
	_, err := io.Copy(h, s)
	if err != nil {
		return nil, err
	}

	return strz.HexEncode(h.Sum(nil)), nil
}

// Sha384Stream returns the SHA-384 checksum of the data.
func Sha384Stream(s io.Reader) ([]byte, error) {
	h := sha512.New384()
	_, err := io.Copy(h, s)
	if err != nil {
		return nil, err
	}

	return strz.HexEncode(h.Sum(nil)), nil
}

// Sha512Stream returns the SHA-512 checksum of the data.
func Sha512Stream(s io.Reader) ([]byte, error) {
	h := sha512.New()
	_, err := io.Copy(h, s)
	if err != nil {
		return nil, err
	}

	return strz.HexEncode(h.Sum(nil)), nil
}
