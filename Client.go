package main

import (
	"encoding/json"
	"log"
	"time"

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

func (c *Client) onMessage() ([]byte, error) {
	_, msg, err := c.ws.ReadMessage()
	if err != nil {
		if websocket.IsCloseError(err, websocket.CloseAbnormalClosure) {
			return nil, nil
		}
		log.Println("Write error :", err)
		c.s.unsub_user <- c
		return nil, err
	}
	return msg, nil
}

func (c *Client) deJson(msg []byte) (string, string, error) {
	var data struct {
		Text   string `json:"text"`
		Sender string `json:"sender"`
	}

	err := json.Unmarshal(msg, &data)
	if err != nil {
		log.Println("JSON decoding error:", err)
		return "", "", err
	}
	return data.Text, data.Sender, nil
}

func (c *Client) write() {
	for {

		msg, err := c.onMessage()
		if err != nil {
			return
		}
		if msg == nil {
			continue
		}

		text, sender, err := c.deJson(msg)
		if err != nil {
			c.s.unsub_user <- c
			return
		}

		if sender == c.username {
			c.s.incoming_msg <- newMessage([]byte(text), c, time.Now())
		}

	}
}

func (c *Client) enJson(msg *Message) ([]byte, error) {
	text, sender, pt := msg.textb, msg.sender.username, msg.ptime
	data := struct {
		Text   string `json:"text"`
		Sender string `json:"sender"`
		Ptime  string `json:"ptime"`
	}{
		Text:   string(text),
		Sender: sender,
		Ptime:  pt.Format(time.DateTime),
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Println("JSON encoding error:", err)
		return nil, err
	}
	return jsonData, nil
}

func (c *Client) read() {
	c.ws.SetReadLimit(int64(messageSize))

	for msg := range c.receive {
		// send message to the user
		jsonData, err := c.enJson(msg)
		if err != nil {
			c.s.unsub_user <- c
			return
		}

		err = c.ws.WriteMessage(websocket.TextMessage, jsonData)
		if err != nil {
			log.Println("Read error :", err)
			c.s.unsub_user <- c
			return
		}
	}
}
