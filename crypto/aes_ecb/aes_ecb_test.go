package aes_ecb

import (
	"fmt"
	"testing"
)

func TestAESEncryptDecrypt(t *testing.T) {
	fmt.Println("AES ECB加密解密测试........")
	key := []byte("1234567890123456")
	blickSize := 16
	tool := NewAesTool(key, blickSize)
	encryptData, _ := tool.Encrypt([]byte("abcdef"))
	fmt.Println(encryptData)
	decryptData, _ := tool.Decrypt(encryptData)
	fmt.Println(string(decryptData))
}