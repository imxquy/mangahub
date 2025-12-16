package udp

import (
	"encoding/json"
	"net"
	"time"
)

type Control struct {
	addr string
}

func NewControl(addr string) *Control {
	return &Control{addr: addr}
}

type BroadcastCmd struct {
	Type    string `json:"type"` // "broadcast"
	MangaID string `json:"manga_id"`
	Message string `json:"message"`
}

func (c *Control) Broadcast(mangaID, message string) error {
	ra, err := net.ResolveUDPAddr("udp", c.addr)
	if err != nil {
		return err
	}
	conn, err := net.DialUDP("udp", nil, ra)
	if err != nil {
		return err
	}
	defer conn.Close()

	_ = conn.SetWriteDeadline(time.Now().Add(2 * time.Second))
	b, err := json.Marshal(BroadcastCmd{Type: "broadcast", MangaID: mangaID, Message: message})
	if err != nil {
		return err
	}
	_, err = conn.Write(b)
	return err
}
