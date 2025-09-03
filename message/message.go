package message

import (
	"io"
	"time"
)

const (
	JoinMessage MessageType = iota // 加入群聊
	ChatMessage                    // 聊天消息
	SetNameMessage
	LeaveMessage
	ErrorMessage
	SystemMessage
	UnknownMessage
)

type MessageType int

type Message struct {
	Type      MessageType
	Owner     string
	Content   string
	Timestamp time.Time
}

type MessageReader interface {
	ReadMessage(r io.Reader) (*Message, error)
}

type MessageWriter interface {
	WriteMessage(writer io.Writer, message *Message) error
}
