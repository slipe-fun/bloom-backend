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
	authapp "github.com/slipe-fun/skid-backend/internal/app/auth"
	chatapp "github.com/slipe-fun/skid-backend/internal/app/chat"
	exchangeapp "github.com/slipe-fun/skid-backend/internal/app/exchange"
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
	authhandler "github.com/slipe-fun/skid-backend/internal/transport/http/auth"
	chathandler "github.com/slipe-fun/skid-backend/internal/transport/http/chat"
	exchangehandler "github.com/slipe-fun/skid-backend/internal/transport/http/exchange"
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
	keysApp := keysapp.NewKeysApp(keysRepo, userRepo)
	authApp := authapp.NewAuthApp(sessionApp, keysApp, userRepo, rdb)
	exchangeApp := exchangeapp.NewExchangeApp(sessionApp, userRepo, rdb)
	userApp := userapp.NewUserApp(userRepo)
	chatApp := chatapp.NewChatApp(chatRepo, messageRepo)
	messageApp := messageapp.NewMessageApp(messageRepo, chatApp)

	hub := types.NewHub(sessionApp, chatApp)

	authHandler := authhandler.NewAuthHandler(authApp)
	exchangeHandler := exchangehandler.NewExchangeHandler(exchangeApp)
	userHandler := userhandler.NewUserHandler(userApp)
	chatHandler := chathandler.NewChatHandler(chatApp, userApp, messageApp, hub)
	messageHandler := messagehandler.NewMessageHandler(chatApp, messageApp, hub)
	sessionHandler := sessionhandler.NewSessionHandler(sessionApp, chatRepo)
	keysHandler := keyshandler.NewKeysHandler(keysApp)

	fiberApp := fiber.New()

	fiberApp.Use(recover.New(recover.Config{
		EnableStackTrace: true,
	}))
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

	authMiddleware := middleware.NewAuthMiddleware(sessionApp, userApp)

	authGroup := fiberApp.Group("/auth")
	exchangeGroup := fiberApp.Group("/exchange")
	userGroup := fiberApp.Group("/user")
	chatGroup := fiberApp.Group("/chat")
	messageGroup := fiberApp.Group("/message")
	sessionGroup := fiberApp.Group("/session")

	authGroup.Post("/register", authHandler.Register)
	authGroup.Get("/login/begin/:auth_lookup_id", authHandler.LoginBegin)
	authGroup.Post("/login/finish", authHandler.LoginFinish)

	exchangeGroup.Post("/session", exchangeHandler.StartSession)

	userGroup.Get("/me", authMiddleware.Handle(), userHandler.GetUser)
	userGroup.Post("/edit", authMiddleware.Handle(), userHandler.EditUser)
	userGroup.Get("/search", userHandler.SearchByUsername)
	userGroup.Get("/:id", userHandler.GetUserByID)
	userGroup.Get("/keys/:type", authMiddleware.Handle(), keysHandler.GetUserKeys)

	fiberApp.Get("/users", userHandler.GetAllUsers)

	chatGroup.Post("/create", authMiddleware.Handle(), chatHandler.CreateChat)
	chatGroup.Get("/:id", authMiddleware.Handle(), chatHandler.GetChatByID)
	chatGroup.Get("/:id/read", authMiddleware.Handle(), chatHandler.GetChatLastReadMessage)
	chatGroup.Get("/:c_id/messages/after/:m_id", authMiddleware.Handle(), chatHandler.GetChatMessagesAfter)
	chatGroup.Get("/:c_id/messages/before/:m_id", authMiddleware.Handle(), chatHandler.GetChatMessagesBefore)

	fiberApp.Get("/chats", authMiddleware.Handle(), chatHandler.GetChatsByUserID)

	messageGroup.Get("/:id", authMiddleware.Handle(), messageHandler.GetMessageByID)
	messageGroup.Post("/send", authMiddleware.Handle(), messageHandler.Send)
	messageGroup.Post("/seen", authMiddleware.Handle(), messageHandler.Seen)

	sessionGroup.Get("/", authMiddleware.Handle(), sessionHandler.GetSessionByToken)
	sessionGroup.Post("/:id/delete", authMiddleware.Handle(), sessionHandler.DeleteSession)

	fiberApp.Get("/sessions", authMiddleware.Handle(), sessionHandler.GetUserSessions)

	fiberApp.Get("/metrics", adaptor.HTTPHandler(promhttp.Handler()))

	fiberApp.Get("/ws", websocket.New(handler.HandleWS(hub)))
	fiberApp.Get("/exchange/ws", websocket.New(exchangehandler.HandleExchangeWS(hub, rdb)))

	log.Fatal(fiberApp.Listen(fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)))
}
