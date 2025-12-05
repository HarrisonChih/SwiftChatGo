package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

// 小写
func Sha256Encode(data string) string {
	h := sha256.New()
	h.Write([]byte(data))

	return hex.EncodeToString(h.Sum(nil))
}

// 大写
func SHA256Encode(data string) string {
	return strings.ToUpper(Sha256Encode(data))
}

// 加密
func MakePassword(plainpwd, salt string) string {
	return Sha256Encode(plainpwd + salt)
}

// 解密
func ValidPassword(plainpwd, salt, password string) bool {
	return Sha256Encode(plainpwd+salt) == password
}
