package model


type Client interface {
	Dial(address string) error
	Setname(name string) error
	SendMessage(message string) error
	ReadMessage() (Message, error)
	Close() error
}
