package server

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"os"

	"github.com/nallerooth/fshare/common"
)

func (s *Server) HandleConnection(c net.Conn) {
	fmt.Println("Got a connection from", c.RemoteAddr())
	msg := common.Message{}

	// TODO: Update size to match common.Message
	// int8 + uint64 + 32 bytes
	msgType := make([]byte, 41)
	buf := bytes.NewBuffer(msgType)

	numBytes, err := c.Read(msgType)

	// convert bytes into struct
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error in connection %s: %s", c.RemoteAddr(), err)
	} else {
		binary.Read(buf, binary.BigEndian, &msg)
		fmt.Printf("%d bytes received from %s: %+v\n", numBytes, c.RemoteAddr(), msg)
		s.processMessage(c, msg)
		c.Close()
	}
}

func (s *Server) processMessage(c net.Conn, msg common.Message) {
	var err error

	switch msg.Type {
	case common.List:
		if err = s.sendList(c); err != nil {
			fmt.Fprintln(os.Stderr, "processMessage:", err)
		}
	case common.File:
		fmt.Println("Process file")
	case common.DeleteFile:
		fmt.Println("process delete file")
	}

	quitMsg := common.Message{
		Type: common.Quit,
	}
	buf := &bytes.Buffer{}
	binary.Write(buf, binary.BigEndian, quitMsg)
	c.Write(buf.Bytes())

	c.Close()
}

func (s *Server) sendList(c net.Conn) error {
	payload := []byte(s.fileListFormatter(s.AvailableFiles(), false))

	msg := common.Message{
		Type:   common.Text,
		Length: uint64(len(payload)),
	}

	buf := bytes.Buffer{}
	err := binary.Write(&buf, binary.BigEndian, msg)
	if err != nil {
		return err
	}

	c.Write(buf.Bytes())
	c.Write(payload)

	fmt.Println("DONE")
	return nil
}
