package protocol

import (
	"bufio"
	"encoding/json"
	"io"
	"tcpchat/model"
)

type MessageReader struct {
	reader *bufio.Reader
}

func NewMessageReader(reader io.Reader) *MessageReader {
	return &MessageReader{
		reader: bufio.NewReader(reader),
	}
}

func (r *MessageReader) ReadMessage() (*model.Message, error) {
	line, err := r.reader.ReadBytes('\n')	
	if err != nil {
		return nil, err
	}
	var message model.Message
	err = json.Unmarshal(line, &message)
	if err != nil {
		return nil, err
	}
	return &message, nil
}