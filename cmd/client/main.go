package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/nallerooth/fshare/client"
	"github.com/nallerooth/fshare/common"
)

var clientConfig = client.Config{
	Passphrase: "",
	RemoteURL:  "localhost",
	RemotePort: 32000,
}

// TODO: Add support for remote alias
func connectToRemote() (net.Conn, error) {
	connStr := fmt.Sprintf("%s:%d", clientConfig.RemoteURL, clientConfig.RemotePort)
	return net.Dial("tcp", connStr)
}

// TODO: Handle timeout of reads
func readBytes(c net.Conn, bytesToRead uint64) ([]byte, error) {
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

func readMessageType(c net.Conn) (common.MessageType, error) {
	resBuf, err := readBytes(c, 1)
	if err != nil {
		return 0, err
	}
	var val common.MessageType
	err = binary.Read(bytes.NewReader(resBuf), binary.BigEndian, &val)
	if err != nil {
		return 0, err
	}

	return val, nil
}

func readUint64(c net.Conn) (uint64, error) {
	resBuf, err := readBytes(c, 8)
	if err != nil {
		return 0, err
	}
	var val uint64
	err = binary.Read(bytes.NewReader(resBuf), binary.BigEndian, &val)
	if err != nil {
		return 0, err
	}

	return val, nil
}

// TODO: readBytes(c net.Conn) ([]byte, error) {}

func connectAndSend(msg *client.InternalClientMessage) error {
	conn, err := connectToRemote()
	if err != nil {
		return err
	}
	defer conn.Close()

	buf := &bytes.Buffer{}
	err = binary.Write(buf, binary.BigEndian, msg.Type)
	if err != nil {
		return fmt.Errorf("Error converting message to bytes: %s", err)
	}

	_, err = conn.Write(buf.Bytes())
	if err != nil {
		return fmt.Errorf("Error writing to connection: %s", err)
	}

	running := true
	for running {

		respType, err := readMessageType(conn)
		if err != nil {
			return err
		}

		switch respType {
		case common.Text:
			// Read length of response
			respLen, err := readUint64(conn)
			if err != nil {
				return err
			}
			if respLen > 0 {
				payload := make([]byte, respLen)
				_, err = conn.Read(payload)
				if err != nil {
					return fmt.Errorf("read common.Text: %s", err)
				}
				fmt.Println(string(payload))
			}

		case common.Quit:
			fmt.Printf("\nReceived end of transmission\n")
			running = false
		}
	}

	return nil
}

// getCommand returns the command section of `program <command> <target>`
// as a lower cased string
func getCommand() (string, error) {
	if len(os.Args) >= 2 {
		return strings.ToLower(os.Args[1]), nil
	}

	return "", fmt.Errorf("Invalid command")
}

// getTarget returns the target section of `program <command> <target>`
func getTarget() (string, error) {
	if len(os.Args) >= 3 {
		return os.Args[2], nil
	}

	return "", fmt.Errorf("Invalid target")
}

func parseCommand() (common.MessageType, string, error) {
	command, err := getCommand()
	if err != nil {
		return -1, "", err
	}

	switch command {
	case "upload":
		file, err := getTarget()
		if err != nil {
			return -1, "", err
		}
		return common.File, file, nil

	case "list":
		return common.List, "", nil

	case "delete":
		name, err := getTarget()
		if err != nil {
			return -1, "", err
		}
		return common.DeleteFile, name, nil

	case "search":
		name, err := getTarget()
		if err != nil {
			return -1, "", err
		}
		return common.Search, name, nil
	}

	return -1, "", fmt.Errorf("Invalid command '%s'", command)
}

func main() {
	command, target, err := parseCommand()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		// TODO: Print usage
	}

	msg := &client.InternalClientMessage{
		Type: command,
	}

	switch command {
	case common.List:
	case common.File:
		msg.LocalFilename = target
	case common.DeleteFile:
		msg.RemoteFilename = target
	case common.DeleteHash:
		msg.Sha256sum = target
	}

	err = connectAndSend(msg)
	if err != nil {
		fmt.Println(err)
	}
}
