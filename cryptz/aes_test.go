package cryptz

import (
	"bytes"
	"testing"

	"github.com/welllog/golib/testz"
)

func TestAESCBCEncrypt(t *testing.T) {
	key := []byte("0123456789abcdef")
	iv := []byte("0123456789abcdef")

	text := []byte("hello world")
	dst := make([]byte, AESCBCEncryptLen(text))
	if err := AESCBCEncrypt(dst, text, key, iv); err != nil {
		t.Fatalf("AESCBCEncrypt error: %s", err)
	}

	n, err := AESCBCDecrypt(dst, dst, key, iv)
	if err != nil {
		t.Fatalf("AESCBCDecrypt error: %s", err)
	}

	if string(dst[:n]) != string(text) {
		t.Fatalf("AESCBCDecrypt(%s) != %s", string(dst[:n]), string(text))
	}
}

func TestAESCBCDecrypt(t *testing.T) {
	key := []byte("0123456789abcdef")
	iv := []byte("0123456789abcdef")

	str := "hello world"
	text := []byte(str)
	dst := append(text, bytes.Repeat([]byte{0}, AESCBCEncryptLen(text)-len(text))...)
	if err := AESCBCEncrypt(dst, text, key, iv); err != nil {
		t.Fatalf("AESCBCEncrypt error: %s", err)
	}

	n, err := AESCBCDecrypt(dst, dst, key, iv)
	if err != nil {
		t.Fatalf("AESCBCDecrypt error: %s", err)
	}

	if string(dst[:n]) != str {
		t.Fatalf("AESCBCDecrypt(%s) != %s", string(dst[:n]), str)
	}
}

func TestAESGCMEncrypt(t *testing.T) {
	key := []byte("0123456789abcdef")

	str := "hello world"
	text := []byte(str)
	nonce := []byte("abc")
	enc := append(text, bytes.Repeat([]byte{0}, AESGCMEncryptLen(text)-len(text))...)
	err := AESGCMEncrypt(enc, text, key, nonce, nil)
	if err != nil {
		t.Fatalf("AESGCMEncrypt error: %s", err)
	}

	err = AESGCMDecrypt(enc[:AESGCMDecryptLen(enc)], enc, key, nonce, nil)
	if err != nil {
		t.Fatalf("AESGCMDecrypt error: %s", err)
	}

	if string(enc[:AESGCMDecryptLen(enc)]) != str {
		t.Fatalf("AESGCMDecrypt(%s) != %s", string(enc[:AESGCMDecryptLen(enc)]), str)
	}
}

func TestAESGCMDecrypt(t *testing.T) {
	key := []byte("0123456789abcdef")

	str := "hello world"
	text := []byte(str)
	nonce := []byte("abc")
	addition := []byte("addition")
	enc := append(text, bytes.Repeat([]byte{0}, AESGCMEncryptLen(text)-len(text))...)
	err := AESGCMEncrypt(enc, text, key, nonce, addition)
	if err != nil {
		t.Fatalf("AESGCMEncrypt error: %s", err)
	}

	err = AESGCMDecrypt(enc[:AESGCMDecryptLen(enc)], enc, key, nonce, addition)
	if err != nil {
		t.Fatalf("AESGCMDecrypt error: %s", err)
	}

	if string(enc[:AESGCMDecryptLen(enc)]) != str {
		t.Fatalf("AESGCMDecrypt(%s) != %s", string(enc[:AESGCMDecryptLen(enc)]), str)
	}
}

func TestAESGCMEncryptDecrypt(t *testing.T) {
	keys := []struct {
		key   []byte
		nonce []byte
	}{
		{[]byte("0123456789abcdef"), []byte("abc")},
		{[]byte("0123456789abcdef"), []byte("0123456789ab")},
		{[]byte("0123456789abcdef01234567"), []byte("0123456789ab")},
		{[]byte("0123456789abcdef0123456789abcdef"), []byte("0123456789ab")},
	}

	texts := [][]byte{
		[]byte("hello world"),
		[]byte("1231"),
		[]byte("ovevvdcq"),
		[]byte("👋，世界"),
		[]byte("what happen"),
		[]byte(""),
		[]byte("0123456789abcdef0123456789abcdef"),
		[]byte("0123456789abcdef0123456789abcdef0123456789abcdef"),
		[]byte("0123456789abcdef+-/;.,,,.'[]123!@#$%^&*()_+0123456789abcdef0123456789abcdef0123456789abcdef"),
	}

	for _, k := range keys {
		for _, text := range texts {
			enc := make([]byte, AESGCMEncryptLen(text))
			err := AESGCMEncrypt(enc, text, k.key, k.nonce, []byte("demo"))
			testz.Nil(t, err)

			err = AESGCMDecrypt(enc[:AESGCMDecryptLen(enc)], enc, k.key, k.nonce, []byte("demo"))
			testz.Nil(t, err)

			testz.Equal(t, string(enc[:AESGCMDecryptLen(enc)]), string(text))
		}
	}
}

func TestAESCTREncryptDecrypt(t *testing.T) {
	key := []byte("0123456789abcdef")
	iv := []byte("abcdef9876543210")

	str := "hello world"
	text := []byte(str)
	err := AESCTREncrypt(text, text, key, iv)
	if err != nil {
		t.Fatalf("AESCTREncrypt error: %s", err)
	}

	err = AESCTRDecrypt(text, text, key, iv)
	if err != nil {
		t.Fatalf("AESCTRDecrypt error: %s", err)
	}

	if string(text) != str {
		t.Fatalf("AESCTRDecrypt(%s) != %s", string(text), str)
	}
}
