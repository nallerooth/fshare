package server

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"os"

	"github.com/nallerooth/fshare/internal/connection"
	"github.com/nallerooth/fshare/internal/message"
)

// HandleConnection is the entry point for a newly established TCP connection
func (s *Server) HandleConnection(c net.Conn) {
	defer c.Close()
	fmt.Println("Handling connection from", c.RemoteAddr()) // TODO: move to logger
	msg := message.Message{}
	err := connection.ReadMessage(c, &msg)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error reading message from client: ", err)
	}

	err = s.processMessage(c, msg)
	if err != nil {
		fmt.Println("Error: ", err)
	}

	fmt.Println("closing conn..") // TODO: remove
}

func (s *Server) processMessage(c net.Conn, msg message.Message) error {
	var err error

	switch msg.Command {
	case message.List:
		fmt.Println("list files")
		if err = s.sendList(c); err != nil {
			return fmt.Errorf("processMessage: %s", err)
		}
	}
	//case message.FileTransfer:
	//fmt.Println("Process file")
	//if err = s.receiveFile(c); err != nil {
	//return fmt.Errorf("processMessage: %s", err)
	//}
	//case message.FileDelete:
	//fmt.Println("process delete file")
	//}

	// Message processing done, send quit command to client
	//c.Write([]byte{byte(connection.Quit)})
	err = s.sendQuit(c)
	if err != nil {
		return err
	}

	return nil
}

func (s *Server) sendQuit(c net.Conn) error {
	msg := message.QuitMessage()
	buf := bytes.Buffer{}
	err := binary.Write(&buf, binary.BigEndian, msg)
	if err != nil {
		return err
	}

	_, err = c.Write(buf.Bytes())
	if err != nil {
		return err
	}

	return nil
}

func (s *Server) sendList(c net.Conn) error {
	payload := []byte(s.fileListFormatter(s.AvailableFiles(), false))

	msg := message.NewText(payload)
	buf := bytes.Buffer{}

	err := binary.Write(&buf, binary.BigEndian, msg)
	if err != nil {
		return fmt.Errorf("error writing msg: %v", err)
	}
	err = binary.Write(&buf, binary.BigEndian, payload)
	if err != nil {
		return fmt.Errorf("error writing payload: %v", err)
	}

	_, err = c.Write(buf.Bytes())
	if err != nil {
		return fmt.Errorf("error writing buffer to socket: %v", err)
	}
	if err != nil {
		return err
	}

	return nil
}

//func (s *Server) receiveFile(c net.Conn) error {
//msg := message.Message{}
//err := msg.ReadFromConn(c)
//if err != nil {
//return err
//}
//_, filename := filepath.Split(string(msg.Filename[:]))

//// Temporary storage
//tmpFile, err := ioutil.TempFile(s.config.Server.Workdir, "upload_*.tmp")
//if err != nil {
//return fmt.Errorf("receiveFile: %s", err)
//}
//bytesToRead := msg.Length
//buf := make([]byte, 0, 1024)

//for bytesToRead > 0 {
//bytesRead, err := c.Read(buf)
//if err != nil {
//return fmt.Errorf("receiveFile: %s", err)
//}
//tmpBytesWritten, err := tmpFile.Write(buf[:bytesRead])
//if tmpBytesWritten != bytesRead {
//return fmt.Errorf("receiveFile: Bytes written to tmp file (%d) did not match number of bytes read from socket (%d)", tmpBytesWritten, bytesRead)
//}

//bytesToRead -= uint64(bytesRead)
//}

//err = os.Rename(tmpFile.Name(), s.config.Server.Workdir+filename)
//if err != nil {
//return fmt.Errorf("receiveFile: unable to move temporary file to given filename: %s", err)
//}

//return nil
//}
