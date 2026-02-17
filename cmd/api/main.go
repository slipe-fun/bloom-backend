package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"

	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	authapp "github.com/slipe-fun/skid-backend/internal/app/auth"
	chatapp "github.com/slipe-fun/skid-backend/internal/app/chat"
	encryptedchatkeysapp "github.com/slipe-fun/skid-backend/internal/app/encrypted_chat_keys"
	friendapp "github.com/slipe-fun/skid-backend/internal/app/friend"
	keysapp "github.com/slipe-fun/skid-backend/internal/app/keys"
	messageapp "github.com/slipe-fun/skid-backend/internal/app/message"
	sessionapp "github.com/slipe-fun/skid-backend/internal/app/session"
	userapp "github.com/slipe-fun/skid-backend/internal/app/user"
	verificationapp "github.com/slipe-fun/skid-backend/internal/app/verification"
	authservice "github.com/slipe-fun/skid-backend/internal/auth"
	"github.com/slipe-fun/skid-backend/internal/config"
	"github.com/slipe-fun/skid-backend/internal/metrics"
	"github.com/slipe-fun/skid-backend/internal/oauth/google"
	"github.com/slipe-fun/skid-backend/internal/pkg/logger"
	"github.com/slipe-fun/skid-backend/internal/repository"
	chatrepo "github.com/slipe-fun/skid-backend/internal/repository/chat"
	encryptedchatkeysrepo "github.com/slipe-fun/skid-backend/internal/repository/encrypted_chat_keys"
	friendrepo "github.com/slipe-fun/skid-backend/internal/repository/friend"
	keysrepo "github.com/slipe-fun/skid-backend/internal/repository/keys"
	messagerepo "github.com/slipe-fun/skid-backend/internal/repository/message"
	sessionrepo "github.com/slipe-fun/skid-backend/internal/repository/session"
	userrepo "github.com/slipe-fun/skid-backend/internal/repository/user"
	verificationrepo "github.com/slipe-fun/skid-backend/internal/repository/verification"
	authhandler "github.com/slipe-fun/skid-backend/internal/transport/http/auth"
	chathandler "github.com/slipe-fun/skid-backend/internal/transport/http/chat"
	encryptedchatkeyshandler "github.com/slipe-fun/skid-backend/internal/transport/http/encrypted_chat_keys"
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

	metrics.Init()

	googleService := google.NewGoogleAuthService(
		cfg.GoogleAuth.ClientID,
		cfg.GoogleAuth.ClientSecret,
		cfg.GoogleAuth.RedirectURL,
	)

	verificationRepo := verificationrepo.NewVerificationRepo(db)
	userRepo := userrepo.NewUserRepo(db)
	chatRepo := chatrepo.NewChatRepo(db, userRepo)
	messageRepo := messagerepo.NewMessageRepo(db)
	sessionRepo := sessionrepo.NewSessionRepo(db, userRepo)
	keysRepo := keysrepo.NewKeysRepo(db)
	encryptedChatKeysRepo := encryptedchatkeysrepo.NewEncryptedChatKeysRepo(db)
	friendRepo := friendrepo.NewFriendRepo(db)

	jwtSvc := authservice.NewJWTService(cfg.JWT.Secret)
	tokenSvc := authservice.NewTokenService(jwtSvc)

	sessionApp := sessionapp.NewSessionApp(sessionRepo, userRepo, jwtSvc, tokenSvc)
	verificationApp := verificationapp.NewAuthApp(verificationRepo)
	authApp := authapp.NewAuthApp(sessionApp, userRepo, verificationRepo, verificationApp, googleService)
	userApp := userapp.NewUserApp(userRepo)
	chatApp := chatapp.NewChatApp(chatRepo)
	messageApp := messageapp.NewMessageApp(messageRepo, chatApp)
	keysApp := keysapp.NewKeysApp(keysRepo)
	encryptedChatKeysApp := encryptedchatkeysapp.NewEncryptedChatKeysApp(encryptedChatKeysRepo, chatRepo, sessionApp)
	friendApp := friendapp.NewFriendApp(friendRepo, userRepo)

	hub := types.NewHub(sessionApp, chatApp)

	authHandler := authhandler.NewAuthHandler(authApp)
	userHandler := userhandler.NewUserHandler(userApp, friendApp)
	chatHandler := chathandler.NewChatHandler(chatApp, userApp, messageApp, hub)
	messageHandler := messagehandler.NewMessageHandler(chatApp, messageApp, hub)
	sessionHandler := sessionhandler.NewSessionHandler(sessionApp, chatRepo)
	keysHandler := keyshandler.NewKeysHandler(keysApp)
	encryptedChatKeysHandler := encryptedchatkeyshandler.NewEncryptedChatKeysHandlerApp(encryptedChatKeysApp)
	friendHandler := friendhandler.NewFriendHandler(friendApp, hub)

	fiberApp := fiber.New()

	fiberApp.Use(recover.New())

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

	fiberApp.Use(middleware.MetricsMiddleware())

	authMiddleware := middleware.NewAuthMiddleware(sessionApp)

	fiberApp.Post("/auth/verify-code", authHandler.VerifyCode)
	fiberApp.Post("/auth/request-code", authHandler.RequestCode)
	fiberApp.Get("/oauth2/google/redirect", authHandler.GoogleRedirect)
	fiberApp.Get("/oauth2/google/exchange-code", authHandler.ExchangeCode)
	fiberApp.Post("/auth/register", authHandler.Register)

	fiberApp.Get("/user/me", authMiddleware.Handle(), userHandler.GetUser)
	fiberApp.Post("/user/edit", authMiddleware.Handle(), userHandler.EditUser)
	fiberApp.Get("/user/search", userHandler.SearchByUsername)
	fiberApp.Get("/user/exists", userHandler.IsUserWithEmailExists)
	fiberApp.Get("/user/:id", userHandler.GetUserByID)
	fiberApp.Get("/user/:id/key-bundle", authMiddleware.Handle(), sessionHandler.GetUserKeyBundle)

	fiberApp.Get("/users", userHandler.GetAllUsers)

	fiberApp.Get("/friends/:status", authMiddleware.Handle(), friendHandler.GetFriends)
	fiberApp.Post("/friend/request", authMiddleware.Handle(), friendHandler.SendRequest)
	fiberApp.Post("/friend/delete", authMiddleware.Handle(), friendHandler.DeleteFriend)

	fiberApp.Post("/chat/create", authMiddleware.Handle(), chatHandler.CreateChat)
	fiberApp.Get("/chats", authMiddleware.Handle(), chatHandler.GetChatsByUserID)
	fiberApp.Get("/chat/:id", authMiddleware.Handle(), chatHandler.GetChatByID)
	fiberApp.Get("/chat/:id/read", authMiddleware.Handle(), chatHandler.GetChatLastReadMessage)
	fiberApp.Get("/chat/:id/messages", authMiddleware.Handle(), chatHandler.GetChatMessages)
	fiberApp.Get("/chat/:c_id/messages/after/:m_id", authMiddleware.Handle(), chatHandler.GetChatMessagesAfter)
	fiberApp.Get("/chat/:c_id/messages/before/:m_id", authMiddleware.Handle(), chatHandler.GetChatMessagesBefore)
	fiberApp.Post("/chat/:id/key", authMiddleware.Handle(), encryptedChatKeysHandler.AddKeys)
	fiberApp.Post("/chat/:id/keys/public", authMiddleware.Handle(), chatHandler.AddChatKeys)

	fiberApp.Post("/chats/keys/private", authMiddleware.Handle(), keysHandler.SaveChatKeys)
	fiberApp.Get("/chats/keys/private", authMiddleware.Handle(), keysHandler.GetUserChatsKeys)

	fiberApp.Get("/message/:id", authMiddleware.Handle(), messageHandler.GetMessageByID)
	fiberApp.Post("/message/send", authMiddleware.Handle(), messageHandler.Send)
	fiberApp.Post("/message/seen", authMiddleware.Handle(), messageHandler.Seen)

	fiberApp.Get("/sessions", authMiddleware.Handle(), sessionHandler.GetUserSessions)
	fiberApp.Get("/session", authMiddleware.Handle(), sessionHandler.GetSessionByToken)
	fiberApp.Post("/session/add-keys", authMiddleware.Handle(), sessionHandler.AddKeys)
	fiberApp.Post("/session/:id/delete", authMiddleware.Handle(), sessionHandler.DeleteSession)

	fiberApp.Get("/metrics", adaptor.HTTPHandler(promhttp.Handler()))

	fiberApp.Get("/ws", websocket.New(handler.HandleWS(hub)))

	log.Fatal(fiberApp.Listen(fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)))
}
