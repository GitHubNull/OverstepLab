package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
)

// ---- AES-256-GCM ----

// Built-in AES key (for challenge purposes only - in real life would be stored securely)
var aesKey = sha256.Sum256([]byte("OverstepLabAES256ChallengeKey!@#"))

// AESEncrypt encrypts plaintext using AES-256-GCM and returns base64-encoded ciphertext.
func AESEncrypt(plaintext []byte) (string, error) {
	block, err := aes.NewCipher(aesKey[:])
	if err != nil {
		return "", fmt.Errorf("aes: failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("aes: failed to create GCM: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return "", fmt.Errorf("aes: failed to generate nonce: %w", err)
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// AESDecrypt decrypts base64-encoded ciphertext using AES-256-GCM.
func AESDecryptFromBase64(encoded string) ([]byte, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, fmt.Errorf("aes: invalid base64: %w", err)
	}

	block, err := aes.NewCipher(aesKey[:])
	if err != nil {
		return nil, fmt.Errorf("aes: failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("aes: failed to create GCM: %w", err)
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, fmt.Errorf("aes: ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("aes: decryption failed: %w", err)
	}

	return plaintext, nil
}

// GetAESKeyBase64 returns the AES key in base64 format so it can be shared with challengers.
func GetAESKeyBase64() string {
	return base64.StdEncoding.EncodeToString(aesKey[:])
}

// ---- RSA ----

// Built-in RSA key pair (for challenge purposes only)
// Generated once at package init
var (
	rsaPrivateKey *rsa.PrivateKey
	rsaPublicKey  *rsa.PublicKey
)

func init() {
	var err error
	rsaPrivateKey, err = rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic("crypto/rsa: failed to generate key: " + err.Error())
	}
	rsaPublicKey = &rsaPrivateKey.PublicKey
}

// RSAEncrypt encrypts plaintext with the built-in RSA public key (OAEP/SHA256).
func RSAEncrypt(plaintext []byte) (string, error) {
	ciphertext, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, rsaPublicKey, plaintext, nil)
	if err != nil {
		return "", fmt.Errorf("rsa: encryption failed: %w", err)
	}
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// RSADecrypt decrypts base64-encoded ciphertext with the built-in RSA private key.
func RSADecryptFromBase64(encoded string) ([]byte, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, fmt.Errorf("rsa: invalid base64: %w", err)
	}

	plaintext, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, rsaPrivateKey, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("rsa: decryption failed: %w", err)
	}
	return plaintext, nil
}

// GetRSAPublicKeyPEM returns the built-in RSA public key in PEM format.
func GetRSAPublicKeyPEM() string {
	pubBytes, err := x509.MarshalPKIXPublicKey(rsaPublicKey)
	if err != nil {
		return ""
	}
	block := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubBytes,
	}
	return string(pem.EncodeToMemory(block))
}
