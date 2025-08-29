package server

import (
	"context"
	"fmt"
	"io"
	"net"
	"sync"
	"tcpchat/logger"
	"tcpchat/message"
	"tcpchat/model"
)

type TCPServer struct {
	clients   map[net.Conn]*ClinetInfo
	clientsMu sync.Mutex
	listener  net.Listener
	log       *logger.Logger
}

type ClinetInfo struct {
	Name   string
	Conn   net.Conn
	writer model.MessageWriter
}

func NewTCPServer() model.Server {
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
	s.log.Info("服务器监听于", "address", address)
	return nil
}

func (s *TCPServer) Start(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			conn, err := s.listener.Accept()
			if err != nil {
				s.log.Error("监听错误", "error", err)
				continue
			}
			s.appendConn(conn)
			go s.serve(conn)
		}
	}
}
func (s *TCPServer) Close() error {
	if s.listener != nil {
		return s.listener.Close()
	}
	return nil
}

func (s *TCPServer) BroadCast(m *model.Message) error {
	s.clientsMu.Lock()
	defer s.clientsMu.Unlock()
	for _, client := range s.clients {
		err := client.writer.WriteMessage(m)
		if err != nil {
			s.log.Error("广播消息错误", "error", err)
		}
	}
	return nil
}

func (s *TCPServer) appendConn(conn net.Conn) {
	cInfo := &ClinetInfo{
		Name:   "",
		Conn:   conn,
		writer: message.NewJsonMessageWriter(conn),
	}
	s.clientsMu.Lock()
	// 这里clientinfo
	s.clients[conn] = cInfo
	s.clientsMu.Unlock()
	s.log.Info("新连接", "remote", conn.RemoteAddr().String())
}
func (s *TCPServer) removeConn(conn net.Conn) {
	s.clientsMu.Lock()
	delete(s.clients, conn)
	s.clientsMu.Unlock()
	conn.Close()
	s.log.Info("连接断开", "remote", conn.RemoteAddr().String())
}

func (s *TCPServer) serve(conn net.Conn) {
	defer s.removeConn(conn)
	reader := message.NewJsonMessageReader(conn)
	for {
		msg, err := reader.ReadMessage()
		if err != nil && err != io.EOF {
			s.log.Error("读取消息错误", "error", err)
			return
		}
		if err == io.EOF {
			leave := &model.Message{
				Type:    model.LeaveMessage,
				Content: fmt.Sprintf("%s 离开了聊天室", s.clients[conn].Name),
			}
			s.BroadCast(leave)
			return
		}
		s.log.Info("收到消息", "remote", conn.RemoteAddr().String(), "message", msg)
		s.hanldedMessage(conn, msg)
	}
}
func (s *TCPServer) setClientName(conn net.Conn, m *model.Message) {
	s.clientsMu.Lock()
	if client, ok := s.clients[conn]; ok {
		client.Name = m.Content
	}
	s.clientsMu.Unlock()
	changeNameMsg := &model.Message{
		Type:    model.SystemMessage,
		Content: fmt.Sprintf("%s 修改了昵称", s.clients[conn].Name),
	}
	s.BroadCast(changeNameMsg)
}

func (s *TCPServer) hanldedMessage(conn net.Conn, m *model.Message) {
	switch m.Type {
	case model.SetNameMessage:
		s.setClientName(conn, m)
	case model.ChatMessage:
		s.BroadCast(m)
	case model.JoinMessage:
		join := &model.Message{
			Type:    model.SystemMessage,
			Content: fmt.Sprintf("%s 加入了聊天室", m.Content),
		}
		s.BroadCast(join)
	default:
		s.log.Warn("未知消息类型", "type", m.Type)
	}
}
