package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/slipe-fun/skid-backend/internal/app/session"
)

type AuthMiddleware struct {
	sessionApp *session.SessionApp
}

func NewAuthMiddleware(sessionApp *session.SessionApp) *AuthMiddleware {
	return &AuthMiddleware{
		sessionApp: sessionApp,
	}
}

func (m *AuthMiddleware) Handle() fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Get("Authorization")
		if token == "" {
			return fiber.ErrUnauthorized
		}

		if len(token) > 7 && token[:7] == "Bearer " {
			token = token[7:]
		}

		session, err := m.sessionApp.GetSession(token)
		if err != nil {
			return fiber.ErrUnauthorized
		}

		c.Locals("session", session)

		return c.Next()
	}
}
