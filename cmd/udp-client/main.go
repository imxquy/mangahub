package main

import (
	"net"
)

func main() {
	addr, _ := net.ResolveUDPAddr("udp", "127.0.0.1:9091")
	c, _ := net.DialUDP("udp", nil, addr)
	defer c.Close()

	c.Write([]byte(`{"type":"register","user_id":"u1"}`))

	buf := make([]byte, 2048)
	for {
		n, _, _ := c.ReadFromUDP(buf)
		println(string(buf[:n]))
	}
}
