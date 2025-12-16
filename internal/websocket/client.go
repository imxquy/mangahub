package websocket

import (
	"encoding/json"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	conn *websocket.Conn
	send chan []byte
}

// ChatMessage là format message broadcast ra cho tất cả client
// đúng theo spec WebSocket chat
type WSChatMessage struct {
	Type      string `json:"type"` // "chat"
	Message   string `json:"message"`
	Timestamp int64  `json:"timestamp"`
}

func NewClient(conn *websocket.Conn) *Client {
	return &Client{
		conn: conn,
		send: make(chan []byte, 64),
	}
}

func (c *Client) ReadLoop(h *Hub) {
	defer func() {
		h.Unregister(c)
		close(c.send)
		_ = c.conn.Close()
	}()

	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			return
		}

		// Wrap raw text into JSON message
		m := WSChatMessage{
			Type:      "chat",
			Message:   string(msg),
			Timestamp: time.Now().Unix(),
		}

		b, _ := json.Marshal(m)
		h.Broadcast(b)
	}
}

func (c *Client) WriteLoop() {
	defer func() { _ = c.conn.Close() }()
	for msg := range c.send {
		if err := c.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			return
		}
	}
}
