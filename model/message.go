package model

import "time"

const (
	JoinMessage MessageType = iota
	ChatMessage
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
	ReadMessage() (*Message, error)
}

type MessageWriter interface {
	WriteMessage(message *Message) error
}
