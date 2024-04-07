package main

type MessageSlice []*Message

type Message struct {
	textb  []byte
	sender *Client
}

func newMessage(msg []byte, c *Client) *Message {
	return &Message{
		textb:  msg,
		sender: c,
	}
}
