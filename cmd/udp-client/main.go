package main

import (
	"fmt"
	"net"
)

func main() {
	server, err := net.ResolveUDPAddr("udp", "127.0.0.1:9091")
	if err != nil {
		panic(err)
	}

	conn, err := net.DialUDP("udp", nil, server)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// Register
	_, _ = conn.Write([]byte(`{"type":"register","user_id":"u1"}`))
	fmt.Println("UDP client registered. Waiting for notifications on same socket...")

	buf := make([]byte, 4096)
	for {
		n, _, err := conn.ReadFromUDP(buf)
		if err != nil {
			continue
		}
		fmt.Println("<<", string(buf[:n]))
	}
}
