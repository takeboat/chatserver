package protocol

import (
	"encoding/json"
	"io"
	"tcpchat/model"
)

type MessageWriter struct {
	writer io.Writer
}

func NewMessageWriter(writer io.Writer) *MessageWriter {
	return &MessageWriter{writer: writer}
}
func (w *MessageWriter) WriteMessage(message *model.Message) error {
	data, err := json.Marshal(message)
	if err != nil {
		return err
	}
	data = append(data, '\n')
	_, err = w.writer.Write(data)
	return err
}
