//go:build go1.20

package cryptz

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdh"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"hash"
	"math"

	"github.com/welllog/golib/strz"
	"github.com/welllog/golib/typez"
)

var (
	// hpkeMagic WLHPKSEG
	hpkeMagic = [magicLen]byte{'W', 'L', 'H', 'P', 'K', 'S', 'E', 'G'}
)

// AEADFactory creates a cipher.AEAD from a key.
type AEADFactory func(key []byte) (cipher.AEAD, error)

// HPKE implements Hybrid Public Key Encryption (HPKE) as per RFC 9180.
// This implementation currently supports:
// - KEM: DHKEM(Curve, HKDF-SHA256)
// - KDF: HKDF-SHA256
// - AEAD: AES-128-GCM (default, but swappable)
// It supports "Base" mode (anonymous) and "Auth" mode (sender authentication).
type HPKE struct {
	curve         ecdh.Curve
	aeadFactory   AEADFactory
	aeadKeyLength int
	aeadOverhead  int
	aeadNonceSize int
	kemEncLen     int
	suiteID       []byte
	kemSuiteID    []byte
	pskIDHash     []byte

	// Current Suite Configuration
	kemID  uint16
	kdfID  uint16
	aeadID uint16
}

const (
	// ModeBase is the mode identifier for Base mode (0x00).
	ModeBase = uint8(0x00)
	// ModeAuth is the mode identifier for Auth mode (0x02).
	ModeAuth = uint8(0x02)

	// Algorithm Identifiers (RFC 9180 Section 7)
	kemP256HKDFSHA256   = uint16(0x0010)
	kemP384HKDFSHA384   = uint16(0x0011)
	kemP521HKDFSHA512   = uint16(0x0012)
	kemX25519HKDFSHA256 = uint16(0x0020)

	kdfHKDFSHA256 = uint16(0x0001)

	// AEAD Algorithm Identifiers (RFC 9180 Section 7.3)
	AeadAES128GCM        = uint16(0x0001)
	AeadAES256GCM        = uint16(0x0002)
	AeadChaCha20Poly1305 = uint16(0x0003)

	defaultBufSize = 256
)

var (
	versionLabel   = []byte("HPKE-v1") // RFC9180
	labelEAEPrk    = []byte("eae_prk")
	labelShared    = []byte("shared_secret")
	labelPSKIDHash = []byte("psk_id_hash")
	labelInfoHash  = []byte("info_hash")
	labelSecret    = []byte("secret")
	labelKey       = []byte("key")
	labelBaseNonce = []byte("base_nonce")

	ErrDHZero = errors.New("hpke: dh shared secret is all-zero")
)

// defaultAESGCM is the default AEAD factory using AES-128-GCM.
func defaultAESGCM(key []byte) (cipher.AEAD, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	return cipher.NewGCM(block)
}

// NewHPKE creates a new HPKE instance with the specified curve.
// Default suite: HKDF-SHA256, AES-128-GCM.
func NewHPKE(curve ecdh.Curve) *HPKE {
	h := &HPKE{
		curve:         curve,
		aeadFactory:   defaultAESGCM,
		aeadKeyLength: 16, // AES-128
		aeadOverhead:  16, // AES-GCM tag size
		aeadNonceSize: 12, // AES-GCM nonce size
		kdfID:         kdfHKDFSHA256,
		aeadID:        AeadAES128GCM,
	}

	// Determine KEM ID
	switch curve {
	case ecdh.P256():
		h.kemID = kemP256HKDFSHA256
		h.kemEncLen = 65
	case ecdh.X25519():
		h.kemID = kemX25519HKDFSHA256
		h.kemEncLen = 32
	case ecdh.P384():
		h.kemID = kemP384HKDFSHA384
		h.kemEncLen = 97
	case ecdh.P521():
		h.kemID = kemP521HKDFSHA512
		h.kemEncLen = 133
	default:
		// For unknown curves (e.g. new ones added to stdlib), try to determine length dynamically.
		prv, err := curve.GenerateKey(rand.Reader)
		if err != nil {
			panic("hpke: failed to generate key to determine kemEncLen")
		}
		h.kemEncLen = len(prv.PublicKey().Bytes())
		h.kemID = kemP256HKDFSHA256 // Best effort fallback
	}

	h.computeSuiteID()
	h.kemSuiteID = make([]byte, 5)
	copy(h.kemSuiteID[0:3], "KEM")
	binary.BigEndian.PutUint16(h.kemSuiteID[3:], h.kemID)

	return h
}

func (h *HPKE) computeSuiteID() {
	// Construct Suite ID: "HPKE" || I2OSP(kem_id, 2) || I2OSP(kdf_id, 2) || I2OSP(aead_id, 2)
	const size = 4 + 2 + 2 + 2
	if len(h.suiteID) < size {
		h.suiteID = make([]byte, size)
	}
	copy(h.suiteID, "HPKE")
	binary.BigEndian.PutUint16(h.suiteID[4:], h.kemID)
	binary.BigEndian.PutUint16(h.suiteID[6:], h.kdfID)
	binary.BigEndian.PutUint16(h.suiteID[8:], h.aeadID)

	h.pskIDHash = labeledExtract(nil, h.suiteID, nil, labelPSKIDHash, nil)
}

// SetAEADFactory sets a custom AEAD factory with specific key length and algorithm ID.
// The nonce length will be determined automatically from the AEAD instance.
// The overhead will be determined by instantiating a dummy AEAD.
func (h *HPKE) SetAEADFactory(factory AEADFactory, keyLen int, aeadID uint16) error {
	h.aeadFactory = factory
	h.aeadKeyLength = keyLen
	h.aeadID = aeadID

	// Recompute Suite ID with new AEAD ID
	h.computeSuiteID()

	// Instantiate a dummy AEAD to determine overhead
	dummyKey := make([]byte, keyLen)
	aead, err := factory(dummyKey)
	if err != nil {
		return err
	}
	h.aeadOverhead = aead.Overhead()
	h.aeadNonceSize = aead.NonceSize()
	return nil
}

// GenerateKey generates a new key pair using the configured curve.
func (h *HPKE) GenerateKey() (*ecdh.PrivateKey, *ecdh.PublicKey, error) {
	prv, err := h.curve.GenerateKey(rand.Reader)
	if err != nil {
		return nil, nil, err
	}
	return prv, prv.PublicKey(), nil
}

// Seal encrypts and authenticates the plaintext, appending the result to dst.
// The result is: ephPub (ephemeral pub key) || ciphertext
// If sendPrv is nil, it uses Base mode (anonymous).
// If sendPrv is provided, it uses Auth mode (authenticated).
func (h *HPKE) Seal(dst []byte, sendPrv *ecdh.PrivateKey, recvPub *ecdh.PublicKey, info, plaintext, aad []byte) ([]byte, error) {

	ephPrv, ephPub, err := h.GenerateKey()
	if err != nil {
		return nil, err
	}

	ephPubBytes := ephPub.Bytes()
	recvPubBytes := recvPub.Bytes()

	// kemContext: ephPubBytes + recvPubBytes + [sendPubBytes]
	kemContextSize := len(ephPubBytes) + len(recvPubBytes)
	var buf, kemContext, dh []byte
	mode := ModeBase

	if sendPrv != nil {
		// Auth Mode
		mode = ModeAuth
		ss1, err := ephPrv.ECDH(recvPub)
		if err != nil {
			return nil, err
		}
		if h.kemID == kemX25519HKDFSHA256 && isAllZero(ss1) {
			return nil, ErrDHZero
		}

		ss2, err := sendPrv.ECDH(recvPub)
		if err != nil {
			return nil, err
		}
		if h.kemID == kemX25519HKDFSHA256 && isAllZero(ss2) {
			return nil, ErrDHZero
		}

		sendPubBytes := sendPrv.PublicKey().Bytes()
		kemContextSize += len(sendPubBytes)
		bufSize := 32 + len(ss1) + len(ss2) + kemContextSize
		if bufSize <= defaultBufSize {
			var tmp [defaultBufSize]byte
			buf = tmp[:]
		} else {
			buf = make([]byte, bufSize)
		}
		dh = buf[32 : 32+len(ss1)+len(ss2)]
		kemContext = buf[32+len(ss1)+len(ss2) : bufSize]
		copy(dh, ss1)
		copy(dh[len(ss1):], ss2)
		copy(kemContext, ephPubBytes)
		copy(kemContext[len(ephPubBytes):], recvPubBytes)
		copy(kemContext[len(ephPubBytes)+len(recvPubBytes):], sendPubBytes)
	} else {
		// Base Mode
		dh, err = ephPrv.ECDH(recvPub)
		if err != nil {
			return nil, err
		}
		if h.kemID == kemX25519HKDFSHA256 && isAllZero(dh) {
			return nil, ErrDHZero
		}

		bufSize := 32 + kemContextSize
		if bufSize <= defaultBufSize {
			var tmp [defaultBufSize]byte
			buf = tmp[:]
		} else {
			buf = make([]byte, bufSize)
		}
		kemContext = buf[32:bufSize]
		copy(kemContext, ephPubBytes)
		copy(kemContext[len(ephPubBytes):], recvPubBytes)
	}

	// Derive Shared Secret
	sharedSecret := h.extractAndExpandDHKEM(buf, dh, kemContext)

	// Derive Keys
	key, baseNonce, err := h.deriveKeys(buf, sharedSecret, info, mode)

	if err != nil {
		return nil, err
	}

	// Encrypt
	aead, err := h.aeadFactory(key)
	if err != nil {
		return nil, err
	}

	ret := append(dst, ephPubBytes...)
	nonce := nonceForSeq(buf[h.aeadKeyLength+h.aeadNonceSize:], baseNonce, 0)
	return aead.Seal(ret, nonce, plaintext, aad), nil
}

// Open decrypts and authenticates the ciphertext, appending the result to dst.
// The input ciphertext is expected to be: ephPub (ephemeral pub key) || ciphertext
// If sendPub is nil, it expects Base mode.
// If sendPub is provided, it expects Auth mode.
//
// To reuse the ciphertext buffer for in-place decryption, use:
//
//	offset := hpke.KEMEncLen()
//	plaintext, err := hpke.Open(ciphertext[offset:offset], recvPrv, sendPub, info, ciphertext, aad)
//
// This writes the plaintext starting at offset, avoiding overwriting unread cipher data.
func (h *HPKE) Open(dst []byte, recvPrv *ecdh.PrivateKey, sendPub *ecdh.PublicKey, info, ciphertext, aad []byte) ([]byte, error) {

	encLen := h.kemEncLen
	if len(ciphertext) < encLen {
		return nil, ErrInvalidCipherText
	}

	ephPubBytes := ciphertext[:encLen]
	cipherData := ciphertext[encLen:]

	// Parse Ephemeral Key
	ephPub, err := h.curve.NewPublicKey(ephPubBytes)
	if err != nil {
		return nil, err
	}

	recvPubBytes := recvPrv.PublicKey().Bytes()
	// kemContext: ephPubBytes + recvPubBytes + [sendPubBytes]
	kemContextSize := len(ephPubBytes) + len(recvPubBytes)
	var buf, kemContext, dh []byte
	mode := ModeBase

	if sendPub != nil {
		// Auth Mode
		mode = ModeAuth
		ss1, err := recvPrv.ECDH(ephPub)
		if err != nil {
			return nil, err
		}
		if h.kemID == kemX25519HKDFSHA256 && isAllZero(ss1) {
			return nil, ErrDHZero
		}
		ss2, err := recvPrv.ECDH(sendPub)
		if err != nil {
			return nil, err
		}
		if h.kemID == kemX25519HKDFSHA256 && isAllZero(ss2) {
			return nil, ErrDHZero
		}

		sendPubBytes := sendPub.Bytes()
		kemContextSize += len(sendPubBytes)
		bufSize := 32 + len(ss1) + len(ss2) + kemContextSize
		if bufSize <= defaultBufSize {
			var tmp [defaultBufSize]byte
			buf = tmp[:]
		} else {
			buf = make([]byte, bufSize)
		}
		dh = buf[32 : 32+len(ss1)+len(ss2)]
		kemContext = buf[32+len(ss1)+len(ss2) : bufSize]
		copy(dh, ss1)
		copy(dh[len(ss1):], ss2)
		copy(kemContext, ephPubBytes)
		copy(kemContext[len(ephPubBytes):], recvPubBytes)
		copy(kemContext[len(ephPubBytes)+len(recvPubBytes):], sendPubBytes)
	} else {
		// Base Mode
		dh, err = recvPrv.ECDH(ephPub)
		if err != nil {
			return nil, err
		}
		if h.kemID == kemX25519HKDFSHA256 && isAllZero(dh) {
			return nil, ErrDHZero
		}

		bufSize := 32 + kemContextSize
		if bufSize <= defaultBufSize {
			var tmp [defaultBufSize]byte
			buf = tmp[:]
		} else {
			buf = make([]byte, bufSize)
		}
		kemContext = buf[32:bufSize]
		copy(kemContext, ephPubBytes)
		copy(kemContext[len(ephPubBytes):], recvPubBytes)
	}

	// Derive Shared Secret
	sharedSecret := h.extractAndExpandDHKEM(buf, dh, kemContext)

	// Derive Keys
	key, baseNonce, err := h.deriveKeys(buf, sharedSecret, info, mode)

	if err != nil {
		return nil, err
	}

	// Decrypt
	aead, err := h.aeadFactory(key)
	if err != nil {
		return nil, err
	}
	nonce := nonceForSeq(buf[h.aeadKeyLength+h.aeadNonceSize:], baseNonce, 0)
	return aead.Open(dst, nonce, cipherData, aad)
}

// CiphertextSize calculates the expected size of the ciphertext for a given plaintext size.
// It includes the encapsulated key length and the AEAD overhead.
func (h *HPKE) CiphertextSize(plaintextSize int) int {
	return h.kemEncLen + plaintextSize + h.aeadOverhead
}

// PlaintextSize calculates the expected size of the plaintext for a given ciphertext size.
// It subtracts the encapsulated key length and the AEAD overhead.
// Returns -1 if the ciphertext is too short.
func (h *HPKE) PlaintextSize(ciphertextSize int) int {
	overhead := h.kemEncLen + h.aeadOverhead
	if ciphertextSize < overhead {
		return -1
	}
	return ciphertextSize - overhead
}

// KEMEncLen returns the length of the encapsulated key (ephemeral public key).
// This is useful for calculating the safe offset when reusing the ciphertext buffer
// for in-place decryption. See Open() documentation for usage example.
func (h *HPKE) KEMEncLen() int {
	return h.kemEncLen
}

type HPKEContext struct {
	key       []byte
	baseNonce []byte
	ephPub    []byte
	nonceBuf  []byte
	seq       uint64
	aead      cipher.AEAD
	kemEncLen int
}

func (h *HPKE) SetupBaseSender(recvPub *ecdh.PublicKey, info []byte) (*HPKEContext, error) {

	ephPrv, ephPub, err := h.GenerateKey()
	if err != nil {
		return nil, err
	}

	dh, err := ephPrv.ECDH(recvPub)
	if err != nil {
		return nil, err
	}
	if h.kemID == kemX25519HKDFSHA256 && isAllZero(dh) {
		return nil, ErrDHZero
	}

	ephPubBytes := ephPub.Bytes()
	recvPubBytes := recvPub.Bytes()

	// kemContext: ephPubBytes + recvPubBytes
	kemContextSize := len(ephPubBytes) + len(recvPubBytes)
	var buf, kemContext []byte

	bufSize := 32 + kemContextSize
	if bufSize <= defaultBufSize {
		var tmp [defaultBufSize]byte
		buf = tmp[:]
	} else {
		buf = make([]byte, bufSize)
	}
	kemContext = buf[32:bufSize]
	copy(kemContext, ephPubBytes)
	copy(kemContext[len(ephPubBytes):], recvPubBytes)

	// Derive Shared Secret
	sharedSecret := h.extractAndExpandDHKEM(buf, dh, kemContext)

	// Derive Keys
	key, baseNonce, err := h.deriveKeys(buf, sharedSecret, info, ModeBase)

	if err != nil {
		return nil, err
	}

	return h.buildHPKEContext(buf, key, baseNonce, ephPubBytes)
}

func (h *HPKE) SetupBaseReceiver(recvPrv *ecdh.PrivateKey, ephPubBytes, info []byte) (*HPKEContext, error) {

	if len(ephPubBytes) != h.kemEncLen {
		return nil, errors.New("hpke: invalid ephemeral public key")
	}

	ephPub, err := h.curve.NewPublicKey(ephPubBytes)
	if err != nil {
		return nil, err
	}

	dh, err := recvPrv.ECDH(ephPub)
	if err != nil {
		return nil, err
	}
	if h.kemID == kemX25519HKDFSHA256 && isAllZero(dh) {
		return nil, ErrDHZero
	}

	recvPubBytes := recvPrv.PublicKey().Bytes()
	// kemContext: ephPubBytes + recvPubBytes
	kemContextSize := len(ephPubBytes) + len(recvPubBytes)
	var buf, kemContext []byte

	bufSize := 32 + kemContextSize
	if bufSize <= defaultBufSize {
		var tmp [defaultBufSize]byte
		buf = tmp[:]
	} else {
		buf = make([]byte, bufSize)
	}
	kemContext = buf[32:bufSize]
	copy(kemContext, ephPubBytes)
	copy(kemContext[len(ephPubBytes):], recvPubBytes)

	// Derive Shared Secret
	sharedSecret := h.extractAndExpandDHKEM(buf, dh, kemContext)

	// Derive Keys
	key, baseNonce, err := h.deriveKeys(buf, sharedSecret, info, ModeBase)

	if err != nil {
		return nil, err
	}

	return h.buildHPKEContext(buf, key, baseNonce, ephPubBytes)
}

func (h *HPKE) SetupAuthSender(recvPub *ecdh.PublicKey, sendPrv *ecdh.PrivateKey, info []byte) (*HPKEContext, error) {

	ephPrv, ephPub, err := h.GenerateKey()
	if err != nil {
		return nil, err
	}

	ss1, err := ephPrv.ECDH(recvPub)
	if err != nil {
		return nil, err
	}
	if h.kemID == kemX25519HKDFSHA256 && isAllZero(ss1) {
		return nil, ErrDHZero
	}

	ss2, err := sendPrv.ECDH(recvPub)
	if err != nil {
		return nil, err
	}
	if h.kemID == kemX25519HKDFSHA256 && isAllZero(ss2) {
		return nil, ErrDHZero
	}

	ephPubBytes := ephPub.Bytes()
	recvPubBytes := recvPub.Bytes()
	sendPubBytes := sendPrv.PublicKey().Bytes()

	// kemContext: ephPubBytes + recvPubBytes + sendPubBytes
	kemContextSize := len(ephPubBytes) + len(recvPubBytes) + len(sendPubBytes)
	var buf, kemContext, dh []byte
	bufSize := 32 + len(ss1) + len(ss2) + kemContextSize
	if bufSize <= defaultBufSize {
		var tmp [defaultBufSize]byte
		buf = tmp[:]
	} else {
		buf = make([]byte, bufSize)
	}
	dh = buf[32 : 32+len(ss1)+len(ss2)]
	kemContext = buf[32+len(ss1)+len(ss2) : bufSize]
	copy(dh, ss1)
	copy(dh[len(ss1):], ss2)
	copy(kemContext, ephPubBytes)
	copy(kemContext[len(ephPubBytes):], recvPubBytes)
	copy(kemContext[len(ephPubBytes)+len(recvPubBytes):], sendPubBytes)

	// Derive Shared Secret
	sharedSecret := h.extractAndExpandDHKEM(buf, dh, kemContext)

	// Derive Keys
	key, baseNonce, err := h.deriveKeys(buf, sharedSecret, info, ModeAuth)

	if err != nil {
		return nil, err
	}

	return h.buildHPKEContext(buf, key, baseNonce, ephPubBytes)
}

func (h *HPKE) SetupAuthReceiver(recvPrv *ecdh.PrivateKey, sendPub *ecdh.PublicKey, ephPubBytes, info []byte) (*HPKEContext, error) {

	if len(ephPubBytes) != h.kemEncLen {
		return nil, errors.New("hpke: invalid ephemeral public key")
	}

	// Parse Ephemeral Key
	ephPub, err := h.curve.NewPublicKey(ephPubBytes)
	if err != nil {
		return nil, err
	}

	ss1, err := recvPrv.ECDH(ephPub)
	if err != nil {
		return nil, err
	}
	if h.kemID == kemX25519HKDFSHA256 && isAllZero(ss1) {
		return nil, ErrDHZero
	}
	ss2, err := recvPrv.ECDH(sendPub)
	if err != nil {
		return nil, err
	}
	if h.kemID == kemX25519HKDFSHA256 && isAllZero(ss2) {
		return nil, ErrDHZero
	}

	recvPubBytes := recvPrv.PublicKey().Bytes()
	sendPubBytes := sendPub.Bytes()
	// kemContext: ephPubBytes + recvPubBytes + sendPubBytes
	kemContextSize := len(ephPubBytes) + len(recvPubBytes) + len(sendPubBytes)
	var buf, kemContext, dh []byte

	bufSize := 32 + len(ss1) + len(ss2) + kemContextSize
	if bufSize <= defaultBufSize {
		var tmp [defaultBufSize]byte
		buf = tmp[:]
	} else {
		buf = make([]byte, bufSize)
	}
	dh = buf[32 : 32+len(ss1)+len(ss2)]
	kemContext = buf[32+len(ss1)+len(ss2) : bufSize]
	copy(dh, ss1)
	copy(dh[len(ss1):], ss2)
	copy(kemContext, ephPubBytes)
	copy(kemContext[len(ephPubBytes):], recvPubBytes)
	copy(kemContext[len(ephPubBytes)+len(recvPubBytes):], sendPubBytes)

	// Derive Shared Secret
	sharedSecret := h.extractAndExpandDHKEM(buf, dh, kemContext)

	// Derive Keys
	key, baseNonce, err := h.deriveKeys(buf, sharedSecret, info, ModeAuth)

	if err != nil {
		return nil, err
	}

	return h.buildHPKEContext(buf, key, baseNonce, ephPubBytes)
}

// Seal seals the plaintext using the sender context.
// It won't prepend the ephemeral public key, as it's already known from context.
func (c *HPKEContext) Seal(dst, plaintext, aad []byte) ([]byte, error) {
	nonce := nonceForSeq(c.nonceBuf, c.baseNonce, c.seq)
	return c.aead.Seal(dst, nonce, plaintext, aad), nil
}

// Open opens the ciphertext using the receiver context.
// The ciphertext don't include the ephemeral public key.
// It could reuse ciphertext as dst for in-place decryption, like ciphertext[:0]
func (c *HPKEContext) Open(dst, ciphertext, aad []byte) ([]byte, error) {
	nonce := nonceForSeq(c.nonceBuf, c.baseNonce, c.seq)
	return c.aead.Open(dst, nonce, ciphertext, aad)
}

func (c *HPKEContext) IncrementSeq() uint64 {
	c.seq++
	return c.seq
}

func (c *HPKEContext) SetSeq(seq uint64) {
	c.seq = seq
}

func (c *HPKEContext) KemEncLen() int {
	return c.kemEncLen
}

func (c *HPKEContext) EphPublicKey() []byte {
	return c.ephPub
}

// HPKEEncrypt encrypts the plaintext using HPKE and returns the base64 encoded ciphertext.
func HPKEEncrypt[T, D typez.StrOrBytes](plainText T, aad D, sendPrv *ecdh.PrivateKey, recvPub *ecdh.PublicKey, hpke *HPKE) ([]byte, error) {

	encLen := encPrefixLen + hpke.CiphertextSize(len(plainText))
	base64EncLen := base64.RawURLEncoding.EncodedLen(encLen)
	ret := make([]byte, base64EncLen)
	enc := ret[base64EncLen-encLen:]

	// magic(8)|version(1)
	copy(enc, hpkeMagic[:])
	enc[magicLen] = 1 // version 1

	_, err := hpke.Seal(enc[:encPrefixLen], sendPrv, recvPub, nil, strz.UnsafeStrOrBytesToBytes(plainText), strz.UnsafeStrOrBytesToBytes(aad))

	if err != nil {
		return nil, err
	}

	base64.RawURLEncoding.Encode(ret, enc)
	return ret, nil
}

// HPKEDecrypt decrypts the base64 encoded ciphertext using HPKE.
func HPKEDecrypt[T, D typez.StrOrBytes](ciphertext T, aad D, recvPrv *ecdh.PrivateKey, sendPub *ecdh.PublicKey, hpke *HPKE) ([]byte, error) {

	enc, err := strz.Base64Decode(ciphertext, base64.RawURLEncoding)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidCipherText, err)
	}

	// magic(8)|version(1)
	if len(enc) < encPrefixLen || !bytes.Equal(enc[:magicLen], hpkeMagic[:]) {
		return nil, ErrInvalidCipherText
	}

	if enc[magicLen] != 1 {
		return nil, ErrInvalidCipherText
	}

	if hpke.PlaintextSize(len(enc)-encPrefixLen) < 0 {
		return nil, ErrInvalidCipherText
	}

	offset := encPrefixLen + hpke.KEMEncLen()
	return hpke.Open(enc[offset:offset], recvPrv, sendPub, nil, enc[encPrefixLen:], strz.UnsafeStrOrBytesToBytes(aad))

}

// --- HKDF Implementation (RFC 5869) with HPKE Labeling ---

// extractAndExpandDHKEM performs the Extract-and-Expand step for DHKEM.
// buf need 32 bytes of space.
func (h *HPKE) extractAndExpandDHKEM(buf, dh, kemContext []byte) []byte {
	// eae_prk = LabeledExtract(salt=0, "eae_prk", dh)  (uses kemSuite in LabeledExtract)
	eaePrk := labeledExtract(buf[:0], h.kemSuiteID, nil, labelEAEPrk, dh)

	// shared_secret = LabeledExpand(eae_prk, "shared_secret", kemContext, Nh)
	mac := hmac.New(sha256.New, eaePrk)
	shared := labeledExpand(buf[:0], mac, h.kemSuiteID, labelShared, kemContext, sha256.Size)
	return shared
}

// deriveKeys derives the AEAD key and baseNonce from the shared secret.
// buf need 1+32+32+keyLen+ceil(nonceLen/32)*32 bytes of space.
func (h *HPKE) deriveKeys(buf, sharedSecret, info []byte, mode uint8) (key, baseNonce []byte, err error) {
	// RFC 9180: secret = LabeledExtract(shared_secret, "secret", psk)
	// For Base mode, psk is empty. shared_secret is the SALT, psk is the IKM!
	secret := labeledExtract(buf[:0], h.suiteID, sharedSecret, labelSecret, nil)

	// Create HMAC instance once and reuse it
	mac := hmac.New(sha256.New, secret)

	// caller ensures buf min len is 256, so we can slice safely.
	// KeySchedule Context: mode || psk_id_hash || info_hash
	// labeledExtract use buf space need hashSize * n
	contextBegin := int(math.Ceil(float64(h.aeadNonceSize)/float64(sha256.Size)))*sha256.Size + h.aeadKeyLength
	contextEnd := contextBegin + 1 + len(h.pskIDHash) + sha256.Size
	context := buf[contextBegin:contextEnd]
	context[0] = mode
	copy(context[1:], h.pskIDHash)
	// info_hash is written into context
	labeledExtract(context[1+len(h.pskIDHash):1+len(h.pskIDHash)], h.suiteID, nil, labelInfoHash, info)

	// Allocate fresh key and nonce
	keyBuf := buf[:0]
	key = labeledExpand(keyBuf, mac, h.suiteID, labelKey, context, h.aeadKeyLength)

	mac.Reset()
	nonceBuf := buf[h.aeadKeyLength:h.aeadKeyLength]
	baseNonce = labeledExpand(nonceBuf, mac, h.suiteID, labelBaseNonce, context, h.aeadNonceSize)

	return key, baseNonce, nil

}

func (h *HPKE) buildHPKEContext(buf, key, baseNonce, ephPubBytes []byte) (*HPKEContext, error) {
	aead, err := h.aeadFactory(key)
	if err != nil {
		return nil, err
	}

	offset1 := h.aeadKeyLength + h.aeadNonceSize
	offset2 := copy(buf[offset1:], ephPubBytes) + offset1
	return &HPKEContext{
		key:       key,
		baseNonce: baseNonce,
		ephPub:    buf[offset1:offset2],
		nonceBuf:  buf[offset2 : offset2+h.aeadNonceSize],
		seq:       0,
		aead:      aead,
		kemEncLen: h.kemEncLen,
	}, nil
}

// dst must have 32 bytes of space at least.
func labeledExtract(dst, suiteID, salt, label, ikm []byte) []byte {
	// RFC 9180: salt should be empty string (zero-length), not zeros
	if salt == nil {
		salt = []byte{} // Empty slice, not zeros[:]
	}
	mac := hmac.New(sha256.New, salt)
	mac.Write(versionLabel)
	mac.Write(suiteID)
	mac.Write(label)
	mac.Write(ikm)
	return mac.Sum(dst)
}

// dst must have ceil(l/HashLen) * HashLen bytes of space at least.
func labeledExpand(dst []byte, mac hash.Hash, suiteID, label, info []byte, l int) []byte {
	// LabeledInfo construction (virtual, we write parts directly)
	// Format: I2OSP(L, 2) || suite_id || label || info

	// mac is already initialized with the key (prk)

	startLen := len(dst)
	targetLen := startLen + l

	var t []byte
	ctr := byte(1)

	// Prepare constant parts of LabeledInfo
	var lenBytes [2]byte
	binary.BigEndian.PutUint16(lenBytes[:], uint16(l))
	var ctrBuf [1]byte

	for len(dst) < targetLen {
		mac.Reset()
		mac.Write(t)
		// Write LabeledInfo parts
		// RFC 9180: I2OSP(L, 2) || "HPKE-v1" || suiteID || label || info
		mac.Write(lenBytes[:])
		mac.Write(versionLabel) // CRITICAL: Was missing!
		mac.Write(suiteID)
		mac.Write(label)
		mac.Write(info)

		ctrBuf[0] = ctr
		mac.Write(ctrBuf[:])

		// Append directly to dst
		dst = mac.Sum(dst)

		// Update t to be the last block (which was just appended)
		// The hash size is 32 for SHA256
		// Note: mac.Sum appends 32 bytes.
		// If we needed less than 32 bytes for the last block, dst will be larger than targetLen.
		// We need to be careful with t.

		// t is the last 32 bytes appended.
		t = dst[len(dst)-32:]

		ctr++
	}
	return dst[:targetLen]
}

func isAllZero(b []byte) bool {
	for _, v := range b {
		if v != 0 {
			return false
		}
	}
	return true
}

// nonceForSeq: baseNonce XOR seq (low-order bytes XOR)
func nonceForSeq(buf, baseNonce []byte, seq uint64) []byte {
	copy(buf, baseNonce)
	n := buf[0:len(baseNonce)]
	// XOR low-order bytes with seq (big-endian)
	for i := 0; i < 8 && i < len(n); i++ {
		b := byte(seq & 0xff)
		n[len(n)-1-i] ^= b
		seq >>= 8
	}
	return n
}
