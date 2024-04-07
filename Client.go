package main

import (
	"log"

	"github.com/gofiber/contrib/websocket"
)

var (
	messageSize = 1024
)

type ClientList map[*Client]bool

type Client struct {
	username string
	// User's web socket connection
	ws *websocket.Conn
	// Server user has been joining
	s *Server
	// message waiting for user to receive
	receive chan []byte
}

func newClient(name string, conn *websocket.Conn, s *Server) *Client {
	return &Client{
		username: name,
		ws:       conn,
		receive:  make(chan []byte),
		s:        s,
	}
}

func (c *Client) write() {
	for {
		_, msg, err := c.ws.ReadMessage()
		if err != nil {
			log.Println("Write error :", err)
			return
		}

		c.s.incoming_msg <- msg
	}
}

func (c *Client) read() {
	c.ws.SetReadLimit(int64(messageSize))
	for msg := range c.receive {
		// send message to the user
		err := c.ws.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			log.Println("Read error :", err)
			return
		}
	}
}
