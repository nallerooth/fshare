package main

import (
	"bytes"
	"encoding/binary"
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

func readMessageType(c net.Conn) (common.MessageType, error) {
	resBuf, err := common.ReadBytes(c, 1)
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
	resBuf, err := common.ReadBytes(c, 8)
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
		return common.FileTransfer, file, nil

	case "list":
		return common.List, "", nil

	case "delete":
		name, err := getTarget()
		if err != nil {
			return -1, "", err
		}
		return common.FileDelete, name, nil

	case "search":
		name, err := getTarget()
		if err != nil {
			return -1, "", err
		}
		return common.FileSearch, name, nil
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
	case common.FileTransfer:
		msg.LocalFilename = target
	case common.FileDelete:
		msg.RemoteFilename = target
	}

	err = connectAndSend(msg)
	if err != nil {
		fmt.Println(err)
	}
}
