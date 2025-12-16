package udp

import (
	"encoding/json"
	"net"
	"sync"
	"time"
)

type Message struct {
	Type    string `json:"type"` // "register" | "broadcast"
	UserID  string `json:"user_id,omitempty"`
	MangaID string `json:"manga_id,omitempty"`
	Message string `json:"message,omitempty"`
}

type Notification struct {
	Type      string `json:"type"` // "chapter_release"
	MangaID   string `json:"manga_id"`
	Message   string `json:"message"`
	Timestamp int64  `json:"timestamp"`
}

type Server struct {
	addr    string
	conn    *net.UDPConn
	mu      sync.Mutex
	clients map[string]*net.UDPAddr
}

func NewServer(addr string) *Server {
	return &Server{
		addr:    addr,
		clients: make(map[string]*net.UDPAddr),
	}
}

func (s *Server) Run() error {
	udpAddr, err := net.ResolveUDPAddr("udp", s.addr)
	if err != nil { return err }

	c, err := net.ListenUDP("udp", udpAddr)
	if err != nil { return err }
	s.conn = c

	buf := make([]byte, 4096)
	for {
		n, clientAddr, err := s.conn.ReadFromUDP(buf)
		if err != nil { continue }

		var m Message
		if err := json.Unmarshal(buf[:n], &m); err != nil {
			continue
		}

		switch m.Type {
		case "register":
			s.mu.Lock()
			s.clients[clientAddr.String()] = clientAddr
			s.mu.Unlock()
			_, _ = s.conn.WriteToUDP([]byte(`{"type":"registered"}`), clientAddr)

		case "broadcast":
			// API gửi lệnh broadcast vào đây
			if m.MangaID != "" && m.Message != "" {
				s.BroadcastChapterRelease(m.MangaID, m.Message)
			}
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
