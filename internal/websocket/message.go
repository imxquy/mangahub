package websocket

type ChatMessage struct {
	Type      string `json:"type"` // "chat" | "system"
	UserID    string `json:"user_id,omitempty"`
	Username  string `json:"username,omitempty"`
	Message   string `json:"message"`
	Timestamp int64  `json:"timestamp"`
}
