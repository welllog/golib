package cryptz

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"

	"github.com/welllog/golib/strz"
	"github.com/welllog/golib/typez"
)

const (
	_SALT_LEN = 8
	_KEY_LEN  = 32
	_CRED_LEN = 48 // BLOCK_LEN(16)+KEY_LEN(32)
)

var (
	// fixedSaltHeader
	fixedSaltHeader = []byte("Salted__")
)

// Encrypt encrypts plainText with secret (openssl aes-256-cbc implementation).
func Encrypt[T, E typez.StrOrBytes](plainText T, secret E) ([]byte, error) {
	var salt [_SALT_LEN]byte
	var cred [_CRED_LEN]byte
	err := fillSaltAndCred(salt[:], cred[:], secret)
	if err != nil {
		return nil, err
	}

	key := cred[:_KEY_LEN] // 32 bytes, 256 / 8
	iv := cred[_KEY_LEN:]  // 16 bytes, same as block size

	/*
		|Salted__(8 byte)|salt(8 byte)|plaintext|
	*/
	dst := make([]byte, aes.BlockSize+AESCBCEncryptLen(plainText))
	copy(dst[0:], fixedSaltHeader)
	copy(dst[8:], salt[:])

	_ = AESCBCEncrypt(dst[aes.BlockSize:], strz.UnsafeStrOrBytesToBytes(plainText), key, iv)

	ret := make([]byte, base64.StdEncoding.EncodedLen(len(dst)))
	base64.StdEncoding.Encode(ret, dst)

	return ret, nil
}

// Decrypt decrypts cipherText with secret (openssl aes-256-cbc implementation).
func Decrypt[T, E typez.StrOrBytes](cipherText T, secret E) ([]byte, error) {
	src, err := strz.Base64Decode(cipherText, base64.StdEncoding)
	if err != nil {
		return nil, err
	}

	if len(src) < 2*aes.BlockSize || len(src)&blockSizeMask != 0 {
		return nil, errors.New("cipherText text length illegal")
	}

	if !bytes.Equal(src[:8], fixedSaltHeader) {
		return nil, errors.New("check cbc fixed header error")
	}

	var cred [_CRED_LEN]byte
	fillCred(cred[:], src[8:aes.BlockSize], secret)

	key := cred[:_KEY_LEN] // 32 bytes, 256 / 8
	iv := cred[_KEY_LEN:]  // 16 bytes, same as block size

	// reuse src
	dst := src[aes.BlockSize:]
	n, err := AESCBCDecrypt(dst, dst, key, iv)
	if err != nil {
		return nil, err
	}

	return dst[:n], nil
}

// GCMEncrypt encrypts plainText with secret and additionalData
func GCMEncrypt[T, E, D typez.StrOrBytes](plainText T, secret E, additionalData D) ([]byte, error) {
	var salt [_SALT_LEN]byte
	var cred [_CRED_LEN]byte
	err := fillSaltAndCred(salt[:], cred[:], secret)
	if err != nil {
		return nil, err
	}

	key := cred[:_KEY_LEN]
	nonce := cred[_KEY_LEN : _KEY_LEN+nonceSize]

	/*
		|Salted__(8 byte)|salt(8 byte)|plaintext|tag(16 byte)|
	*/
	dst := make([]byte, aes.BlockSize+AESGCMEncryptLen(plainText))
	copy(dst[0:], fixedSaltHeader)
	copy(dst[8:], salt[:])

	_ = AESGCMEncrypt(
		dst[aes.BlockSize:],
		strz.UnsafeStrOrBytesToBytes(plainText),
		key,
		nonce,
		strz.UnsafeStrOrBytesToBytes(additionalData),
	)

	ret := make([]byte, hex.EncodedLen(len(dst)))
	hex.Encode(ret, dst)
	return ret, nil
}

// GCMDecrypt decrypts cipherText with secret and additionalData
func GCMDecrypt[T, E, D typez.StrOrBytes](cipherText T, secret E, additionalData D) ([]byte, error) {
	src, err := strz.HexDecode(cipherText)
	if err != nil {
		return nil, fmt.Errorf("hex decode error: %w", err)
	}

	if len(cipherText) < aes.BlockSize {
		return nil, errors.New("cipherText text length illegal")
	}

	if !bytes.Equal(src[:8], fixedSaltHeader) {
		return nil, errors.New("check fixed header error")
	}

	var cred [_CRED_LEN]byte
	fillCred(cred[:], src[8:aes.BlockSize], secret)

	key := cred[:_KEY_LEN] // 32 bytes, 256 / 8
	nonce := cred[_KEY_LEN : _KEY_LEN+nonceSize]

	// reuse src
	dst := src[aes.BlockSize:]
	err = AESGCMDecrypt(dst, dst, key, nonce, strz.UnsafeStrOrBytesToBytes(additionalData))
	if err != nil {
		return nil, err
	}
	return dst[:AESGCMDecryptLen(dst)], nil
}

// EncryptStreamTo encrypts stream to out with secret
func EncryptStreamTo[E typez.StrOrBytes](out io.Writer, stream io.Reader, secret E) error {
	var salt [_SALT_LEN]byte
	var cred [_CRED_LEN]byte
	err := fillSaltAndCred(salt[:], cred[:], secret)
	if err != nil {
		return err
	}

	key := cred[:_KEY_LEN] // 32 bytes, 256 / 8
	iv := cred[_KEY_LEN:]

	block, err := aes.NewCipher(key)
	if err != nil {
		return fmt.Errorf("NewCipher error: %w", err)
	}

	_, err = out.Write(fixedSaltHeader)
	if err != nil {
		return fmt.Errorf("write fixed salt header error: %w", err)
	}

	_, err = out.Write(salt[:])
	if err != nil {
		return fmt.Errorf("write salt error: %w", err)
	}

	encStream := cipher.NewCTR(block, iv)
	writer := &cipher.StreamWriter{S: encStream, W: out}
	_, err = io.Copy(writer, stream)
	if err != nil {
		return fmt.Errorf("copy stream error: %w", err)
	}

	return nil
}

// DecryptStreamTo decrypts stream to out with secret
func DecryptStreamTo[E typez.StrOrBytes](out io.Writer, stream io.Reader, secret E) error {
	saltHeader := make([]byte, aes.BlockSize)

	n, err := stream.Read(saltHeader)
	if err != nil {
		return fmt.Errorf("read header error: %w", err)
	}

	if n != aes.BlockSize {
		return fmt.Errorf("read header less error: n=%d", n)
	}

	if !bytes.Equal(saltHeader[:8], fixedSaltHeader) {
		return errors.New("check fixed header error")
	}

	var cred [_CRED_LEN]byte
	fillCred(cred[:], saltHeader[8:], secret)

	key := cred[:_KEY_LEN] // 32 bytes, 256 / 8
	iv := cred[_KEY_LEN:]

	block, err := aes.NewCipher(key)
	if err != nil {
		return fmt.Errorf("NewCipher error: %w", err)
	}

	decStream := cipher.NewCTR(block, iv)
	reader := &cipher.StreamReader{S: decStream, R: stream}

	_, err = io.Copy(out, reader)
	if err != nil {
		return fmt.Errorf("copy stream error: %w", err)
	}

	return nil
}

func fillSaltAndCred[E typez.StrOrBytes](salt, cred []byte, secret E) error {
	_, err := io.ReadFull(rand.Reader, salt)
	if err != nil {
		return fmt.Errorf("generate random salt error: %w", err)
	}

	fillCred(cred, salt, secret)

	return nil
}

func fillCred[E typez.StrOrBytes](cred []byte, salt []byte, secret E) {
	buf := make([]byte, 0, 16+len(secret)+len(salt))
	var prevSum [16]byte
	for i := 0; i < 3; i++ { // salted 48byte, md5 16byte, three times could fill
		n := 0 // first prevSum length is zero,so n must be zero
		if i > 0 {
			n = 16
		}
		buf = buf[:n+len(secret)+len(salt)]
		copy(buf, prevSum[:])
		copy(buf[n:], secret)
		copy(buf[n+len(secret):], salt)
		prevSum = md5.Sum(buf)        // md5(prevSum + secret + salt)
		copy(cred[i*16:], prevSum[:]) // concat every md5
	}
}
