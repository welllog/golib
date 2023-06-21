package cryptz

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"fmt"

	"github.com/welllog/golib/typez"
)

const (
	// blockSizeMask
	blockSizeMask = aes.BlockSize - 1
	// gcmTagSize default tag size
	gcmTagSize = 16
	// nonceSize default nonce size
	nonceSize = 12
)

// PrePadPatterns
var prePadPatterns [aes.BlockSize + 1][]byte

func init() {
	for i := 0; i < len(prePadPatterns); i++ {
		prePadPatterns[i] = bytes.Repeat([]byte{byte(i)}, i)
	}
	/*
		[]
		[1]
		[2 2]
		[3 3 3]
		[4 4 4 4]
		[5 5 5 5 5]
		[6 6 6 6 6 6]
		[7 7 7 7 7 7 7]
		[8 8 8 8 8 8 8 8]
		[9 9 9 9 9 9 9 9 9]
		[10 10 10 10 10 10 10 10 10 10]
		[11 11 11 11 11 11 11 11 11 11 11]
		[12 12 12 12 12 12 12 12 12 12 12 12]
		[13 13 13 13 13 13 13 13 13 13 13 13 13]
		[14 14 14 14 14 14 14 14 14 14 14 14 14 14]
		[15 15 15 15 15 15 15 15 15 15 15 15 15 15 15]
		[16 16 16 16 16 16 16 16 16 16 16 16 16 16 16 16]
	*/
}

// AESCBCEncryptLen returns the length of the encrypted data
func AESCBCEncryptLen[T typez.StrOrBytes](plainText T) int {
	return len(plainText) + aes.BlockSize - (len(plainText) & blockSizeMask)
}

// AESCBCDecryptLen returns the length of the decrypted data
func AESCBCDecryptLen[T typez.StrOrBytes](cipherText T) int {
	return len(cipherText)
}

// AESGCMEncryptLen returns the length of the encrypted data
func AESGCMEncryptLen[T typez.StrOrBytes](plainText T) int {
	return len(plainText) + gcmTagSize
}

// AESGCMDecryptLen returns the length of the decrypted data
func AESGCMDecryptLen[T typez.StrOrBytes](cipherText T) int {
	return len(cipherText) - gcmTagSize
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
func AESCBCDecrypt(dst, cipherText, key, iv []byte) (int, error) {
	if len(cipherText) < aes.BlockSize || len(cipherText)&blockSizeMask != 0 {
		return 0, errors.New("cipherText length illegal")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return 0, fmt.Errorf("NewCipher error: %w", err)
	}

	cbc := cipher.NewCBCDecrypter(block, iv)
	cbc.CryptBlocks(dst, cipherText)

	return pkcs7UnPadding(dst)
}

// AESGCMEncrypt encrypts plainText with key and additionalData
// key length must be 16, 24 or 32 bytes to select AES-128, AES-192 or AES-256.
// plainText could pre grow tagSize(default 16) so dst could reuse plainText memory
func AESGCMEncrypt(dst, plainText, key, nonce, additionalData []byte) error {
	block, err := aes.NewCipher(key)
	if err != nil {
		return fmt.Errorf("NewCipher error: %w", err)
	}

	gcm, err := cipher.NewGCMWithNonceSize(block, len(nonce))
	if err != nil {
		return fmt.Errorf("NewGCM error: %w", err)
	}

	gcm.Seal(dst[:0], nonce, plainText, additionalData)
	return nil
}

// AESGCMDecrypt decrypts encryptText with key and additionalData
// key length must be 16, 24 or 32 bytes to select AES-128, AES-192 or AES-256.
// dst could reuse encryptText memory, like encryptText[:AESGCMDecryptLen(encryptText)]
func AESGCMDecrypt(dst, cipherText, key, nonce, additionalData []byte) error {
	block, err := aes.NewCipher(key)
	if err != nil {
		return fmt.Errorf("NewCipher error: %w", err)
	}

	gcm, err := cipher.NewGCMWithNonceSize(block, len(nonce))
	if err != nil {
		return fmt.Errorf("NewGCM error: %w", err)
	}

	_, err = gcm.Open(dst[:0], nonce, cipherText, additionalData)
	if err != nil {
		return fmt.Errorf("GCM Open error: %w", err)
	}

	return nil
}

// PKCS5Padding adds padding to the input data according to the PKCS#5 standard
func PKCS5Padding(data []byte) ([]byte, error) {
	return PKCS7Padding(data, 8)
}

// PKCS5UnPadding removes padding from the input data according to the PKCS#5 standard
func PKCS5UnPadding(data []byte) ([]byte, error) {
	return PKCS7UnPadding(data, 8)
}

// PKCS7Padding adds padding to the input data according to the PKCS#7 standard
func PKCS7Padding(data []byte, blockSize int) ([]byte, error) {
	// Check if input parameters are valid
	if len(data) == 0 {
		return nil, errors.New("input data cannot be empty")
	}

	if blockSize <= 0 {
		return nil, errors.New("block size must be a positive integer")
	}

	// Calculate the padding length
	paddingLen := blockSize - (len(data) % blockSize)

	// Create the padding bytes
	paddingBytes := bytes.Repeat([]byte{byte(paddingLen)}, paddingLen)

	// Append the padding bytes to the input data
	return append(data, paddingBytes...), nil
}

// PKCS7UnPadding removes padding from the input data according to the PKCS#7 standard
func PKCS7UnPadding(data []byte, blockSize int) ([]byte, error) {
	// Check if input parameters are valid
	if len(data) == 0 {
		return nil, errors.New("input data cannot be empty")
	}

	if blockSize <= 0 {
		return nil, errors.New("block size must be a positive integer")
	}

	if len(data)%blockSize != 0 {
		return nil, errors.New("input data length must be a multiple of block size")
	}

	// Get the padding length
	paddingLen := int(data[len(data)-1])

	// Check if the padding length is valid
	if paddingLen <= 0 || paddingLen > blockSize {
		return nil, errors.New("invalid padding length")
	}

	// Check if the padding bytes are correct
	paddingBytes := bytes.Repeat([]byte{byte(paddingLen)}, paddingLen)
	if !bytes.Equal(data[len(data)-paddingLen:], paddingBytes) {
		return nil, errors.New("invalid padding bytes")
	}

	// Remove the padding
	return data[:len(data)-paddingLen], nil
}

func pkcs7UnPadding(data []byte) (int, error) {
	paddingLen := int(data[len(data)-1])
	if paddingLen > aes.BlockSize || paddingLen <= 0 {
		return 0, errors.New("invalid padding length")
	}
	if !bytes.Equal(prePadPatterns[paddingLen], data[len(data)-paddingLen:]) {
		return 0, errors.New("invalid padding bytes")
	}
	return len(data) - paddingLen, nil
}
