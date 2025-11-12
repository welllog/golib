package cryptz

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"hash"
)

const (
	SHA256 HashAlgo = iota
	SHA512
	MD5
	SHA1
)

type HashAlgo uint8

func (h HashAlgo) Factory() func() hash.Hash {
	switch h {
	case SHA256:
		return sha256.New
	case SHA512:
		return sha512.New
	case MD5:
		return md5.New
	case SHA1:
		return sha1.New
	default:
		return sha256.New
	}
}
