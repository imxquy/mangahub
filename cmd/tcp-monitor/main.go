package main

import (
	"bufio"
	"fmt"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:9090")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	fmt.Println("TCP monitor connected to 127.0.0.1:9090")
	fmt.Println("Waiting for broadcasts...")

	sc := bufio.NewScanner(conn)
	for sc.Scan() {
		fmt.Println("<<", sc.Text())
	}
	fmt.Println("TCP monitor disconnected.")
}
