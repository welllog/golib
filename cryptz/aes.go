package cryptz

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"

	"github.com/welllog/golib/typez"
)

// AESCBCEncryptLen returns the length of the encrypted data
func AESCBCEncryptLen[T typez.StrOrBytes](plainText T) int {
	return len(plainText) + aes.BlockSize - (len(plainText) & blockSizeMask)
}

// AESCBCEncrypt encrypts plainText with key and iv
// key length must be 16, 24 or 32 bytes to select AES-128, AES-192 or AES-256.
// iv length must be 16 bytes, iv should be random to ensure safety
// plainText could pre grow padding length, so dst could reuse plainText memory
func AESCBCEncrypt(dst, plainText, key, iv []byte) error {
	block, err := aes.NewCipher(key)
	if err != nil {
		return fmt.Errorf("NewCipher error: %w", err)
	}

	paddingLen := aes.BlockSize - (len(plainText) & blockSizeMask)
	copy(dst, plainText)
	copy(dst[len(plainText):], prePadPatterns[paddingLen])

	cbc := cipher.NewCBCEncrypter(block, iv)
	cbc.CryptBlocks(dst, dst)

	return nil
}

// AESCBCDecrypt decrypts encryptText with key and iv
// key length must be 16, 24 or 32 bytes to select AES-128, AES-192 or AES-256.
// iv length must be 16 bytes, iv should be random to ensure safety
// dst could reuse encryptText memory
func AESCBCDecrypt(dst, encryptText, key, iv []byte) (int, error) {
	if len(encryptText) < aes.BlockSize || len(encryptText)&blockSizeMask != 0 {
		return 0, fmt.Errorf("encrypt text length illegal: len=%d", len(encryptText))
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return 0, fmt.Errorf("NewCipher error: %w", err)
	}

	cbc := cipher.NewCBCDecrypter(block, iv)
	cbc.CryptBlocks(dst, encryptText)

	return pkcs7UnPadding(dst)
}

// AESGCMEncrypt encrypts plainText with key and additionalData
// key length must be 16, 24 or 32 bytes to select AES-128, AES-192 or AES-256.
// dst could reuse plainText memory
func AESGCMEncrypt(dst, plainText, key, additionalData []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("NewCipher error: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("NewGCM error: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("ReadFull error: %w", err)
	}

	return gcm.Seal(dst, nonce, plainText, additionalData), nil
}

// AESGCMDecrypt decrypts encryptText with key and additionalData
// key length must be 16, 24 or 32 bytes to select AES-128, AES-192 or AES-256.
// dst could reuse encryptText memory
func AESGCMDecrypt(dst, encryptText, key, additionalData []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("NewCipher error: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("NewGCM error: %w", err)
	}

	nonceSize := gcm.NonceSize()
	if len(encryptText) < nonceSize {
		return nil, fmt.Errorf("encrypt text length illegal: len=%d", len(encryptText))
	}

	// to reuse encryptText storage as dst, we need check if dst is encryptText
	if bytes.Equal(dst, encryptText) {
		dst = dst[nonceSize:]
	}
	return gcm.Open(dst, encryptText[:nonceSize], encryptText[nonceSize:], additionalData)
}
