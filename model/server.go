package model

import (
	"context"
)

type Server interface {
	Listen(port string) error
	BroadCast(message *Message) error
	Close() error
	Start(ctx context.Context)
}
