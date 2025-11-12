package cryptz

import (
	"os"
	"testing"

	"github.com/welllog/golib/testz"
)

func TestRsaHybridEncryptDecrypt(t *testing.T) {
	var str = "dadsadsajwbd1wibdiw1bidw1"

	pub, err := ParseRsaPublicKey(pubKey)
	testz.Nil(t, err)

	prv, err := ParseRsaPrivateKey(prvKey)
	testz.Nil(t, err)

	enc, err := RsaHybridEncrypt(str, "", pub)
	testz.Nil(t, err)

	dec, err := RsaHybridDecrypt(enc, "", prv)
	testz.Nil(t, err)

	testz.Equal(t, str, string(dec))
}

func TestRsaHybridEncryptStreamTo(t *testing.T) {
	pub, err := ParseRsaPublicKey(pubKey)
	testz.Nil(t, err)

	prv, err := ParseRsaPrivateKey(prvKey)
	testz.Nil(t, err)

	text := "hello, this is a test!!!"
	fileName := "STREAM_HYBRID.txt"
	if err := os.WriteFile(fileName, []byte(text), 0644); err != nil {
		t.Fatalf("os.WriteFile(%s) error: %s", fileName, err)
	}

	f1, err := os.Open(fileName)
	if err != nil {
		t.Fatalf("os.Open(%s) error: %s", fileName, err)
	}

	f2, err := os.Create(fileName + ".enc")
	if err != nil {
		t.Fatalf("os.Create(%s.enc) error: %s", fileName, err)
	}

	if err := RsaHybridEncryptStream(f2, f1, pub); err != nil {
		t.Fatalf("RsaHybridEncryptStream error: %s", err)
	}

	_ = f1.Close()
	_ = f2.Close()

	f3, err := os.Open(fileName + ".enc")
	if err != nil {
		t.Fatalf("os.Open(%s.enc) error: %s", fileName, err)
	}

	f4, err := os.Create(fileName)
	if err != nil {
		t.Fatalf("os.Create(%s) error: %s", fileName, err)
	}

	if err := RsaHybridDecryptStream(f4, f3, prv); err != nil {
		t.Fatalf("RsaHybridDecryptStream error: %s", err)
	}

	_ = f3.Close()
	_ = f4.Close()

	b, err := os.ReadFile(fileName)
	if err != nil {
		t.Fatalf("os.ReadFile(%s) error: %s", fileName, err)
	}
	if string(b) != text {
		t.Fatalf("RsaHybridDecryptStream error, expected: %s, actual: %s", text, string(b))
	}

	_ = os.Remove(fileName)
	_ = os.Remove(fileName + ".enc")
}
