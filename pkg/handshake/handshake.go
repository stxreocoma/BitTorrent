package handshake

import (
	"fmt"
	"io"
)

type Handshake struct {
	Pstr     string
	InfoHash [20]byte
	PeerID   [20]byte
}

const reservedBytes = 8

// infoHash length is 19 but we need the next index so it is 19+1
const infoHashLen = 20
const peerIDLen = 20

func (h *Handshake) Serialize() []byte {
	pstrlen := len(h.Pstr)
	bufLen := 49 + pstrlen
	buf := make([]byte, bufLen)
	buf[0] = byte(pstrlen)
	copy(buf[1:], h.Pstr)
	//Leave 8 reserved bytes
	copy(buf[1+pstrlen+reservedBytes:], h.InfoHash[:])
	copy(buf[1+pstrlen+reservedBytes+infoHashLen:], h.PeerID[:])
	return buf
}

func Read(r io.Reader) (*Handshake, error) {
	var h Handshake
	bufLen := make([]byte, 1)
	_, err := io.ReadFull(r, bufLen)
	if err != nil {
		return nil, err
	}
	pstrLen := int(bufLen[0])
	if pstrLen == 0 {
		err := fmt.Errorf("pstr length is zero")
		return nil, err
	}
	buf := make([]byte, pstrLen+reservedBytes+infoHashLen+peerIDLen)
	_, err = io.ReadFull(r, buf)
	if err != nil {
		return nil, err
	}
	var infoHash, peerID [20]byte
	copy(infoHash[:], buf[pstrLen+reservedBytes:pstrLen+reservedBytes+infoHashLen])
	copy(peerID[:], buf[pstrLen+reservedBytes+infoHashLen:])
	h = Handshake{
		Pstr:     string(buf[0:pstrLen]),
		InfoHash: infoHash,
		PeerID:   peerID,
	}
	return &h, nil
}

func New(infoHash, peerID [20]byte) *Handshake {
	return &Handshake{
		Pstr:     "BitTorrent protocol",
		InfoHash: infoHash,
		PeerID:   peerID,
	}
}
