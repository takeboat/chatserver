package client

import (
	"net"
	"tcpchat/logger"
	"tcpchat/message"
	"tcpchat/model"
	"time"
)

type TCPClient struct {
	reader  model.MessageReader
	writer  model.MessageWriter
	Name    string
	conn    net.Conn
	log     *logger.Logger
	msgChan chan *model.Message
}
type ClientOptions func(*TCPClient)

func NewTCPClient(opts ...ClientOptions) model.Client {
	c := &TCPClient{
		log:     logger.NewLogger(logger.WithGroup("tcp_client")),
		msgChan: make(chan *model.Message, 100),
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
	c.log.Info("开启客户端消息")
	defer c.log.Info("客户端消息结束")
	// wait for server to be ready
	time.Sleep(time.Second * 4)
	for {
		msg, err := c.reader.ReadMessage()
		if err != nil {
			c.log.Error("读取消息错误", "error", err)
			return
		}
		c.log.Info("收到消息", "message", msg)
		// todo handel message
		c.msgChan <- msg
		c.handelMessage(msg)
	}
}
func (c *TCPClient) handelMessage(msg *model.Message) {
	switch msg.Type {
	case model.ChatMessage:
		c.log.Info("收到聊天消息", "message", msg.Content)
	case model.SystemMessage:
		c.log.Info("收到系统消息", "message", msg.Content)
	default:
		c.log.Warn("收到未知类型消息", "type", msg.Type)
	}
}
