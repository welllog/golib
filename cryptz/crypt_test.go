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
	secret := "test"
	text := "hello, this is a test!!!"
	fileName := "STREAM.txt"
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

	if err := EncryptStreamTo(f2, f1, secret); err != nil {
		t.Fatalf("EncryptStreamTo error: %s", err)
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

	if err := DecryptStreamTo(f4, f3, secret); err != nil {
		t.Fatalf("DecryptStreamTo error: %s", err)
	}

	_ = f3.Close()
	_ = f4.Close()

	b, err := os.ReadFile(fileName)
	if err != nil {
		t.Fatalf("os.ReadFile(%s) error: %s", fileName, err)
	}
	if string(b) != text {
		t.Fatalf("DecryptStreamTo error, expected: %s, actual: %s", text, string(b))
	}

	_ = os.Remove(fileName)
	_ = os.Remove(fileName + ".enc")
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
