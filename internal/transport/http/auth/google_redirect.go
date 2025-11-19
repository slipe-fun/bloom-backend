package auth

import (
	"fmt"
	"net/url"

	"github.com/gofiber/fiber/v2"
	"github.com/slipe-fun/skid-backend/internal/config"
)

func (h *AuthHandler) GoogleRedirect(c *fiber.Ctx) error {
	if c.Query("state") != "random-state" {
		return c.Status(400).SendString("invalid state")
	}

	code := c.Query("code")
	if code == "" {
		return c.Status(400).SendString("no code")
	}

	cfg := config.LoadConfig("configs/config.yaml")

	redirectURI := fmt.Sprintf("%s://oauth2redirect/google?code=%s", cfg.GoogleAuth.BundleID, url.QueryEscape(code))

	return c.Redirect(redirectURI, fiber.StatusFound)
}
