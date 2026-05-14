package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"

	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	chatapp "github.com/slipe-fun/skid-backend/internal/app/chat"
	keysapp "github.com/slipe-fun/skid-backend/internal/app/keys"
	messageapp "github.com/slipe-fun/skid-backend/internal/app/message"
	sessionapp "github.com/slipe-fun/skid-backend/internal/app/session"
	userapp "github.com/slipe-fun/skid-backend/internal/app/user"
	authservice "github.com/slipe-fun/skid-backend/internal/auth"
	"github.com/slipe-fun/skid-backend/internal/config"
	"github.com/slipe-fun/skid-backend/internal/metrics"
	"github.com/slipe-fun/skid-backend/internal/pkg/logger"
	"github.com/slipe-fun/skid-backend/internal/redis"
	"github.com/slipe-fun/skid-backend/internal/repository"
	chatrepo "github.com/slipe-fun/skid-backend/internal/repository/chat"
	keysrepo "github.com/slipe-fun/skid-backend/internal/repository/keys"
	messagerepo "github.com/slipe-fun/skid-backend/internal/repository/message"
	sessionrepo "github.com/slipe-fun/skid-backend/internal/repository/session"
	userrepo "github.com/slipe-fun/skid-backend/internal/repository/user"
	chathandler "github.com/slipe-fun/skid-backend/internal/transport/http/chat"
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

	rdb, err := redis.InitRedis(cfg)
	if err != nil {
		log.Fatalf("Redis error: %v", err)
	}
	defer rdb.Close()

	if err := logger.Init("logs/app.log"); err != nil {
		panic(err)
	}

	metrics.Init()

	userRepo := userrepo.NewUserRepo(db)
	chatRepo := chatrepo.NewChatRepo(db, userRepo)
	messageRepo := messagerepo.NewMessageRepo(db)
	sessionRepo := sessionrepo.NewSessionRepo(db, userRepo)
	keysRepo := keysrepo.NewKeysRepo(db)
	jwtSvc := authservice.NewJWTService(cfg.JWT.Secret)
	tokenSvc := authservice.NewTokenService(jwtSvc)

	sessionApp := sessionapp.NewSessionApp(sessionRepo, userRepo, jwtSvc, tokenSvc)
	// authApp := authapp.NewAuthApp(sessionApp)
	userApp := userapp.NewUserApp(userRepo)
	chatApp := chatapp.NewChatApp(chatRepo, messageRepo)
	messageApp := messageapp.NewMessageApp(messageRepo, chatApp)
	keysApp := keysapp.NewKeysApp(keysRepo)

	hub := types.NewHub(sessionApp, chatApp)

	// authHandler := authhandler.NewAuthHandler(authApp)
	userHandler := userhandler.NewUserHandler(userApp)
	chatHandler := chathandler.NewChatHandler(chatApp, userApp, messageApp, hub)
	messageHandler := messagehandler.NewMessageHandler(chatApp, messageApp, hub)
	sessionHandler := sessionhandler.NewSessionHandler(sessionApp, chatRepo)
	keysHandler := keyshandler.NewKeysHandler(keysApp)

	fiberApp := fiber.New()

	fiberApp.Use(recover.New())
	fiberApp.Use(cors.New())

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

	fiberApp.Get("/user/me", authMiddleware.Handle(), userHandler.GetUser)
	fiberApp.Post("/user/edit", authMiddleware.Handle(), userHandler.EditUser)
	fiberApp.Get("/user/search", userHandler.SearchByUsername)
	fiberApp.Get("/user/:id", userHandler.GetUserByID)

	fiberApp.Get("/users", userHandler.GetAllUsers)

	fiberApp.Post("/chat/create", authMiddleware.Handle(), chatHandler.CreateChat)
	fiberApp.Get("/chats", authMiddleware.Handle(), chatHandler.GetChatsByUserID)
	fiberApp.Get("/chat/:id", authMiddleware.Handle(), chatHandler.GetChatByID)
	fiberApp.Get("/chat/:id/read", authMiddleware.Handle(), chatHandler.GetChatLastReadMessage)
	fiberApp.Get("/chat/:c_id/messages/after/:m_id", authMiddleware.Handle(), chatHandler.GetChatMessagesAfter)
	fiberApp.Get("/chat/:c_id/messages/before/:m_id", authMiddleware.Handle(), chatHandler.GetChatMessagesBefore)

	fiberApp.Get("/user/keys/:type", authMiddleware.Handle(), keysHandler.GetUserKeys)

	fiberApp.Get("/message/:id", authMiddleware.Handle(), messageHandler.GetMessageByID)
	fiberApp.Post("/message/send", authMiddleware.Handle(), messageHandler.Send)
	fiberApp.Post("/message/seen", authMiddleware.Handle(), messageHandler.Seen)

	fiberApp.Get("/sessions", authMiddleware.Handle(), sessionHandler.GetUserSessions)
	fiberApp.Get("/session", authMiddleware.Handle(), sessionHandler.GetSessionByToken)
	fiberApp.Post("/session/:id/delete", authMiddleware.Handle(), sessionHandler.DeleteSession)

	fiberApp.Get("/metrics", adaptor.HTTPHandler(promhttp.Handler()))

	fiberApp.Get("/ws", websocket.New(handler.HandleWS(hub)))

	log.Fatal(fiberApp.Listen(fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)))
}
