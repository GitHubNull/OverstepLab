package common

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateOrderNo() string {
	return fmt.Sprintf("ORD%s%d", time.Now().Format("20060102150405"), randInt(1000, 9999))
}

func GenerateAPIKey() string {
	b := make([]byte, 32)
	rand.Read(b)
	return "sk_" + hex.EncodeToString(b)
}

func randInt(min, max int) int {
	if min >= max {
		return min
	}
	return min + int(time.Now().UnixNano())%(max-min+1)
}
