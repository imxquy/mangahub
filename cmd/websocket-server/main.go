package main

import (
	"log"
	"net/http"

	ws "mangahub/internal/websocket"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func main() {
	hub := ws.NewHub()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		c := ws.NewClient(conn)
		hub.Register(c)
		go c.WriteLoop()
		go c.ReadLoop(hub)
	})

	log.Println("WebSocket Chat listening on :9093 (/ws)")
	log.Fatal(http.ListenAndServe(":9093", nil))
}
