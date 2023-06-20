package cryptz

import (
	"fmt"
	"os"
	"testing"

	"github.com/welllog/golib/testz"
)

func TestEncryptToBase64String(t *testing.T) {
	tests1 := []struct {
		text string
		pass string
	}{
		{"hello, this is a test!!!", "whaterror"},
		{"👋，世界", "测试"},
	}
	for _, tt := range tests1 {
		enc, err := EncryptToBase64String(tt.text, tt.pass)
		if err != nil {
			t.Fatalf("EncryptToString(%s, %s) error: %s", tt.text, tt.pass, err)
		}

		dec, err := DecryptBase64ToString(enc, tt.pass)
		if err != nil {
			t.Fatalf("DecryptFromString(%s, %s) error: %s", enc, tt.pass, err)
		}

		testz.Equal(t, tt.text, dec,
			fmt.Sprintf("EncryptToString(%s, %s) != DecryptFromString(%s, %s)", tt.text, tt.pass, enc, tt.pass),
		)
	}

	tests2 := []struct {
		text []byte
		pass []byte
	}{
		{[]byte("hello, this is a test!!!"), []byte("whaterror")},
		{[]byte("??，世界"), []byte("测试")},
	}
	for _, tt := range tests2 {
		enc, err := EncryptToBase64String(tt.text, tt.pass)
		if err != nil {
			t.Fatalf("Encrypt(%s, %s) error: %s", tt.text, tt.pass, err)
		}

		dec, err := DecryptBase64ToString(enc, tt.pass)
		if err != nil {
			t.Fatalf("Decrypt(%s, %s) error: %s", enc, tt.pass, err)
		}

		testz.Equal(t, string(tt.text), dec,
			fmt.Sprintf("EncryptToString(%s, %s) != DecryptFromString(%s, %s)", string(tt.text), string(tt.pass), enc, string(tt.pass)),
		)
	}

	tests3 := []struct {
		text string
		pass []byte
	}{
		{"hello, this is a test!!!", []byte("whaterror")},
		{"??，世界", []byte("测试")},
	}
	for _, tt := range tests3 {
		enc, err := EncryptToBase64String([]byte(tt.text), tt.pass)
		if err != nil {
			t.Fatalf("Encrypt(%s, %s) error: %s", tt.text, tt.pass, err)
		}

		dec, err := DecryptBase64ToString(enc, tt.pass)
		if err != nil {
			t.Fatalf("Decrypt(%s, %s) error: %s", enc, tt.pass, err)
		}

		testz.Equal(t, tt.text, dec,
			fmt.Sprintf("EncryptToString(%s, %s) != DecryptFromString(%s, %s)", tt.text, string(tt.pass), enc, string(tt.pass)),
		)
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
		enc, err := EncryptToBase64String(tt.text, []byte(tt.pass))
		if err != nil {
			t.Fatalf("Encrypt(%s, %s) error: %s", tt.text, tt.pass, err)
		}

		dec, err := DecryptBase64ToString(enc, []byte(tt.pass))
		if err != nil {
			t.Fatalf("Decrypt(%s, %s) error: %s", enc, tt.pass, err)
		}

		testz.Equal(t, string(tt.text), dec,
			fmt.Sprintf("EncryptToString(%s, %s) != DecryptFromString(%s, %s)", tt.text, tt.pass, enc, tt.pass),
		)
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
