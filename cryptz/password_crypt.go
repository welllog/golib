package cryptz

import (
	"bufio"
	"bytes"
	"crypto/aes"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"

	"github.com/welllog/golib/strz"
	"github.com/welllog/golib/typez"
)

const (
	magicLen     = 8
	versionLen   = 1
	encPrefixLen = magicLen + versionLen
	saltLen16    = 16
	//   text pass v1 layout: magic(8)|version(1)|salt(16)|nonce(12)|keyDeriverHeader|cipherText|tag(16)
	// stream pass v1 layout: magic(8)|version(1)|salt(16)|iv(16)|keyDeriverHeader|cipherStream
	passV1PrefixLen = encPrefixLen + saltLen16 + nonceSize
)

var (
	ErrInvalidCipherText       = errors.New("invalid cipher text")
	ErrInvalidKeyDeriverHeader = errors.New("invalid header")
	ErrInvalidCipherStream     = errors.New("invalid cipher stream")

	passMagic       = [magicLen]byte{'W', 'L', 'P', 'A', 'S', 'S', 'W', 'D'}
	passStreamMagic = [magicLen]byte{'W', 'L', 'P', 'W', 'D', 'S', 'T', 'R'}

	keyDeriverRegistry = make(map[[IDLen]byte]KeyDeriver, 4)

	defKeyDeriver = PBKDF2KeyDeriver{
		Iter: 10_000,
		Hash: SHA256,
	}
)

func init() {
	RegisterKeyDeriver(PBKDF2KeyDeriver{})
}

func RegisterKeyDeriver(deriver KeyDeriver) {
	keyDeriverRegistry[deriver.ID()] = deriver
}

// PasswordEncrypt using password to encrypt plainText with additional data ad.
// It generates a random salt and nonce internally. It uses the provided keyDeriver
// to derive the encryption key from the password and salt.
// If keyDeriver is nil, it uses the default PBKDF2 with SHA256 and 10,000 iterations.
// If using custom key deriver, make sure to call RegisterKeyDeriver register it before decryption.
// The output is base64 URL encoded cipher text.
func PasswordEncrypt[T, P, D typez.StrOrBytes](plainText T, password P, ad D, keyDeriver KeyDeriver) ([]byte, error) {
	if keyDeriver == nil {
		keyDeriver = defKeyDeriver
	}

	// version 1 ------------
	// magic(8)|version(1)|salt(16)|nonce(12)|keyDeriverHeader|cipherText|tag(16)
	deriverHeader := keyDeriver.Header()
	encLen := passV1PrefixLen + len(deriverHeader) + len(plainText) + gcmTagSize
	base64EncLen := base64.RawURLEncoding.EncodedLen(encLen)

	ret := make([]byte, base64EncLen)
	enc := ret[base64EncLen-encLen:]

	copy(enc, passMagic[:])
	enc[magicLen] = 1 // version 1

	_, err := rand.Read(enc[encPrefixLen:passV1PrefixLen])
	if err != nil {
		return nil, err
	}
	salt := enc[encPrefixLen : encPrefixLen+saltLen16]
	nonce := enc[encPrefixLen+saltLen16 : passV1PrefixLen]

	offset := passV1PrefixLen
	offset += copy(enc[offset:], deriverHeader)

	key := keyDeriver.Key(strz.UnsafeStrOrBytesToBytes(password), salt, keyLen)

	err = AESGCMEncrypt(
		enc[offset:],
		strz.UnsafeStrOrBytesToBytes(plainText),
		key,
		nonce,
		strz.UnsafeStrOrBytesToBytes(ad),
	)
	if err != nil {
		return nil, err
	}
	// version 1 ------------

	base64.RawURLEncoding.Encode(ret, enc)
	return ret, nil
}

// PasswordDecrypt using password to decrypt cipherText with additional data ad.
// It uses the key deriver header in the cipherText to restore the key deriver
// and derive the encryption key from the password and salt.
// If using custom key deriver, make sure to call RegisterKeyDeriver register it before decryption.
// The input cipherText is expected to be base64 URL encoded.
func PasswordDecrypt[T, P, D typez.StrOrBytes](cipherText T, password P, ad D) ([]byte, error) {
	enc, err := strz.Base64Decode(cipherText, base64.RawURLEncoding)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidCipherText, err)
	}

	// magic(8)|version(1)
	if len(enc) < encPrefixLen || !bytes.Equal(enc[:magicLen], passMagic[:]) {
		return nil, ErrInvalidCipherText
	}

	switch enc[magicLen] {
	case 1:
		return passwordDecryptV1(enc, password, ad)
	default:
		return nil, ErrInvalidCipherText
	}
}

// PasswordEncryptStream using password to encrypt data from stream and write to dst.
// It generates a random salt and iv internally. It uses the provided keyDeriver
// to derive the encryption key from the password and salt.
// If keyDeriver is nil, it uses the default PBKDF2 with SHA256 and 10,000 iterations.
// If using custom key deriver, make sure to call RegisterKeyDeriver register it before decryption.
func PasswordEncryptStream[P typez.StrOrBytes](dst io.Writer, stream io.Reader, password P, keyDeriver KeyDeriver) error {
	if keyDeriver == nil {
		keyDeriver = defKeyDeriver
	}

	// version 1 ------------
	// magic(8)|version(1)|salt(16)|iv(16)|keyDeriverHeader|cipherStream
	deriverHeader := keyDeriver.Header()

	saltAndIv := make([]byte, saltLen16+aes.BlockSize)
	_, err := rand.Read(saltAndIv)
	if err != nil {
		return err
	}

	w := bufio.NewWriter(dst)

	_, _ = w.Write(passStreamMagic[:])
	_ = w.WriteByte(1) // version 1
	_, _ = w.Write(saltAndIv)
	_, _ = w.Write(deriverHeader)

	salt := saltAndIv[:saltLen16]
	iv := saltAndIv[saltLen16:]
	key := keyDeriver.Key(strz.UnsafeStrOrBytesToBytes(password), salt, keyLen)

	err = AESCTRStreamEncrypt(w, stream, key, iv)
	if err != nil {
		return err
	}

	return w.Flush()
}

// PasswordDecryptStream using password to decrypt data from stream and write to dst.
// It uses the key deriver header in the stream to restore the key deriver
// and derive the encryption key from the password and salt.
// If using custom key deriver, make sure to call RegisterKeyDeriver register it before decryption.
func PasswordDecryptStream[P typez.StrOrBytes](dst io.Writer, stream io.Reader, password P) error {
	buf := make([]byte, saltLen16+aes.BlockSize+IDLen)
	r := bufio.NewReader(stream)

	_, err := io.ReadFull(r, buf[:encPrefixLen])
	if err != nil {
		return fmt.Errorf("%w: %s", ErrInvalidCipherStream, err)
	}

	if !bytes.Equal(buf[:magicLen], passStreamMagic[:]) {
		return fmt.Errorf("%w: %s", ErrInvalidCipherStream, "invalid header")
	}

	if buf[magicLen] != 1 {
		return fmt.Errorf("%w: %s", ErrInvalidCipherStream, "unsupported version")
	}

	_, err = io.ReadFull(r, buf)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrInvalidCipherStream, err)
	}

	salt := buf[:saltLen16]
	iv := buf[saltLen16 : saltLen16+aes.BlockSize]
	deriverID := buf[saltLen16+aes.BlockSize : saltLen16+aes.BlockSize+IDLen]

	var cmpID [8]byte
	copy(cmpID[:], deriverID)
	deriver, ok := keyDeriverRegistry[cmpID]
	if !ok {
		return ErrInvalidKeyDeriverHeader
	}

	deriverHeaderLen := deriver.HeaderLen()
	if deriverHeaderLen < IDLen {
		return fmt.Errorf("%w: %s", ErrInvalidCipherStream, "deriver header length too short")
	}
	deriverHeader := make([]byte, deriverHeaderLen)
	_, err = io.ReadFull(r, deriverHeader[IDLen:])
	if err != nil {
		return fmt.Errorf("%w: %s", ErrInvalidCipherStream, err)
	}

	copy(deriverHeader, deriverID)
	encDeriver, err := deriver.Restore(deriverHeader)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidKeyDeriverHeader, err)
	}
	key := encDeriver.Key(strz.UnsafeStrOrBytesToBytes(password), salt, keyLen)

	return AESCTRStreamDecrypt(dst, r, key, iv)
}

func passwordDecryptV1[P, D typez.StrOrBytes](enc []byte, password P, ad D) ([]byte, error) {
	// magic(8)|version(1)|salt(16)|nonce(12)|keyDeriverHeader|cipherText|tag(16)
	if len(enc) < passV1PrefixLen+IDLen+gcmTagSize {
		return nil, ErrInvalidCipherText
	}

	salt := enc[encPrefixLen : encPrefixLen+saltLen16]
	nonce := enc[encPrefixLen+saltLen16 : passV1PrefixLen]
	deriverID := enc[passV1PrefixLen : passV1PrefixLen+IDLen]

	var cmpID [8]byte
	copy(cmpID[:], deriverID)
	deriver, ok := keyDeriverRegistry[cmpID]
	if !ok {
		return nil, ErrInvalidKeyDeriverHeader
	}

	headerLen := deriver.HeaderLen()
	if len(enc) < passV1PrefixLen+headerLen+gcmTagSize {
		return nil, ErrInvalidCipherText
	}

	encDeriver, err := deriver.Restore(enc[passV1PrefixLen : passV1PrefixLen+headerLen])
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidKeyDeriverHeader, err)
	}

	key := encDeriver.Key(strz.UnsafeStrOrBytesToBytes(password), salt, keyLen)
	cipherData := enc[passV1PrefixLen+headerLen:]
	plainText := cipherData[:len(cipherData)-gcmTagSize]

	err = AESGCMDecrypt(plainText, cipherData, key, nonce, strz.UnsafeStrOrBytesToBytes(ad))
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidCipherText, err)
	}

	return plainText, nil
}
