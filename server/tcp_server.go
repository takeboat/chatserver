package server

import (
	"log/slog"
	"net"
	"sync"
	"tcpchat/logger"
	"tcpchat/message"
	"tcpchat/model"
)

type TCPServer struct {
	clients   map[net.Conn]*ClinetInfo
	clientsMu sync.RWMutex
	listener  net.Listener
	log       *logger.Logger
}

type ClinetInfo struct {
	Name   string
	Conn   net.Conn
	writer model.MessageWriter
}

func NewTCPServer() *TCPServer {
	s := &TCPServer{
		clients: make(map[net.Conn]*ClinetInfo),
		log:     logger.NewLogger(logger.WithGroup("tcp_server")),
	}
	return s
}
func (s *TCPServer) Listen(address string) error {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	s.listener = listener
	slog.Info("服务器监听于地址", address)
	// todo 这里添加东西
	return nil
}

func (s *TCPServer) Start() {

}
func (s *TCPServer) Close() error {
	if s.listener != nil {
		return s.listener.Close()
	}
	return nil
}

func (s *TCPServer) BroadCast(m model.Message) error {

	return nil
}
func (s *TCPServer) acceptConn() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			continue
		}
		s.handleConn(conn)
	}
}
func (s *TCPServer) handleConn(conn net.Conn) {
	cInfo := &ClinetInfo{
		Name:   "",
		Conn:   conn,
		writer: message.NewJsonMessageWriter(conn),
	}
	s.clientsMu.Lock()
	// 这里clientinfo
	s.clients[conn] = cInfo
	s.clientsMu.Unlock()
}
func (s *TCPServer) removeConn(conn net.Conn) {
	s.clientsMu.Lock()
	delete(s.clients, conn)
	s.clientsMu.Unlock()
	conn.Close()
}

func (s *TCPServer) serve() {

}
