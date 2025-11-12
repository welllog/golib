package cryptz

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
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
	saltLen = 8
	keyLen  = 32
	credLen = 48 // BLOCK_LEN(16)+KEY_LEN(32)
)

var (
	// fixedSaltHeader
	fixedSaltHeader = []byte("Salted__")
)

// Encrypt encrypts plainText with secret (openssl aes-256-cbc implementation).
func Encrypt[T, E typez.StrOrBytes](plainText T, secret E) ([]byte, error) {
	enc, err := SaltBySecretCBCEncrypt(plainText, secret)
	if err != nil {
		return nil, err
	}

	ret := make([]byte, base64.StdEncoding.EncodedLen(len(enc)))
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
// Deprecated: use GCMEncryptV2 instead
func GCMEncrypt[T, E, D typez.StrOrBytes](plainText T, secret E, additionalData D) ([]byte, error) {
	enc, err := SaltBySecretGCMEncrypt(plainText, secret, additionalData)
	if err != nil {
		return nil, err
	}

	ret := make([]byte, hex.EncodedLen(len(enc)))
	hex.Encode(ret, enc)
	return ret, nil
}

// GCMDecrypt decrypts cipherText with secret and additionalData
// Deprecated: use GCMDecryptV2 instead
func GCMDecrypt[T, E, D typez.StrOrBytes](cipherText T, secret E, additionalData D) ([]byte, error) {
	enc, err := strz.HexDecode(cipherText)
	if err != nil {
		return nil, fmt.Errorf("hex decode error: %w", err)
	}

	return SaltBySecretGCMDecrypt(enc, secret, additionalData, true)
}

// GCMEncryptV2 encrypts plainText with secret and additionalData
func GCMEncryptV2[T, E, D typez.StrOrBytes](plainText T, secret E, additionalData D) ([]byte, error) {
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}

	derived := pbkdf2(strz.UnsafeStrOrBytesToBytes(secret), salt, 10000, keyLen+nonceSize)
	key := derived[:keyLen]
	nonce := derived[keyLen:]

	enc := make([]byte, 16+AESGCMEncryptLen(plainText))
	copy(enc[0:16], salt)

	err = AESGCMEncrypt(
		enc[16:], strz.UnsafeStrOrBytesToBytes(plainText), key, nonce, strz.UnsafeStrOrBytesToBytes(additionalData),
	)
	if err != nil {
		return nil, err
	}

	ret := make([]byte, base64.URLEncoding.EncodedLen(len(enc)))
	base64.URLEncoding.Encode(ret, enc)
	return ret, nil
}

// GCMDecryptV2 decrypts cipherText with secret and additionalData
func GCMDecryptV2[T, E, D typez.StrOrBytes](cipherText T, secret E, additionalData D) ([]byte, error) {
	enc, err := strz.Base64Decode(cipherText, base64.URLEncoding)
	if err != nil {
		return nil, err
	}

	if len(enc) < 32 {
		return nil, errors.New("cipherText too short")
	}

	salt := enc[:16]
	derived := pbkdf2(strz.UnsafeStrOrBytesToBytes(secret), salt, 10000, keyLen+nonceSize)
	key := derived[:keyLen]
	nonce := derived[keyLen:]

	cipherData := enc[16:]
	decLen := AESGCMDecryptLen(cipherData)
	if decLen < 0 {
		return nil, errors.New("cipherText length illegal")
	}

	plainText := cipherData[:decLen]

	err = AESGCMDecrypt(plainText, cipherData, key, nonce, strz.UnsafeStrOrBytesToBytes(additionalData))
	if err != nil {
		return nil, err
	}

	return plainText, nil
}

// EncryptStreamTo encrypts stream to out with secret
func EncryptStreamTo[E typez.StrOrBytes](out io.Writer, stream io.Reader, secret E) error {
	var salt [saltLen]byte
	var cred [credLen]byte
	err := fillSaltAndCred(salt[:], cred[:], secret)
	if err != nil {
		return err
	}

	key := cred[:keyLen] // 32 bytes, 256 / 8
	iv := cred[keyLen:]

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

	_, err := io.ReadFull(stream, saltHeader)
	if err != nil {
		return fmt.Errorf("read header error: %w", err)
	}

	if !bytes.Equal(saltHeader[:8], fixedSaltHeader) {
		return errors.New("check fixed header error")
	}

	var cred [credLen]byte
	fillCred(cred[:], saltHeader[8:], secret)

	key := cred[:keyLen] // 32 bytes, 256 / 8
	iv := cred[keyLen:]

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

	enc := make([]byte, 2+len(encKey)+nonceSize+AESGCMEncryptLen(plaintext))
	binary.BigEndian.PutUint16(enc[0:2], uint16(len(encKey)))
	offset := 2
	offset += copy(enc[offset:], encKey)
	offset += copy(enc[offset:], nonce)

	err = AESGCMEncrypt(enc[offset:], strz.UnsafeStrOrBytesToBytes(plaintext), aesKey, nonce, nil)
	if err != nil {
		return nil, err
	}

	ret := make([]byte, base64.URLEncoding.EncodedLen(len(enc)))
	base64.URLEncoding.Encode(ret, enc)

	return ret, nil
}

// HybridDecrypt use RSA-OAEP + AES-GCM to decrypt ciphertext
func HybridDecrypt[T typez.StrOrBytes](ciphertext T, pri *rsa.PrivateKey) ([]byte, error) {
	enc, err := strz.Base64Decode(ciphertext, base64.URLEncoding)
	if err != nil {
		return nil, err
	}

	if len(enc) < 2 {
		return nil, errors.New("ciphertext too short")
	}

	keySize := int(binary.BigEndian.Uint16(enc[:2]))
	expectedMinLen := 2 + keySize + nonceSize
	if len(enc) < expectedMinLen {
		return nil, errors.New("ciphertext too short")
	}

	encKey := enc[2 : 2+keySize]
	nonce := enc[2+keySize : 2+keySize+nonceSize]
	encData := enc[2+keySize+nonceSize:]

	aesKey, err := rsa.DecryptOAEP(sha256.New(), nil, pri, encKey, nil)
	if err != nil {
		return nil, err
	}

	plaintext := encData[:AESGCMDecryptLen(encData)]
	err = AESGCMDecrypt(plaintext, encData, aesKey, nonce, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

// HybridEncryptStreamTo use RSA-OAEP + AES-CTR to encrypt stream to out
func HybridEncryptStreamTo(out io.Writer, stream io.Reader, pub *rsa.PublicKey) error {
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

	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return fmt.Errorf("NewCipher error: %w", err)
	}

	bs := make([]byte, 2)
	binary.BigEndian.PutUint16(bs, uint16(len(encKey)))
	_, err = out.Write(bs)
	if err != nil {
		return err
	}
	_, err = out.Write(encKey)
	if err != nil {
		return err
	}
	_, err = out.Write(iv)
	if err != nil {
		return err
	}

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

	keySize := int(binary.BigEndian.Uint16(bs))
	encKey := make([]byte, keySize)
	_, err = io.ReadFull(stream, encKey)
	if err != nil {
		return fmt.Errorf("read enc key error: %w", err)
	}

	iv := make([]byte, aes.BlockSize)
	_, err = io.ReadFull(stream, iv)
	if err != nil {
		return fmt.Errorf("read iv error: %w", err)
	}

	aesKey, err := rsa.DecryptOAEP(sha256.New(), nil, pri, encKey, nil)
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
	enc := make([]byte, aes.BlockSize+AESCBCEncryptLen(plainText))
	copy(enc[0:], fixedSaltHeader)
	copy(enc[8:], salt[:])

	_ = AESCBCEncrypt(enc[aes.BlockSize:], strz.UnsafeStrOrBytesToBytes(plainText), key, iv)

	return enc, nil
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
func SaltBySecretGCMEncrypt[T, E, D typez.StrOrBytes](plainText T, secret E, additionalData D) ([]byte, error) {
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
	enc := make([]byte, aes.BlockSize+AESGCMEncryptLen(plainText))
	copy(enc[0:], fixedSaltHeader)
	copy(enc[8:], salt[:])

	_ = AESGCMEncrypt(
		enc[aes.BlockSize:],
		strz.UnsafeStrOrBytesToBytes(plainText),
		key,
		nonce,
		strz.UnsafeStrOrBytesToBytes(additionalData),
	)

	return enc, nil
}

// SaltBySecretGCMDecrypt
func SaltBySecretGCMDecrypt[E, D typez.StrOrBytes](cipherText []byte, secret E, additionalData D, reuseCipherText bool) ([]byte, error) {
	if len(cipherText) < aes.BlockSize {
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

func pbkdf2(password, salt []byte, iter, keyLen int) []byte {
	hashLen := sha256.Size
	numBlocks := (keyLen + hashLen - 1) / hashLen
	var out []byte
	for block := 1; block <= numBlocks; block++ {
		// U1 = HMAC(password, salt || INT(block))
		h := hmac.New(sha256.New, password)
		h.Write(salt)
		h.Write([]byte{
			byte(block >> 24), byte(block >> 16), byte(block >> 8), byte(block),
		})
		u := h.Sum(nil)
		t := make([]byte, len(u))
		copy(t, u)

		// U2..Uc
		for i := 1; i < iter; i++ {
			h = hmac.New(sha256.New, password)
			h.Write(u)
			u = h.Sum(nil)
			for j := range t {
				t[j] ^= u[j]
			}
		}
		out = append(out, t...)
	}
	return out[:keyLen]
}
