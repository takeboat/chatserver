package client

type Client interface {
	Dial(address string) error
	Setname(name string) error
	SendMessage(message string) error
	ReadMessage() (string, error)
	Close() error
}