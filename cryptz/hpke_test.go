//go:build go1.20

package cryptz

import (
	"bytes"
	"crypto/cipher"
	"crypto/ecdh"
	"crypto/rand"
	"errors"
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
		t.Errorf("dst and cipherText should share the same underlying array. dst: %p, cipherText: %p",
			func() interface{} {
				if len(dst) > 0 {
					return &dst[0]
				}
				return nil
			}(),
			func() interface{} {
				if len(cipherText) > 0 {
					return &cipherText[0]
				}
				return nil
			}(),
		)
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

	sctx.SetSeq(0)
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

func TestHPKEContext_NonceReuseDefense(t *testing.T) {
	hpke := NewHPKE(ecdh.X25519())
	prv, pub, err := hpke.GenerateKey()
	testz.Nil(t, err)

	sctx, err := hpke.SetupBaseSender(pub, nil)
	testz.Nil(t, err)

	plaintext1 := []byte("secret message 1")
	plaintext2 := []byte("secret message 2")

	// 1. 第一次加密应成功
	testz.Equal(t, uint64(0), sctx.Seq())
	ct1, err := sctx.Seal(nil, plaintext1, nil)
	testz.Nil(t, err)

	// 2. 同一 seq 下未更新 seq 再次调用 Seal，必须被拦截并返回 ErrNonceReuse
	_, err = sctx.Seal(nil, plaintext2, nil)
	testz.Equal(t, ErrNonceReuse, err)

	// 3. 调用 IncrementSeq 后，可以正常加密下一条
	sctx.IncrementSeq()
	testz.Equal(t, uint64(1), sctx.Seq())
	ct2, err := sctx.Seal(nil, plaintext2, nil)
	testz.Nil(t, err)
	testz.Assert(t, len(ct2) > 0)

	// 4. 调用 SetSeq 显式指定序号（如业务重试），可以正常重新加密
	sctx.SetSeq(1)
	ct2Retry, err := sctx.Seal(nil, plaintext2, nil)
	testz.Nil(t, err)
	testz.Assert(t, len(ct2Retry) > 0)

	// 5. 验证解密端 Open 允许对同一密文进行多次幂等解密（不受防呆拦截）
	rctx, err := hpke.SetupBaseReceiver(prv, sctx.EphPublicKey(), nil)
	testz.Nil(t, err)

	// 解密第一条
	pt1, err := rctx.Open(nil, ct1, nil)
	testz.Nil(t, err)
	testz.Equal(t, plaintext1, pt1)

	// 再次解密同一条（验证幂等重试解密）
	pt1Retry, err := rctx.Open(nil, ct1, nil)
	testz.Nil(t, err)
	testz.Equal(t, plaintext1, pt1Retry)
}

func TestHPKE_API_P384(t *testing.T) {
	hpke := NewHPKE(ecdh.P384())
	testz.Equal(t, 97, hpke.KEMEncLen())
	testz.Equal(t, 48, hpke.KEMSecretLen())

	recvPrv, recvPub, err := hpke.GenerateKey()
	testz.Nil(t, err)

	sendPrv, sendPub, err := hpke.GenerateKey()
	testz.Nil(t, err)

	msg := []byte("hello p384 with sha384 DHKEM")
	ad := []byte("auth header")

	// Auth mode
	ct, err := hpke.Seal(nil, sendPrv, recvPub, nil, msg, ad)
	testz.Nil(t, err)

	pt, err := hpke.Open(nil, recvPrv, sendPub, nil, ct, ad)
	testz.Nil(t, err)
	testz.Equal(t, msg, pt)

	// Base mode
	ctBase, err := hpke.Seal(nil, nil, recvPub, nil, msg, ad)
	testz.Nil(t, err)

	ptBase, err := hpke.Open(nil, recvPrv, nil, nil, ctBase, ad)
	testz.Nil(t, err)
	testz.Equal(t, msg, ptBase)

	// Context API
	sctx, err := hpke.SetupBaseSender(recvPub, nil)
	testz.Nil(t, err)
	rctx, err := hpke.SetupBaseReceiver(recvPrv, sctx.EphPublicKey(), nil)
	testz.Nil(t, err)

	ctCtx, err := sctx.Seal(nil, msg, ad)
	testz.Nil(t, err)
	ptCtx, err := rctx.Open(nil, ctCtx, ad)
	testz.Nil(t, err)
	testz.Equal(t, msg, ptCtx)
}

func TestHPKE_API_P521(t *testing.T) {
	hpke := NewHPKE(ecdh.P521())
	testz.Equal(t, 133, hpke.KEMEncLen())
	testz.Equal(t, 64, hpke.KEMSecretLen())

	recvPrv, recvPub, err := hpke.GenerateKey()
	testz.Nil(t, err)

	sendPrv, sendPub, err := hpke.GenerateKey()
	testz.Nil(t, err)

	msg := []byte("hello p521 with sha512 DHKEM")
	ad := []byte("auth header 521")

	// Auth mode
	ct, err := hpke.Seal(nil, sendPrv, recvPub, nil, msg, ad)
	testz.Nil(t, err)

	pt, err := hpke.Open(nil, recvPrv, sendPub, nil, ct, ad)
	testz.Nil(t, err)
	testz.Equal(t, msg, pt)

	// Base mode
	ctBase, err := hpke.Seal(nil, nil, recvPub, nil, msg, ad)
	testz.Nil(t, err)

	ptBase, err := hpke.Open(nil, recvPrv, nil, nil, ctBase, ad)
	testz.Nil(t, err)
	testz.Equal(t, msg, ptBase)

	// Context API
	sctx, err := hpke.SetupBaseSender(recvPub, nil)
	testz.Nil(t, err)
	rctx, err := hpke.SetupBaseReceiver(recvPrv, sctx.EphPublicKey(), nil)
	testz.Nil(t, err)

	ctCtx, err := sctx.Seal(nil, msg, ad)
	testz.Nil(t, err)
	ptCtx, err := rctx.Open(nil, ctCtx, ad)
	testz.Nil(t, err)
	testz.Equal(t, msg, ptCtx)
}

func TestHPKE_KEM_Parameters(t *testing.T) {
	cases := []struct {
		curve     ecdh.Curve
		encLen    int
		secretLen int
	}{
		{ecdh.P256(), 65, 32},
		{ecdh.X25519(), 32, 32},
		{ecdh.P384(), 97, 48},
		{ecdh.P521(), 133, 64},
	}

	for _, tc := range cases {
		h := NewHPKE(tc.curve)
		if got := h.KEMEncLen(); got != tc.encLen {
			t.Errorf("curve %v KEMEncLen got %d, want %d", tc.curve, got, tc.encLen)
		}
		if got := h.KEMSecretLen(); got != tc.secretLen {
			t.Errorf("curve %v KEMSecretLen got %d, want %d", tc.curve, got, tc.secretLen)
		}
	}
}

func TestHPKE_Seal_PlaintextOverlap(t *testing.T) {
	h := NewHPKE(ecdh.P256())
	_, recvPub, err := h.GenerateKey()
	if err != nil {
		t.Fatal(err)
	}

	pt := bytes.Repeat([]byte("A"), 200)
	_, err = h.Seal(pt[:0], nil, recvPub, nil, pt, nil)
	if !errors.Is(err, ErrBufferOverlap) {
		t.Fatalf("expected ErrBufferOverlap when dst and plaintext overlap, got: %v", err)
	}
}
