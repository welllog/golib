package cryptz

import (
	"bytes"
	"crypto"
	"crypto/md5"
	_ "crypto/sha512"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/welllog/golib/testz"
)

type testKeyDeriver struct{}

func (t testKeyDeriver) ID() [8]byte {
	return [8]byte{'t', 'e', 's', 't', '_', 'k', 'd', 'r'}
}

func (t testKeyDeriver) Key(password, salt []byte, keyLen int) []byte {
	tmp := make([]byte, len(password)+len(salt))
	copy(tmp, password)
	copy(tmp[len(password):], salt)

	hash := md5.Sum(tmp)
	key := make([]byte, keyLen)
	offset := 0
	for {
		offset += copy(key[offset:], hash[:])
		if offset >= keyLen {
			break
		}
		hash = md5.Sum(hash[:])
	}
	return key
}

func (t testKeyDeriver) Header() []byte {
	id := t.ID()
	return id[:]
}

func (t testKeyDeriver) HeaderLen() int {
	return len(t.Header())
}

func (t testKeyDeriver) Restore(deriverHeader []byte) (KeyDeriver, error) {
	if len(deriverHeader) != t.HeaderLen() {
		return nil, fmt.Errorf("invalid testKeyDeriver header length")
	}

	id := t.ID()
	if !bytes.Equal(deriverHeader, id[:]) {
		return nil, fmt.Errorf("invalid testKeyDeriver header id")
	}

	return t, nil
}

func TestPasswordEncryptDecrypt(t *testing.T) {
	RegisterKeyDeriver(testKeyDeriver{})

	tests := []struct {
		plainText []byte
		pass      []byte
		addData   []byte
		deriver   KeyDeriver
	}{
		{
			[]byte("hello world"), []byte("im a pass"), []byte("im a additional data"),
			nil,
		},
		{
			[]byte("hello world"), []byte("im a pass"), nil,
			PBKDF2KeyDeriver{0, 0},
		},
		{
			[]byte("hello world"), []byte("im a pass"), []byte(""),
			PBKDF2KeyDeriver{1, crypto.SHA512},
		},
		{
			[]byte("👋,world"), []byte("test"), nil,
			PBKDF2KeyDeriver{10, crypto.MD5},
		},
		{
			[]byte("hello world"), []byte("im a pass"), nil,
			testKeyDeriver{},
		},
		{
			[]byte("hello world"), []byte("im a pass"), []byte("asdsa"),
			testKeyDeriver{},
		},
		{
			[]byte("👋,world"), []byte("test"), nil,
			testKeyDeriver{},
		},
	}

	for _, tt := range tests {
		enc, err := PasswordEncrypt(tt.plainText, tt.pass, tt.addData, tt.deriver)
		testz.Nil(t, err)

		dec, err := PasswordDecrypt(enc, tt.pass, tt.addData)
		testz.Nil(t, err)

		testz.Equal(t, string(tt.plainText), string(dec))
	}

	begin := time.Now()
	enc, err := PasswordEncrypt([]byte("hello world"), []byte("test"), []byte("test"),
		PBKDF2KeyDeriver{100_0000, crypto.SHA256})
	testz.Nil(t, err)
	fmt.Printf("PasswordEncrypt pbkf2-sha256 iter 100_0000 cost: %d ms \n", time.Since(begin).Milliseconds())

	begin = time.Now()
	dec, err := PasswordDecrypt(enc, []byte("test"), []byte("test"))
	testz.Nil(t, err)
	testz.Equal(t, "hello world", string(dec))
	fmt.Printf("PasswordDecrypt pbkf2-sha256 iter 100_0000 cost: %d ms \n", time.Since(begin).Milliseconds())
}

func TestPasswordEncryptDecryptLarge(t *testing.T) {
	pwd := []byte("test")
	ad := []byte("additional data")
	plainText := make([]byte, 512*1024)

	cipherText, err := PasswordEncrypt(plainText, pwd, ad, PBKDF2KeyDeriver{10000, crypto.SHA256})
	testz.Nil(t, err)

	decrypted, err := PasswordDecrypt(cipherText, pwd, ad)
	testz.Nil(t, err)
	testz.Equal(t, len(plainText), len(decrypted))

}

func TestPasswordEncryptDecryptStream(t *testing.T) {
	secret := "test"
	text := "hello, this is a test!!!"
	fileName := "PWD_STREAM.txt"
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

	if err := PasswordEncryptStream(f2, f1, secret, nil); err != nil {
		t.Fatalf("PasswordEncryptStream error: %s", err)
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

	if err := PasswordDecryptStream(f4, f3, secret); err != nil {
		t.Fatalf("PasswordDecryptStream error: %s", err)
	}

	_ = f3.Close()
	_ = f4.Close()

	b, err := os.ReadFile(fileName)
	if err != nil {
		t.Fatalf("os.ReadFile(%s) error: %s", fileName, err)
	}
	if string(b) != text {
		t.Fatalf("PasswordDecryptStream error, expected: %s, actual: %s", text, string(b))
	}

	_ = os.Remove(fileName)
	_ = os.Remove(fileName + ".enc")
}

func TestPasswordEncryptDecryptStream2(t *testing.T) {
	RegisterKeyDeriver(testKeyDeriver{})

	secret := "test"
	text := "hello, this is a test!!!"
	fileName := "PWD_STREAM2.txt"
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

	if err := PasswordEncryptStream(f2, f1, secret, testKeyDeriver{}); err != nil {
		t.Fatalf("PasswordEncryptStream error: %s", err)
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

	if err := PasswordDecryptStream(f4, f3, secret); err != nil {
		t.Fatalf("PasswordDecryptStream error: %s", err)
	}

	_ = f3.Close()
	_ = f4.Close()

	b, err := os.ReadFile(fileName)
	if err != nil {
		t.Fatalf("os.ReadFile(%s) error: %s", fileName, err)
	}
	if string(b) != text {
		t.Fatalf("PasswordDecryptStream error, expected: %s, actual: %s", text, string(b))
	}

	_ = os.Remove(fileName)
	_ = os.Remove(fileName + ".enc")
}
