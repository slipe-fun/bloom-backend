package validations

import (
	"encoding/base64"
	"fmt"
)

func ValidateCryptoLength(encoded string, expectedSize int) error {
	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return fmt.Errorf("invalid base64 encoding")
	}
	if len(data) != expectedSize {
		return fmt.Errorf("invalid length: expected %d, got %d", expectedSize, len(data))
	}
	return nil
}
