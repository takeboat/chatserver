package message

import (
	"encoding/json"
	"io"
	"tcpchat/model"
)

type JsonMessageWriter struct {
	writer io.Writer
}

func NewJsonMessageWriter(writer io.Writer) *JsonMessageWriter {
	return &JsonMessageWriter{writer: writer}
}
func (w *JsonMessageWriter) WriteMessage(message *model.Message) error {
	data, err := json.Marshal(message)
	if err != nil {
		return err
	}
	data = append(data, '\n')
	_, err = w.writer.Write(data)
	return err
}
