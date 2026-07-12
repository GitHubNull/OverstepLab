package crypto

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// ---- Hash functions ----

func MD5Hash(data []byte) string {
	h := md5.Sum(data)
	return hex.EncodeToString(h[:])
}

func SHA1Hash(data []byte) string {
	h := sha1.Sum(data)
	return hex.EncodeToString(h[:])
}

func SHA256Hash(data []byte) string {
	h := sha256.Sum256(data)
	return hex.EncodeToString(h[:])
}

// ---- HMAC-SHA256 Signature ----

var hmacKey = []byte("OverstepLabHMACSecretKey!@#2024")

// HMACSHA256Sign creates an HMAC-SHA256 signature for the given data.
func HMACSHA256Sign(data []byte) string {
	mac := hmac.New(sha256.New, hmacKey)
	mac.Write(data)
	return hex.EncodeToString(mac.Sum(nil))
}

// HMACSHA256Verify verifies the HMAC-SHA256 signature.
func HMACSHA256Verify(data []byte, signature string) bool {
	expected := HMACSHA256Sign(data)
	return hmac.Equal([]byte(expected), []byte(signature))
}

// GetHMACKey returns the HMAC key as base64 (for challenge hints).
func GetHMACKeyBase64() string {
	return base64.StdEncoding.EncodeToString(hmacKey)
}

// ---- Signed Parameter Encoding ----

// EncodeSignedParam encodes a value with an HMAC signature: base64(data) + "." + hex(hmac(data))
func EncodeSignedParam(value string) string {
	data := []byte(value)
	encoded := base64.StdEncoding.EncodeToString(data)
	sig := HMACSHA256Sign(data)
	return encoded + "." + sig
}

// DecodeSignedParam decodes a signed parameter and verifies the signature.
func DecodeSignedParam(encoded string) ([]byte, error) {
	parts := strings.SplitN(encoded, ".", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("signed: invalid format, expected base64data.signature")
	}

	data, err := base64.StdEncoding.DecodeString(parts[0])
	if err != nil {
		return nil, fmt.Errorf("signed: invalid base64 data: %w", err)
	}

	signature := parts[1]
	if !HMACSHA256Verify(data, signature) {
		return nil, fmt.Errorf("signed: HMAC signature verification failed")
	}

	return data, nil
}

// ---- Simple Hash Signature (MD5/SHA256) ----
// Real-world style: pack all params into a sorted string, compute MD5+salt, put in header X-Hash-Sign
// Used for E-09 challenge: simple hash integrity check bypass

var hashSalt = []byte("OverstepLabHashSalt2024")

// SerializeParams packs params into a sorted key=value&key2=value2 string (same as frontend)
func SerializeParams(data map[string]interface{}) string {
	var entries []string
	for key, value := range data {
		var valStr string
		switch v := value.(type) {
		case string:
			valStr = v
		case float64:
			valStr = strconv.FormatFloat(v, 'f', -1, 64)
		case int, int64, uint, uint64:
			valStr = fmt.Sprintf("%v", v)
		case bool:
			valStr = strconv.FormatBool(v)
		default:
			valStr = fmt.Sprintf("%v", v)
		}
		entries = append(entries, key+"="+valStr)
	}
	sort.Strings(entries)
	return strings.Join(entries, "&")
}

// ComputeHashSign computes MD5(payload + salt) — same algorithm as frontend
func ComputeHashSign(payload string) string {
	h := md5.Sum([]byte(payload + string(hashSalt)))
	return hex.EncodeToString(h[:])
}

// ComputeHmacSign computes HMAC-SHA256(payload) — same algorithm as frontend
func ComputeHmacSign(payload string) string {
	mac := hmac.New(sha256.New, hmacKey)
	mac.Write([]byte(payload))
	return hex.EncodeToString(mac.Sum(nil))
}

// VerifyHashSign verifies the X-Hash-Sign header against the given params
func VerifyHashSign(params map[string]interface{}, sign string) bool {
	expected := ComputeHashSign(SerializeParams(params))
	return strings.EqualFold(sign, expected)
}

// GetHashSalt returns the hash salt for challenge participants.
func GetHashSalt() []byte {
	return hashSalt
}

// SignID signs a numeric ID for use in signed challenges.
func SignID(id uint) string {
	return EncodeSignedParam(fmt.Sprintf("%d", id))
}
