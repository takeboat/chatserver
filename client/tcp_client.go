package client

import (
	"fmt"
	"net"
	"tcpchat/logger"
	"tcpchat/message"
	"time"
)

type TCPClient struct {
	reader    message.MessageReader
	writer    message.MessageWriter
	Name      string
	conn      net.Conn
	log       *logger.Logger
	onMessage func(message *message.Message)
}

type ClientOptions func(*TCPClient)

func (c *TCPClient) WithOnMessage(onMessage func(message *message.Message)) ClientOptions {
	return func(c *TCPClient) {
		c.onMessage = onMessage
	}
}

func NewTCPClient(opts ...ClientOptions) Client {
	c := &TCPClient{
		log:  logger.NewLogger(logger.WithGroup("tcp_client")),
		Name: fmt.Sprintf("client_%d", time.Now().UnixNano()),
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func (c *TCPClient) Dial(address string) error {
	c.log.Info("dialing", "address", address)
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return err
	}
	c.conn = conn
	c.reader = message.NewJsonMessageReader()
	c.writer = message.NewJsonMessageWriter()
	c.log.Info("dialed", "address", address)
	return nil
}

func (c *TCPClient) Setname(name string) error {
	message := message.Message{
		Content: name,
		Owner:   c.Name,
		Type:    message.SetNameMessage,
	}
	c.log.Info("set name", "name", name)
	c.Name = name
	return c.writer.WriteMessage(c.conn, &message)
}

func (c *TCPClient) SendMessage(content string) error {
	message := message.Message{
		Content: content,
		Owner:   c.Name,
		Type:    message.ChatMessage,
	}
	c.log.Info("send message", "message", content)
	return c.writer.WriteMessage(c.conn, &message)
}

func (c *TCPClient) Close() error {
	return c.conn.Close()
}

func (c *TCPClient) Start() {
	// 确保 reader 和 writer 不为 nil
	if c.reader == nil || c.writer == nil {
		c.log.Error("reader or writer is nil")
		return // 防止 panic
	}

	c.log.Info("开启客户端消息")
	defer c.log.Info("客户端消息结束")

	// 启动读取消息的 goroutine
	go func() {
		for {
			message, err := c.reader.ReadMessage(c.conn)
			if err != nil {
				c.log.Error("读取消息错误", "error", err)
				return
			}
			if c.onMessage != nil {
				c.onMessage(message)
			}
		}
	}()
}
