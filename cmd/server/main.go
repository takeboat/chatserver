package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"tcpchat/server"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	port := flag.String("port", "8080", "server port")
	flag.Parse()
	s := server.NewTCPServer()
	slog.Info("启动服务", "port", *port)
	go func() {
		s.Listen(fmt.Sprintf(":%s", *port))
		s.Start(ctx)
	}()

	<-interrupt
	cancel()
	slog.Info("服务已关闭")
}
