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

	r := s.searchRoom(room_name)

	client := newClient(username, conn, r)

	r.serve(client)
}

func (s *Server) searchRoom(roomname string) *Room {
	r, ok := s.rooms[roomname]
	if !ok {
		r = newRoom(roomname)
		s.rooms[roomname] = r
		go r.running()
	}
	return r
}
