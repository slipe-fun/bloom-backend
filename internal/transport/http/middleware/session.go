package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/slipe-fun/skid-backend/internal/app/session"
	"github.com/slipe-fun/skid-backend/internal/app/user"
)

type AuthMiddleware struct {
	sessionApp *session.SessionApp
	userApp    *user.UserApp
}

func NewAuthMiddleware(sessionApp *session.SessionApp, userApp *user.UserApp) *AuthMiddleware {
	return &AuthMiddleware{
		sessionApp: sessionApp,
		userApp:    userApp,
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

		user, err := m.userApp.GetUserByID(session.UserID)
		if err != nil {
			return fiber.ErrUnauthorized
		}

		c.Locals("session", session)
		c.Locals("session_user", user)
		c.Locals("token", token)

		return c.Next()
	}
}
