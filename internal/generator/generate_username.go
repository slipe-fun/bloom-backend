package generator

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"

	"github.com/slipe-fun/skid-backend/internal/config"
)

func GenerateUsername() (string, error) {
	adjIdx, err := rand.Int(rand.Reader, big.NewInt(int64(len(config.Adjectives))))
	if err != nil {
		return "", fmt.Errorf("failed to select adjective: %w", err)
	}
	adj := strings.ToLower(strings.TrimSpace(config.Adjectives[adjIdx.Int64()]))

	nounIdx, err := rand.Int(rand.Reader, big.NewInt(int64(len(config.Nouns))))
	if err != nil {
		return "", fmt.Errorf("failed to select noun: %w", err)
	}
	noun := strings.ToLower(strings.TrimSpace(config.Nouns[nounIdx.Int64()]))

	suffixBytes := make([]byte, 4)
	if _, err := rand.Read(suffixBytes); err != nil {
		return "", fmt.Errorf("failed to generate random suffix: %w", err)
	}
	suffix := hex.EncodeToString(suffixBytes)

	username := fmt.Sprintf("%s-%s-%s", adj, noun, suffix)

	return username, nil
}
