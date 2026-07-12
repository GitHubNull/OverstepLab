package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/oversteplab/oversteplab/internal/crypto"
)

// EncodingMiddleware 透明编码中间件
// 当请求包含 X-Encoding-Type 头时，自动解码 query 参数和 body 字段
func EncodingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		encType := c.GetHeader("X-Encoding-Type")
		if encType == "" {
			c.Set("encoding_type", "none")
			c.Next()
			return
		}

		encType = strings.ToLower(strings.TrimSpace(encType))
		c.Set("encoding_type", encType)

		// 解码 query 参数
		if len(c.Request.URL.RawQuery) > 0 {
			query := c.Request.URL.Query()
			for key, values := range query {
				decoded := make([]string, len(values))
				for i, v := range values {
					decoded[i] = DecodeQueryValue(encType, v)
				}
				query[key] = decoded
			}
			c.Request.URL.RawQuery = query.Encode()
		}

		// 解码 POST/PUT body（JSON）
		if c.Request.Method == "POST" || c.Request.Method == "PUT" || c.Request.Method == "DELETE" {
			if c.Request.Body != nil && c.Request.ContentLength > 0 {
				bodyBytes, err := io.ReadAll(c.Request.Body)
				if err == nil && len(bodyBytes) > 0 {
					decoded := DecodeBodyFields(encType, bodyBytes)
					c.Request.Body = io.NopCloser(bytes.NewBuffer(decoded))
					c.Request.ContentLength = int64(len(decoded))
				}
			}
		}

		c.Next()
	}
}

// DecodeQueryValue 根据编码类型解码单个query值
func DecodeQueryValue(encType string, value string) string {
	if value == "" {
		return value
	}
	switch encType {
	case "base64":
		decoded, err := crypto.Base64Decode(value)
		if err != nil {
			return value
		}
		return string(decoded)
	case "base32":
		decoded, err := crypto.Base32Decode(value)
		if err != nil {
			return value
		}
		return string(decoded)
	case "caesar":
		return crypto.CaesarDecodeShift(value, 3)
	case "custom-base64":
		decoded, err := crypto.CustomBase64Decode(value)
		if err != nil {
			return value
		}
		return string(decoded)
	case "multi":
		// Multi-layer: Base64 -> Base32
		decoded, err := crypto.Base32Decode(value)
		if err != nil {
			return value
		}
		decoded2, err := crypto.Base64Decode(string(decoded))
		if err != nil {
			return value
		}
		return string(decoded2)
	case "aes":
		// Key: oversteplab-aes-secret-key-32b (hashed to 32 bytes internally)
		decoded, err := crypto.AESDecrypt(value, []byte("oversteplab-aes-secret-key-32b"))
		if err != nil {
			return value
		}
		return string(decoded)
	case "sm4":
		// Key: oversteplab-sm4-secret-key (hashed to 16 bytes internally)
		decoded, err := crypto.SM4Decrypt(value, []byte("oversteplab-sm4-secret-key"))
		if err != nil {
			return value
		}
		return string(decoded)
	case "hash-sign":
		// Format: value:md5(value|salt) — extract value part
		parts := strings.SplitN(value, ":", 2)
		return parts[0]
	default:
		return value
	}
}

// DecodeBodyFields 解码JSON body中所有字符串字段值
func DecodeBodyFields(encType string, body []byte) []byte {
	if encType == "none" || encType == "" || encType == "hmac" || encType == "hash-sign" {
		// hmac/hash-sign: body stays plaintext, signature is in header
		return body
	}
	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return body
	}
	for key, val := range data {
		if s, ok := val.(string); ok && s != "" {
			data[key] = DecodeQueryValue(encType, s)
		}
	}
	result, err := json.Marshal(data)
	if err != nil {
		return body
	}
	return result
}

// DecodeUintParam 从query中提取并解码uint类型参数
// 用于handler中替换 c.Query("vpsId") 等
func DecodeUintParam(c *gin.Context, key string) uint {
	encType := c.GetString("encoding_type")
	if encType == "" {
		encType = "none"
	}
	val := c.Query(key)
	if val == "" {
		// 尝试从body中获取
		if c.Request.Body != nil {
			bodyBytes, _ := io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			var data map[string]interface{}
			if json.Unmarshal(bodyBytes, &data) == nil {
				// 重新设置body供后续使用
				c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
				if v, ok := data[key]; ok {
					switch vv := v.(type) {
					case string:
						val = vv
					case float64:
						return uint(vv)
					}
				}
			}
		}
	}
	if val == "" {
		return 0
	}
	decoded := DecodeQueryValue(encType, val)
	n, err := strconv.ParseUint(decoded, 10, 64)
	if err != nil {
		return 0
	}
	return uint(n)
}

// DecodeUintBodyField 从已解析的map中解码uint字段
func DecodeUintBodyField(rawData map[string]interface{}, key string, encType string) uint {
	if v, ok := rawData[key]; ok {
		switch vv := v.(type) {
		case float64:
			return uint(vv)
		case string:
			decoded := DecodeQueryValue(encType, vv)
			n, err := strconv.ParseUint(decoded, 10, 64)
			if err == nil {
				return uint(n)
			}
		}
	}
	return 0
}

// DecodeIntParam 从query中解码int参数
func DecodeIntParam(c *gin.Context, key string) int {
	encType := c.GetString("encoding_type")
	if encType == "" {
		encType = "none"
	}
	val := c.Query(key)
	if val == "" {
		return 0
	}
	decoded := DecodeQueryValue(encType, val)
	n, err := strconv.Atoi(decoded)
	if err != nil {
		return 0
	}
	return n
}

// GetEncodingType 从请求上下文中获取编码类型
func GetEncodingType(c *gin.Context) string {
	encType, exists := c.Get("encoding_type")
	if !exists {
		return "none"
	}
	return encType.(string)
}

// verifyHashSign 验证E-09哈希签名（中间件级别）
func verifyHashSign(c *gin.Context, salt string) bool {
	signHeader := c.GetHeader("X-Hash-Sign")
	if signHeader == "" {
		// 无签名头，放行（业务由E-09判定）
		return true
	}

	// 读取body用于签名验证
	var bodyStr string
	if c.Request.Body != nil {
		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err == nil {
			bodyStr = string(bodyBytes)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}
	}
	// 也包含query参数
	if c.Request.URL.RawQuery != "" {
		bodyStr += "?" + c.Request.URL.RawQuery
	}

	expected := crypto.ComputeBodyHashSign(bodyStr, salt)
	return signHeader == expected
}

// verifyHmacSign 验证E-07 HMAC签名（中间件级别）
func verifyHmacSign(c *gin.Context, key string) bool {
	signHeader := c.GetHeader("X-HMAC-Sign")
	if signHeader == "" {
		return true
	}

	var bodyStr string
	if c.Request.Body != nil {
		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err == nil {
			bodyStr = string(bodyBytes)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}
	}
	if c.Request.URL.RawQuery != "" {
		bodyStr += "?" + c.Request.URL.RawQuery
	}

	return crypto.HMACSHA256Verify(key, bodyStr, signHeader)
}

// HandleEncodingSignatures 处理E-07/E-09签名验证
// 作为中间件在需要签名验证的路由上使用
func HandleEncodingSignatures() gin.HandlerFunc {
	return func(c *gin.Context) {
		encType := GetEncodingType(c)

		switch encType {
		case "hmac":
			key := "oversteplab-hmac-secret-key"
			if !verifyHmacSign(c, key) {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
					"code":    403,
					"message": "HMAC签名验证失败",
				})
				return
			}
		case "hash-sign":
			salt := "oversteplab-hash-salt"
			if !verifyHashSign(c, salt) {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
					"code":    403,
					"message": "哈希签名验证失败",
				})
				return
			}
		}
		c.Next()
	}
}
