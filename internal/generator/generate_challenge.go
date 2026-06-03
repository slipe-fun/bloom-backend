package generator

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

func GenerateChallenge() (string, error) {
	bytes := make([]byte, 32)

	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to read random bytes: %w", err)
	}

	return base64.RawURLEncoding.EncodeToString(bytes), nil
}
