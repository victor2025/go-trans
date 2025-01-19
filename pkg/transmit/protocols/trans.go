package protocols

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"go-trans/utils"
	"io"
	"net"
)

const (
	magic                    = 151 // magic number
	headSize                 = 4   // size of head
	ByteType       TransType = 0
	StrType        TransType = 1
	NumType        TransType = 2
	NewFileType    TransType = 3
	EndType        TransType = 4
	DisconnectType TransType = 5
)

type TRANS struct {
	Head    *THead
	Content []byte
}

type THead struct {
	Magic     uint8
	Type      TransType
	TotalSize uint16
}

type TransType uint8

// ByteTransMsg 字节类型消息体
func ByteTransMsg(data []byte) *TRANS {
	// 配置数据
	size := uint16(headSize + len(data)) // head + content
	head := &THead{
		Magic:     magic,
		Type:      ByteType,
		TotalSize: size,
	}
	return &TRANS{
		Head:    head,
		Content: data,
	}
}

// StrTransMsg 字符串类型消息体
func StrTransMsg(data string) *TRANS {
	trans := ByteTransMsg([]byte(data))
	trans.Head.Type = StrType
	return trans
}

// NumTransMsg 数值类型消息体
func NumTransMsg(data int64) *TRANS {
	byteArr := make([]byte, 64)
	n := binary.PutVarint(byteArr, data)
	trans := ByteTransMsg(byteArr[:n])
	trans.Head.Type = NumType
	return trans
}

// EndTransMsg 终止消息体
func EndTransMsg(md5 []byte) *TRANS {
	trans := ByteTransMsg(md5)
	trans.Head.Type = EndType
	return trans
}

// EmptyBodyTransMsg 空消息体
func EmptyBodyTransMsg(transType TransType) *TRANS {
	trans := ByteTransMsg(make([]byte, 0))
	trans.Head.Type = transType
	return trans
}

// ReceiveNextTrans 接收下一个消息体
func ReceiveNextTrans(conn net.Conn) (*TRANS, error) {
	var err error
	// read head
	headBytes := make([]byte, headSize)
	_, err = io.ReadFull(conn, headBytes)
	utils.HandleError(err)
	if err != nil {
		return nil, err
	}
	head, err := parseHead(headBytes)
	utils.HandleError(err)

	// read content
	content := make([]byte, head.TotalSize-headSize)
	_, err = io.ReadFull(conn, content)
	utils.HandleError(err)
	return &TRANS{
		Head:    head,
		Content: content,
	}, nil
}

// 解析头部
func parseHead(bytes []byte) (*THead, error) {
	if len(bytes) != headSize {
		return nil, errors.New("invalid head byte data")
	}
	m := uint8(bytes[0])
	if m != magic {
		return nil, fmt.Errorf("invalid magic number, expect %d, got %d", magic, m)
	}
	t := uint8(bytes[1])
	size := utils.Bytes2Uint64(bytes[2:4])
	return &THead{
		Magic:     magic,
		Type:      TransType(t),
		TotalSize: uint16(size),
	}, nil
}

func (t *TRANS) Bytes() []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.BigEndian, t.Head)
	buf.Write(t.Content)
	return buf.Bytes()
}
