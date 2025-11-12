package cryptz

import (
	"os"
	"testing"

	"github.com/welllog/golib/testz"
)

func TestEncrypt(t *testing.T) {
	tests1 := []struct {
		text string
		pass string
	}{
		{"hello, this is a test!!!", "whaterror"},
		{"👋，世界", "测试"},
	}
	for _, tt := range tests1 {
		enc, err := Encrypt(tt.text, tt.pass)
		if err != nil {
			t.Fatalf("EncryptToString(%s, %s) error: %s", tt.text, tt.pass, err)
		}

		dec, err := Decrypt(enc, tt.pass)
		if err != nil {
			t.Fatalf("DecryptFromString(%s, %s) error: %s", enc, tt.pass, err)
		}

		testz.Equal(t, tt.text, string(dec))
	}

	tests2 := []struct {
		text []byte
		pass []byte
	}{
		{[]byte("hello, this is a test!!!"), []byte("whaterror")},
		{[]byte("??，世界"), []byte("测试")},
	}
	for _, tt := range tests2 {
		enc, err := Encrypt(tt.text, tt.pass)
		if err != nil {
			t.Fatalf("Encrypt(%s, %s) error: %s", tt.text, tt.pass, err)
		}

		dec, err := Decrypt(enc, tt.pass)
		if err != nil {
			t.Fatalf("Decrypt(%s, %s) error: %s", enc, tt.pass, err)
		}

		testz.Equal(t, string(tt.text), string(dec))
	}

	tests3 := []struct {
		text string
		pass []byte
	}{
		{"hello, this is a test!!!", []byte("whaterror")},
		{"??，世界", []byte("测试")},
	}
	for _, tt := range tests3 {
		enc, err := Encrypt([]byte(tt.text), tt.pass)
		if err != nil {
			t.Fatalf("Encrypt(%s, %s) error: %s", tt.text, tt.pass, err)
		}

		dec, err := Decrypt(enc, tt.pass)
		if err != nil {
			t.Fatalf("Decrypt(%s, %s) error: %s", enc, tt.pass, err)
		}

		testz.Equal(t, tt.text, string(dec))
	}

	tests4 := []struct {
		text []byte
		pass string
	}{
		{[]byte("hello, this is a test!!!"), "whaterror"},
		{[]byte("??，世界"), "测试"},
		{[]byte("this is no pass"), ""},
	}
	for _, tt := range tests4 {
		enc, err := Encrypt(tt.text, []byte(tt.pass))
		if err != nil {
			t.Fatalf("Encrypt(%s, %s) error: %s", tt.text, tt.pass, err)
		}

		dec, err := Decrypt(enc, []byte(tt.pass))
		if err != nil {
			t.Fatalf("Decrypt(%s, %s) error: %s", enc, tt.pass, err)
		}

		testz.Equal(t, string(tt.text), string(dec))
	}
}

func TestEncryptStreamTo(t *testing.T) {
	text := "hello, this is a test!!!"
	if err := os.WriteFile("1.txt", []byte(text), 0644); err != nil {
		t.Fatalf("os.WriteFile(1.txt) error: %s", err)
	}

	f1, err := os.Open("1.txt")
	if err != nil {
		t.Fatalf("os.Open(1.txt) error: %s", err)
	}

	f2, err := os.Create("1.txt.enc")
	if err != nil {
		t.Fatalf("os.Create(1.txt.enc) error: %s", err)
	}

	if err := EncryptStreamTo(f2, f1, "123456"); err != nil {
		t.Fatalf("EncryptStreamTo error: %s", err)
	}

	_ = f1.Close()
	_ = f2.Close()

	f3, err := os.Open("1.txt.enc")
	if err != nil {
		t.Fatalf("os.Open(1.txt.enc) error: %s", err)
	}

	f4, err := os.Create("1.txt")
	if err != nil {
		t.Fatalf("os.Create(1.txt) error: %s", err)
	}

	if err := DecryptStreamTo(f4, f3, []byte("123456")); err != nil {
		t.Fatalf("DecryptStreamTo error: %s", err)
	}

	_ = f3.Close()
	_ = f4.Close()

	b, err := os.ReadFile("1.txt")
	if err != nil {
		t.Fatalf("os.ReadFile(1.txt) error: %s", err)
	}
	if string(b) != text {
		t.Fatalf("DecryptStreamTo error, expected: %s, actual: %s", text, string(b))
	}

	_ = os.Remove("1.txt")
	_ = os.Remove("1.txt.enc")
}

func TestGCMEncrypt(t *testing.T) {
	plainText := "hello world"
	pass := "im a pass"
	additionalData := "im a additional data"
	enc, err := GCMEncrypt(plainText, pass, additionalData)
	if err != nil {
		t.Fatalf("GCMEncrypt error: %s", err)
	}

	dec, err := GCMDecrypt(enc, pass, additionalData)
	if err != nil {
		t.Fatalf("GCMDecrypt error: %s", err)
	}

	if string(dec) != plainText {
		t.Fatalf("GCMDecrypt error, expected: %s, actual: %s", plainText, string(dec))
	}

	enc, err = GCMEncrypt(plainText, pass, "")
	if err != nil {
		t.Fatalf("GCMEncrypt error: %s", err)
	}

	dec, err = GCMDecrypt(enc, pass, []byte{})
	if err != nil {
		t.Fatalf("GCMDecrypt error: %s", err)
	}

	if string(dec) != plainText {
		t.Fatalf("GCMDecrypt error, expected: %s, actual: %s", plainText, string(dec))
	}
}

func TestHybridEncryptDecrypt(t *testing.T) {
	var str = "dadsadsajwbd1wibdiw1bidw1"

	pub, err := ParseRsaPublicKey(pubKey)
	testz.Nil(t, err)

	pri, err := ParseRsaPrivateKey(prvKey)
	testz.Nil(t, err)

	enc, err := HybridEncrypt(str, pub)
	testz.Nil(t, err)

	dec, err := HybridDecrypt(enc, pri)
	testz.Nil(t, err)

	testz.Equal(t, str, string(dec))
}

func TestHybridEncryptStreamTo(t *testing.T) {
	pub, err := ParseRsaPublicKey(pubKey)
	testz.Nil(t, err)

	pri, err := ParseRsaPrivateKey(prvKey)
	testz.Nil(t, err)

	text := "hello, this is a test!!!"
	if err := os.WriteFile("2.txt", []byte(text), 0644); err != nil {
		t.Fatalf("os.WriteFile(2.txt) error: %s", err)
	}

	f1, err := os.Open("2.txt")
	if err != nil {
		t.Fatalf("os.Open(2.txt) error: %s", err)
	}

	f2, err := os.Create("2.txt.enc")
	if err != nil {
		t.Fatalf("os.Create(2.txt.enc) error: %s", err)
	}

	if err := HybridEncryptStreamTo(f2, f1, pub); err != nil {
		t.Fatalf("HybridEncryptStreamTo error: %s", err)
	}

	_ = f1.Close()
	_ = f2.Close()

	f3, err := os.Open("2.txt.enc")
	if err != nil {
		t.Fatalf("os.Open(2.txt.enc) error: %s", err)
	}

	f4, err := os.Create("2.txt")
	if err != nil {
		t.Fatalf("os.Create(2.txt) error: %s", err)
	}

	if err := HybridDecryptStreamTo(f4, f3, pri); err != nil {
		t.Fatalf("DecryptStreamTo error: %s", err)
	}

	_ = f3.Close()
	_ = f4.Close()

	b, err := os.ReadFile("2.txt")
	if err != nil {
		t.Fatalf("os.ReadFile(1.txt) error: %s", err)
	}
	if string(b) != text {
		t.Fatalf("DecryptStreamTo error, expected: %s, actual: %s", text, string(b))
	}

	_ = os.Remove("2.txt")
	_ = os.Remove("2.txt.enc")
}
