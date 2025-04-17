package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

// HashPassword Returns the hash of password with our credentials
func HashPassword(password string) (string, error) {
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}

	combined := append(salt, []byte(password)...)

	hash := sha256.New()
	_, err = hash.Write(combined)
	if err != nil {
		return "", err
	}

	hashedPassword := hex.EncodeToString(hash.Sum(nil))
	saltString := hex.EncodeToString(salt)
	result := fmt.Sprintf("%s:%s", saltString, hashedPassword)
	return result, nil
}

func CheckPassword(password, storedValue string) bool {
	parts := strings.Split(storedValue, ":")
	if len(parts) != 2 {
		return false
	}

	saltString := parts[0]
	hashedPassword := parts[1]

	salt, err := hex.DecodeString(saltString)
	if err != nil {
		return false
	}

	combined := append(salt, []byte(password)...)

	hash := sha256.New()
	_, err = hash.Write(combined)
	if err != nil {
		return false
	}

	return hex.EncodeToString(hash.Sum(nil)) == hashedPassword
}
