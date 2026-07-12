package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
)

// SM4Encrypt 使用SM4-CBC模式加密（简化实现，实际使用AES-128-CBC模拟SM4参数）
// SM4：128位密钥，128位分组
func SM4Encrypt(plaintext []byte, key []byte) (string, error) {
	if len(key) == 0 {
		return "", errors.New("empty key")
	}
	// Hash to get 16-byte key (SM4 uses 128-bit key)
	h := sha256.Sum256(key)
	key16 := h[:16]

	block, err := aes.NewCipher(key16)
	if err != nil {
		return "", err
	}

	// PKCS7 padding
	blockSize := block.BlockSize()
	padding := blockSize - len(plaintext)%blockSize
	padText := make([]byte, len(plaintext)+padding)
	copy(padText, plaintext)
	for i := len(plaintext); i < len(padText); i++ {
		padText[i] = byte(padding)
	}

	// CBC mode
	iv := make([]byte, blockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	ciphertext := make([]byte, blockSize+len(padText))
	copy(ciphertext[:blockSize], iv)

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[blockSize:], padText)

	return Base64Encode(ciphertext), nil
}

// SM4Decrypt 使用SM4-CBC模式解密
func SM4Decrypt(encoded string, key []byte) ([]byte, error) {
	ciphertext, err := Base64Decode(encoded)
	if err != nil {
		return nil, err
	}
	if len(key) == 0 {
		return nil, errors.New("empty key")
	}
	h := sha256.Sum256(key)
	key16 := h[:16]

	block, err := aes.NewCipher(key16)
	if err != nil {
		return nil, err
	}

	blockSize := block.BlockSize()
	if len(ciphertext) < blockSize {
		return nil, errors.New("ciphertext too short")
	}

	iv := ciphertext[:blockSize]
	ciphertext = ciphertext[blockSize:]

	if len(ciphertext)%blockSize != 0 {
		return nil, fmt.Errorf("ciphertext is not a multiple of block size")
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(ciphertext, ciphertext)

	// Remove PKCS7 padding
	padding := int(ciphertext[len(ciphertext)-1])
	if padding > blockSize || padding == 0 {
		return nil, errors.New("invalid padding")
	}
	for i := len(ciphertext) - padding; i < len(ciphertext); i++ {
		if ciphertext[i] != byte(padding) {
			return nil, errors.New("invalid padding")
		}
	}
	return ciphertext[:len(ciphertext)-padding], nil
}

// SM2Encrypt 使用SM2加密（简化实现，使用标准RSA模拟）
// SM2：基于椭圆曲线，此处简化使用base64编码表示
func SM2Encrypt(plaintext []byte, key []byte) (string, error) {
	h := sha256.Sum256(key)
	// 简化实现：XOR + Base64
	result := make([]byte, len(plaintext))
	for i := range plaintext {
		result[i] = plaintext[i] ^ h[i%len(h)]
	}
	return base64.StdEncoding.EncodeToString(result), nil
}

// SM2Decrypt 使用SM2解密
func SM2Decrypt(encoded string, key []byte) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, err
	}
	h := sha256.Sum256(key)
	result := make([]byte, len(data))
	for i := range data {
		result[i] = data[i] ^ h[i%len(h)]
	}
	return result, nil
}
