package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	// "github.com/gofiber/websocket/v2"
	"github.com/slipe-fun/skid-backend/internal/app"
	"github.com/slipe-fun/skid-backend/internal/config"
	"github.com/slipe-fun/skid-backend/internal/repository"
	"github.com/slipe-fun/skid-backend/internal/service"
	"github.com/slipe-fun/skid-backend/internal/transport/http"
)

func main() {
	cfg := config.LoadConfig("configs/config.yaml")

	db := repository.InitDB(cfg)
	defer db.Close()

	userRepo := repository.NewUserRepo(db)

	jwtSvc := service.NewJWTService(cfg.JWT.Secret)

	authApp := app.NewAuthApp(userRepo, jwtSvc)
	userApp := app.NewUserApp(userRepo, jwtSvc)

	authHandler := http.NewAuthHandler(authApp)
	userHandler := http.NewUserHandler(userApp)

	fiberApp := fiber.New()

	fiberApp.Post("/login", authHandler.Login)
	fiberApp.Post("/register", authHandler.Register)
	fiberApp.Get("/user/me", userHandler.GetUser)
	fiberApp.Get("/user/:id", userHandler.GetUserById)

	// fiberApp.Get("/ws", websocket.New(func(c *websocket.Conn) {
	// 	defer c.Close()
	// 	for {
	// 		mt, msg, err := c.ReadMessage()
	// 		if err != nil {
	// 			break
	// 		}
	// 		c.WriteMessage(mt, msg)
	// 	}
	// }))

	log.Fatal(fiberApp.Listen(fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)))
}
