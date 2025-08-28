package http

import (
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func ExtractBearerToken(c *fiber.Ctx) (string, error) {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("missing authorization header")
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == "" {
		return "", errors.New("invalid token format")
	}

	return token, nil
}
