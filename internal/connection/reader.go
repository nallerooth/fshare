package connection

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/nallerooth/fshare/internal/message"
)

// ReadBytes reads N bytes from a connection
// TODO: Handle timeout of reads
func ReadBytes(c net.Conn, bytesToRead uint64) ([]byte, error) {
	if bytesToRead == 0 {
		return nil, errors.New("attempted to read 0 bytes")
	}
	fmt.Printf("attempting to read %d bytes\n", bytesToRead)
	fmt.Printf("local: %v\nremote: %v", c.LocalAddr(), c.RemoteAddr())
	buf := make([]byte, 0, bytesToRead)
	bRead := uint64(0)

	c.SetReadDeadline(time.Now())
	for bRead < bytesToRead {
		n, err := c.Read(buf)
		if err != nil {
			return nil, err
		}
		bRead += uint64(n)
	}

	return buf, nil
}

// ReadMessage reads the sizeof Message from a connection
func ReadMessage(c net.Conn, msg *message.Message) error {
	err := binary.Read(c, binary.BigEndian, msg)
	if err != nil {
		if errors.Is(err, io.EOF) {
			return fmt.Errorf("could not read from socket: %w", err)
		}
		if errors.Is(err, io.ErrUnexpectedEOF) {
			return fmt.Errorf("could not read full message: %w", err)
		}
		return fmt.Errorf("unknown error while reading message: %w", err)
	}

	return nil
}
