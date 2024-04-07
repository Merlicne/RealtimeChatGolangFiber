package main

import (
	"github.com/gofiber/contrib/websocket"
)

type Server struct {
	rooms RoomList
}

func newServer() *Server {
	return &Server{
		rooms: make(RoomList),
	}
}

func (s *Server) serveWS(conn *websocket.Conn) {
	room_name := conn.Params("room")
	username := conn.Query("username")
	r, ok := s.rooms[room_name]
	if !ok {
		r = newRoom(room_name)
		s.rooms[room_name] = r
		go r.running()
	}

	client := newClient(username, conn, r)

	r.serve(client)
}
