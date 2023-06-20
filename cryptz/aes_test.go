package cryptz

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"testing"
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
	enc, err := AESGCMEncrypt(text, text, key, nonce, nil)
	if err != nil {
		t.Fatalf("AESGCMEncrypt error: %s", err)
	}

	dec, err := AESGCMDecrypt(enc, enc, key, nonce, nil)
	if err != nil {
		t.Fatalf("AESGCMDecrypt error: %s", err)
	}

	if string(dec) != str {
		t.Fatalf("AESGCMDecrypt(%s) != %s", string(dec), str)
	}
}

func TestAESGCMDecrypt(t *testing.T) {
	key := []byte("0123456789abcdef")

	str := "hello world"
	text := []byte(str)
	nonce := []byte("abc")
	addition := []byte("addition")
	enc, err := AESGCMEncrypt(nil, text, key, nonce, addition)
	if err != nil {
		t.Fatalf("AESGCMEncrypt error: %s", err)
	}

	dec, err := AESGCMDecrypt(nil, enc, key, nonce, nil)
	if err == nil {
		t.Fatal("AESGCMDecrypt additionalData missing should fail")
	}

	dec, err = AESGCMDecrypt(nil, enc, key, nonce, addition)
	if err != nil {
		t.Fatalf("AESGCMDecrypt error: %s", err)
	}

	if string(dec) != str {
		t.Fatalf("AESGCMDecrypt(%s) != %s", string(dec), str)
	}
}

func TestExamp(t *testing.T) {
	// 需要加密的原始数据
	plaintext := []byte("Hello, World!")

	// 32字节的AES密钥，可以是16、24或32字节
	key := []byte("01234567890123456789012345678901")

	// 附加数据
	additionalData := []byte("Additional Data")

	// 加密数据
	ciphertext, err := encrypt(plaintext, key, additionalData)
	if err != nil {
		fmt.Println("加密失败:", err)
		return
	}

	fmt.Printf("加密后的数据: %x\n", ciphertext)

	// 解密数据
	decryptedText, err := decrypt(ciphertext, key, additionalData)
	if err != nil {
		fmt.Println("解密失败:", err)
		return
	}

	fmt.Printf("解密后的数据: %s\n", decryptedText)
}

func encrypt(plaintext []byte, key []byte, additionalData []byte) ([]byte, error) {
	// 创建AES加密器
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// 创建GCM实例，使用block作为底层加密器
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// 生成随机的nonce（初始化向量）
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	// 加密数据，同时传入additionalData
	ciphertext := aesGCM.Seal(nil, nonce, plaintext, additionalData)

	// 将nonce与密文拼接在一起返回
	ciphertext = append(nonce, ciphertext...)

	return ciphertext, nil
}

// 解密函数
func decrypt(ciphertext []byte, key []byte, additionalData []byte) ([]byte, error) {
	// 创建AES解密器
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// 创建GCM实例，使用block作为底层解密器
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// 验证密文长度是否足够包含nonce
	if len(ciphertext) < aesGCM.NonceSize() {
		return nil, errors.New("invalid ciphertext")
	}

	// 提取nonce
	nonce := ciphertext[:aesGCM.NonceSize()]

	// 解密数据，同时传入additionalData
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext[aesGCM.NonceSize():], additionalData)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}
