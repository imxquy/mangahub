package tcp

import (
	"encoding/json"
	"net"
	"time"
)

type Notifier struct {
	addr string
}

func NewNotifier(addr string) *Notifier {
	return &Notifier{addr: addr}
}

func (n *Notifier) SendProgressUpdate(p ProgressUpdate) error {
	conn, err := net.DialTimeout("tcp", n.addr, 2*time.Second)
	if err != nil {
		return err
	}
	defer conn.Close()

	b, err := json.Marshal(p)
	if err != nil {
		return err
	}
	b = append(b, '\n')
	_, err = conn.Write(b)
	return err
}
