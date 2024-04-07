package main

type MessageSlice []*Message

type Message struct {
	textb  []byte
	sender *Client
	ptime  string
}

func newMessage(msg []byte, c *Client, pt string) *Message {
	return &Message{
		textb:  msg,
		sender: c,
		ptime:  pt,
	}
}
