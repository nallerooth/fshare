package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"

	"github.com/nallerooth/fshare/internal/config"
	"github.com/nallerooth/fshare/internal/connection"
	"github.com/nallerooth/fshare/internal/message"
)

func readMessage(c net.Conn) (*message.Message, error) {
	msg := &message.Message{}
	err := connection.ReadMessage(c, msg)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

func readUint64(c net.Conn) (uint64, error) {
	resBuf, err := connection.ReadBytes(c, 8)
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

func connectAndSend(conf config.Config, msg *message.Message) error {
	conn, err := connection.NewClientConnection(conf.Client)
	if err != nil {
		return err
	}
	defer conn.Close()

	buf := bytes.Buffer{}
	err = binary.Write(&buf, binary.BigEndian, msg)
	if err != nil {
		return fmt.Errorf("Error writing to connection: %s", err)
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		running := true
		for running {
			msg, err := readMessage(conn)
			if err != nil {
				fmt.Println("go routine error: ", err)
				return
			}

			switch msg.Command {
			case message.Text:
				// Read length of response
				if msg.DataLength > 0 {
					payload := make([]byte, msg.DataLength)
					_, err = conn.Read(payload)
					if err != nil {
						return
					}
					fmt.Println(string(payload))
				}

			case message.Quit:
				fmt.Printf("\nReceived end of transmission\n")
				running = false
				wg.Done()

			default:
				fmt.Println("unknown message type")
			}
		}

	}()

	_, err = conn.Write(buf.Bytes())
	if err != nil {
		return err
	}

	wg.Wait()

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

func createMessageFromArguments() (*message.Message, error) {
	var msg message.Message
	name := ""

	command, err := getCommand()
	if err != nil {
		return nil, err
	}

	switch command {
	case "upload":
		name, err = getTarget()
		if err != nil {
			return nil, err
		}
		msg.Command = message.FileTransfer

	case "list":
		msg = message.NewList()

	case "delete":
		name, err = getTarget()
		if err != nil {
			return nil, err
		}
		msg.Command = message.FileDelete

	case "search":
		name, err = getTarget()
		if err != nil {
			return nil, err
		}
		msg.Command = message.FileSearch

	default:
		return nil, fmt.Errorf("Invalid command '%s'", command)
	}

	// Set msg.Name = name
	// TODO: Make sure to grab the file extension, if available
	copy(msg.Name[:], []byte(name))
	return &msg, nil
}

func main() {
	config := config.Config{
		Passphrase: "",
		Client: config.ClientConfig{
			ServerURL:  "localhost",
			ServerPort: 32000,
		},
	}

	msg, err := createMessageFromArguments()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		// TODO: Print usage
	}

	err = connectAndSend(config, msg)
	if err != nil {
		fmt.Println(err)
	}
}
