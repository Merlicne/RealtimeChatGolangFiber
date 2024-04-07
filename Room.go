package main

import (
	"log"
	"sync"
	"time"
)

type RoomList map[string]*Room

type Room struct {
	room_name string
	// CLient collection
	Clients ClientList

	MessageHistory MessageSlice

	// hold incoming message in channel
	incoming_msg chan *Message

	sub_user chan *Client

	unsub_user chan *Client

	sync.RWMutex
}

func newRoom(name string) *Room {
	return &Room{
		room_name:    name,
		Clients:      make(ClientList),
		incoming_msg: make(chan *Message),
		sub_user:     make(chan *Client),
		unsub_user:   make(chan *Client),
	}
}

func (s *Room) running() {

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

func (s *Room) broadcast(msg *Message) {
	for c := range s.Clients {
		c.receive <- msg
	}
	s.MessageHistory = append(s.MessageHistory, msg)
}

func (s *Room) reChat(c *Client) {
	for _, msg := range s.MessageHistory {
		c.receive <- msg
	}
}

func (s *Room) user_connect(c *Client) {
	s.Lock()
	s.Clients[c] = true
	s.Unlock()
	msg := " has connected "
	log.Println(c.username + msg + "to " + s.room_name)
	s.reChat(c)
	go s.broadcast(newMessage([]byte(msg), c, time.Now().Format(time.DateTime)))
}

func (s *Room) user_disconnect(c *Client) {
	s.Lock()
	if _, ok := s.Clients[c]; ok {
		delete(s.Clients, c)
		c.ws.Close()
	}
	s.Unlock()
	msg := " has disconnected"
	log.Println(c.username + msg + " from " + s.room_name)
	go s.broadcast(newMessage([]byte(msg), c, time.Now().Format(time.DateTime)))
}

func (s *Room) serve(c *Client) {
	s.sub_user <- c
	go c.write()
	c.read()
}
