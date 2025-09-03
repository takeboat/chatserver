package message

import (
	"bufio"
	"encoding/json"
	"io"
)

type JsonMessageReader struct {
	reader *bufio.Reader
}

func NewJsonMessageReader() *JsonMessageReader {
	return &JsonMessageReader{}
}

func (jmr *JsonMessageReader) ReadMessage(r io.Reader) (*Message, error) {
	jmr.reader = bufio.NewReader(r)

	line, err := jmr.reader.ReadBytes('\n')
	if err != nil {
		return nil, err
	}
	var message Message
	err = json.Unmarshal(line, &message)
	if err != nil {
		return nil, err
	}
	return &message, nil
}

// TLV: Type-Length-Value
// 这是一种编码格式，经常用于通信协议之中
// T: Type, 1 or 4字节，表示消息类型
// L: Length, 表示消息长度 通常为固定长度为4字节 or 2字节
// V: Value, 表示消息内容 长度由Length来决定
type TLVMessageReader struct {
	reader *bufio.Reader
}

func NewTLVMessageReader(reader io.Reader) *TLVMessageReader {
	return &TLVMessageReader{
		reader: bufio.NewReader(reader),
	}
}

// func (tlv *TLVMessageReader) ReadMessage() (*model.Message, error) {
// 	// 读取Type 假设1 字节
// 	typeByte, err := tlv.reader.ReadByte()
// 	if err != nil {
// 		return nil, err
// 	}
// 	// 读取Lenght 假设4字节 大端
// 	lengthBytes := make([]byte, 4)
// 	_, err = io.ReadFull(tlv.reader, lengthBytes)
// 	if err != nil {
// 		return nil, err
// 	}
// 	// 读取Value
// 	var message model.Message
// 	// 根据 Type来解析 Value
// 	return &message, nil
// }
