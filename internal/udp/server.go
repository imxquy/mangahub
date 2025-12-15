package udp

import (
	"encoding/json"
	"net"
	"sync"
	"time"
)

type RegisterMessage struct {
	Type   string `json:"type"`   
	UserID string `json:"user_id"`
}

type Notification struct {
	Type      string `json:"type"` 
	MangaID   string `json:"manga_id"`
	Message   string `json:"message"`
	Timestamp int64  `json:"timestamp"`
}

type Server struct {
	addr    string
	conn    *net.UDPConn
	mu      sync.Mutex
	clients map[string]*net.UDPAddr // key = addr.String()
}

func NewServer(addr string) *Server {
	return &Server{
		addr:    addr,
		clients: make(map[string]*net.UDPAddr),
	}
}

func (s *Server) Run() error {
	udpAddr, err := net.ResolveUDPAddr("udp", s.addr)
	if err != nil {
		return err
	}
	c, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return err
	}
	s.conn = c

	buf := make([]byte, 2048)
	for {
		n, clientAddr, err := s.conn.ReadFromUDP(buf)
		if err != nil {
			continue
		}

		var reg RegisterMessage
		if err := json.Unmarshal(buf[:n], &reg); err != nil {
			continue
		}
		if reg.Type == "register" {
			s.mu.Lock()
			s.clients[clientAddr.String()] = clientAddr
			s.mu.Unlock()

			// optional ack
			ack := []byte(`{"type":"registered"}`)
			_, _ = s.conn.WriteToUDP(ack, clientAddr)
		}
	}
}

func (s *Server) BroadcastChapterRelease(mangaID, message string) {
	noti := Notification{
		Type:      "chapter_release",
		MangaID:   mangaID,
		Message:   message,
		Timestamp: time.Now().Unix(),
	}
	b, _ := json.Marshal(noti)

	s.mu.Lock()
	defer s.mu.Unlock()
	for _, a := range s.clients {
		_, _ = s.conn.WriteToUDP(b, a)
	}
}
