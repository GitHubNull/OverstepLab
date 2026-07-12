package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"net/url"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	cryptopkg "github.com/oversteplab/oversteplab/internal/crypto"
)

// EncodingMiddleware transparently decodes request data when X-Encoding-Type header is present.
// It runs after AuthMiddleware on all authenticated routes, so handlers always receive
// already-decoded data without any awareness of the encoding layer.
//
// Supported encoding types: base64, base32, caesar, custom_base64, custom_base32,
// hex, base58, base85, multi, aes, rsa, sm2, sm4, signed
func EncodingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip crypto utility endpoints and legacy encoded routes
		if strings.HasPrefix(c.Request.URL.Path, "/api/v1/crypto") ||
			strings.HasPrefix(c.Request.URL.Path, "/api/v1/encoded") {
			c.Next()
			return
		}

		encType := strings.TrimSpace(c.GetHeader("X-Encoding-Type"))
		if encType == "" || encType == "none" {
			c.Next()
			return
		}

		// Special handling for hash type (E-09): verify X-Hash-Sign header
		if encType == "hash" {
			if !verifyHashSign(c) {
				c.AbortWithStatusJSON(400, gin.H{"code": 1, "message": "hash signature verification failed"})
				return
			}
			c.Next()
			return
		}

		// Special handling for signed type (E-07): verify X-HMAC-Sign header
		if encType == "signed" {
			if !verifyHmacSign(c) {
				c.AbortWithStatusJSON(400, gin.H{"code": 1, "message": "HMAC signature verification failed"})
				return
			}
			c.Next()
			return
		}

		// Decode query parameters for GET/HEAD requests
		if c.Request.Method == "GET" || c.Request.Method == "HEAD" {
			decodeQueryParams(c, encType)
			c.Next()
			return
		}

		// Decode JSON body string fields for POST/PUT/DELETE/PATCH requests
		contentType := c.GetHeader("Content-Type")
		if !strings.Contains(contentType, "application/json") {
			c.Next()
			return
		}

		body, err := io.ReadAll(c.Request.Body)
		if err != nil || len(body) == 0 {
			c.Next()
			return
		}

		decodedBody, err := decodeJSONBody(body, encType)
		if err != nil {
			// Body is not valid JSON or decode failed, pass through unchanged
			c.Request.Body = io.NopCloser(bytes.NewReader(body))
			c.Next()
			return
		}

		c.Request.Body = io.NopCloser(bytes.NewReader(decodedBody))
		c.Request.ContentLength = int64(len(decodedBody))
		c.Next()
	}
}

// decodeQueryParams decodes all query parameter values using the given encoding type.
func decodeQueryParams(c *gin.Context, encType string) {
	rawQuery := c.Request.URL.RawQuery
	if rawQuery == "" {
		return
	}

	// Parse and rebuild query string with decoded values
	var pairs []string
	// Use split on "&" for simple parsing that preserves order
	for _, pair := range strings.Split(rawQuery, "&") {
		eqIdx := strings.Index(pair, "=")
		if eqIdx < 0 {
			pairs = append(pairs, pair)
			continue
		}
		key := pair[:eqIdx]
		value := pair[eqIdx+1:]

		// URL-decode the value first (query params may be percent-encoded)
		if decodedValue, err := url.QueryUnescape(value); err == nil {
			value = decodedValue
		}

		if decoded, err := cryptopkg.DecodeParam(value, encType); err == nil {
			pairs = append(pairs, key+"="+url.QueryEscape(decoded))
		} else {
			pairs = append(pairs, pair)
		}
	}

	c.Request.URL.RawQuery = strings.Join(pairs, "&")
}

// decodeJSONBody decodes all string values in a JSON body recursively.
func decodeJSONBody(body []byte, encType string) ([]byte, error) {
	var data interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	decoded := decodeValue(data, encType)
	return json.Marshal(decoded)
}

// decodeValue recursively decodes string and number values in arbitrary JSON structures.
func decodeValue(v interface{}, encType string) interface{} {
	switch val := v.(type) {
	case string:
		if decoded, err := cryptopkg.DecodeParam(val, encType); err == nil {
			// Try to convert back to number if the decoded value is numeric
			if num, err := strconv.ParseFloat(decoded, 64); err == nil {
				// Check if it's an integer
				if float64(int64(num)) == num {
					return int64(num)
				}
				return num
			}
			return decoded
		}
		return val
	case float64:
		// Numbers in JSON are parsed as float64. Try decoding as string representation.
		if decoded, err := cryptopkg.DecodeParam(strconv.FormatFloat(val, 'f', -1, 64), encType); err == nil {
			// Try to convert back to number if the decoded value is numeric
			if num, err := strconv.ParseFloat(decoded, 64); err == nil {
				// Check if it's an integer
				if float64(int64(num)) == num {
					return int64(num)
				}
				return num
			}
			return decoded
		}
		return v
	case map[string]interface{}:
		result := make(map[string]interface{}, len(val))
		for k, child := range val {
			result[k] = decodeValue(child, encType)
		}
		return result
	case []interface{}:
		result := make([]interface{}, len(val))
		for i, child := range val {
			result[i] = decodeValue(child, encType)
		}
		return result
	default:
		return v
	}
}

// verifyHashSign verifies the X-Hash-Sign header for E-09 challenge.
// For GET: params sorted as key=value&key2=value2 → MD5+salt
// For POST/PUT: JSON body raw string → MD5+salt
// If no params/body and no X-Hash-Sign header, allow through (no signature needed for empty requests)
func verifyHashSign(c *gin.Context) bool {
	sign := strings.TrimSpace(c.GetHeader("X-Hash-Sign"))

	var payload string
	if c.Request.Method == "GET" || c.Request.Method == "HEAD" {
		// GET: sort query params as key=value&key2=value2
		params := make(map[string]interface{})
		for key, values := range c.Request.URL.Query() {
			if len(values) > 0 {
				params[key] = values[0]
			}
		}
		payload = cryptopkg.SerializeParams(params)
		// No params and no sign header → allow
		if payload == "" && sign == "" {
			return true
		}
	} else {
		// POST/PUT: read JSON body as raw string
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			return false
		}
		// Restore body for downstream handlers
		c.Request.Body = io.NopCloser(bytes.NewReader(body))
		payload = string(body)
		// No body and no sign header → allow
		if payload == "" && sign == "" {
			return true
		}
	}

	if sign == "" {
		return false
	}

	expected := cryptopkg.ComputeHashSign(payload)
	return strings.EqualFold(sign, expected)
}

// verifyHmacSign verifies the X-HMAC-Sign header for E-07 challenge.
// Same logic as verifyHashSign but uses HMAC-SHA256 with a secret key.
func verifyHmacSign(c *gin.Context) bool {
	sign := strings.TrimSpace(c.GetHeader("X-HMAC-Sign"))

	var payload string
	if c.Request.Method == "GET" || c.Request.Method == "HEAD" {
		params := make(map[string]interface{})
		for key, values := range c.Request.URL.Query() {
			if len(values) > 0 {
				params[key] = values[0]
			}
		}
		payload = cryptopkg.SerializeParams(params)
		if payload == "" && sign == "" {
			return true
		}
	} else {
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			return false
		}
		c.Request.Body = io.NopCloser(bytes.NewReader(body))
		payload = string(body)
		if payload == "" && sign == "" {
			return true
		}
	}

	if sign == "" {
		return false
	}

	expected := cryptopkg.ComputeHmacSign(payload)
	return strings.EqualFold(sign, expected)
}
