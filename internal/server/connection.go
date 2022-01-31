package server

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"

	"github.com/nallerooth/fshare/internal/common"
	message "github.com/nallerooth/fshare/internal/common/messages"
)

// HandleConnection is the entry point for a newly established TCP connection
func (s *Server) HandleConnection(c net.Conn) {
	fmt.Println("Got a connection from", c.RemoteAddr())
	msg := common.Message{}
	err := common.ReadMessage(c, &msg)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error reading message from client: ", err)
	}
	fmt.Printf("Got a message of type: %+v\n", msg)
	s.processMessage(c, msg)
	c.Close()
	/*
		var msg common.MessageType

		// TODO: Update size to match common.Message
		// int8 + uint64 + [32 bytes]filename
		msgType := make([]byte, 41)
		buf := bytes.NewBuffer(msgType)

		numBytes, err := c.Read(msgType)

		// Convert bytes into struct
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error in connection %s: %s", c.RemoteAddr(), err)
		} else {
			binary.Read(buf, binary.BigEndian, &msg)
			fmt.Printf("%d bytes received from %s: %+v\n", numBytes, c.RemoteAddr(), msg)
			s.processMessage(c, msg)
			c.Close()
		}
	*/
}

func (s *Server) processMessage(c net.Conn, msg common.Message) error {
	var err error

	switch msg.Command {
	case common.List:
		if err = s.sendList(c); err != nil {
			return fmt.Errorf("processMessage: %s", err)
		}
	case common.FileTransfer:
		fmt.Println("Process file")
		if err := s.receiveFile(c); err != nil {
			return fmt.Errorf("processMessage: %s", err)
		}
	case common.FileDelete:
		fmt.Println("process delete file")
	}

	// Message processing done, send quit command to client
	c.Write([]byte{byte(common.Quit)})

	c.Close()

	return nil
}

func (s *Server) sendList(c net.Conn) error {
	payload := []byte(s.fileListFormatter(s.AvailableFiles(), false))

	msg := message.TextMessage{
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

func (s *Server) receiveFile(c net.Conn) error {
	msg := message.FileTransferMessage{}
	err := msg.ReadFromConn(c)
	if err != nil {
		return err
	}
	_, filename := filepath.Split(string(msg.Filename[:]))

	// Temporary storage
	tmpFile, err := ioutil.TempFile(s.config.Workdir, "upload_*.tmp")
	if err != nil {
		return fmt.Errorf("receiveFile: %s", err)
	}
	bytesToRead := msg.Length
	buf := make([]byte, 0, 1024)

	for bytesToRead > 0 {
		bytesRead, err := c.Read(buf)
		if err != nil {
			return fmt.Errorf("receiveFile: %s", err)
		}
		tmpBytesWritten, err := tmpFile.Write(buf[:bytesRead])
		if tmpBytesWritten != bytesRead {
			return fmt.Errorf("receiveFile: Bytes written to tmp file (%d) did not match number of bytes read from socket (%d)", tmpBytesWritten, bytesRead)
		}

		bytesToRead -= uint64(bytesRead)
	}

	err = os.Rename(tmpFile.Name(), s.config.Workdir+filename)
	if err != nil {
		return fmt.Errorf("receiveFile: unable to move temporary file to given filename: %s", err)
	}

	return nil
}
