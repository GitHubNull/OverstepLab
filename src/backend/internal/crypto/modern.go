package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"
)

// AES-256-GCM 加密
func AESEncrypt(plaintext []byte, key []byte) (string, error) {
	// Key must be 32 bytes for AES-256
	if len(key) != 32 {
		// Hash to get 32-byte key
		h := sha256.Sum256(key)
		key = h[:]
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return Base64Encode(ciphertext), nil
}

// AES-256-GCM 解密
func AESDecrypt(encoded string, key []byte) ([]byte, error) {
	ciphertext, err := Base64Decode(encoded)
	if err != nil {
		return nil, err
	}
	if len(key) != 32 {
		h := sha256.Sum256(key)
		key = h[:]
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}

// RSAEncrypt 使用公钥加密
func RSAEncrypt(plaintext []byte, pubKey *rsa.PublicKey) (string, error) {
	encrypted, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, pubKey, plaintext, nil)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(encrypted), nil
}

// RSADecrypt 使用私钥解密
func RSADecrypt(encoded string, privKey *rsa.PrivateKey) ([]byte, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, err
	}
	return rsa.DecryptOAEP(sha256.New(), rand.Reader, privKey, ciphertext, nil)
}
