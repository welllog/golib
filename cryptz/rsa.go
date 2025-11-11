package cryptz

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"hash"

	"github.com/welllog/golib/strz"
	"github.com/welllog/golib/typez"
)

// ParseRsaPublicKey parses a PEM encoded RSA public key
func ParseRsaPublicKey[E typez.StrOrBytes](pemData E) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(strz.UnsafeStrOrBytesToBytes(pemData))
	if block == nil {
		return nil, errors.New("invalid public key PEM data")
	}

	switch block.Type {
	case "PUBLIC KEY":
		key, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		rsaPub, ok := key.(*rsa.PublicKey)
		if !ok {
			return nil, errors.New("not an RSA public key")
		}
		return rsaPub, nil
	case "RSA PUBLIC KEY":
		return x509.ParsePKCS1PublicKey(block.Bytes)
	default:
		return nil, errors.New("unsupported public key type: " + block.Type)
	}
}

// ParseRsaPrivateKey parses a PEM encoded RSA private key
func ParseRsaPrivateKey[E typez.StrOrBytes](pemData E) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(strz.UnsafeStrOrBytesToBytes(pemData))
	if block == nil {
		return nil, errors.New("invalid private key PEM data")
	}

	switch block.Type {
	case "PRIVATE KEY":
		key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		rsaPri, ok := key.(*rsa.PrivateKey)
		if !ok {
			return nil, errors.New("not an RSA private key")
		}
		return rsaPri, nil
	case "RSA PRIVATE KEY":
		return x509.ParsePKCS1PrivateKey(block.Bytes)
	default:
		return nil, errors.New("unsupported private key type: " + block.Type)
	}
}

// RsaOAEPEncrypt encrypts plaintext using RSA-OAEP
// The maximum length of plaintext less than pub.Size() - 2*hash.Size() - 2
// For example, with a 2048-bit key and SHA-256, the maximum plaintext length is 190 bytes.
func RsaOAEPEncrypt[T, L typez.StrOrBytes](plaintext T, label L, pub *rsa.PublicKey, hash hash.Hash) ([]byte, error) {
	maxLen := pub.Size() - 2*hash.Size() - 2
	pt := strz.UnsafeStrOrBytesToBytes(plaintext)
	if len(pt) > maxLen {
		return nil, errors.New("plaintext length exceeds single encryption limit")
	}
	return rsa.EncryptOAEP(hash, rand.Reader, pub, pt, strz.UnsafeStrOrBytesToBytes(label))
}

// RsaOAEPDecrypt decrypts ciphertext using RSA-OAEP
func RsaOAEPDecrypt[T, L typez.StrOrBytes](ciphertext T, label L, pri *rsa.PrivateKey, hash hash.Hash) ([]byte, error) {
	return rsa.DecryptOAEP(hash, rand.Reader, pri, strz.UnsafeStrOrBytesToBytes(ciphertext), strz.UnsafeStrOrBytesToBytes(label))
}

// RsaPKCS1v15Encrypt encrypts plaintext using RSA PKCS#1 v1.5
func RsaPKCS1v15Encrypt[T typez.StrOrBytes](plaintext T, pub *rsa.PublicKey) ([]byte, error) {
	maxLen := pub.Size() - 11
	pt := strz.UnsafeStrOrBytesToBytes(plaintext)
	if len(pt) > maxLen {
		return nil, errors.New("plaintext length exceeds single encryption limit")
	}
	return rsa.EncryptPKCS1v15(rand.Reader, pub, pt)
}

// RsaPKCS1v15Decrypt decrypts ciphertext using RSA PKCS#1 v1.5
func RsaPKCS1v15Decrypt[T typez.StrOrBytes](ciphertext T, pri *rsa.PrivateKey) ([]byte, error) {
	return rsa.DecryptPKCS1v15(rand.Reader, pri, strz.UnsafeStrOrBytesToBytes(ciphertext))
}
