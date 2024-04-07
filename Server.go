package main

import (
	"log"
	"sync"

	"github.com/gofiber/contrib/websocket"
)

type Server struct {
	// CLient collection
	Clients ClientList

	// hold incoming message in channel
	incoming_msg chan []byte

	sub_user chan *Client

	unsub_user chan *Client

	sync.RWMutex
}

func newServer() *Server {
	return &Server{
		Clients:      make(ClientList),
		incoming_msg: make(chan []byte),
		sub_user:     make(chan *Client),
		unsub_user:   make(chan *Client),
	}
}

func (s *Server) running() {
	for {
		select {
		case c := <-s.sub_user:
			s.user_connect(c)
		case c := <-s.unsub_user:
			s.user_disconnect(c)
		case msg := <-s.incoming_msg:
			go s.broadcast(msg)
		}
	}
}

func (s *Server) broadcast(msg []byte) {
	for c := range s.Clients {
		c.receive <- msg
	}
}

func (s *Server) user_connect(c *Client) {
	s.Lock()
	s.Clients[c] = true
	msg := c.username + " has connected"
	log.Println(msg)
	go s.broadcast([]byte(msg))
	s.Unlock()
}

func (s *Server) user_disconnect(c *Client) {
	s.Lock()
	delete(s.Clients, c)
	c.ws.Close()
	msg := c.username + " has disconnected"
	log.Println(msg)
	go s.broadcast([]byte(msg))
	s.Unlock()
}

func (s *Server) serveWS(conn *websocket.Conn) {
	client := newClient(conn.Query("username"), conn, s)
	s.sub_user <- client
	defer func() { s.unsub_user <- client }()
	go client.write()
	client.read()
}
