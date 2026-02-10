package validations

import (
	"encoding/base64"
	"errors"
)

func decodeBase64(s string) ([]byte, error) {
	b, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return nil, errors.New("некорректный base64")
	}
	return b, nil
}

func validateLength(b []byte, expected int) error {
	if len(b) != expected {
		return errors.New("неверная длина поля")
	}
	return nil
}

func ValidateCEKFields(
	cekWrap, cekWrapIV, cekWrapSalt,
	encapsulatedKeySender, cekWrapSender, cekWrapSenderIV, cekWrapSenderSalt string,
) error {
	fields := []struct {
		value    string
		expected int
		name     string
	}{
		{cekWrap, 48, "CEKWrap"},
		{cekWrapIV, 12, "CEKWrapIV"},
		{cekWrapSalt, 32, "CEKWrapSalt"},
		{encapsulatedKeySender, 1088, "EncapsulatedKeySender"},
		{cekWrapSender, 48, "CEKWrapSender"},
		{cekWrapSenderIV, 12, "CEKWrapSenderIV"},
		{cekWrapSenderSalt, 32, "CEKWrapSenderSalt"},
	}

	for _, f := range fields {
		b, err := decodeBase64(f.value)
		if err != nil {
			return errors.New(f.name + ": некорректный base64")
		}
		if err := validateLength(b, f.expected); err != nil {
			return errors.New(f.name + ": " + err.Error())
		}
	}

	return nil
}
