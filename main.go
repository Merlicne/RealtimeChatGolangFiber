package main

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func main() {

	app := fiber.New()
	server := newServer()
	app.Static("/", "./webpage", fiber.Static{})
	app.Get("/", serveHome)
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
