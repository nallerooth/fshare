package common

import (
	"errors"
	"net"
)

// TODO: Handle timeout of reads
func ReadBytes(c net.Conn, bytesToRead uint64) ([]byte, error) {
	if bytesToRead == 0 {
		return nil, errors.New("attempted to read 0 bytes")
	}
	buf := make([]byte, 0, bytesToRead)
	bRead := uint64(0)

	for bRead < bytesToRead {
		n, err := c.Read(buf)
		if err != nil {
			return nil, err
		}
		bRead += uint64(n)
	}

	return buf, nil
}
