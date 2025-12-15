package main

import (
	"log"
	"mangahub/internal/udp"
)

func main() {
	s := udp.NewServer(":9091")
	log.Println("UDP Notify listening on :9091")
	if err := s.Run(); err != nil {
		log.Fatalf("udp server: %v", err)
	}
}
