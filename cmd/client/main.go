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

func connectAndSend(msg *common.Message, target *string) error {
	connStr := fmt.Sprintf("%s:%d", clientConfig.RemoteURL, clientConfig.RemotePort)

	conn, err := net.Dial("tcp", connStr)
	if err != nil {
		return err
	}
	defer conn.Close()

	buf := &bytes.Buffer{}
	err = binary.Write(buf, binary.BigEndian, *msg)
	if err != nil {
		return fmt.Errorf("Error converting message to bytes: %s", err)
	}

	_, err = conn.Write(buf.Bytes())
	if err != nil {
		return fmt.Errorf("Error writing to connection: %s", err)
	}

	if target != nil {
		fmt.Printf("Received target: %s\n", *target)
	}

	resp := common.Message{}
	responseBuffer := make([]byte, 41)

	running := true
	for running {

		_, err = conn.Read(responseBuffer)
		if err != nil {
			return err
		}
		err = binary.Read(bytes.NewReader(responseBuffer), binary.BigEndian, &resp)
		if err != nil {
			return fmt.Errorf("Error in binary.Read: %s", err)
		}

		switch resp.Type {
		case common.Text:
			if resp.Length > 0 {
				payload := make([]byte, resp.Length)
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

	var msg *client.InternalClientMessage

	switch command {
	case common.List:
		msg = &client.InternalClientMessage{
			Type: command,
		}
	case common.File:
		msg = &client.InternalClientMessage{
			Type:          common.File,
			LocalFilename: target,
		}
	}

	err = connectAndSend(msg, target)
	if err != nil {
		fmt.Println(err)
	}
}
