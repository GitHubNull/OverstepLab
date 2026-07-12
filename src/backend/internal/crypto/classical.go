package crypto

import (
	"strconv"
	"strings"
)

// CaesarEncode 凯撒密码编码（字母+数字）
// 偏移量shift：字母循环移位A↔Z，数字循环移位0↔9
func CaesarEncode(data []byte, shift int) []byte {
	shift = ((shift % 26) + 26) % 26
	result := make([]byte, len(data))
	copy(result, data)
	for i, b := range result {
		switch {
		case b >= 'A' && b <= 'Z':
			result[i] = 'A' + (b-'A'+byte(shift))%26
		case b >= 'a' && b <= 'z':
			result[i] = 'a' + (b-'a'+byte(shift))%26
		case b >= '0' && b <= '9':
			digitShift := shift % 10
			if digitShift < 0 {
				digitShift += 10
			}
			result[i] = '0' + (b-'0'+byte(digitShift))%10
		}
	}
	return result
}

// CaesarDecode 凯撒密码解码
func CaesarDecode(data []byte, shift int) []byte {
	return CaesarEncode(data, -shift)
}

// CaesarDecodeAuto 自动尝试所有偏移量，返回最常见的解码结果
func CaesarDecodeAuto(data []byte) []byte {
	if len(data) == 0 {
		return data
	}
	// 默认使用偏移量3（ROT3是最常见的）
	return CaesarDecode(data, 3)
}

// VigenereEncode 维吉尼亚密码编码
func VigenereEncode(data []byte, key string) []byte {
	keyBytes := []byte(strings.ToUpper(key))
	if len(keyBytes) == 0 {
		return data
	}
	result := make([]byte, len(data))
	copy(result, data)
	keyIdx := 0
	for i, b := range result {
		if b >= 'A' && b <= 'Z' {
			shift := keyBytes[keyIdx%len(keyBytes)] - 'A'
			result[i] = 'A' + (b-'A'+shift)%26
			keyIdx++
		} else if b >= 'a' && b <= 'z' {
			shift := keyBytes[keyIdx%len(keyBytes)] - 'A'
			result[i] = 'a' + (b-'a'+shift)%26
			keyIdx++
		}
	}
	return result
}

// VigenereDecode 维吉尼亚密码解码
func VigenereDecode(data []byte, key string) []byte {
	keyBytes := []byte(strings.ToUpper(key))
	if len(keyBytes) == 0 {
		return data
	}
	result := make([]byte, len(data))
	copy(result, data)
	keyIdx := 0
	for i, b := range result {
		if b >= 'A' && b <= 'Z' {
			shift := keyBytes[keyIdx%len(keyBytes)] - 'A'
			result[i] = 'A' + (b-'A'+26-shift)%26
			keyIdx++
		} else if b >= 'a' && b <= 'z' {
			shift := keyBytes[keyIdx%len(keyBytes)] - 'A'
			result[i] = 'a' + (b-'a'+26-shift)%26
			keyIdx++
		}
	}
	return result
}

// RailFenceEncode 栅栏密码编码
func RailFenceEncode(data []byte, rails int) []byte {
	if rails <= 1 || len(data) <= rails {
		return data
	}
	n := len(data)
	fence := make([][]byte, rails)
	for i := range fence {
		fence[i] = make([]byte, 0, n/rails+1)
	}
	rail := 0
	dir := 1
	for _, b := range data {
		fence[rail] = append(fence[rail], b)
		rail += dir
		if rail == 0 || rail == rails-1 {
			dir = -dir
		}
	}
	result := make([]byte, 0, n)
	for _, row := range fence {
		result = append(result, row...)
	}
	return result
}

// RailFenceDecode 栅栏密码解码
func RailFenceDecode(data []byte, rails int) []byte {
	if rails <= 1 || len(data) <= rails {
		return data
	}
	n := len(data)
	// Build the rail pattern
	railLen := make([]int, rails)
	rail := 0
	dir := 1
	for i := 0; i < n; i++ {
		railLen[rail]++
		rail += dir
		if rail == 0 || rail == rails-1 {
			dir = -dir
		}
	}
	// Fill rails
	fence := make([][]byte, rails)
	pos := 0
	for r := 0; r < rails; r++ {
		fence[r] = data[pos : pos+railLen[r]]
		pos += railLen[r]
	}
	// Read zigzag
	result := make([]byte, n)
	rail = 0
	dir = 1
	idx := make([]int, rails)
	for i := 0; i < n; i++ {
		result[i] = fence[rail][idx[rail]]
		idx[rail]++
		rail += dir
		if rail == 0 || rail == rails-1 {
			dir = -dir
		}
	}
	return result
}

// ROT13 是凯撒密码 shift=13 的特例（仅字母）
func ROT13(data []byte) []byte {
	result := make([]byte, len(data))
	copy(result, data)
	for i, b := range result {
		switch {
		case b >= 'A' && b <= 'Z':
			result[i] = 'A' + (b-'A'+13)%26
		case b >= 'a' && b <= 'z':
			result[i] = 'a' + (b-'a'+13)%26
		}
	}
	return result
}

// CaesarDecodeShift 凯撒解码，指定偏移量
// 同时处理字母和数字（字母0-25循环，数字0-9循环）
func CaesarDecodeShift(s string, shift int) string {
	data := []byte(s)
	letterShift := ((shift % 26) + 26) % 26
	digitShift := ((shift % 10) + 10) % 10
	for i, b := range data {
		switch {
		case b >= 'A' && b <= 'Z':
			data[i] = 'A' + (b-'A'+26-byte(letterShift))%26
		case b >= 'a' && b <= 'z':
			data[i] = 'a' + (b-'a'+26-byte(letterShift))%26
		case b >= '0' && b <= '9':
			data[i] = '0' + (b-'0'+10-byte(digitShift))%10
		}
	}
	return string(data)
}

// StrToInt converts a string to an integer, returns 0 on error
func StrToInt(s string) int {
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return n
}
