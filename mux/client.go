package mux

import (
	"github.com/gorilla/websocket"
	"github.com/komly/chatik/types"
)

type Client struct {
	out   chan *types.Response
	close chan interface{}
	conn  *websocket.Conn
}

func NewClient(conn *websocket.Conn) *Client {
	c := &Client{
		out:   make(chan *types.Response),
		close: make(chan interface{}),
		conn:  conn,
	}
	return c
}

func (c Client) WritePump() {
	for {
		select {
		case r := <-c.out:
			c.conn.WriteJSON(r)
			break
		case <-c.close:
			c.conn.Close()
			return
		}

	}
}

func (c Client) Write(resp *types.Response) {
	c.out <- resp
}

func (c Client) Close() {
	c.close <- struct{}{}
}
