package model


type Client interface {
	Dial(address string) error
	Setname(name string) error
	SendMessage(content string) error
	Close() error
}
