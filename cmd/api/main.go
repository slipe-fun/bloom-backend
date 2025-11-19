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
	SessionApp "github.com/slipe-fun/skid-backend/internal/app/session"
	UserApp "github.com/slipe-fun/skid-backend/internal/app/user"
	VerificationApp "github.com/slipe-fun/skid-backend/internal/app/verification"
	"github.com/slipe-fun/skid-backend/internal/config"
	"github.com/slipe-fun/skid-backend/internal/repository"
	ChatRepo "github.com/slipe-fun/skid-backend/internal/repository/chat"
	MessageRepo "github.com/slipe-fun/skid-backend/internal/repository/message"
	SessionRepo "github.com/slipe-fun/skid-backend/internal/repository/session"
	UserRepo "github.com/slipe-fun/skid-backend/internal/repository/user"
	VerificationRepo "github.com/slipe-fun/skid-backend/internal/repository/verification"
	"github.com/slipe-fun/skid-backend/internal/service"
	"github.com/slipe-fun/skid-backend/internal/service/oauth2"
	"github.com/slipe-fun/skid-backend/internal/transport/http/auth"
	"github.com/slipe-fun/skid-backend/internal/transport/http/chat"
	"github.com/slipe-fun/skid-backend/internal/transport/http/message"
	"github.com/slipe-fun/skid-backend/internal/transport/http/middleware"
	"github.com/slipe-fun/skid-backend/internal/transport/http/session"
	"github.com/slipe-fun/skid-backend/internal/transport/http/user"
	"github.com/slipe-fun/skid-backend/internal/transport/ws/handler"
	"github.com/slipe-fun/skid-backend/internal/transport/ws/types"
)

func main() {
	cfg := config.LoadConfig("configs/config.yaml")

	db := repository.InitDB(cfg)
	defer db.Close()

	googleService := oauth2.NewGoogleAuthService(
		cfg.GoogleAuth.ClientId,
		cfg.GoogleAuth.ClientSecret,
		cfg.GoogleAuth.RedirectURL,
	)

	verificationRepo := VerificationRepo.NewVerificationRepo(db)
	userRepo := UserRepo.NewUserRepo(db, verificationRepo)
	chatRepo := ChatRepo.NewChatRepo(db, userRepo)
	messageRepo := MessageRepo.NewMessageRepo(db)
	sessionRepo := SessionRepo.NewSessionRepo(db, userRepo)

	jwtSvc := service.NewJWTService(cfg.JWT.Secret)
	tokenSvc := service.NewTokenService(jwtSvc)

	sessionApp := SessionApp.NewSessionApp(sessionRepo, userRepo, jwtSvc, tokenSvc)
	verificationApp := VerificationApp.NewAuthApp(verificationRepo)
	authApp := AuthApp.NewAuthApp(sessionApp, userRepo, verificationRepo, verificationApp, jwtSvc, googleService)
	userApp := UserApp.NewUserApp(sessionApp, userRepo, jwtSvc, tokenSvc)
	chatApp := ChatApp.NewChatApp(sessionApp, chatRepo, tokenSvc)
	messageApp := MessageApp.NewMessageApp(sessionApp, messageRepo, chatApp, tokenSvc)

	authHandler := auth.NewAuthHandler(authApp, (*oauth2.GoogleAuthService)(googleService))
	userHandler := user.NewUserHandler(userApp)
	chatHandler := chat.NewChatHandler(chatApp, userApp, messageApp)
	messageHandler := message.NewMessageHandler(chatApp, userApp, messageApp)
	sessionHandler := session.NewSessionHandler(sessionApp)

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
	fiberApp.Get("/oauth2/google/redirect", authHandler.GoogleRedirect)
	fiberApp.Get("/oauth2/google/exchange-code", authHandler.ExchangeCode)
	fiberApp.Post("/auth/register", authHandler.Register)

	fiberApp.Get("/user/me", userHandler.GetUser)
	fiberApp.Post("/user/edit", userHandler.EditUser)
	fiberApp.Get("/user/search", userHandler.SearchByUsername)
	fiberApp.Get("/user/exists", userHandler.IsUserWithEmailExists)
	fiberApp.Get("/user/:id", userHandler.GetUserById)

	fiberApp.Post("/chat/create", chatHandler.CreateChat)
	fiberApp.Get("/chats", chatHandler.GetChatsByUserId)
	fiberApp.Get("/chat/:id", chatHandler.GetChatById)
	fiberApp.Get("/chat/:id/messages", chatHandler.GetChatMessages)
	fiberApp.Get("/chat/:c_id/messages/after/:m_id", chatHandler.GetChatMessagesAfter)
	fiberApp.Post("/chat/:id/addkeys", chatHandler.AddChatKeys)

	fiberApp.Get("/message/:id", messageHandler.GetMessageById)

	fiberApp.Get("/sessions", sessionHandler.GetUserSessions)
	fiberApp.Get("/session", sessionHandler.GetSessionByToken)
	fiberApp.Post("/session/:id/delete", sessionHandler.DeleteSession)

	hub := types.NewHub(sessionApp, chatApp, messageApp, userApp, jwtSvc, tokenSvc)
	fiberApp.Get("/ws", websocket.New(handler.HandleWS(hub)))

	log.Fatal(fiberApp.Listen(fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)))
}
