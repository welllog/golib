package cryptz

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
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
	cipherText, err := SaltBySecretCBCEncrypt(plainText, secret)
	if err != nil {
		return nil, err
	}

	ret := make([]byte, base64.StdEncoding.EncodedLen(len(cipherText)))
	base64.StdEncoding.Encode(ret, cipherText)

	return ret, nil
}

// Decrypt decrypts cipherText with secret (openssl aes-256-cbc implementation).
func Decrypt[T, E typez.StrOrBytes](cipherText T, secret E) ([]byte, error) {
	src, err := strz.Base64Decode(cipherText, base64.StdEncoding)
	if err != nil {
		return nil, err
	}

	return SaltBySecretCBCDecrypt(src, secret, true)
}

// GCMEncrypt encrypts plainText with secret and additionalData
func GCMEncrypt[T, E, D typez.StrOrBytes](plainText T, secret E, additionalData D) ([]byte, error) {
	cipherText, err := SaltBySecretGCMEncrypt(plainText, secret, additionalData)
	if err != nil {
		return nil, err
	}

	ret := make([]byte, hex.EncodedLen(len(cipherText)))
	hex.Encode(ret, cipherText)
	return ret, nil
}

// GCMDecrypt decrypts cipherText with secret and additionalData
func GCMDecrypt[T, E, D typez.StrOrBytes](cipherText T, secret E, additionalData D) ([]byte, error) {
	src, err := strz.HexDecode(cipherText)
	if err != nil {
		return nil, fmt.Errorf("hex decode error: %w", err)
	}

	return SaltBySecretGCMDecrypt(src, secret, additionalData, true)
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

// HybridEncrypt use RSA-OAEP + AES-GCM to encrypt plaintext
func HybridEncrypt[T typez.StrOrBytes](plaintext T, pub *rsa.PublicKey) ([]byte, error) {
	var rb [_KEY_LEN + nonceSize]byte
	_, err := rand.Read(rb[:])
	if err != nil {
		return nil, err
	}

	aesKey := rb[:_KEY_LEN] // AES-256
	encKey, err := RsaOAEPEncrypt(aesKey, []byte(nil), pub, sha256.New())
	if err != nil {
		return nil, err
	}

	nonce := rb[_KEY_LEN:] // nonceSize bytes

	buf := make([]byte, 2+len(encKey)+nonceSize+AESGCMEncryptLen(plaintext))
	binary.BigEndian.PutUint16(buf[0:2], uint16(len(encKey)))
	offset := 2
	offset += copy(buf[offset:], encKey)
	offset += copy(buf[offset:], nonce)

	err = AESGCMEncrypt(buf[offset:], strz.UnsafeStrOrBytesToBytes(plaintext), aesKey, nonce, nil)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

// HybridDecrypt use RSA-OAEP + AES-GCM to decrypt ciphertext
func HybridDecrypt[T typez.StrOrBytes](ciphertext T, pri *rsa.PrivateKey) ([]byte, error) {
	ciphertextBytes := strz.UnsafeStrOrBytesToBytes(ciphertext)
	if len(ciphertextBytes) < 2 {
		return nil, errors.New("ciphertext too short")
	}

	keyLen := int(binary.BigEndian.Uint16(ciphertextBytes[:2]))
	expectedMinLen := 2 + keyLen + nonceSize
	if len(ciphertextBytes) < expectedMinLen {
		return nil, errors.New("ciphertext too short")
	}

	encKey := ciphertextBytes[2 : 2+keyLen]
	nonce := ciphertextBytes[2+keyLen : 2+keyLen+nonceSize]
	encData := ciphertextBytes[2+keyLen+nonceSize:]

	aesKey, err := RsaOAEPDecrypt(encKey, []byte(nil), pri, sha256.New())
	if err != nil {
		return nil, err
	}

	plaintext := make([]byte, AESGCMDecryptLen(encData))
	err = AESGCMDecrypt(plaintext, encData, aesKey, nonce, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

// HybridEncryptStreamTo use RSA-OAEP + AES-CTR to encrypt stream to out
func HybridEncryptStreamTo(out io.Writer, stream io.Reader, pub *rsa.PublicKey) error {
	var rb [_KEY_LEN + aes.BlockSize]byte
	_, err := rand.Read(rb[:])
	if err != nil {
		return err
	}

	aesKey := rb[:_KEY_LEN] // AES-256
	encKey, err := RsaOAEPEncrypt(aesKey, []byte(nil), pub, sha256.New())
	if err != nil {
		return err
	}
	iv := rb[_KEY_LEN:]

	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return fmt.Errorf("NewCipher error: %w", err)
	}

	bs := make([]byte, 2)
	binary.BigEndian.PutUint16(bs, uint16(len(encKey)))
	_, _ = out.Write(bs)
	_, _ = out.Write(encKey)
	_, _ = out.Write(iv)

	encStream := cipher.NewCTR(block, iv)
	writer := &cipher.StreamWriter{S: encStream, W: out}
	_, err = io.Copy(writer, stream)
	if err != nil {
		return fmt.Errorf("copy stream error: %w", err)
	}

	return nil
}

// HybridDecryptStreamTo use RSA-OAEP + AES-CTR to decrypt stream to out
func HybridDecryptStreamTo(out io.Writer, stream io.Reader, pri *rsa.PrivateKey) error {
	bs := make([]byte, 2)
	_, err := io.ReadFull(stream, bs)
	if err != nil {
		return fmt.Errorf("read key length error: %w", err)
	}

	keyLen := int(binary.BigEndian.Uint16(bs))
	encKey := make([]byte, keyLen)
	_, err = io.ReadFull(stream, encKey)
	if err != nil {
		return fmt.Errorf("read enc key error: %w", err)
	}

	iv := make([]byte, aes.BlockSize)
	_, err = io.ReadFull(stream, iv)
	if err != nil {
		return fmt.Errorf("read iv error: %w", err)
	}

	aesKey, err := RsaOAEPDecrypt(encKey, []byte(nil), pri, sha256.New())
	if err != nil {
		return err
	}

	block, err := aes.NewCipher(aesKey)
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

// SaltBySecretCBCEncrypt
func SaltBySecretCBCEncrypt[T, E typez.StrOrBytes](plainText T, secret E) ([]byte, error) {
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

	return dst, nil
}

// SaltBySecretCBCDecrypt
func SaltBySecretCBCDecrypt[E typez.StrOrBytes](cipherText []byte, secret E, reuseCipherText bool) ([]byte, error) {
	if len(cipherText) < 2*aes.BlockSize || len(cipherText)&blockSizeMask != 0 {
		return nil, errors.New("cipherText text length illegal")
	}

	if !bytes.Equal(cipherText[:8], fixedSaltHeader) {
		return nil, errors.New("check cbc fixed header error")
	}

	var cred [_CRED_LEN]byte
	fillCred(cred[:], cipherText[8:aes.BlockSize], secret)

	key := cred[:_KEY_LEN] // 32 bytes, 256 / 8
	iv := cred[_KEY_LEN:]  // 16 bytes, same as block size

	cipherText = cipherText[aes.BlockSize:]
	dst := cipherText
	if !reuseCipherText {
		dst = make([]byte, len(dst))
	}

	n, err := AESCBCDecrypt(dst, cipherText, key, iv)
	if err != nil {
		return nil, err
	}

	return dst[:n], nil
}

// SaltBySecretGCMEncrypt
func SaltBySecretGCMEncrypt[T, E, D typez.StrOrBytes](plainText T, secret E, additionalData D) ([]byte, error) {
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

	return dst, nil
}

// SaltBySecretGCMDecrypt
func SaltBySecretGCMDecrypt[E, D typez.StrOrBytes](cipherText []byte, secret E, additionalData D, reuseCipherText bool) ([]byte, error) {
	if len(cipherText) < aes.BlockSize {
		return nil, errors.New("cipherText text length illegal")
	}

	if !bytes.Equal(cipherText[:8], fixedSaltHeader) {
		return nil, errors.New("check fixed header error")
	}

	var cred [_CRED_LEN]byte
	fillCred(cred[:], cipherText[8:aes.BlockSize], secret)

	key := cred[:_KEY_LEN] // 32 bytes, 256 / 8
	nonce := cred[_KEY_LEN : _KEY_LEN+nonceSize]

	cipherText = cipherText[aes.BlockSize:]
	dst := cipherText
	if !reuseCipherText {
		dst = make([]byte, len(dst))
	}
	err := AESGCMDecrypt(dst, cipherText, key, nonce, strz.UnsafeStrOrBytesToBytes(additionalData))
	if err != nil {
		return nil, err
	}
	return dst[:AESGCMDecryptLen(dst)], nil
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
