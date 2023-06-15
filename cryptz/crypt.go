package cryptz

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"

	"github.com/welllog/golib/strz"
	"github.com/welllog/golib/typez"
)

const (
	CBC_SALT_LEN = 8
	CBC_KEY_LEN  = 32
	CBC_CRED_LEN = 48 // CBC_BLOCK_LEN(16)+CBC_KEY_LEN(32)
)

var (
	// PrePadPatterns
	prePadPatterns [aes.BlockSize + 1][]byte
	// fixedSaltHeader
	fixedSaltHeader = []byte("Salted__")
	// blockSizeMask
	blockSizeMask = aes.BlockSize - 1
)

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

func EncryptLen[T typez.StrOrBytes](s T) int {
	n := len(s) + aes.BlockSize
	return n + aes.BlockSize - (n & blockSizeMask)
}

func DecryptLen[T typez.StrOrBytes](s T) int {
	return len(s) - aes.BlockSize
}

// EncryptToBase64String openssl AES-256-CBC implementation
func EncryptToBase64String[T, E typez.StrOrBytes](plainText T, pass E) (string, error) {
	dst := make([]byte, EncryptLen(plainText))
	err := Encrypt(dst, plainText, pass)
	if err != nil {
		return "", err
	}

	return strz.Base64EncodeToString(dst, base64.StdEncoding), nil
}

// DecryptBase64ToString openssl AES-256-CBC implementation
func DecryptBase64ToString[E typez.StrOrBytes](encryptText string, pass E) (string, error) {
	src, err := strz.Base64Decode(encryptText, base64.StdEncoding)
	if err != nil {
		return "", err
	}

	if len(src) < aes.BlockSize {
		return "", errors.New("encrypt text length illegal")
	}

	// reuse src
	dst := src[aes.BlockSize:]
	n, err := Decrypt(dst, src, pass)
	return strz.UnsafeString(dst[:n]), err
}

// Encrypt encrypts plainText with pass
func Encrypt[T, E typez.StrOrBytes](dst []byte, plainText T, pass E) error {
	var salt [CBC_SALT_LEN]byte
	var cred [CBC_CRED_LEN]byte
	err := fillSaltAndCred(salt[:], cred[:], pass)
	if err != nil {
		return err
	}

	key := cred[:CBC_KEY_LEN] // 32 bytes, 256 / 8
	iv := cred[CBC_KEY_LEN:]  // 16 bytes, same as block size

	block, err := aes.NewCipher(key)
	if err != nil {
		return fmt.Errorf("NewCipher error: %w", err)
	}

	/*
		|Salted__(8 byte)|salt(8 byte)|plaintext|
	*/
	dataLen := len(plainText) + aes.BlockSize
	copy(dst[0:], fixedSaltHeader)
	copy(dst[8:], salt[:])
	copy(dst[aes.BlockSize:], plainText)

	paddingLen := aes.BlockSize - (dataLen & blockSizeMask)
	copy(dst[dataLen:], prePadPatterns[paddingLen])

	cbc := cipher.NewCBCEncrypter(block, iv)
	// encrypt from plaintext position to end
	cbc.CryptBlocks(dst[aes.BlockSize:], dst[aes.BlockSize:])
	return nil
}

// Decrypt decrypts enc with pass
func Decrypt[T, E typez.StrOrBytes](dst []byte, encryptText T, pass E) (int, error) {
	/*
		|Salted__(8 byte)|salt(8 byte)|encrypt_text|
	*/
	if len(encryptText) < aes.BlockSize || len(encryptText)&blockSizeMask != 0 {
		return 0, fmt.Errorf("encrypt text length illegal: len=%d", len(encryptText))
	}

	b := strz.UnsafeStrOrBytesToBytes(encryptText)
	saltHeader := b[:aes.BlockSize]
	if !bytes.Equal(saltHeader[:8], fixedSaltHeader) {
		return 0, errors.New("check cbc fixed header error")
	}

	var cred [CBC_CRED_LEN]byte
	fillCred(cred[:], saltHeader[8:], pass)

	key := cred[:CBC_KEY_LEN] // 32 bytes, 256 / 8
	iv := cred[CBC_KEY_LEN:]  // 16 bytes, same as block size

	block, err := aes.NewCipher(key)
	if err != nil {
		return 0, fmt.Errorf("NewCipher error: %w", err)
	}

	cbc := cipher.NewCBCDecrypter(block, iv)
	cbc.CryptBlocks(dst, b[aes.BlockSize:])
	return pkcs7UnPadding(dst)
}

// EncryptStreamTo encrypts stream to out with pass
func EncryptStreamTo[E typez.StrOrBytes](out io.Writer, stream io.Reader, pass E) error {
	var salt [CBC_SALT_LEN]byte
	var cred [CBC_CRED_LEN]byte
	err := fillSaltAndCred(salt[:], cred[:], pass)
	if err != nil {
		return err
	}

	key := cred[:CBC_KEY_LEN] // 32 bytes, 256 / 8
	iv := cred[CBC_KEY_LEN:]

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

	encStream := cipher.NewCFBEncrypter(block, iv)
	writer := &cipher.StreamWriter{S: encStream, W: out}
	_, err = io.Copy(writer, stream)
	if err != nil {
		return fmt.Errorf("copy stream error: %w", err)
	}

	return nil
}

// DecryptStreamTo decrypts stream to out with pass
func DecryptStreamTo[E typez.StrOrBytes](out io.Writer, stream io.Reader, pass E) error {
	saltHeader := make([]byte, aes.BlockSize)

	n, err := stream.Read(saltHeader)
	if err != nil {
		return fmt.Errorf("read header error: %w", err)
	}

	if n != aes.BlockSize {
		return fmt.Errorf("read header less error: n=%d", n)
	}

	var cred [CBC_CRED_LEN]byte
	fillCred(cred[:], saltHeader[8:], pass)

	key := cred[:CBC_KEY_LEN] // 32 bytes, 256 / 8
	iv := cred[CBC_KEY_LEN:]

	block, err := aes.NewCipher(key)
	if err != nil {
		return fmt.Errorf("NewCipher error: %w", err)
	}

	decStream := cipher.NewCFBDecrypter(block, iv)
	reader := &cipher.StreamReader{S: decStream, R: stream}

	_, err = io.Copy(out, reader)
	if err != nil {
		return fmt.Errorf("copy stream error: %w", err)
	}

	return nil
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

func fillSaltAndCred[E typez.StrOrBytes](salt, cred []byte, pass E) error {
	_, err := io.ReadFull(rand.Reader, salt)
	if err != nil {
		return fmt.Errorf("generate random salt error: %w", err)
	}

	fillCred(cred, salt, pass)

	return nil
}

func fillCred[E typez.StrOrBytes](cred []byte, salt []byte, pass E) {
	buf := make([]byte, 0, 16+len(pass)+len(salt))
	var prevSum [16]byte
	for i := 0; i < 3; i++ { // salted 48byte, md5 16byte, three times could fill
		n := 0 // first prevSum length is zero,so n must be zero
		if i > 0 {
			n = 16
		}
		buf = buf[:n+len(pass)+len(salt)]
		copy(buf, prevSum[:])
		copy(buf[n:], pass)
		copy(buf[n+len(pass):], salt)
		prevSum = md5.Sum(buf)        // md5(prevSum + pass + salt)
		copy(cred[i*16:], prevSum[:]) // concat every md5
	}
}
