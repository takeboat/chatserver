package client

import (
	"net"
	"tcpchat/logger"
	"tcpchat/message"
	"tcpchat/model"
)

type TCPClient struct {
	reader    model.MessageReader
	writer    model.MessageWriter
	Name      string
	conn      net.Conn
	log       *logger.Logger
	onMessage func(message *model.Message)
}
type ClientOptions func(*TCPClient)

func (c *TCPClient) WithOnMessage(onMessage func(message *model.Message)) ClientOptions {
	return func(c *TCPClient) {
		c.onMessage = onMessage
	}
}
func NewTCPClient(opts ...ClientOptions) model.Client {
	c := &TCPClient{
		log: logger.NewLogger(logger.WithGroup("tcp_client")),
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
	c.reader = message.NewJsonMessageReader(conn)
	c.writer = message.NewJsonMessageWriter(conn)
	c.log.Info("dialed", "address", address)
	return nil
}

func (c *TCPClient) Setname(name string) error {
	message := model.Message{
		Content: name,
		Type:    model.SetNameMessage,
	}
	c.log.Info("set name", "name", name)
	return c.writer.WriteMessage(&message)
}
func (c *TCPClient) SendMessage(content string) error {
	message := model.Message{
		Content: content,
		Type:    model.ChatMessage,
	}
	c.log.Info("send message", "message", content)
	return c.writer.WriteMessage(&message)
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
	
	// 确保 onMessage 回调函数不为 nil
	if c.onMessage == nil {
		c.log.Error("onMessage callback is nil")
		return
	}

	c.log.Info("开启客户端消息")
	defer c.log.Info("客户端消息结束")
	
	// 启动读取消息的 goroutine
	go func() {
		for {
			message, err := c.reader.ReadMessage()
			if err != nil {
				c.log.Error("读取消息错误", "error", err)
				return
			}
			c.onMessage(message)
		}
	}()
}