package crypto

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
)

// MD5Hash 计算MD5哈希
func MD5Hash(data string) string {
	h := md5.Sum([]byte(data))
	return hex.EncodeToString(h[:])
}

// SHA256Hash 计算SHA256哈希
func SHA256Hash(data string) string {
	h := sha256.Sum256([]byte(data))
	return hex.EncodeToString(h[:])
}

// HMACSHA256Sign 计算HMAC-SHA256签名
func HMACSHA256Sign(key string, payload string) string {
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(payload))
	return hex.EncodeToString(mac.Sum(nil))
}

// HMACSHA256Verify 验证HMAC-SHA256签名
func HMACSHA256Verify(key string, payload string, signature string) bool {
	expected := HMACSHA256Sign(key, payload)
	return hmac.Equal([]byte(expected), []byte(signature))
}

// ComputeHashSign 计算 value:md5(value|salt) 格式的哈希签名
// 这是旧式签名，仅用于E-09挑战教学
func ComputeHashSign(value string, salt string) string {
	hashStr := value + "|" + salt
	return fmt.Sprintf("%s:%s", value, MD5Hash(hashStr))
}

// VerifyHashSign 验证 value:md5格式的哈希签名
func VerifyHashSign(signedValue string, salt string) bool {
	idx := strings.LastIndex(signedValue, ":")
	if idx < 0 {
		return false
	}
	value := signedValue[:idx]
	expected := ComputeHashSign(value, salt)
	return signedValue == expected
}

// SortParams 将map参数排序后拼接为字符串，用于签名计算
func SortParams(params map[string]string) string {
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var parts []string
	for _, k := range keys {
		parts = append(parts, fmt.Sprintf("%s=%s", k, params[k]))
	}
	return strings.Join(parts, "&")
}

// ComputeBodyHashSign 计算请求体的哈希签名
// 用于E-09挑战的整体签名机制
func ComputeBodyHashSign(body string, salt string) string {
	return MD5Hash(body + "|" + salt)
}
