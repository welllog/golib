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
	"fmt"
	"hash"

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
	infoHash      []byte

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
	zeros          [sha256.Size]byte
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
	h.infoHash = labeledExtract(nil, h.suiteID, nil, labelInfoHash, nil)
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
func (h *HPKE) Seal(dst []byte, sendPrv *ecdh.PrivateKey, recvPub *ecdh.PublicKey, plaintext, aad []byte) ([]byte, error) {
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
		ss2, err := sendPrv.ECDH(recvPub)
		if err != nil {
			return nil, err
		}

		sendPubBytes := sendPrv.PublicKey().Bytes()
		kemContextSize += len(sendPubBytes)
		bufSize := 32 + len(ss1) + len(ss2) + kemContextSize
		if bufSize <= 512 {
			var tmp [512]byte
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

		bufSize := 32 + kemContextSize
		if bufSize <= 512 {
			var tmp [512]byte
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
	key, baseNonce, err := h.deriveKeys(buf, sharedSecret, mode)
	if err != nil {
		return nil, err
	}

	// Encrypt
	aead, err := h.aeadFactory(key)
	if err != nil {
		return nil, err
	}

	ret := append(dst, ephPubBytes...)

	return aead.Seal(ret, baseNonce, plaintext, aad), nil
}

// Open decrypts and authenticates the ciphertext, appending the result to dst.
// The input ciphertext is expected to be: ephPub (ephemeral pub key) || ciphertext
// If sendPub is nil, it expects Base mode.
// If sendPub is provided, it expects Auth mode.
//
// To reuse the ciphertext buffer for in-place decryption, use:
//
//	offset := hpke.KEMEncLen()
//	plaintext, err := hpke.Open(ciphertext[offset:offset], recvPrv, sendPub, ciphertext, aad)
//
// This writes the plaintext starting at offset, avoiding overwriting unread cipher data.
func (h *HPKE) Open(dst []byte, recvPrv *ecdh.PrivateKey, sendPub *ecdh.PublicKey, ciphertext, aad []byte) ([]byte, error) {
	encLen := h.kemEncLen
	if len(ciphertext) < encLen {
		return nil, ErrInvalidCipherText
	}

	ephPubBytes := ciphertext[:encLen]
	cipherData := ciphertext[encLen:]

	// 1. Parse Ephemeral Key
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
		ss2, err := recvPrv.ECDH(sendPub)
		if err != nil {
			return nil, err
		}

		sendPubBytes := sendPub.Bytes()
		kemContextSize += len(sendPubBytes)
		bufSize := 32 + len(ss1) + len(ss2) + kemContextSize
		if bufSize <= 512 {
			var tmp [512]byte
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

		bufSize := 32 + kemContextSize
		if bufSize <= 512 {
			var tmp [512]byte
			buf = tmp[:]
		} else {
			buf = make([]byte, bufSize)
		}
		kemContext = buf[32:bufSize]
		copy(kemContext, ephPubBytes)
		copy(kemContext[len(ephPubBytes):], recvPubBytes)
	}

	// 4. Derive Shared Secret
	sharedSecret := h.extractAndExpandDHKEM(buf, dh, kemContext)

	// 5. Derive Keys
	key, baseNonce, err := h.deriveKeys(buf, sharedSecret, mode)
	if err != nil {
		return nil, err
	}

	// 6. Decrypt
	aead, err := h.aeadFactory(key)
	if err != nil {
		return nil, err
	}

	return aead.Open(dst, baseNonce, cipherData, aad)
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
// buf need 1+32+32+keyLen+nonceLen bytes of space.
func (h *HPKE) deriveKeys(buf, sharedSecret []byte, mode uint8) (key, baseNonce []byte, err error) {
	// secret = LabeledExtract("", "secret", shared_secret) using HPKE suite id
	secret := labeledExtract(buf[:0], h.suiteID, nil, labelSecret, sharedSecret)

	// Create HMAC instance once and reuse it
	mac := hmac.New(sha256.New, secret)

	// KeySchedule Context: mode || psk_id_hash || info_hash
	// Reuse buf from the beginning (kemContext/sharedSecret/secret are no longer needed)
	contextEnd := 1 + len(h.pskIDHash) + len(h.infoHash)
	context := buf[0:contextEnd]
	context[0] = mode
	copy(context[1:], h.pskIDHash)
	copy(context[1+len(h.pskIDHash):], h.infoHash)

	// Reuse buf for key and baseNonce
	// labeledExpand appends to dst, so we need to pass zero-length slices
	keyBuf := buf[contextEnd:contextEnd]
	key = labeledExpand(keyBuf, mac, h.suiteID, labelKey, context, h.aeadKeyLength)

	nonceBuf := buf[contextEnd+h.aeadKeyLength : contextEnd+h.aeadKeyLength]
	baseNonce = labeledExpand(nonceBuf, mac, h.suiteID, labelBaseNonce, context, h.aeadNonceSize)

	return key, baseNonce, nil
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

	_, err := hpke.Seal(enc[:encPrefixLen], sendPrv, recvPub, strz.UnsafeStrOrBytesToBytes(plainText), strz.UnsafeStrOrBytesToBytes(aad))
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
	ret, err := hpke.Open(enc[offset:offset], recvPrv, sendPub, enc[encPrefixLen:], strz.UnsafeStrOrBytesToBytes(aad))
	return ret, err
}

func labeledExtract(dst, suiteID, salt, label, ikm []byte) []byte {
	if salt == nil {
		salt = zeros[:]
	}
	mac := hmac.New(sha256.New, salt)
	mac.Write(versionLabel)
	mac.Write(suiteID)
	mac.Write(label)
	mac.Write(ikm)
	return mac.Sum(dst)
}

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
		mac.Write(lenBytes[:])
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
