package handlers

import (
	"github.com/gorilla/websocket"
	"github.com/komly/chatik/mux"
	"github.com/komly/chatik/types"
	"log"
	"net/http"
)

type ws struct {
	u   websocket.Upgrader
	mux *mux.Mux
}

func NewWS(mux *mux.Mux) *ws {
	w := &ws{
		u: websocket.Upgrader{
			ReadBufferSize:  2048,
			WriteBufferSize: 2048,
		},
		mux: mux,
	}
	return w
}

func (wsk ws) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := wsk.u.Upgrade(w, r, nil)
	if err != nil {
		log.Print(err)
		http.Error(w, "Can't upgrade connection", http.StatusInternalServerError)
	}
	defer conn.Close()

	c := mux.NewClient(conn)
	go c.WritePump()
	defer c.Close()

	wsk.mux.Add(c)
	defer wsk.mux.Remove(c)

	req := types.Request{}
	for {
		err := conn.ReadJSON(&req)
		if err != nil {
			log.Print(err)
			return
		}
		log.Print(req)
		if err := wsk.mux.Process(c, &req); err != nil {
			log.Print(err)
			return
		}
	}
}
