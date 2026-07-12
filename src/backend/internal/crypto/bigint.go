package crypto

import (
	"math/big"
)

// BigIntToBytes 将 big.Int 转换为字节数组
func BigIntToBytes(x *big.Int) []byte {
	return x.Bytes()
}

// BytesToBigInt 将字节数组转换为 big.Int
func BytesToBigInt(data []byte) *big.Int {
	return new(big.Int).SetBytes(data)
}

// BigIntToString 将 big.Int 转换为十进制字符串
func BigIntToString(x *big.Int) string {
	return x.String()
}

// StringToBigInt 将十进制字符串转换为 big.Int
func StringToBigInt(s string) (*big.Int, bool) {
	x, ok := new(big.Int).SetString(s, 10)
	return x, ok
}

// BigIntToHex 将 big.Int 转换为十六进制字符串
func BigIntToHex(x *big.Int) string {
	return x.Text(16)
}

// HexToBigInt 将十六进制字符串转换为 big.Int
func HexToBigInt(s string) (*big.Int, bool) {
	x, ok := new(big.Int).SetString(s, 16)
	return x, ok
}

// ModExp 计算 base^exp mod m
func ModExp(base, exp, m *big.Int) *big.Int {
	return new(big.Int).Exp(base, exp, m)
}

// ModMul 计算 (a * b) mod m
func ModMul(a, b, m *big.Int) *big.Int {
	return new(big.Int).Mul(a, b).Mod(new(big.Int).Mul(a, b), m)
}
