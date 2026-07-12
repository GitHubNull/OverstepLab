package crypto

import (
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"

	"github.com/emmansun/gmsm/sm2"
	"github.com/emmansun/gmsm/sm4"
)

// ---- SM2 (国密非对称加密) ----

var (
	sm2PrivateKey *sm2.PrivateKey
	sm2PublicKey  *ecdsa.PublicKey
)

func init() {
	var err error
	sm2PrivateKey, err = sm2.GenerateKey(rand.Reader)
	if err != nil {
		panic("crypto/sm2: failed to generate key: " + err.Error())
	}
	sm2PublicKey = &sm2PrivateKey.PublicKey
}

// SM2Encrypt encrypts plaintext with the built-in SM2 public key.
func SM2Encrypt(plaintext []byte) (string, error) {
	ciphertext, err := sm2.Encrypt(rand.Reader, sm2PublicKey, plaintext, nil)
	if err != nil {
		return "", fmt.Errorf("sm2: encryption failed: %w", err)
	}
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// SM2Decrypt decrypts base64-encoded ciphertext with the built-in SM2 private key.
func SM2DecryptFromBase64(encoded string) ([]byte, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, fmt.Errorf("sm2: invalid base64: %w", err)
	}

	plaintext, err := sm2.Decrypt(sm2PrivateKey, ciphertext)
	if err != nil {
		return nil, fmt.Errorf("sm2: decryption failed: %w", err)
	}
	return plaintext, nil
}

// GetSM2PublicKeyHex returns the SM2 public key in hex format.
func GetSM2PublicKeyHex() string {
	pubBytes := sm2PublicKey.X.Bytes()
	pubBytes = append(pubBytes, sm2PublicKey.Y.Bytes()...)
	return hex.EncodeToString(pubBytes)
}

// ---- SM4 (国密对称加密) ----

var sm4Key = []byte("OverstepLabSM4!@") // 16 bytes for SM4-128

func init() {
	// SM4 key must be exactly 16 bytes
	if len(sm4Key) != 16 {
		panic("crypto/sm4: key must be 16 bytes")
	}
}

// SM4Encrypt encrypts plaintext with SM4-CBC using the built-in key.
func SM4Encrypt(plaintext []byte) (string, error) {
	block, err := sm4.NewCipher(sm4Key)
	if err != nil {
		return "", fmt.Errorf("sm4: failed to create cipher: %w", err)
	}

	// Generate random IV
	iv := make([]byte, sm4.BlockSize)
	if _, err := rand.Read(iv); err != nil {
		return "", fmt.Errorf("sm4: failed to generate IV: %w", err)
	}

	// PKCS7 padding
	plaintext = pkcs7Pad(plaintext, sm4.BlockSize)

	ciphertext := make([]byte, len(plaintext))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, plaintext)

	// Prepend IV to ciphertext for storage
	result := make([]byte, len(iv)+len(ciphertext))
	copy(result, iv)
	copy(result[len(iv):], ciphertext)

	return base64.StdEncoding.EncodeToString(result), nil
}

// SM4Decrypt decrypts base64-encoded ciphertext with SM4-CBC.
func SM4DecryptFromBase64(encoded string) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, fmt.Errorf("sm4: invalid base64: %w", err)
	}

	block, err := sm4.NewCipher(sm4Key)
	if err != nil {
		return nil, fmt.Errorf("sm4: failed to create cipher: %w", err)
	}

	if len(data) < sm4.BlockSize {
		return nil, fmt.Errorf("sm4: ciphertext too short")
	}

	iv := data[:sm4.BlockSize]
	ciphertext := data[sm4.BlockSize:]

	if len(ciphertext)%sm4.BlockSize != 0 {
		return nil, fmt.Errorf("sm4: ciphertext is not a multiple of block size")
	}

	plaintext := make([]byte, len(ciphertext))
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(plaintext, ciphertext)

	// Remove PKCS7 padding
	plaintext, err = pkcs7Unpad(plaintext, sm4.BlockSize)
	if err != nil {
		return nil, fmt.Errorf("sm4: invalid padding: %w", err)
	}

	return plaintext, nil
}

// GetSM4KeyBase64 returns the SM4 key in base64 format.
func GetSM4KeyBase64() string {
	return base64.StdEncoding.EncodeToString(sm4Key)
}

// ---- PKCS7 Padding helpers ----

func pkcs7Pad(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padText := make([]byte, padding)
	for i := range padText {
		padText[i] = byte(padding)
	}
	return append(data, padText...)
}

func pkcs7Unpad(data []byte, blockSize int) ([]byte, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("pkcs7: data is empty")
	}
	if len(data)%blockSize != 0 {
		return nil, fmt.Errorf("pkcs7: data not multiple of block size")
	}

	padding := int(data[len(data)-1])
	if padding > blockSize || padding == 0 {
		return nil, fmt.Errorf("pkcs7: invalid padding size %d", padding)
	}

	// Verify all padding bytes
	for i := len(data) - padding; i < len(data); i++ {
		if data[i] != byte(padding) {
			return nil, fmt.Errorf("pkcs7: invalid padding byte at position %d", i)
		}
	}

	return data[:len(data)-padding], nil
}
