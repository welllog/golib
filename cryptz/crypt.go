package cryptz

import (
	"bytes"
	"crypto/aes"
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
	saltLen = 8
	keyLen  = 32
	credLen = 48 // IV(16)+KEY_LEN(32)
)

var (
	// fixedSaltHeader
	fixedSaltHeader = []byte("Salted__")
)

// Encrypt encrypts plainText with secret (openssl aes-256-cbc implementation).
func Encrypt[T, E typez.StrOrBytes](plainText T, secret E) ([]byte, error) {
	encLen := aes.BlockSize + AESCBCEncryptLen(plainText)
	base64EncLen := base64.StdEncoding.EncodedLen(encLen)
	ret := make([]byte, base64EncLen)
	// get the tail buffer to write encrypted data, then use base64 to encode it to the head; this could avoid extra buffer allocation
	enc := ret[base64EncLen-encLen:]

	_, err := SaltBySecretCBCEncrypt(plainText, secret, enc)
	if err != nil {
		return nil, err
	}

	base64.StdEncoding.Encode(ret, enc)
	return ret, nil
}

// Decrypt decrypts cipherText with secret (openssl aes-256-cbc implementation).
func Decrypt[T, E typez.StrOrBytes](cipherText T, secret E) ([]byte, error) {
	enc, err := strz.Base64Decode(cipherText, base64.StdEncoding)
	if err != nil {
		return nil, err
	}

	return SaltBySecretCBCDecrypt(enc, secret, true)
}

// GCMEncrypt encrypts plainText with secret and additionalData
// This is not enough secure for high security requirements.
// Use PasswordEncrypt instead.
func GCMEncrypt[T, E, D typez.StrOrBytes](plainText T, secret E, additionalData D) ([]byte, error) {
	encLen := aes.BlockSize + AESGCMEncryptLen(plainText)
	hexEncLen := hex.EncodedLen(encLen)
	ret := make([]byte, hexEncLen)
	// get the tail buf to write encrypted data, then using hex encode it to head, this could avoid extra buf
	enc := ret[hexEncLen-encLen:]

	_, err := SaltBySecretGCMEncrypt(plainText, secret, additionalData, enc)
	if err != nil {
		return nil, err
	}

	hex.Encode(ret, enc)
	return ret, nil
}

// GCMDecrypt decrypts cipherText with secret and additionalData
// This is not enough secure for high security requirements.
// Use PasswordDecrypt instead.
func GCMDecrypt[T, E, D typez.StrOrBytes](cipherText T, secret E, additionalData D) ([]byte, error) {
	enc, err := strz.HexDecode(cipherText)
	if err != nil {
		return nil, fmt.Errorf("hex decode error: %w", err)
	}

	return SaltBySecretGCMDecrypt(enc, secret, additionalData, true)
}

// EncryptStreamTo encrypts stream to dst with secret
func EncryptStreamTo[E typez.StrOrBytes](dst io.Writer, stream io.Reader, secret E) error {
	var salt [saltLen]byte
	var cred [credLen]byte
	err := fillSaltAndCred(salt[:], cred[:], secret)
	if err != nil {
		return err
	}

	key := cred[:keyLen] // 32 bytes, 256 / 8
	iv := cred[keyLen:]

	_, err = dst.Write(fixedSaltHeader)
	if err != nil {
		return fmt.Errorf("write header error: %w", err)
	}

	_, err = dst.Write(salt[:])
	if err != nil {
		return fmt.Errorf("write header error: %w", err)
	}

	return AESCTRStreamEncrypt(dst, stream, key, iv)
}

// DecryptStreamTo decrypts stream to dst with secret
func DecryptStreamTo[E typez.StrOrBytes](dst io.Writer, stream io.Reader, secret E) error {
	saltHeader := make([]byte, aes.BlockSize)

	_, err := io.ReadFull(stream, saltHeader)
	if err != nil {
		return fmt.Errorf("read header error: %w", err)
	}

	if !bytes.Equal(saltHeader[:8], fixedSaltHeader) {
		return errors.New("invalid encrypted stream")
	}

	var cred [credLen]byte
	fillCred(cred[:], saltHeader[8:], secret)

	key := cred[:keyLen] // 32 bytes, 256 / 8
	iv := cred[keyLen:]

	return AESCTRStreamDecrypt(dst, stream, key, iv)
}

// SaltBySecretCBCEncrypt
func SaltBySecretCBCEncrypt[T, E typez.StrOrBytes](plainText T, secret E, buf []byte) ([]byte, error) {
	var salt [saltLen]byte
	var cred [credLen]byte
	err := fillSaltAndCred(salt[:], cred[:], secret)
	if err != nil {
		return nil, err
	}

	key := cred[:keyLen] // 32 bytes, 256 / 8
	iv := cred[keyLen:]  // 16 bytes, same as block size

	/*
		|Salted__(8 byte)|salt(8 byte)|plaintext|
	*/
	enc := buf
	encLen := aes.BlockSize + AESCBCEncryptLen(plainText)
	if len(buf) < encLen {
		enc = make([]byte, encLen)
	}
	copy(enc[0:], fixedSaltHeader)
	copy(enc[8:], salt[:])

	err = AESCBCEncrypt(enc[aes.BlockSize:], strz.UnsafeStrOrBytesToBytes(plainText), key, iv)
	if err != nil {
		return nil, err
	}

	return enc[:encLen], nil
}

// SaltBySecretCBCDecrypt
func SaltBySecretCBCDecrypt[E typez.StrOrBytes](cipherText []byte, secret E, reuseCipherText bool) ([]byte, error) {
	if len(cipherText) < 2*aes.BlockSize || len(cipherText)&blockSizeMask != 0 {
		return nil, errors.New("cipherText text length illegal")
	}

	if !bytes.Equal(cipherText[:8], fixedSaltHeader) {
		return nil, errors.New("check cbc fixed header error")
	}

	var cred [credLen]byte
	fillCred(cred[:], cipherText[8:aes.BlockSize], secret)

	key := cred[:keyLen] // 32 bytes, 256 / 8
	iv := cred[keyLen:]  // 16 bytes, same as block size

	cipherData := cipherText[aes.BlockSize:]
	plainText := cipherData
	if !reuseCipherText {
		plainText = make([]byte, len(cipherData))
	}

	n, err := AESCBCDecrypt(plainText, cipherData, key, iv)
	if err != nil {
		return nil, err
	}

	return plainText[:n], nil
}

// SaltBySecretGCMEncrypt
func SaltBySecretGCMEncrypt[T, E, D typez.StrOrBytes](plainText T, secret E, additionalData D, buf []byte) ([]byte, error) {
	var salt [saltLen]byte
	var cred [credLen]byte
	err := fillSaltAndCred(salt[:], cred[:], secret)
	if err != nil {
		return nil, err
	}

	key := cred[:keyLen]
	nonce := cred[keyLen : keyLen+nonceSize]

	/*
		|Salted__(8 byte)|salt(8 byte)|plaintext|tag(16 byte)|
	*/
	enc := buf
	encLen := aes.BlockSize + AESGCMEncryptLen(plainText)
	if len(buf) < encLen {
		enc = make([]byte, encLen)
	}
	copy(enc[0:], fixedSaltHeader)
	copy(enc[8:], salt[:])

	err = AESGCMEncrypt(
		enc[aes.BlockSize:],
		strz.UnsafeStrOrBytesToBytes(plainText),
		key,
		nonce,
		strz.UnsafeStrOrBytesToBytes(additionalData),
	)
	if err != nil {
		return nil, err
	}

	return enc[:encLen], nil
}

// SaltBySecretGCMDecrypt
func SaltBySecretGCMDecrypt[E, D typez.StrOrBytes](cipherText []byte, secret E, additionalData D, reuseCipherText bool) ([]byte, error) {
	// min: 16(salt header) + 16(tag)
	if len(cipherText) < 32 {
		return nil, errors.New("cipherText text length illegal")
	}

	if !bytes.Equal(cipherText[:8], fixedSaltHeader) {
		return nil, errors.New("check fixed header error")
	}

	var cred [credLen]byte
	fillCred(cred[:], cipherText[8:aes.BlockSize], secret)

	key := cred[:keyLen] // 32 bytes, 256 / 8
	nonce := cred[keyLen : keyLen+nonceSize]

	cipherData := cipherText[aes.BlockSize:]
	decLen := AESGCMDecryptLen(cipherData)
	if decLen < 0 {
		return nil, errors.New("cipherText length illegal")
	}

	plainText := cipherData[:decLen]
	if !reuseCipherText {
		plainText = make([]byte, decLen)
	}
	err := AESGCMDecrypt(plainText, cipherData, key, nonce, strz.UnsafeStrOrBytesToBytes(additionalData))
	if err != nil {
		return nil, err
	}
	return plainText, nil
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
	for i := 0; i < 3; i++ { // cred 48byte, md5 16byte, three times could fill
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
