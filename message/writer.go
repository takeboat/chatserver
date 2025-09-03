package message

import (
	"bufio"
	"encoding/json"
	"io"
)

type JsonMessageWriter struct {
	writer io.Writer
}

func NewJsonMessageWriter() *JsonMessageWriter {
	return &JsonMessageWriter{}
}
func (w *JsonMessageWriter) WriteMessage(writer io.Writer, message *Message) error {
	w.writer = bufio.NewWriter(writer)

	data, err := json.Marshal(message)
	if err != nil {
		return err
	}
	data = append(data, '\n')
	_, err = w.writer.Write(data)
	return err
}
