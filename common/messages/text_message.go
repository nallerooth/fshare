package message

import (
	"bytes"
	"encoding/binary"
	"net"

	"github.com/nallerooth/fshare/common"
)

type TextMessage struct {
	Length uint64
}

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
