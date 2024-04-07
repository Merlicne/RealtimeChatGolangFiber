package main

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	// server start
	server := newServer()
	// CORS middleware
	app.Use(func(c *fiber.Ctx) error {
		c.Set("Access-Control-Allow-Origin", "*")
		c.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
		c.Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		return c.Next()
	})
	app.Static("/", "./webpage", fiber.Static{})
	app.Get("/", serveHome)

	// serve Web socket for user
	app.Get("/ws/:room", websocket.New(server.serveWS))
	app.Listen(":3000")
}

func serveHome(c *fiber.Ctx) error {
	return c.SendFile("webpage")
}

func IsWebSocket(c *fiber.Ctx) error {
	if websocket.IsWebSocketUpgrade(c) {
		c.Locals("allowed", true)
		return c.Next()
	}
	return fiber.ErrUpgradeRequired
}
