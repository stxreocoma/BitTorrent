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
		return nil, nil
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

func FormatRequest(index, begin, length int) *Message {
	payload := make([]byte, 12)
	binary.BigEndian.PutUint32(payload[0:4], uint32(index))
	binary.BigEndian.PutUint32(payload[4:8], uint32(begin))
	binary.BigEndian.PutUint32(payload[8:12], uint32(length))
	return &Message{ID: MsgRequest, Payload: payload}
}

func FormatHave(index int) *Message {
	payload := make([]byte, 4)
	binary.BigEndian.PutUint32(payload[0:4], uint32(index))
	return &Message{ID: MsgHave, Payload: payload}
}
func ParsePiece(index int, buf []byte, msg *Message) (int, error) {
	if msg.ID != MsgPiece {
		return 0, fmt.Errorf("expected PIECE (ID %d), got ID %d", MsgPiece, msg.ID)
	}
	if len(msg.Payload) < 8 {
		return 0, fmt.Errorf("payload too short. %d < 8", len(msg.Payload))
	}
	parsedIndex := int(binary.BigEndian.Uint32(msg.Payload[0:4]))
	if parsedIndex != index {
		return 0, fmt.Errorf("expected index of %d, got %d", index, parsedIndex)
	}
	begin := int(binary.BigEndian.Uint32(msg.Payload[4:8]))
	if begin >= len(buf) {
		return 0, fmt.Errorf("begin offset too high. %d >= %d", begin, len(buf))
	}
	data := msg.Payload[8:]
	if begin+len(data) > len(buf) {
		return 0, fmt.Errorf("data too long [%d] for offset %d with length %d", len(data), begin, len(buf))
	}
	copy(buf[begin:], data)
	return len(data), nil
}

func ParseHave(msg *Message) (int, error) {
	if msg.ID != MsgHave {
		return 0, fmt.Errorf("expected HAVE (ID %d), got ID %d", MsgHave, msg.ID)
	}
	if len(msg.Payload) != 4 {
		return 0, fmt.Errorf("expected payload length 4, got length %d", len(msg.Payload))
	}
	index := int(binary.BigEndian.Uint32(msg.Payload))
	return index, nil
}

func (m *Message) name() string {
	if m == nil {
		return "KeepAlive"
	}
	switch m.ID {
	case MsgChoke:
		return "Choke"
	case MsgUnchoke:
		return "Unchoke"
	case MsgInterested:
		return "Interested"
	case MsgNotInterested:
		return "NotInterested"
	case MsgHave:
		return "Have"
	case MsgBitfield:
		return "Bitfield"
	case MsgRequest:
		return "Request"
	case MsgPiece:
		return "Piece"
	case MsgCancel:
		return "Cancel"
	default:
		return fmt.Sprintf("Unknown#%d", m.ID)
	}
}

func (m *Message) String() string {
	if m == nil {
		return m.name()
	}
	return fmt.Sprintf("%s [%d]", m.name(), len(m.Payload))
}
