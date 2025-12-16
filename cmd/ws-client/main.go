package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/websocket"
)

func main() {
	url := "ws://127.0.0.1:9093/ws"
	c, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		panic(err)
	}
	defer c.Close()

	fmt.Println("WS connected:", url)
	fmt.Println("Type a message and press Enter. Ctrl+C to quit.")

	// read loop
	go func() {
		for {
			_, msg, err := c.ReadMessage()
			if err != nil {
				fmt.Println("WS read error:", err)
				return
			}
			fmt.Println("<<", string(msg))
		}
	}()

	// graceful exit
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	sc := bufio.NewScanner(os.Stdin)
	for {
		select {
		case <-sig:
			_ = c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "bye"))
			return
		default:
		}

		if !sc.Scan() {
			return
		}
		line := sc.Text()
		if line == "" {
			continue
		}
		if err := c.WriteMessage(websocket.TextMessage, []byte(line)); err != nil {
			fmt.Println("WS write error:", err)
			return
		}
	}
}
