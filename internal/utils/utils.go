package utils

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
)

// GenerateSecureCode генерирует криптографически безопасный код заданной длины
func GenerateSecureCode(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// GenerateSessionID генерирует криптографически безопасный ID сессии
func GenerateSessionID() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}
