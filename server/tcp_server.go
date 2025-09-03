package server

import (
	"context"
	"fmt"
	"io"
	"net"
	"sync"
	"sync/atomic"
	"tcpchat/logger"
	"tcpchat/message"
)

// 使用原子操作代替全局变量
var globalID int64

type TCPServer struct {
	clients   map[int64]*ClientInfo
	clientsMu sync.RWMutex
	listener  net.Listener
	log       *logger.Logger
	writer    message.MessageWriter
}

type ClientInfo struct {
	ID   int64
	Name string
	Conn net.Conn
}

func NewTCPServer() Server {
	s := &TCPServer{
		clients: make(map[int64]*ClientInfo),
		log:     logger.NewLogger(logger.WithGroup("tcp_server")),
		writer:  &message.JsonMessageWriter{}, // 初始化writer
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
	defer s.Close()
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

func (s *TCPServer) Broadcast(m *message.Message) error {
	s.clientsMu.RLock()
	defer s.clientsMu.RUnlock()

	for _, client := range s.clients {
		err := s.writer.WriteMessage(client.Conn, m)
		if err != nil {
			s.log.Error("广播消息错误", "error", err, "client", client.Conn.RemoteAddr().String())
		}
	}
	return nil
}

func (s *TCPServer) appendConn(conn net.Conn) int64 {
	// 使用原子操作生成唯一ID
	id := atomic.AddInt64(&globalID, 1)

	cInfo := &ClientInfo{
		ID:   id,
		Name: fmt.Sprintf("用户%d", id),
		Conn: conn,
	}

	s.clientsMu.Lock()
	s.clients[id] = cInfo
	s.clientsMu.Unlock()

	s.log.Info("新连接", "remote", conn.RemoteAddr().String(), "id", id)
	return id
}

func (s *TCPServer) removeConn(id int64) {
	s.clientsMu.Lock()
	clientInfo, exists := s.clients[id]
	if exists {
		delete(s.clients, id)
	}
	s.clientsMu.Unlock()

	if exists {
		clientInfo.Conn.Close()
		s.log.Info("连接断开", "remote", clientInfo.Conn.RemoteAddr().String(), "id", id)
	}
}

func (s *TCPServer) getClientIDByConn(conn net.Conn) (int64, bool) {
	s.clientsMu.RLock()
	defer s.clientsMu.RUnlock()

	for id, client := range s.clients {
		if client.Conn == conn {
			return id, true
		}
	}
	return 0, false
}

func (s *TCPServer) serve(conn net.Conn) {
	// 获取客户端ID
	clientID, exists := s.getClientIDByConn(conn)
	if !exists {
		s.log.Error("无法找到客户端ID", "remote", conn.RemoteAddr().String())
		conn.Close()
		return
	}

	defer func() {
		s.removeConn(clientID)
	}()

	reader := message.NewJsonMessageReader()
	for {
		msg, err := reader.ReadMessage(conn)
		if err != nil {
			if err != io.EOF {
				s.log.Error("读取消息错误", "error", err, "client_id", clientID)
			} else {
				// 客户端正常断开连接
				s.clientsMu.RLock()
				clientName := s.clients[clientID].Name
				s.clientsMu.RUnlock()

				leave := &message.Message{
					Type:    message.LeaveMessage,
					Owner:   clientName,
					Content: fmt.Sprintf("%s 离开了聊天室", clientName),
				}
				s.Broadcast(leave)
			}
			return
		}

		s.log.Info("收到消息", "client_id", clientID, "message", msg)
		s.handleMessage(clientID, msg)
	}
}

func (s *TCPServer) setClientName(clientID int64, m *message.Message) {
	s.clientsMu.Lock()
	if client, ok := s.clients[clientID]; ok {
		oldName := client.Name
		client.Name = m.Content
		s.clientsMu.Unlock()

		changeNameMsg := &message.Message{
			Type:    message.SystemMessage,
			Content: fmt.Sprintf("%s 修改了昵称为 %s", oldName, m.Content),
		}
		s.Broadcast(changeNameMsg)
	} else {
		s.clientsMu.Unlock()
	}
}

func (s *TCPServer) handleMessage(clientID int64, m *message.Message) {
	switch m.Type {
	case message.SetNameMessage:
		s.setClientName(clientID, m)
	case message.ChatMessage:
		// 设置消息发送者
		s.clientsMu.RLock()
		if client, ok := s.clients[clientID]; ok {
			m.Owner = client.Name
		}
		s.clientsMu.RUnlock()
		s.Broadcast(m)
	case message.JoinMessage:
		s.clientsMu.RLock()
		clientName := s.clients[clientID].Name
		s.clientsMu.RUnlock()

		join := &message.Message{
			Type:    message.SystemMessage,
			Content: fmt.Sprintf("%s 加入了聊天室", clientName),
		}
		s.Broadcast(join)
	default:
		s.log.Warn("未知消息类型", "type", m.Type)
	}
}
