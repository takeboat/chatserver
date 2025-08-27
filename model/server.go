package model

type Server interface {
	Listen(port string) error
	BroadCast(message *Message) error
	Close() error 
	Start() 
}
