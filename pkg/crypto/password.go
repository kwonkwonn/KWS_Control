package crypto

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"strings"
)

func HashPasswordWithSalt(password string, salt []byte) []byte {
	hash := sha256.New()
	temp := strings.ToUpper(hex.EncodeToString(salt))
	hash.Write([]byte(password))
	hash.Write([]byte(temp))
	return hash.Sum(nil)
}

func GenerateRandomSalt(length int) ([]byte, error) {
	salt := make([]byte, length)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}
	return salt, nil
}

func GenerateRandomPassword(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(bytes)[:length], nil
}
