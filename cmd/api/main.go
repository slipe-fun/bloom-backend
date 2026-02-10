package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"

	authapp "github.com/slipe-fun/skid-backend/internal/app/auth"
	chatapp "github.com/slipe-fun/skid-backend/internal/app/chat"
	friendapp "github.com/slipe-fun/skid-backend/internal/app/friend"
	keysapp "github.com/slipe-fun/skid-backend/internal/app/keys"
	messageapp "github.com/slipe-fun/skid-backend/internal/app/message"
	sessionapp "github.com/slipe-fun/skid-backend/internal/app/session"
	userapp "github.com/slipe-fun/skid-backend/internal/app/user"
	verificationapp "github.com/slipe-fun/skid-backend/internal/app/verification"
	"github.com/slipe-fun/skid-backend/internal/config"
	"github.com/slipe-fun/skid-backend/internal/repository"
	chatrepo "github.com/slipe-fun/skid-backend/internal/repository/chat"
	friendrepo "github.com/slipe-fun/skid-backend/internal/repository/friend"
	keysrepo "github.com/slipe-fun/skid-backend/internal/repository/keys"
	messagerepo "github.com/slipe-fun/skid-backend/internal/repository/message"
	sessionrepo "github.com/slipe-fun/skid-backend/internal/repository/session"
	userrepo "github.com/slipe-fun/skid-backend/internal/repository/user"
	verificationrepo "github.com/slipe-fun/skid-backend/internal/repository/verification"
	"github.com/slipe-fun/skid-backend/internal/service"
	"github.com/slipe-fun/skid-backend/internal/service/logger"
	"github.com/slipe-fun/skid-backend/internal/service/oauth2"
	authhandler "github.com/slipe-fun/skid-backend/internal/transport/http/auth"
	chathandler "github.com/slipe-fun/skid-backend/internal/transport/http/chat"
	friendhandler "github.com/slipe-fun/skid-backend/internal/transport/http/friend"
	keyshandler "github.com/slipe-fun/skid-backend/internal/transport/http/keys"
	messagehandler "github.com/slipe-fun/skid-backend/internal/transport/http/message"
	"github.com/slipe-fun/skid-backend/internal/transport/http/middleware"
	sessionhandler "github.com/slipe-fun/skid-backend/internal/transport/http/session"
	userhandler "github.com/slipe-fun/skid-backend/internal/transport/http/user"
	"github.com/slipe-fun/skid-backend/internal/transport/ws/handler"
	"github.com/slipe-fun/skid-backend/internal/transport/ws/types"
)

func main() {
	cfg := config.LoadConfig("configs/config.yaml")

	db := repository.InitDB(cfg)
	defer db.Close()

	if err := logger.Init("logs/app.log"); err != nil {
		panic(err)
	}

	googleService := oauth2.NewGoogleAuthService(
		cfg.GoogleAuth.ClientId,
		cfg.GoogleAuth.ClientSecret,
		cfg.GoogleAuth.RedirectURL,
	)

	verificationRepo := verificationrepo.NewVerificationRepo(db)
	userRepo := userrepo.NewUserRepo(db, verificationRepo)
	chatRepo := chatrepo.NewChatRepo(db, userRepo)
	messageRepo := messagerepo.NewMessageRepo(db)
	sessionRepo := sessionrepo.NewSessionRepo(db, userRepo)
	keysRepo := keysrepo.NewKeysRepo(db, chatRepo)
	friendRepo := friendrepo.NewFriendRepo(db)

	jwtSvc := service.NewJWTService(cfg.JWT.Secret)
	tokenSvc := service.NewTokenService(jwtSvc)

	sessionApp := sessionapp.NewSessionApp(sessionRepo, userRepo, jwtSvc, tokenSvc)
	verificationApp := verificationapp.NewAuthApp(verificationRepo)
	authApp := authapp.NewAuthApp(sessionApp, userRepo, verificationRepo, verificationApp, jwtSvc, googleService)
	userApp := userapp.NewUserApp(sessionApp, userRepo, jwtSvc, tokenSvc)
	chatApp := chatapp.NewChatApp(sessionApp, chatRepo, tokenSvc)
	messageApp := messageapp.NewMessageApp(sessionApp, messageRepo, chatApp, tokenSvc)
	keysApp := keysapp.NewKeysApp(sessionApp, keysRepo, userApp, chatApp)
	friendApp := friendapp.NewFriendApp(sessionApp, friendRepo, userRepo, tokenSvc)

	hub := types.NewHub(sessionApp, chatApp, messageApp, userApp, jwtSvc, tokenSvc)

	authHandler := authhandler.NewAuthHandler(authApp, (*oauth2.GoogleAuthService)(googleService))
	userHandler := userhandler.NewUserHandler(userApp, friendApp)
	chatHandler := chathandler.NewChatHandler(chatApp, userApp, messageApp, hub)
	messageHandler := messagehandler.NewMessageHandler(chatApp, userApp, messageApp, hub)
	sessionHandler := sessionhandler.NewSessionHandler(sessionApp)
	keysHandler := keyshandler.NewKeysHandler(keysApp, chatApp)
	friendHandler := friendhandler.NewFriendHandler(friendApp)

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

	fiberApp.Get("/friends/:status", friendHandler.GetFriends)
	fiberApp.Post("/friend/request", friendHandler.SendRequest)
	fiberApp.Post("/friend/delete", friendHandler.DeleteFriend)

	fiberApp.Post("/chat/create", chatHandler.CreateChat)
	fiberApp.Get("/chats", chatHandler.GetChatsByUserId)
	fiberApp.Get("/chat/:id", chatHandler.GetChatById)
	fiberApp.Get("/chat/:id/read", chatHandler.GetChatLastReadMessage)
	fiberApp.Get("/chat/:id/messages", chatHandler.GetChatMessages)
	fiberApp.Get("/chat/:c_id/messages/after/:m_id", chatHandler.GetChatMessagesAfter)
	fiberApp.Get("/chat/:c_id/messages/before/:m_id", chatHandler.GetChatMessagesBefore)
	fiberApp.Post("/chat/:id/keys/public", chatHandler.AddChatKeys)

	fiberApp.Post("/chats/keys/private", keysHandler.SaveChatKeys)
	fiberApp.Get("/chats/keys/private", keysHandler.GetUserChatsKeys)

	fiberApp.Get("/message/:id", messageHandler.GetMessageById)
	fiberApp.Post("/message/send", messageHandler.Send)
	fiberApp.Post("/message/seen", messageHandler.Seen)

	fiberApp.Get("/sessions", sessionHandler.GetUserSessions)
	fiberApp.Get("/session", sessionHandler.GetSessionByToken)
	fiberApp.Post("/session/:id/delete", sessionHandler.DeleteSession)

	fiberApp.Get("/ws", websocket.New(handler.HandleWS(hub)))

	log.Fatal(fiberApp.Listen(fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)))
}
