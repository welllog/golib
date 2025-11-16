package cryptz

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"fmt"
	"io"

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

// AESCTREncryptLen returns the length of the encrypted data
func AESCTREncryptLen[T typez.StrOrBytes](plainText T) int {
	return len(plainText)
}

// AESCTRDecryptLen returns the length of the decrypted data
func AESCTRDecryptLen[T typez.StrOrBytes](cipherText T) int {
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

// AESCBCDecrypt decrypts cipherText with key and iv
// key length must be 16, 24 or 32 bytes to select AES-128, AES-192 or AES-256.
// iv length must be 16 bytes, iv should be random to ensure safety
// dst could reuse cipherText memory
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

// AESCTREncrypt encrypts plainText with key and iv
// key length must be 16, 24 or 32 bytes to select AES-128, AES-192 or AES-256.
// iv length must be 16 bytes, iv should be random to ensure safety
// dst could reuse plainText memory
func AESCTREncrypt(dst, plainText, key, iv []byte) error {
	block, err := aes.NewCipher(key)
	if err != nil {
		return fmt.Errorf("NewCipher error: %w", err)
	}

	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(dst, plainText)
	return nil
}

// AESCTRDecrypt decrypts cipherText with key and iv
// key length must be 16, 24 or 32 bytes to select AES-128, AES-192 or AES-256.
// iv length must be 16 bytes, iv should be random to ensure safety
// dst could reuse cipherText memory
func AESCTRDecrypt(dst, cipherText, key, iv []byte) error {
	block, err := aes.NewCipher(key)
	if err != nil {
		return fmt.Errorf("NewCipher error: %w", err)
	}

	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(dst, cipherText)
	return nil
}

// AESCTRStreamEncrypt encrypts data from src reader to dst writer with key and iv
// key length must be 16, 24 or 32 bytes to select AES-128, AES-192 or AES-256.
// iv length must be 16 bytes, iv should be random to ensure safety
func AESCTRStreamEncrypt(dst io.Writer, stream io.Reader, key, iv []byte) error {
	block, err := aes.NewCipher(key)
	if err != nil {
		return fmt.Errorf("NewCipher error: %w", err)
	}

	encStream := cipher.NewCTR(block, iv)
	writer := &cipher.StreamWriter{S: encStream, W: dst}
	_, err = io.Copy(writer, stream)
	if err != nil {
		return fmt.Errorf("encrypt stream error: %w", err)
	}

	return nil
}

// AESCTRStreamDecrypt decrypts data from src reader to dst writer with key and iv
// key length must be 16, 24 or 32 bytes to select AES-128, AES-192 or AES-256.
// iv length must be 16 bytes, iv should be random to ensure safety
func AESCTRStreamDecrypt(dst io.Writer, stream io.Reader, key, iv []byte) error {
	block, err := aes.NewCipher(key)
	if err != nil {
		return fmt.Errorf("NewCipher error: %w", err)
	}

	decStream := cipher.NewCTR(block, iv)
	reader := &cipher.StreamReader{S: decStream, R: stream}
	_, err = io.Copy(dst, reader)
	if err != nil {
		return fmt.Errorf("decrypt stream error: %w", err)
	}

	return nil
}

// AESGCMEncrypt encrypts plainText with key and additionalData
// key length must be 16, 24 or 32 bytes to select AES-128, AES-192 or AES-256.
// nonce recommend 12 bytes length for better performance.
// additionalData could be nil
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

// AESGCMDecrypt decrypts cipherText with key and additionalData
// key length must be 16, 24 or 32 bytes to select AES-128, AES-192 or AES-256.
// nonce recommend 12 bytes length for better performance.
// additionalData could be nil
// dst could reuse cipherText memory, like cipherText[:AESGCMDecryptLen(cipherText)]
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

	if blockSize <= 0 || blockSize > 255 {
		return nil, errors.New("block size must be between 1 and 255")
	}

	// Calculate the padding length
	paddingLen := blockSize - (len(data) % blockSize)

	// Create the padding bytes
	var paddingBytes []byte
	if paddingLen <= aes.BlockSize {
		paddingBytes = prePadPatterns[paddingLen]
	} else {
		paddingBytes = bytes.Repeat([]byte{byte(paddingLen)}, paddingLen)
	}

	// Append the padding bytes to the input data
	return append(data, paddingBytes...), nil
}

// PKCS7UnPadding removes padding from the input data according to the PKCS#7 standard
func PKCS7UnPadding(data []byte, blockSize int) ([]byte, error) {
	// Check if input parameters are valid
	if len(data) == 0 {
		return nil, errors.New("input data cannot be empty")
	}

	if blockSize <= 0 || blockSize > 255 {
		return nil, errors.New("block size must be between 1 and 255")
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
	paddingBytes := data[len(data)-paddingLen:]
	for i := 0; i < paddingLen; i++ {
		if paddingBytes[i] != byte(paddingLen) {
			return nil, errors.New("invalid padding bytes")
		}
	}

	// Remove the padding
	return data[:len(data)-paddingLen], nil
}

func pkcs7UnPadding(data []byte) (int, error) {
	paddingLen := int(data[len(data)-1])
	if paddingLen > aes.BlockSize || paddingLen <= 0 {
		return 0, errors.New("invalid padding length")
	}
	if len(data) < paddingLen {
		return 0, errors.New("invalid padding length")
	}
	if !bytes.Equal(prePadPatterns[paddingLen], data[len(data)-paddingLen:]) {
		return 0, errors.New("invalid padding bytes")
	}
	return len(data) - paddingLen, nil
}
