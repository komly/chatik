package main

import (
	"github.com/komly/chatik/handlers"
	"github.com/komly/chatik/mux"
	"log"
	"net/http"
)

func main() {
	mux := mux.NewMux()
	go mux.Run()

	http.Handle("/ws", handlers.NewWS(mux))
	http.Handle("/", handlers.NewHome())
	log.Fatal(http.ListenAndServe(":3000", nil))
}
