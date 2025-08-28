package http

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/slipe-fun/skid-backend/internal/app"
)

type AuthHandler struct {
	authApp *app.AuthApp
}

func NewAuthHandler(authApp *app.AuthApp) *AuthHandler {
	return &AuthHandler{authApp: authApp}
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}

	token, user, err := h.authApp.Login(req.Username, req.Password, time.Hour*8760)

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"token": token,
		"user": fiber.Map{
			"id":       user.ID,
			"username": user.Username,
			"date":     user.Date,
		},
	})
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}

	token, user, err := h.authApp.Register(req.Username, req.Password, time.Hour*8760)

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"token": token,
		"user": fiber.Map{
			"id":       user.ID,
			"username": user.Username,
			"date":     user.Date,
		},
	})
}
