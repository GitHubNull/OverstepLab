package crypto

import (
	"encoding/base32"
	"encoding/base64"
	"encoding/hex"
	"math/big"
	"strings"
)

var (
	// CustomBase64Table 自定义Base64编码表（打乱标准表顺序）
	CustomBase64Table = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz+/"
	customBase64      = base64.NewEncoding(CustomBase64Table).WithPadding(base64.NoPadding)

	// Base58 alphabet (Bitcoin style, no 0/O/I/l)
	base58Alphabet = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"
)

func Base64Encode(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

func Base64Decode(s string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(s)
}

func Base64URLEncode(data []byte) string {
	return base64.URLEncoding.EncodeToString(data)
}

func Base64URLDecode(s string) ([]byte, error) {
	return base64.URLEncoding.DecodeString(s)
}

func Base32Encode(data []byte) string {
	return base32.StdEncoding.EncodeToString(data)
}

func Base32Decode(s string) ([]byte, error) {
	s = strings.ToUpper(s)
	s = strings.TrimRight(s, "=")
	// Pad to multiple of 8
	if pad := len(s) % 8; pad != 0 {
		s += strings.Repeat("=", 8-pad)
	}
	return base32.StdEncoding.DecodeString(s)
}

func Base58Encode(data []byte) string {
	x := new(big.Int).SetBytes(data)
	zero := big.NewInt(0)
	base := big.NewInt(58)
	mod := new(big.Int)
	var result []byte

	for x.Cmp(zero) > 0 {
		x.DivMod(x, base, mod)
		result = append([]byte{base58Alphabet[mod.Int64()]}, result...)
	}

	// Add leading zeros
	for _, b := range data {
		if b == 0 {
			result = append([]byte{base58Alphabet[0]}, result...)
		} else {
			break
		}
	}
	return string(result)
}

func Base58Decode(s string) ([]byte, error) {
	x := big.NewInt(0)
	base := big.NewInt(58)
	for _, c := range []byte(s) {
		idx := strings.IndexByte(base58Alphabet, c)
		if idx < 0 {
			return nil, ErrInvalidEncoding
		}
		x.Mul(x, base)
		x.Add(x, big.NewInt(int64(idx)))
	}

	// Count leading zeros
	leadingZeros := 0
	for _, c := range []byte(s) {
		if c == base58Alphabet[0] {
			leadingZeros++
		} else {
			break
		}
	}

	result := x.Bytes()
	result = append(make([]byte, leadingZeros), result...)
	return result, nil
}

func HexEncode(data []byte) string {
	return hex.EncodeToString(data)
}

func HexDecode(s string) ([]byte, error) {
	return hex.DecodeString(s)
}

func CustomBase64Encode(data []byte) string {
	return customBase64.EncodeToString(data)
}

func CustomBase64Decode(s string) ([]byte, error) {
	return customBase64.DecodeString(s)
}

// ErrInvalidEncoding is returned when encoded input is invalid
var ErrInvalidEncoding = errorString("crypto: invalid encoding")

type errorString string

func (e errorString) Error() string { return string(e) }
