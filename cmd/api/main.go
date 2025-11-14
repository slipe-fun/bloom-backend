package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"

	AuthApp "github.com/slipe-fun/skid-backend/internal/app/auth"
	ChatApp "github.com/slipe-fun/skid-backend/internal/app/chat"
	MessageApp "github.com/slipe-fun/skid-backend/internal/app/message"
	UserApp "github.com/slipe-fun/skid-backend/internal/app/user"
	VerificationApp "github.com/slipe-fun/skid-backend/internal/app/verification"
	"github.com/slipe-fun/skid-backend/internal/config"
	"github.com/slipe-fun/skid-backend/internal/repository"
	ChatRepo "github.com/slipe-fun/skid-backend/internal/repository/chat"
	MessageRepo "github.com/slipe-fun/skid-backend/internal/repository/message"
	UserRepo "github.com/slipe-fun/skid-backend/internal/repository/user"
	VerificationRepo "github.com/slipe-fun/skid-backend/internal/repository/verification"
	"github.com/slipe-fun/skid-backend/internal/service"
	"github.com/slipe-fun/skid-backend/internal/transport/http/auth"
	"github.com/slipe-fun/skid-backend/internal/transport/http/chat"
	"github.com/slipe-fun/skid-backend/internal/transport/http/message"
	"github.com/slipe-fun/skid-backend/internal/transport/http/middleware"
	"github.com/slipe-fun/skid-backend/internal/transport/http/user"
	"github.com/slipe-fun/skid-backend/internal/transport/ws/handler"
	"github.com/slipe-fun/skid-backend/internal/transport/ws/types"
)

func main() {
	cfg := config.LoadConfig("configs/config.yaml")

	db := repository.InitDB(cfg)
	defer db.Close()

	verificationRepo := VerificationRepo.NewVerificationRepo(db)
	userRepo := UserRepo.NewUserRepo(db, verificationRepo)
	chatRepo := ChatRepo.NewChatRepo(db, userRepo)
	messageRepo := MessageRepo.NewMessageRepo(db)

	jwtSvc := service.NewJWTService(cfg.JWT.Secret)
	tokenSvc := service.NewTokenService(jwtSvc)

	verificationApp := VerificationApp.NewAuthApp(verificationRepo)
	authApp := AuthApp.NewAuthApp(userRepo, verificationRepo, verificationApp, jwtSvc)
	userApp := UserApp.NewUserApp(userRepo, jwtSvc, tokenSvc)
	chatApp := ChatApp.NewChatApp(chatRepo, tokenSvc)
	messageApp := MessageApp.NewMessageApp(messageRepo, chatApp, tokenSvc)

	authHandler := auth.NewAuthHandler(authApp)
	userHandler := user.NewUserHandler(userApp)
	chatHandler := chat.NewChatHandler(chatApp, userApp, messageApp)
	messageHandler := message.NewMessageHandler(chatApp, userApp, messageApp)

	fiberApp := fiber.New()

	if cfg.RateLimit.Enabled {
		rateLimiter := middleware.NewAdaptiveRateLimiter(cfg.RateLimitWindow())

		rateLimiter.SetLimit("auth", cfg.RateLimit.AuthRequestsPerMinute)
		rateLimiter.SetLimit("api", cfg.RateLimit.GeneralRequestsPerMinute)

		fiberApp.Use(rateLimiter.RateLimit())

		go func() {
			ticker := time.NewTicker(5 * time.Minute)
			defer ticker.Stop()
			for range ticker.C {
				rateLimiter.Cleanup()
			}
		}()
	}

	fiberApp.Post("/auth/verify-code", authHandler.VerifyCode)
	fiberApp.Post("/auth/request-code", authHandler.RequestCode)
	fiberApp.Post("/auth/register", authHandler.Register)

	fiberApp.Get("/user/me", userHandler.GetUser)
	fiberApp.Get("/user/search", userHandler.SearchByUsername)
	fiberApp.Get("/user/:id", userHandler.GetUserById)

	fiberApp.Post("/chat/create", chatHandler.CreateChat)
	fiberApp.Get("/chats", chatHandler.GetChatsByUserId)
	fiberApp.Get("/chat/:id", chatHandler.GetChatById)
	fiberApp.Get("/chat/:id/messages", chatHandler.GetChatMessages)
	fiberApp.Get("/chat/:c_id/messages/after/:m_id", chatHandler.GetChatMessagesAfter)
	fiberApp.Post("/chat/:id/addkeys", chatHandler.AddChatKeys)

	fiberApp.Get("/message/:id", messageHandler.GetMessageById)

	hub := types.NewHub(chatApp, messageApp, userApp, jwtSvc, tokenSvc)
	fiberApp.Get("/ws", websocket.New(handler.HandleWS(hub)))

	log.Fatal(fiberApp.Listen(fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)))
}
