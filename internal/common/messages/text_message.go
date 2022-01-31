package message

import (
	"bytes"
	"encoding/binary"
	"net"

	"github.com/nallerooth/fshare/internal/common"
)

// TextMessage is a small header describing a range of text bytes
type TextMessage struct {
	Length uint64
}

// ReadFromConn will read the header Length number of bytes from the connection
func (msg *TextMessage) ReadFromConn(c net.Conn) error {
	msgLen := uint64(8)

	buf, err := common.ReadBytes(c, msgLen)
	if err != nil {
		return err
	}

	err = binary.Read(bytes.NewBuffer(buf), binary.BigEndian, msg)
	if err != nil {
		return err
	}

	return nil
}
