package crypto

import (
	"fmt"
	"strconv"
	"strings"
)

// ---- Caesar Cipher (凯撒密码) ----

// CaesarEncode applies a Caesar shift to a numeric string.
// For non-numeric strings, it shifts ASCII letters.
func CaesarEncode(plain string, shift int) string {
	var result strings.Builder
	for _, c := range plain {
		result.WriteRune(caesarShift(c, shift))
	}
	return result.String()
}

// CaesarDecode attempts to decode a Caesar-encoded string.
// Since the shift is unknown, it returns the result using shift=3 (most common).
// For challenge purposes, the encoded_handler will use a fixed shift=3.
func CaesarDecode(encoded string) (string, error) {
	// Try to detect if this is a pure numeric string with Caesar shift
	// Default shift = 3 for the challenge
	return CaesarDecodeWithShift(encoded, 3)
}

// CaesarDecodeWithShift decodes with a known shift.
func CaesarDecodeWithShift(encoded string, shift int) (string, error) {
	// Reverse shift: letters use mod 26, digits use mod 10
	reverseShiftLetters := (26 - (shift % 26)) % 26
	reverseShiftDigits := (10 - (shift % 10)) % 10

	var result strings.Builder
	for _, c := range encoded {
		if c >= '0' && c <= '9' {
			// Numeric Caesar: shift digits
			digit := int(c - '0')
			shifted := (digit + reverseShiftDigits) % 10
			result.WriteRune(rune('0' + shifted))
		} else if c >= 'A' && c <= 'Z' {
			shifted := ((int(c-'A') + reverseShiftLetters) % 26)
			result.WriteRune(rune('A' + shifted))
		} else if c >= 'a' && c <= 'z' {
			shifted := ((int(c-'a') + reverseShiftLetters) % 26)
			result.WriteRune(rune('a' + shifted))
		} else {
			result.WriteRune(c)
		}
	}
	return result.String(), nil
}

// CaesarEncodeWithShift applies Caesar shift to a string.
func CaesarEncodeWithShift(plain string, shift int) string {
	var result strings.Builder
	for _, c := range plain {
		result.WriteRune(caesarShift(c, shift))
	}
	return result.String()
}

func caesarShift(c rune, shift int) rune {
	s := shift % 26
	if s < 0 {
		s += 26
	}
	if c >= '0' && c <= '9' {
		return rune('0' + (int(c-'0')+s)%10)
	}
	if c >= 'A' && c <= 'Z' {
		return rune('A' + (int(c-'A')+s)%26)
	}
	if c >= 'a' && c <= 'z' {
		return rune('a' + (int(c-'a')+s)%26)
	}
	return c
}

// ---- Vigenère Cipher (维吉尼亚密码) ----

func VigenereEncode(plain, key string) string {
	key = strings.ToUpper(key)
	if len(key) == 0 {
		return plain
	}

	var result strings.Builder
	keyIdx := 0
	for _, c := range plain {
		if c >= 'A' && c <= 'Z' {
			shift := int(key[keyIdx%len(key)] - 'A')
			shifted := (int(c-'A') + shift) % 26
			result.WriteRune(rune('A' + shifted))
			keyIdx++
		} else if c >= 'a' && c <= 'z' {
			shift := int(key[keyIdx%len(key)] - 'A')
			shifted := (int(c-'a') + shift) % 26
			result.WriteRune(rune('a' + shifted))
			keyIdx++
		} else if c >= '0' && c <= '9' {
			shift := int(key[keyIdx%len(key)] - 'A')
			shifted := (int(c-'0') + shift) % 10
			result.WriteRune(rune('0' + shifted))
			keyIdx++
		} else {
			result.WriteRune(c)
		}
	}
	return result.String()
}

func VigenereDecode(encoded, key string) string {
	key = strings.ToUpper(key)
	if len(key) == 0 {
		return encoded
	}

	var result strings.Builder
	keyIdx := 0
	for _, c := range encoded {
		if c >= 'A' && c <= 'Z' {
			shift := int(key[keyIdx%len(key)] - 'A')
			shifted := (int(c-'A') - shift + 26) % 26
			result.WriteRune(rune('A' + shifted))
			keyIdx++
		} else if c >= 'a' && c <= 'z' {
			shift := int(key[keyIdx%len(key)] - 'A')
			shifted := (int(c-'a') - shift + 26) % 26
			result.WriteRune(rune('a' + shifted))
			keyIdx++
		} else if c >= '0' && c <= '9' {
			shift := int(key[keyIdx%len(key)] - 'A')
			shifted := (int(c-'0') - shift + 10) % 10
			result.WriteRune(rune('0' + shifted))
			keyIdx++
		} else {
			result.WriteRune(c)
		}
	}
	return result.String()
}

// ---- Rail Fence Cipher (栅栏密码) ----

func RailFenceEncode(plain string, rails int) string {
	if rails <= 1 || len(plain) <= rails {
		return plain
	}

	fence := make([][]rune, rails)
	for i := range fence {
		fence[i] = make([]rune, 0, len(plain)/rails+1)
	}

	rail := 0
	dir := 1
	for _, c := range plain {
		fence[rail] = append(fence[rail], c)
		rail += dir
		if rail == 0 || rail == rails-1 {
			dir = -dir
		}
	}

	var result strings.Builder
	for _, row := range fence {
		result.WriteString(string(row))
	}
	return result.String()
}

func RailFenceDecode(encoded string, rails int) (string, error) {
	if rails <= 1 || len(encoded) <= rails {
		return encoded, nil
	}

	// Build the rail pattern
	pattern := make([]int, len(encoded))
	rail := 0
	dir := 1
	for i := range pattern {
		pattern[i] = rail
		rail += dir
		if rail == 0 || rail == rails-1 {
			dir = -dir
		}
	}

	// Count chars per rail
	counts := make([]int, rails)
	for _, r := range pattern {
		counts[r]++
	}

	// Fill rails
	rails_content := make([][]rune, rails)
	pos := 0
	for i := 0; i < rails; i++ {
		rails_content[i] = []rune(encoded[pos : pos+counts[i]])
		pos += counts[i]
	}

	// Read off in pattern order
	result := make([]rune, len(encoded))
	indices := make([]int, rails)
	for i, r := range pattern {
		result[i] = rails_content[r][indices[r]]
		indices[r]++
	}

	return string(result), nil
}

// EncodeID encodes a numeric ID with a specific encoding scheme.
func EncodeID(id uint, encType string) string {
	idStr := strconv.FormatUint(uint64(id), 10)
	switch encType {
	case "base64":
		return Base64Encode([]byte(idStr))
	case "base32":
		return Base32Encode([]byte(idStr))
	case "base58":
		return Base58Encode([]byte(idStr))
	case "base85":
		return Base85Encode([]byte(idStr))
	case "custom_base64":
		return CustomBase64Encode([]byte(idStr))
	case "custom_base32":
		return CustomBase32Encode([]byte(idStr))
	case "caesar":
		return CaesarEncodeWithShift(idStr, 3)
	case "multi":
		// Base64 -> Base32 nested encoding
		return Base32Encode([]byte(Base64Encode([]byte(idStr))))
	default:
		return idStr
	}
}

// ---- Utility for challenges: detect if a string looks like a particular encoding ----

func DetectEncoding(s string) string {
	// Try common encodings
	if _, err := Base64Decode(s); err == nil && looksLikeBase64(s) {
		return "base64"
	}
	if _, err := Base32Decode(s); err == nil && looksLikeBase32(s) {
		return "base32"
	}
	return "unknown"
}

func looksLikeBase64(s string) bool {
	return len(s)%4 == 0 && len(s) >= 4
}

func looksLikeBase32(s string) bool {
	s = strings.ToUpper(s)
	for _, c := range s {
		if !((c >= 'A' && c <= 'Z') || (c >= '2' && c <= '7') || c == '=') {
			return false
		}
	}
	return len(s) >= 4
}

// Ensure fmt is used (imported above)
var _ = fmt.Sprintf
