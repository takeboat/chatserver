package server

import (
	"context"
)

type Server interface {
	Listen(address string) error
	Start(ctx context.Context)
}
