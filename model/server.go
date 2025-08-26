package model

type Server interface {
	Listen(port string) error
	Broadcast(message string) error
	Close() error
}