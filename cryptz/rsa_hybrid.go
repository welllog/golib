package cryptz

import (
	"bufio"
	"bytes"
	"crypto/aes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/welllog/golib/strz"
	"github.com/welllog/golib/typez"
)

const (
	encKeyLenLen = 2 // 2 bytes to store encKey length
	//   text hybrid v1 layout: magic(8)|version(1)|nonce(12)|encKeyLen(2)|encKey(encKeyLen)|cipherText|tag(16)
	// stream hybrid v1 layout: magic(8)|version(1)|iv(16)|encKeyLen(2)|encKey(encKeyLen)|cipherStream
	hybridV1PrefixLen = encPrefixLen + nonceSize + encKeyLenLen
)

var (
	// hybridMagic WLHBRSEG
	hybridMagic = [magicLen]byte{'W', 'L', 'H', 'B', 'R', 'S', 'E', 'G'}
	// hybridStreamMagic WLHBRSTR
	hybridStreamMagic = [magicLen]byte{'W', 'L', 'H', 'B', 'R', 'S', 'T', 'R'}
)

// RsaHybridEncrypt use RSA-OAEP + AES-GCM to encrypt plaintext
// It returns base64 URL encoded cipher text
func RsaHybridEncrypt[T, D typez.StrOrBytes](plaintext T, ad D, pub *rsa.PublicKey) ([]byte, error) {
	var rb [keyLen + nonceSize]byte
	_, err := rand.Read(rb[:])
	if err != nil {
		return nil, err
	}

	aesKey := rb[:keyLen] // AES-256
	nonce := rb[keyLen:]  // nonceSize bytes

	encKey, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, pub, aesKey, nil)
	if err != nil {
		return nil, err
	}

	// version 1 ------------
	// magic(8)|version(1)|nonce(12)|encKeyLen(2)|encKey(encKeyLen)|cipherText|tag(16)
	encLen := hybridV1PrefixLen + len(encKey) + len(plaintext) + gcmTagSize
	base64EncLen := base64.RawURLEncoding.EncodedLen(encLen)

	ret := make([]byte, base64EncLen)
	enc := ret[base64EncLen-encLen:]

	copy(enc, hybridMagic[:])
	enc[magicLen] = 1 // version 1
	copy(enc[encPrefixLen:], nonce)

	binary.BigEndian.PutUint16(enc[encPrefixLen+nonceSize:hybridV1PrefixLen], uint16(len(encKey)))
	offset := hybridV1PrefixLen
	offset += copy(enc[hybridV1PrefixLen:], encKey)

	err = AESGCMEncrypt(enc[offset:], strz.UnsafeStrOrBytesToBytes(plaintext), aesKey, nonce,
		strz.UnsafeStrOrBytesToBytes(ad))
	if err != nil {
		return nil, err
	}

	base64.RawURLEncoding.Encode(ret, enc)
	return ret, nil
}

// RsaHybridDecrypt use RSA-OAEP + AES-GCM to decrypt ciphertext
// The input ciphertext is expected to be base64 URL encoded
func RsaHybridDecrypt[T, D typez.StrOrBytes](ciphertext T, ad D, prv *rsa.PrivateKey) ([]byte, error) {
	enc, err := strz.Base64Decode(ciphertext, base64.RawURLEncoding)
	if err != nil {
		return nil, err
	}

	if len(enc) < encPrefixLen || !bytes.Equal(enc[:magicLen], hybridMagic[:]) {
		return nil, ErrInvalidCipherText
	}

	switch enc[magicLen] {
	case 1:
		return rsaHybridDecryptV1(enc, ad, prv)
	default:
		return nil, ErrInvalidCipherText
	}
}

// RsaHybridEncryptStream encrypts data from stream and writes to dst using RSA-OAEP + AES-CTR
func RsaHybridEncryptStream(dst io.Writer, stream io.Reader, pub *rsa.PublicKey) error {
	// magic(8)|version(1)|iv(16)|encKeyLen(2)|encKey(encKeyLen)|cipherStream
	var rb [keyLen + aes.BlockSize]byte
	_, err := rand.Read(rb[:])
	if err != nil {
		return err
	}

	aesKey := rb[:keyLen] // AES-256
	iv := rb[keyLen:]

	encKey, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, pub, aesKey, nil)
	if err != nil {
		return err
	}

	encKeyLenBytes := make([]byte, encKeyLenLen)
	binary.BigEndian.PutUint16(encKeyLenBytes, uint16(len(encKey)))

	w := bufio.NewWriter(dst)
	_, _ = w.Write(hybridStreamMagic[:])
	_ = w.WriteByte(1)
	_, _ = w.Write(iv)
	_, _ = w.Write(encKeyLenBytes)
	_, _ = w.Write(encKey)
	err = AESCTRStreamEncrypt(w, stream, aesKey, iv)
	if err != nil {
		return err
	}

	return w.Flush()
}

// RsaHybridDecryptStream decrypts data from stream and writes to dst using RSA-OAEP + AES-CTR
func RsaHybridDecryptStream(dst io.Writer, stream io.Reader, prv *rsa.PrivateKey) error {
	// magic(8)|version(1)|iv(16)|encKeyLen(2)|encKey(encKeyLen)|cipherStream
	buf := make([]byte, aes.BlockSize+prv.Size()) // prv.Size min 128 bytes for 1024-bit key
	r := bufio.NewReader(stream)

	_, err := io.ReadFull(r, buf[:encPrefixLen])
	if err != nil {
		return fmt.Errorf("%w: %s", ErrInvalidCipherStream, err)
	}

	if !bytes.Equal(buf[:magicLen], hybridStreamMagic[:]) {
		return fmt.Errorf("%w: %s", ErrInvalidCipherStream, "invalid header")
	}

	if buf[magicLen] != 1 {
		return fmt.Errorf("%w: %s", ErrInvalidCipherStream, "unsupported version")
	}

	_, err = io.ReadFull(r, buf[:aes.BlockSize+encKeyLenLen])
	if err != nil {
		return fmt.Errorf("%w: %s", ErrInvalidCipherStream, err)
	}

	iv := buf[:aes.BlockSize]
	buf = buf[aes.BlockSize:]

	encKeyLen := int(binary.BigEndian.Uint16(buf[:encKeyLenLen]))
	if encKeyLen <= 0 || encKeyLen > prv.Size() {
		return fmt.Errorf("%w: %s", ErrInvalidCipherStream, "invalid header")
	}
	_, err = io.ReadFull(r, buf[:encKeyLen])
	if err != nil {
		return fmt.Errorf("%w: %s", ErrInvalidCipherStream, err)
	}

	aesKey, err := rsa.DecryptOAEP(sha256.New(), nil, prv, buf[:encKeyLen], nil)
	if err != nil {
		return err
	}

	return AESCTRStreamDecrypt(dst, r, aesKey, iv)
}

func rsaHybridDecryptV1[D typez.StrOrBytes](enc []byte, ad D, prv *rsa.PrivateKey) ([]byte, error) {
	// magic(8)|version(1)|nonce(12)|encKeyLen(2)|encKey(encKeyLen)|cipherText|tag(16)
	if len(enc) < hybridV1PrefixLen {
		return nil, ErrInvalidCipherText
	}

	nonce := enc[encPrefixLen : encPrefixLen+nonceSize]
	encKeyLen := int(binary.BigEndian.Uint16(enc[encPrefixLen+nonceSize : hybridV1PrefixLen]))
	if encKeyLen <= 0 || encKeyLen > prv.Size() {
		return nil, ErrInvalidCipherText
	}

	if len(enc) < hybridV1PrefixLen+encKeyLen+gcmTagSize {
		return nil, ErrInvalidCipherText
	}

	encKey := enc[hybridV1PrefixLen : hybridV1PrefixLen+encKeyLen]
	encData := enc[hybridV1PrefixLen+encKeyLen:]

	aesKey, err := rsa.DecryptOAEP(sha256.New(), nil, prv, encKey, nil)
	if err != nil {
		return nil, err
	}

	plaintext := encData[:len(encData)-gcmTagSize]
	err = AESGCMDecrypt(plaintext, encData, aesKey, nonce, strz.UnsafeStrOrBytesToBytes(ad))
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}
