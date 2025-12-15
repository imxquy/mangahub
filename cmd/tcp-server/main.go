package main

import (
	"log"
	"mangahub/internal/tcp"
)

func main() {
	s := tcp.NewProgressSyncServer(":9090")
	log.Println("TCP Sync listening on :9090")
	if err := s.Run(); err != nil {
		log.Fatalf("tcp server: %v", err)
	}
}
