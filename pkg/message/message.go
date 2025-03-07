package message

import (
	"encoding/binary"
	"fmt"
	"io"
)

type messageID uint8

const (
	MsgChoke messageID = iota
	MsgUnchoke
	MsgInterested
	MsgNotInterested
	MsgHave
	MsgBitfield
	MsgRequest
	MsgPiece
	MsgCancel
)

type Message struct {
	ID      messageID
	Payload []byte
}

// id length is 1 byte
const idLen = 1

// message length is uint32 number, so its length is 4 bytes
const messageLen = 4

func (m *Message) Serialize() []byte {
	if m == nil {
		return make([]byte, 4)
	}
	length := uint32(len(m.Payload) + idLen)
	buf := make([]byte, length+messageLen)
	binary.BigEndian.PutUint32(buf[0:messageLen], length)
	buf[messageLen] = byte(m.ID)
	copy(buf[messageLen+idLen:], m.Payload)
	return buf
}

func Read(r io.Reader) (*Message, error) {
	var m Message
	BufLen := make([]byte, messageLen)
	_, err := io.ReadFull(r, BufLen)
	if err != nil {
		return nil, err
	}
	length := binary.BigEndian.Uint32(BufLen)
	if length == 0 {
		err := fmt.Errorf("message length is zero")
		return nil, err
	}
	buf := make([]byte, length)
	_, err = io.ReadFull(r, buf)
	if err != nil {
		return nil, err
	}
	m = Message{
		ID:      messageID(buf[0]),
		Payload: buf[1:],
	}
	return &m, nil
}
