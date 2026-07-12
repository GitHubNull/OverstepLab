package crypto

import (
	"math/big"
)

// ---- Big integer helpers for Base58 ----

func bigIntZero() *big.Int {
	return new(big.Int).SetInt64(0)
}

func bigIntFromInt(v int64) *big.Int {
	return new(big.Int).SetInt64(v)
}

func bigIntFromBytes(data []byte) *big.Int {
	return new(big.Int).SetBytes(data)
}

func bigIntDiv58(n *big.Int) (*big.Int, int64) {
	base := big.NewInt(58)
	q := new(big.Int)
	r := new(big.Int)
	q.DivMod(n, base, r)
	return q, r.Int64()
}

func bigIntMul(a, b *big.Int) *big.Int {
	return new(big.Int).Mul(a, b)
}

func bigIntAdd(a, b *big.Int) *big.Int {
	return new(big.Int).Add(a, b)
}
