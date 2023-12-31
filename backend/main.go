package main

import (
	"backend/api"
	"backend/ws"
	"backend/storage"
	"errors"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	frecover "github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/contrib/websocket"
)

func main() {
	// Load the .env variables
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	// Setup storage
	storage.Setup()

	// Create a new fiber app
	app := fiber.New(fiber.Config{
		// Handle the error and return a suitable http code
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
	
			var e *fiber.Error
			if errors.As(err, &e) {
				code = e.Code
			}

			log.Println(ctx.Path(), err.Error())
	
			return ctx.SendStatus(code)
		},
	})

	// Middleware allowing for panics to be recovered from
	app.Use(frecover.New())
	
	cfg := websocket.Config{RecoverHandler: func(conn *websocket.Conn) {
		if err := recover(); err != nil {
            log.Println(err)
        }
	}}
	
	app.Use(cors.New(cors.Config{
        AllowHeaders:     "Origin, Content-Type, Accept, Content-Length, Accept-Language, Accept-Encoding, Connection, Access-Control-Allow-Origin",
        AllowOrigins:     os.Getenv("ORIGIN"),
        AllowCredentials: true,
        AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
    }))

	// setup routes

	app.Post("/signup", api.HandlePostSignup)
	app.Post("/signup/email", api.HandlePostEmail)
	app.Post("/login", api.HandlePostLogin)
	app.Post("/logout", api.Auth, api.HandlePostLogout)

	app.Post("/authed", api.Auth, api.Authed)

	app.Post("/friend/request", api.Auth, api.HandleSendFriendRequest)
	app.Post("/friend/request/decline", api.Auth, api.HandleDeclineFriendRequest)
	app.Post("/friend/request/accept", api.Auth, api.HandleAcceptFriendRequest)
	
	app.Get("/friends", api.Auth, api.HandleGetFriends)
	app.Get("/messages/:id", api.Auth, api.HandleGetMessages)

	app.Delete("/friend/:id", api.Auth, api.HandleRemoveFriend)
	app.Delete("/friend/request/:id", api.Auth, api.HandleRemoveFriendRequest)

	app.Post("/friend", api.Auth, api.HandleSendFriendRequest)

	app.Get("/ws", api.Auth, websocket.New(ws.Conn, cfg))

	// Listen on the port defined in .env
    panic(app.Listen(":" + os.Getenv("PORT")))
}
