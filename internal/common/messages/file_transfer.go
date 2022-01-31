package message

import (
	"bytes"
	"encoding/binary"
	"net"

	"github.com/nallerooth/fshare/internal/common"
)

type FileTransferMessage struct {
	Length   uint64
	Filename [64]byte
	Hash     [64]byte
}

func (msg *FileTransferMessage) ReadFromConn(c net.Conn) error {
	msgLen := uint64(8 + 64 + 64)

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
