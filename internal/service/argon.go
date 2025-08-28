package service

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

var (
	argonTime    uint32 = 1
	argonMemory  uint32 = 16 * 1024
	argonThreads uint8  = 4
	argonKeyLen  uint32 = 32
)

func generateSalt(n int) ([]byte, error) {
	salt := make([]byte, n)
	_, err := rand.Read(salt)
	return salt, err
}

func HashPassword(password string) (string, error) {
	salt, err := generateSalt(16)
	if err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, argonTime, argonMemory, argonThreads, argonKeyLen)

	encodedSalt := base64.RawStdEncoding.EncodeToString(salt)
	encodedHash := base64.RawStdEncoding.EncodeToString(hash)

	return fmt.Sprintf("%s$%s", encodedSalt, encodedHash), nil
}

func VerifyPassword(password, encoded string) (bool, error) {
	parts := strings.Split(encoded, "$")
	if len(parts) != 2 {
		return false, fmt.Errorf("invalid hash format")
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[0])
	if err != nil {
		return false, err
	}

	hash, err := base64.RawStdEncoding.DecodeString(parts[1])
	if err != nil {
		return false, err
	}

	testHash := argon2.IDKey([]byte(password), salt, argonTime, argonMemory, argonThreads, uint32(len(hash)))

	if len(testHash) != len(hash) {
		return false, nil
	}

	var diff byte
	for i := 0; i < len(hash); i++ {
		diff |= hash[i] ^ testHash[i]
	}

	return diff == 0, nil
}
