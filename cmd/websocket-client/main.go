package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/gorilla/websocket"
)

func main() {
	c, _, err := websocket.DefaultDialer.Dial("ws://localhost:9093/ws", nil)
	if err != nil {
		panic(err)
	}
	defer c.Close()

	fmt.Println("Connected to WebSocket server")

	// đọc message từ server
	go func() {
		for {
			_, msg, err := c.ReadMessage()
			if err != nil {
				fmt.Println("read error:", err)
				return
			}
			fmt.Println("recv:", string(msg))
		}
	}()

	// nhập message từ terminal
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()
		if err := c.WriteMessage(websocket.TextMessage, []byte(text)); err != nil {
			fmt.Println("write error:", err)
			return
		}
	}
}
