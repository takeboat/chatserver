package client

import (
	"context"
	"net"
	"tcpchat/model"
)

type TCPClient struct {
	reader model.MessageReader
	writer model.MessageWriter
	Name   string
	conn   net.Conn
}
type ClientOptions func(*TCPClient)

func WithReader(reader model.MessageReader) ClientOptions {
	return func(c *TCPClient) {
		c.reader = reader
	}
}
func WithWriter(writer model.MessageWriter) ClientOptions {
	return func(c *TCPClient) {
		c.writer = writer
	}
}
func NewTCPClient(ctx context.Context, opts ...ClientOptions) *TCPClient {
	c := &TCPClient{}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func (c *TCPClient) Dial(address string) error {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return err
	}
	c.conn = conn
	return nil
}

func (c *TCPClient) Setname(name string) error {
	message := model.Message{
		Content: name,
		Type:    model.SetName_Message,
	}
	return c.writer.WriteMessage(&message)
}
func (c *TCPClient) SendMessage(content string) error {
	message := model.Message{
		Content: content,
		Type: model.Chat_Message,
	}
	return c.writer.WriteMessage(&message)
}
func (c *TCPClient) ReadMessage() (*model.Message, error) {
	return c.reader.ReadMessage()
}

func (c *TCPClient) Close() error {
	return c.conn.Close()
}