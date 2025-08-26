package server

import (
	"net"
	"sync"
	"tcpchat/model"
)

type TCPServer struct {
	clients   map[net.Conn]*ClinetInfo
	clientsMu sync.RWMutex
	listener  net.Listener
	writer    model.MessageWriter
}

type ClinetInfo struct {
	Name string
	Conn net.Conn
}

func NewTCPServer() *TCPServer {
	s := &TCPServer{
		clients: make(map[net.Conn]*ClinetInfo),
	}
	return s
}
func (s *TCPServer) Listen(address string) error {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	s.listener = listener
	// todo
	return nil
}

func (s *TCPServer) Close() error {
	if s.listener != nil {
		return s.listener.Close()
	}
	return nil
}

func (s *TCPServer) acceptConn() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			continue
		}
		go s.handleConn(conn)

	}
}
func (s *TCPServer) handleConn(conn net.Conn) {
}
