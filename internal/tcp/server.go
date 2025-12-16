package tcp

import (
	"bufio"
	"encoding/json"
	"net"
	"sync"
)

type ProgressUpdate struct {
	UserID    string `json:"user_id"`
	MangaID   string `json:"manga_id"`
	Chapter   int    `json:"chapter"`
	Timestamp int64  `json:"timestamp"`
}

type Server struct {
	addr      string
	mu        sync.Mutex
	conns     map[net.Conn]struct{}
	Broadcast chan ProgressUpdate
}

func NewProgressSyncServer(addr string) *Server {
	return &Server{
		addr:      addr,
		conns:     make(map[net.Conn]struct{}),
		Broadcast: make(chan ProgressUpdate, 100),
	}
}

func (s *Server) Run() error {
	ln, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}
	go s.broadcastLoop()

	for {
		c, err := ln.Accept()
		if err != nil {
			continue
		}
		s.mu.Lock()
		s.conns[c] = struct{}{}
		s.mu.Unlock()

		go s.handleConn(c)
	}
}

func (s *Server) handleConn(c net.Conn) {
	defer func() {
		s.mu.Lock()
		delete(s.conns, c)
		s.mu.Unlock()
		_ = c.Close()
	}()

	sc := bufio.NewScanner(c)
	for sc.Scan() {
		line := sc.Bytes()

		var pu ProgressUpdate
		if err := json.Unmarshal(line, &pu); err == nil && pu.UserID != "" && pu.MangaID != "" {
			s.Broadcast <- pu
		}
	}
}

func (s *Server) broadcastLoop() {
	for msg := range s.Broadcast {
		b, _ := json.Marshal(msg)
		b = append(b, '\n')

		s.mu.Lock()
		for c := range s.conns {
			_, _ = c.Write(b)
		}
		s.mu.Unlock()
	}
}
