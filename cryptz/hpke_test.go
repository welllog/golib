//go:build go1.20

package cryptz

import (
	"bytes"
	"crypto/cipher"
	"crypto/ecdh"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"reflect"
	"testing"

	"github.com/welllog/golib/testz"
)

func TestHPKE_API_P256(t *testing.T) {
	hpke := NewHPKE(ecdh.P256())

	// Receiver generates keys
	recvPrv, recvPub, err := hpke.GenerateKey()
	if err != nil {
		t.Fatal(err)
	}

	// Sender generates keys
	sendPrv, sendPub, err := hpke.GenerateKey()
	if err != nil {
		t.Fatal(err)
	}

	// Message
	msg := []byte("hello world")
	ad := []byte("header")

	// Sender Encrypts
	// Seal(dst, senderPrv, receiverPub, info, msg, ad)
	ct, err := hpke.Seal(nil, sendPrv, recvPub, nil, msg, ad)

	if err != nil {
		t.Fatal(err)
	}

	// Receiver Decrypts
	// Open(dst, receiverPrv, senderPub, info, msg, ad)
	pt, err := hpke.Open(nil, recvPrv, sendPub, nil, ct, ad)

	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(msg, pt) {
		t.Errorf("decrypted text does not match original: got %x, want %x", pt, msg)
	}
}

func TestHPKE_API_X25519(t *testing.T) {
	// Test with X25519 curve
	hpke := NewHPKE(ecdh.X25519())

	recvPrv, recvPub, err := hpke.GenerateKey()
	if err != nil {
		t.Fatal(err)
	}

	sendPrv, sendPub, err := hpke.GenerateKey()
	if err != nil {
		t.Fatal(err)
	}

	msg := []byte("hello x25519")
	ad := []byte("header")

	ct, err := hpke.Seal(nil, sendPrv, recvPub, nil, msg, ad)

	if err != nil {
		t.Fatal(err)
	}

	pt, err := hpke.Open(nil, recvPrv, sendPub, nil, ct, ad)

	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(msg, pt) {
		t.Errorf("decrypted text does not match original: got %x, want %x", pt, msg)
	}
}

func TestHPKE_API_WrongSender(t *testing.T) {
	hpke := NewHPKE(ecdh.P256())

	recvPrv, recvPub, _ := hpke.GenerateKey()
	sendPrv, _, _ := hpke.GenerateKey()
	_, pkFake, _ := hpke.GenerateKey()

	msg := []byte("secret")
	ct, err := hpke.Seal(nil, sendPrv, recvPub, nil, msg, nil)

	if err != nil {
		t.Fatal(err)
	}

	// Decrypt with wrong sender public key
	_, err = hpke.Open(nil, recvPrv, pkFake, nil, ct, nil)

	if err == nil {
		t.Fatal("expected decryption error due to wrong sender public key, got nil")
	}
}

func TestHPKE_API_LargeData(t *testing.T) {
	hpke := NewHPKE(ecdh.P256())
	recvPrv, recvPub, _ := hpke.GenerateKey()
	sendPrv, sendPub, _ := hpke.GenerateKey()

	data := make([]byte, 1024*1024) // 1MB
	rand.Read(data)

	ct, err := hpke.Seal(nil, sendPrv, recvPub, nil, data, nil)

	if err != nil {
		t.Fatal(err)
	}

	pt, err := hpke.Open(nil, recvPrv, sendPub, nil, ct, nil)

	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(data, pt) {
		t.Error("large data mismatch")
	}
}

func TestHPKE_API_DstBuffer(t *testing.T) {
	hpke := NewHPKE(ecdh.P256())
	recvPrv, recvPub, _ := hpke.GenerateKey()
	sendPrv, sendPub, _ := hpke.GenerateKey()

	msg := []byte("buffer test")

	// Pre-allocate dst
	dst := make([]byte, 0, 100)
	dst = append(dst, []byte("prefix")...)

	ct, err := hpke.Seal(dst, sendPrv, recvPub, nil, msg, nil)

	if err != nil {
		t.Fatal(err)
	}

	if !bytes.HasPrefix(ct, []byte("prefix")) {
		t.Error("ciphertext should contain prefix")
	}

	// Decrypt
	// Remove prefix for decryption input (Open expects enc||ct, but our ct has prefix)
	// Wait, Open expects the FULL message passed to it to be enc||ct.
	// But here 'ct' includes 'prefix'.
	// We need to pass the actual enc||ct part to Open.
	actualCT := ct[len("prefix"):]

	dst2 := make([]byte, 0, 100)
	pt, err := hpke.Open(dst2, recvPrv, sendPub, nil, actualCT, nil)

	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(pt, msg) {
		t.Error("decrypted text mismatch")
	}
}

func TestHPKE_CustomAEAD_Called(t *testing.T) {
	hpke := NewHPKE(ecdh.P256())

	called := false
	hpke.SetAEADFactory(func(key []byte) (cipher.AEAD, error) {
		called = true
		return defaultAESGCM(key)
	}, 16, AeadChaCha20Poly1305)

	_, recvPub, _ := hpke.GenerateKey()
	sendPrv, _, _ := hpke.GenerateKey()

	hpke.Seal(nil, sendPrv, recvPub, nil, []byte("msg"), nil)

	if !called {
		t.Error("Custom AEAD factory was not called")
	}
}

func TestHPKE_API_BaseMode(t *testing.T) {
	hpke := NewHPKE(ecdh.P256())

	// Receiver keys only
	recvPrv, recvPub, err := hpke.GenerateKey()
	if err != nil {
		t.Fatal(err)
	}

	msg := []byte("hello base mode")
	ad := []byte("header")

	// Encrypt (Seal) with nil sender keys -> Base Mode
	ct, err := hpke.Seal(nil, nil, recvPub, nil, msg, ad)

	if err != nil {
		t.Fatal(err)
	}

	// Decrypt (Open) with nil sender public key -> Base Mode
	pt, err := hpke.Open(nil, recvPrv, nil, nil, ct, ad)

	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(msg, pt) {
		t.Errorf("decrypted text does not match original: got %x, want %x", pt, msg)
	}
}

func TestHPKE_SizePrediction(t *testing.T) {
	hpke := NewHPKE(ecdh.P256())

	plaintextSize := 100
	ciphertextSize := hpke.CiphertextSize(plaintextSize)

	// P-256 Enc length = 65
	// AES-GCM Overhead = 16
	expectedCiphertextSize := 65 + plaintextSize + 16

	if ciphertextSize != expectedCiphertextSize {
		t.Errorf("CiphertextSize: got %d, want %d", ciphertextSize, expectedCiphertextSize)
	}

	calculatedPlaintextSize := hpke.PlaintextSize(ciphertextSize)
	if calculatedPlaintextSize != plaintextSize {
		t.Errorf("PlaintextSize: got %d, want %d", calculatedPlaintextSize, plaintextSize)
	}

	// Test with too short ciphertext
	if hpke.PlaintextSize(10) != -1 {
		t.Error("PlaintextSize should return -1 for short ciphertext")
	}
}

func TestHPKE_MemoryReuse(t *testing.T) {
	hpke := NewHPKE(ecdh.P256())
	_, recvPub, err := hpke.GenerateKey()
	if err != nil {
		t.Fatal(err)
	}
	sendPrv, _, err := hpke.GenerateKey()
	if err != nil {
		t.Fatal(err)
	}

	msg := []byte("hello memory reuse")
	ad := []byte("header")

	// Pre-allocate buffer
	ciphertextSize := hpke.CiphertextSize(len(msg))
	dst := make([]byte, 0, ciphertextSize)

	// Encrypt
	ct, err := hpke.Seal(dst, sendPrv, recvPub, nil, msg, ad)

	if err != nil {
		t.Fatal(err)
	}

	if len(ct) != ciphertextSize {
		t.Errorf("Ciphertext length mismatch: got %d, want %d", len(ct), ciphertextSize)
	}

	// Check if dst was reused (capacity should be same if no reallocation happened)
	if cap(ct) != cap(dst) {
		t.Logf("Note: Buffer might have been reallocated if not enough capacity or implementation detail changed. Cap dst: %d, Cap ct: %d", cap(dst), cap(ct))
	}
}

func TestHPKE_DecryptReuseBuffer(t *testing.T) {
	hpke := NewHPKE(ecdh.P256())
	recvPrv, recvPub, err := hpke.GenerateKey()
	if err != nil {
		t.Fatal(err)
	}
	sendPrv, _, err := hpke.GenerateKey()
	if err != nil {
		t.Fatal(err)
	}

	msg := []byte("hello in-place decryption")
	ad := []byte("header")

	// Encrypt
	ct, err := hpke.Seal(nil, sendPrv, recvPub, nil, msg, ad)

	if err != nil {
		t.Fatal(err)
	}

	// Decrypt in-place: use ct's buffer as dst
	// We pass ct[:0] as dst, so it appends to the start of ct's backing array.
	// Since ct contains [enc || ciphertext], and we write plaintext (which is shorter) starting at 0,
	// while reading starts at len(enc), this should be safe and efficient.
	pt, err := hpke.Open(ct[:0], recvPrv, sendPrv.PublicKey(), nil, ct, ad)

	if err != nil {
		t.Fatal(err)
	}

	if string(pt) != string(msg) {
		t.Errorf("got %s, want %s", pt, msg)
	}

	// Verify memory reuse
	// Check if pt and ct share the same backing array
	if cap(pt) != cap(ct) {
		t.Logf("Note: Capacities differ, might have reallocated. Cap pt: %d, Cap ct: %d", cap(pt), cap(ct))
	}
}

func TestHPKE_CiphertextSize(t *testing.T) {
	hpke := NewHPKE(ecdh.P256())
	plainText := []byte("sample text")
	expectedSize := hpke.CiphertextSize(len(plainText))

	sprv, spub, _ := hpke.GenerateKey()
	rpriv, rpub, _ := hpke.GenerateKey()

	dst := make([]byte, 0, expectedSize)
	cipherText, err := hpke.Seal(dst, sprv, rpub, nil, plainText, nil)

	if err != nil {
		t.Fatal(err)
	}

	if reflect.ValueOf(dst).Pointer() != reflect.ValueOf(cipherText).Pointer() {
		t.Errorf("dst and cipherText should share the same underlying array. dst: %p, cipherText: %p", &dst[0], &cipherText[0])
	}

	ret, err := hpke.Open(dst, rpriv, spub, nil, cipherText, nil)

	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(ret, plainText) {
		t.Errorf("Decrypted text mismatch: got %s, want %s", ret, plainText)
	}
}

func TestHPKE_EncryptDecrypt(t *testing.T) {
	hpke := NewHPKE(ecdh.P521())
	plainText := []byte("sample text")

	sprv, spub, _ := hpke.GenerateKey()
	rprv, rpub, _ := hpke.GenerateKey()

	cipherText, err := HPKEEncrypt(plainText, "hello", sprv, rpub, hpke)

	testz.Nil(t, err)

	ret, err := HPKEDecrypt(cipherText, "hello", rprv, spub, hpke)

	testz.Nil(t, err)

	testz.Equal(t, plainText, ret)

	cipherText, err = HPKEEncrypt(plainText, "hello2", nil, rpub, hpke)

	testz.Nil(t, err)

	ret, err = HPKEDecrypt(cipherText, "hello2", rprv, nil, hpke)

	testz.Nil(t, err)

	testz.Equal(t, plainText, ret)
}

func TestHPKE_Decrypt_LargePlaintext_Overlap(t *testing.T) {
	hpke := NewHPKE(ecdh.P256())
	// P256 encLen = 65. Prefix = 9. Total offset = 74.
	// We need plaintext > 74 bytes to trigger overlap.
	plainText := make([]byte, 100)
	for i := range plainText {
		plainText[i] = byte(i)
	}

	sprv, spub, _ := hpke.GenerateKey()
	rprv, rpub, _ := hpke.GenerateKey()

	cipherText, err := HPKEEncrypt(plainText, "header", sprv, rpub, hpke)

	testz.Nil(t, err)

	// This should panic if overlap check is triggered
	ret, err := HPKEDecrypt(cipherText, "header", rprv, spub, hpke)

	testz.Nil(t, err)
	testz.Equal(t, plainText, ret)
}

func TestHPKEContext_Seal(t *testing.T) {
	hpke := NewHPKE(ecdh.X25519())
	prv, pub, err := hpke.GenerateKey()
	testz.Nil(t, err)

	sctx, err := hpke.SetupBaseSender(pub, nil)

	testz.Nil(t, err)

	plaintext := []byte("plaintext1")
	aad := []byte("ad1")
	ret, err := sctx.Seal(nil, plaintext, aad)
	testz.Nil(t, err)

	ret, err = sctx.Open(ret[:0], ret, aad)
	testz.Nil(t, err)
	testz.Equal(t, plaintext, ret)

	rctx, err := hpke.SetupBaseReceiver(prv, sctx.EphPublicKey(), nil)

	testz.Nil(t, err)

	ret, err = sctx.Seal(nil, plaintext, aad)
	testz.Nil(t, err)

	ret, err = rctx.Open(ret[:0], ret, aad)
	testz.Nil(t, err)
	testz.Equal(t, plaintext, ret)

	sctx.IncrementSeq()
	rctx.IncrementSeq()
	ret, err = sctx.Seal(nil, plaintext, aad)
	testz.Nil(t, err)

	ret, err = rctx.Open(ret[:0], ret, aad)
	testz.Nil(t, err)
	testz.Equal(t, plaintext, ret)

	sprv, spub, err := hpke.GenerateKey()
	testz.Nil(t, err)
	sctx, err = hpke.SetupAuthSender(pub, sprv, nil)

	testz.Nil(t, err)
	rctx, err = hpke.SetupAuthReceiver(prv, spub, sctx.EphPublicKey(), nil)

	testz.Nil(t, err)

	plaintext = make([]byte, 256)
	for i := 0; i < 100; i++ {
		ret, err = sctx.Seal(nil, plaintext, nil)
		testz.Nil(t, err)

		ret, err = rctx.Open(ret[:0], ret, nil)
		testz.Nil(t, err)
		testz.Equal(t, plaintext, ret)

		sctx.IncrementSeq()
		rctx.IncrementSeq()
	}

}

func TestDeriveKeys(t *testing.T) {
	sharedSecret := []byte("shared secret")
	info := []byte("context info")
	mode := uint8(0)

	buf1 := make([]byte, 256)
	buf2 := make([]byte, 256)
	key1, nonce1, err := deriveKeys1(buf1, sharedSecret, info, mode)
	testz.Nil(t, err)

	key2, nonce2, err := deriveKeys2(buf2, sharedSecret, info, mode)
	testz.Nil(t, err)

	fmt.Println(bytes.Equal(key1, key2))
	fmt.Println(bytes.Equal(nonce1, nonce2))
	fmt.Println(len(nonce1), len(nonce2))
}

func deriveKeys1(buf, sharedSecret, info []byte, mode uint8) (key, baseNonce []byte, err error) {
	pskIDHash := make([]byte, 32)
	// RFC 9180: secret = LabeledExtract(shared_secret, "secret", psk)
	// For Base mode, psk is empty. shared_secret is the SALT, psk is the IKM!
	secret := labeledExtract(buf[:0], nil, sharedSecret, labelSecret, nil)

	// Create HMAC instance once and reuse it
	mac := hmac.New(sha256.New, secret)

	infoHash := labeledExtract(nil, nil, nil, labelInfoHash, info)

	// KeySchedule Context: mode || psk_id_hash || info_hash
	// Allocate fresh context to avoid overlap issues
	context := make([]byte, 1+len(pskIDHash)+len(infoHash))
	context[0] = mode
	copy(context[1:], pskIDHash)
	copy(context[1+len(pskIDHash):], infoHash)

	// Allocate fresh key and nonce
	key = labeledExpand(nil, mac, nil, labelKey, context, 16)
	mac.Reset()
	baseNonce = labeledExpand(nil, mac, nil, labelBaseNonce, context, 12)

	return key, baseNonce, nil
}

func deriveKeys2(buf, sharedSecret, info []byte, mode uint8) (key, baseNonce []byte, err error) {
	pskIDHash := make([]byte, 32)
	// RFC 9180: secret = LabeledExtract(shared_secret, "secret", psk)
	// For Base mode, psk is empty. shared_secret is the SALT, psk is the IKM!
	secret := labeledExtract(buf[:0], nil, sharedSecret, labelSecret, nil)

	// Create HMAC instance once and reuse it
	mac := hmac.New(sha256.New, secret)

	// KeySchedule Context: mode || psk_id_hash || info_hash
	contextBegin := 16 + 12
	contextEnd := contextBegin + 1 + len(pskIDHash) + sha256.Size
	context := buf[contextBegin:contextEnd]
	context[0] = mode
	copy(context[1:], pskIDHash)
	// info_hash = LabeledExtract(nil, suite_id, "info_hash", info)
	// info_hash is written into context
	labeledExtract(context[1+len(pskIDHash):1+len(pskIDHash)], nil, nil, labelInfoHash, info)

	// Allocate fresh key and nonce
	keyBuf := buf[:0]
	key = labeledExpand(keyBuf, mac, nil, labelKey, context, 16)

	mac.Reset()
	nonceBuf := buf[16:16]
	baseNonce = labeledExpand(nonceBuf, mac, nil, labelBaseNonce, context, 12)

	return key, baseNonce, nil

}
