package message

import (
	"bufio"
	"encoding/json"
	"io"
	"tcpchat/model"
)

type JsonMessageReader struct {
	reader *bufio.Reader
}

func NewMessageReader(reader io.Reader) *JsonMessageReader {
	return &JsonMessageReader{
		reader: bufio.NewReader(reader),
	}
}

func (r *JsonMessageReader) ReadMessage() (*model.Message, error) {
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

type TLVMessageReader struct {
}
