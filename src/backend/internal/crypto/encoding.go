package crypto

import (
	"encoding/base32"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"
)

// ---- Base64 ----

func Base64Encode(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

func Base64Decode(s string) ([]byte, error) {
	s = strings.TrimSpace(s)
	return base64.StdEncoding.DecodeString(s)
}

// ---- Base32 (RFC 4648 standard, no padding) ----

func Base32Encode(data []byte) string {
	return base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(data)
}

func Base32Decode(s string) ([]byte, error) {
	s = strings.TrimSpace(s)
	s = strings.ToUpper(s)
	return base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(s)
}

// ---- Base32 Hex ----

func Base32HexEncode(data []byte) string {
	return base32.HexEncoding.WithPadding(base32.NoPadding).EncodeToString(data)
}

func Base32HexDecode(s string) ([]byte, error) {
	s = strings.TrimSpace(s)
	s = strings.ToUpper(s)
	return base32.HexEncoding.WithPadding(base32.NoPadding).DecodeString(s)
}

// ---- Hex ----

func HexEncode(data []byte) string {
	return hex.EncodeToString(data)
}

func HexDecode(s string) ([]byte, error) {
	s = strings.TrimSpace(s)
	return hex.DecodeString(s)
}

// ---- Custom Base64 (swapped charset: A<->Z, a<->z, 0<->9, +<->/) ----

const customBase64Charset = "ZYXWVUTSRQPONMLKJIHGFEDCBAzyxwvutsrqponmlkjihgfedcba9876543210/+"

var customBase64Encoding = base64.NewEncoding(customBase64Charset)

func CustomBase64Encode(data []byte) string {
	return customBase64Encoding.EncodeToString(data)
}

func CustomBase64Decode(s string) ([]byte, error) {
	s = strings.TrimSpace(s)
	return customBase64Encoding.DecodeString(s)
}

// ---- Custom Base32 (reversed RFC 4648 charset) ----

const customBase32Charset = "ZYXWVUTSRQPONMLKJIHGFEDCBA765432"

var customBase32Encoding = base32.NewEncoding(customBase32Charset).WithPadding(base32.NoPadding)

func CustomBase32Encode(data []byte) string {
	return customBase32Encoding.EncodeToString(data)
}

func CustomBase32Decode(s string) ([]byte, error) {
	s = strings.TrimSpace(s)
	s = strings.ToUpper(s)
	return customBase32Encoding.DecodeString(s)
}

// ---- Base58 (Bitcoin-style) ----

const base58Alphabet = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

func Base58Encode(data []byte) string {
	// Count leading zeros
	zeros := 0
	for zeros < len(data) && data[zeros] == 0 {
		zeros++
	}

	// Convert to big integer in base 58
	num := bigIntFromBytes(data)
	result := make([]byte, 0, len(data)*136/100+1)
	for num.Sign() > 0 {
		var rem int64
		num, rem = bigIntDiv58(num)
		result = append(result, base58Alphabet[rem])
	}

	// Add leading zeros as '1's
	for i := 0; i < zeros; i++ {
		result = append(result, '1')
	}

	// Reverse
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}

	return string(result)
}

func Base58Decode(s string) ([]byte, error) {
	s = strings.TrimSpace(s)
	num := bigIntZero()
	base := bigIntFromInt(58)

	for _, c := range []byte(s) {
		idx := strings.IndexByte(base58Alphabet, c)
		if idx < 0 {
			return nil, fmt.Errorf("invalid base58 character: %c", c)
		}
		num = bigIntMul(num, base)
		num = bigIntAdd(num, bigIntFromInt(int64(idx)))
	}

	// Count leading '1's
	zeros := 0
	for zeros < len(s) && s[zeros] == '1' {
		zeros++
	}

	result := num.Bytes()
	decoded := make([]byte, zeros+len(result))
	copy(decoded[zeros:], result)
	return decoded, nil
}

// ---- Base85 (Ascii85/Z85 variant) ----

const base85Alphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz!#$%&()*+-;<=>?@^_`{|}~"

func Base85Encode(data []byte) string {
	if len(data) == 0 {
		return ""
	}

	var result strings.Builder
	for i := 0; i < len(data); i += 4 {
		chunk := data[i:]
		if len(chunk) > 4 {
			chunk = chunk[:4]
		}
		padding := 4 - len(chunk)
		var val uint32
		for j, b := range chunk {
			val |= uint32(b) << (24 - 8*j)
		}

		tmp := make([]byte, 5)
		for j := 4; j >= 0; j-- {
			tmp[j] = base85Alphabet[val%85]
			val /= 85
		}
		result.Write(tmp[:5-padding])
	}
	return result.String()
}

func Base85Decode(s string) ([]byte, error) {
	s = strings.TrimSpace(s)
	if len(s) == 0 {
		return nil, nil
	}

	var result []byte
	for i := 0; i < len(s); i += 5 {
		chunk := s[i:]
		if len(chunk) > 5 {
			chunk = chunk[:5]
		}
		padding := 5 - len(chunk)
		var val uint32
		for _, c := range []byte(chunk) {
			idx := strings.IndexByte(base85Alphabet, c)
			if idx < 0 {
				return nil, fmt.Errorf("invalid base85 character: %c", c)
			}
			val = val*85 + uint32(idx)
		}
		for j := 0; j < padding; j++ {
			val = val*85 + 84
		}

		tmp := make([]byte, 4)
		for j := 3; j >= 0; j-- {
			tmp[j] = byte(val & 0xFF)
			val >>= 8
		}
		result = append(result, tmp[:4-padding]...)
	}
	return result, nil
}

// DecodeByType dispatches to the correct decoder based on encoding type string.
func DecodeByType(encoded string, encType string) ([]byte, error) {
	switch encType {
	case "base64":
		return Base64Decode(encoded)
	case "base32":
		return Base32Decode(encoded)
	case "base32hex":
		return Base32HexDecode(encoded)
	case "base58":
		return Base58Decode(encoded)
	case "base85":
		return Base85Decode(encoded)
	case "custom_base64":
		return CustomBase64Decode(encoded)
	case "custom_base32":
		return CustomBase32Decode(encoded)
	case "hex":
		return HexDecode(encoded)
	case "caesar":
		raw, err := CaesarDecode(encoded)
		if err != nil {
			return nil, err
		}
		return []byte(raw), nil
	case "aes":
		raw, err := AESDecryptFromBase64(encoded)
		if err != nil {
			return nil, err
		}
		return raw, nil
	case "rsa":
		raw, err := RSADecryptFromBase64(encoded)
		if err != nil {
			return nil, err
		}
		return raw, nil
	case "sm4":
		raw, err := SM4DecryptFromBase64(encoded)
		if err != nil {
			return nil, err
		}
		return raw, nil
	case "signed":
		raw, err := DecodeSignedParam(encoded)
		if err != nil {
			return nil, err
		}
		return raw, nil
	case "multi":
		// Base32 -> Base64 nested decoding
		base32Decoded, err := Base32Decode(encoded)
		if err != nil {
			return nil, fmt.Errorf("multi: base32 decode failed: %w", err)
		}
		base64Decoded, err := Base64Decode(string(base32Decoded))
		if err != nil {
			return nil, fmt.Errorf("multi: base64 decode failed: %w", err)
		}
		return base64Decoded, nil
	default:
		return nil, fmt.Errorf("unsupported encoding type: %s", encType)
	}
}

// DecodeParam decodes a single encoded parameter string using the specified encoding type,
// and returns the decoded value as a string (most IDs are numeric strings).
func DecodeParam(encoded string, encType string) (string, error) {
	decoded, err := DecodeByType(encoded, encType)
	if err != nil {
		return "", err
	}
	return string(decoded), nil
}
