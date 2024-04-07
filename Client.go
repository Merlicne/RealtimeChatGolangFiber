package main

import (
	"encoding/json"
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
	// Room user has been joining
	s *Room
	// message waiting for user to receive
	receive chan *Message
}

func newClient(name string, conn *websocket.Conn, s *Room) *Client {
	return &Client{
		username: name,
		ws:       conn,
		receive:  make(chan *Message),
		s:        s,
	}
}

func (c *Client) write() {
	for {
		_, msg, err := c.ws.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseAbnormalClosure) {
				continue
			}
			log.Println("Write error :", err)
			return
		}

		var data struct {
			Text   string `json:"text"`
			Sender string `json:"sender"`
		}

		err = json.Unmarshal(msg, &data)
		if err != nil {
			log.Println("JSON decoding error:", err)
			return
		}

		if data.Sender == c.username {
			c.s.incoming_msg <- newMessage([]byte(data.Text), c)
		}

	}
}

func (c *Client) read() {
	c.ws.SetReadLimit(int64(messageSize))

	for msg := range c.receive {
		// send message to the user
		text, sender := msg.textb, msg.sender.username
		data := struct {
			Text   string `json:"text"`
			Sender string `json:"sender"`
		}{
			Text:   string(text),
			Sender: sender,
		}

		jsonData, err := json.Marshal(data)
		if err != nil {
			log.Println("JSON encoding error:", err)
			return
		}

		err = c.ws.WriteMessage(websocket.TextMessage, []byte(jsonData))
		if err != nil {
			log.Println("Read error :", err)
			return
		}
	}
}
