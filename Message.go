package main

import "time"

type MessageSlice []*Message

type Message struct {
	textb  []byte
	sender *Client
	ptime  time.Time
}

func newMessage(msg []byte, c *Client, pt time.Time) *Message {
	return &Message{
		textb:  msg,
		sender: c,
		ptime:  pt,
	}
}
