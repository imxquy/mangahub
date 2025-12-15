package websocket

import "github.com/gorilla/websocket"

type Client struct {
	conn *websocket.Conn
	send chan []byte
}

func NewClient(conn *websocket.Conn) *Client {
	return &Client{conn: conn, send: make(chan []byte, 64)}
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
		h.Broadcast(msg)
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
