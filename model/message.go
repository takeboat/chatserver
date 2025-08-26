package model

import "time"

const (
	Join_Message MessageType = iota
	Chat_Message
	SetName_Message
	Leave_Message
	Error_Message
	System_Message
	Unknown_Message
)

type MessageType int

type Message struct {
	Type      MessageType
	Content   string
	Timestamp time.Time
}

type MessageReader interface {
	ReadMessage() (*Message, error)
}
type MessageWriter interface {
	WriteMessage(message *Message) error
}
