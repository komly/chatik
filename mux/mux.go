package mux

import (
	"fmt"
	"github.com/komly/chatik/types"
	"log"
)

type Mux struct {
	clients map[*Client]interface{}
	ops     chan func(map[*Client]interface{})
}

func NewMux() *Mux {
	m := Mux{
		clients: make(map[*Client]interface{}),
		ops:     make(chan func(map[*Client]interface{})),
	}
	return &m
}

func (m *Mux) Add(c *Client) {
	m.ops <- func(clients map[*Client]interface{}) {
		clients[c] = struct{}{}
		log.Print("connected")
	}
}

func (m *Mux) Remove(c *Client) {
	m.ops <- func(clients map[*Client]interface{}) {
		delete(clients, c)
		log.Print("disconnected")
	}
}

func (m *Mux) Broadcast(c *Client, req *types.Request, ch chan error) {
	m.ops <- func(clients map[*Client]interface{}) {
		for cl := range clients {
			cl.Write(&types.Response{
				Type: req.Type,
			})
		}
		ch <- nil
	}
}

func (m *Mux) Process(c *Client, req *types.Request) error {
	ch := make(chan error, 1)

	switch req.Type {
	case "broadcast":
		m.Broadcast(c, req, ch)
		break
	default:
		m.ops <- func(clients map[*Client]interface{}) {
			c.Write(&types.Response{
				Type: "invalid_command",
			})
			c.Close()
			ch <- fmt.Errorf("Invalid command from user")
		}
	}
	return <-ch
}

func (m *Mux) Run() {
	for op := range m.ops {
		op(m.clients)
	}
}
