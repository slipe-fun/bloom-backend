package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"

	AuthApp "github.com/slipe-fun/skid-backend/internal/app/auth"
	ChatApp "github.com/slipe-fun/skid-backend/internal/app/chat"
	MessageApp "github.com/slipe-fun/skid-backend/internal/app/message"
	UserApp "github.com/slipe-fun/skid-backend/internal/app/user"
	"github.com/slipe-fun/skid-backend/internal/config"
	"github.com/slipe-fun/skid-backend/internal/repository"
	ChatRepo "github.com/slipe-fun/skid-backend/internal/repository/chat"
	MessageRepo "github.com/slipe-fun/skid-backend/internal/repository/message"
	UserRepo "github.com/slipe-fun/skid-backend/internal/repository/user"
	"github.com/slipe-fun/skid-backend/internal/service"
	"github.com/slipe-fun/skid-backend/internal/transport/http/auth"
	"github.com/slipe-fun/skid-backend/internal/transport/http/chat"
	"github.com/slipe-fun/skid-backend/internal/transport/http/message"
	"github.com/slipe-fun/skid-backend/internal/transport/http/user"
	"github.com/slipe-fun/skid-backend/internal/transport/ws/handler"
	"github.com/slipe-fun/skid-backend/internal/transport/ws/types"
)

func main() {
	cfg := config.LoadConfig("configs/config.yaml")

	db := repository.InitDB(cfg)
	defer db.Close()

	userRepo := UserRepo.NewUserRepo(db)
	chatRepo := ChatRepo.NewChatRepo(db)
	messageRepo := MessageRepo.NewMessageRepo(db)

	jwtSvc := service.NewJWTService(cfg.JWT.Secret)
	tokenSvc := service.NewTokenService(jwtSvc)

	authApp := AuthApp.NewAuthApp(userRepo, jwtSvc)
	userApp := UserApp.NewUserApp(userRepo, jwtSvc, tokenSvc)
	chatApp := ChatApp.NewChatApp(chatRepo, tokenSvc)
	messageApp := MessageApp.NewMessageApp(messageRepo, chatApp, tokenSvc)

	authHandler := auth.NewAuthHandler(authApp)
	userHandler := user.NewUserHandler(userApp)
	chatHandler := chat.NewChatHandler(chatApp, userApp, messageApp)
	messageHandler := message.NewMessageHandler(chatApp, userApp, messageApp)

	fiberApp := fiber.New()

	fiberApp.Post("/auth/login", authHandler.Login)
	fiberApp.Post("/auth/register", authHandler.Register)

	fiberApp.Get("/user/me", userHandler.GetUser)
	fiberApp.Get("/user/:id", userHandler.GetUserById)

	fiberApp.Post("/chat/create", chatHandler.CreateChat)
	fiberApp.Get("/chats", chatHandler.GetChatsByUserId)
	fiberApp.Get("/chat/:id", chatHandler.GetChatById)
	fiberApp.Get("/chat/:id/messages", chatHandler.GetChatMessages)
	fiberApp.Post("/chat/:id/addkeys", chatHandler.AddChatKeys)

	fiberApp.Get("/message/:id", messageHandler.GetMessageById)

	hub := types.NewHub()
	fiberApp.Get("/ws", websocket.New(handler.HandleWS(hub)))

	log.Fatal(fiberApp.Listen(fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)))
}
